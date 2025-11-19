# gRPC authentication

> В gRPC аутентификация делится на два уровня: **транспортная защита** (TLS/mTLS) и **прикладочная авторизация** (токены, JWT, API‑ключи и т.д.). Оба уровня обычно комбинируются.

---

## Уровни безопасности в gRPC

В типичном gRPC‑сервисе есть два слоя:

1. **Transport security**  
   - шифрование канала (TLS);
   - аутентификация сервера (клиент убеждается, что говорит с «правильным» хостом);
   - опционально — аутентификация клиента (mTLS).
2. **Application‑level auth**  
   - кто именно пользователь/сервис;
   - какие у него права (roles/permissions);
   - правила доступа к методам/ресурсам.

Transport‑слой закрывает трафик от MITM и «подслушивания», а прикладочный слой отвечает за [translate:who are you] и [translate:what are you allowed to do].

---

## TLS и mTLS

### TLS (аутентификация сервера)

Минимальный вариант — просто защищённый канал:

- сервер предъявляет сертификат, подписанный доверенным CA;
- клиент проверяет сертификат (CN/SAN → host, цепочку доверия);
- весь gRPC‑трафик шифруется.

Go‑пример (сервер, сильно упрощённо):

```go
creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
if err != nil {
log.Fatal(err)
}

s := grpc.NewServer(grpc.Creds(creds))
```

Клиент:

```go
creds, err := credentials.NewClientTLSFromFile("ca.pem", "server.domain")
conn, err := grpc.Dial(
"server.domain:443",
grpc.WithTransportCredentials(creds),
)
```

### mTLS (двусторонняя аутентификация)

В mTLS обе стороны имеют сертификат:

- сервер проверяет клиентский сертификат (кто это);
- клиент проверяет сервер.

На практике:

- сертификаты клиентов часто выдаёт отдельный CA;
- идентичность клиента (service name, spiffe‑id и т.д.) вытаскивается из поля Subject/SubjectAltName;
- эта идентичность используется как «subject» для authorisation (RBAC).

mTLS хорошо подходит для **внутри‑кластерного сервис‑to‑service** трафика.

---

## Токены, JWT и метаданные

gRPC опирается на HTTP/2, а значит есть **метаданные** (headers). Типичные варианты:

- заголовок [translate:authorization]: [translate:Bearer <jwt>];
- кастомные заголовки (metadata key): [translate:x-api-key], [translate:x-user-id], [translate:x-service-token].

### Клиент: как передать токен

Пример на Go с JWT токеном:

```go
md := metadata.Pairs("authorization", "Bearer "+token)
ctx := metadata.NewOutgoingContext(ctx, md)

resp, err := client.GetUser(ctx, &usersv1.GetUserRequest{Id: 1})
```

Либо через пер‑call credentials (обёртка, которая автоматически добавляет токен ко всем вызовам клиента).

### Сервер: как вытащить токен

На сервере токен читается из `metadata.MD`:

```go
func (s *userServer) GetUser(ctx context.Context, req *usersv1.GetUserRequest) (*usersv1.GetUserResponse, error) {
md, ok := metadata.FromIncomingContext(ctx)
if !ok {
return nil, status.Error(codes.Unauthenticated, "missing metadata")
}

    authHeader := md.Get("authorization")
    if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Bearer ") {
        return nil, status.Error(codes.Unauthenticated, "invalid authorization header")
    }

    token := strings.TrimPrefix(authHeader, "Bearer ")

    // дальше парс/верификация JWT / opaque токена
    // subject, roles, expiry и т.д.

    return s.handleGetUser(ctx, req, /*subject*/ )
}
```

В реальном проекте эту логику лучше вынести в **interceptor**, чтобы не копипастить по всем методам.

---

## Interceptors для auth

### Unary‑интерсептор

Общий паттерн:

- на входе:
  - достать токен/метаданные;
  - валидировать (подпись, expiry, audience);
  - построить объект [translate:principal]/[translate:subject] и положить в `context.Context`;
- если всё плохо — вернуть `Unauthenticated` или `PermissionDenied`;
- вызвать `handler(ctx, req)` с обогащённым контекстом.

Условный пример:

```go
func authUnaryInterceptor(
ctx context.Context,
req any,
info *grpc.UnaryServerInfo,
handler grpc.UnaryHandler,
) (any, error) {
md, ok := metadata.FromIncomingContext(ctx)
if !ok {
return nil, status.Error(codes.Unauthenticated, "missing metadata")
}

    token := extractToken(md) // достать Bearer / api-key

    principal, err := validateToken(token) // JWT / opaque / custom
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid token")
    }

    // положить principal в контекст
    ctx = context.WithValue(ctx, principalKey{}, principal)

    return handler(ctx, req)
}
```

Дальше в бизнес‑коде можно достать principal и делать проверку прав.

---

## Auth vs AuthZ (RBAC/ABAC)

Важно разделять:

- **Authentication** — кто это? (subject: user/service, id, provider);
- **Authorization** — что ему можно? (roles, permissions, ACL, policy).

Типичный pipeline:

1. Интерсептор аутентифицирует:
   - разбирает токен / клиентский сертификат;
   - формирует subject (user id, service id, scopes).
2. Дальше:
   - либо тот же интерсептор делает RBAC (по таблице «метод → роль/permission»);
   - либо отдельный слой авторизации (policy‑engine).

Например, можно сделать mapping:

- `/user.v1.UserService/GetUser` → роль [translate:USER_READ];
- `/user.v1.UserService/DeleteUser` → роль [translate:ADMIN].

И в интерсепторе авторизации проверять:

- есть ли у principal нужная роль/permission для `info.FullMethod`.

---

## Варианты аутентификации

### 1. mTLS между сервисами

Используется, когда:

- все клиенты — сервисы внутри кластера;
- сертификаты выдаёт internal CA;
- идентичность сервиса → из сертификата.

Плюсы:

- нет токенов, только сертификаты;
- хорошо интегрируется с сервис‑мэшем (Istio, Linkerd, etc.).

Минусы:

- сложность управления сертификатами (ротация, revocation);
- неудобно для «человеческих» клиентов (браузеры, мобильные).

---

### 2. JWT/OAuth2 токены

Классическая схема:

- отдельный auth‑сервер выдаёт JWT access‑token;
- клиент gRPC передаёт его как [translate:authorization: Bearer <token>];
- gRPC‑сервер:
  - валидирует подпись (по ключу/ключам IdP);
  - проверяет [translate:exp], [translate:aud], [translate:iss], [translate:sub]/[translate:scope]/roles.

Плюсы:

- не нужно ходить в storage для каждого запроса (self‑contained токен);
- легко отдавать идентичность/роли сервисам.

Минусы:

- забота о ротации ключей;
- внимание к размеру и содержимому токена.

---

### 3. API‑ключи / статические токены

Просто:

- клиент передаёт [translate:x-api-key];
- сервер проверяет ключ по таблице/хранилищу.

Подходит для:

- простых внутренних интеграций;
- машин‑to‑машин без сложного протокола.

Минусы:

- сложно управлять правами/ролями;
- revoke/rotation часто делаются «вручную».

---

## Практические рекомендации

Кратко:

- **Всегда** включать TLS хотя бы на внешней границе (ingress, публичный API).
- Для **внутрисервисного** трафика:
  - рассмотреть mTLS через сервис‑мэш или вручную;
  - использовать идентичность из сертификатов для authZ.
- Для **user‑level auth**:
  - брать токены от внешнего IdP (OAuth2/OIDC) и проверять их в gRPC‑интерсепторе;
  - не хранить пароли прямо в gRPC‑сервисах.
- Разделять обязанности:
  - один интерсептор отвечает за аутентификацию и наполнение контекста;
  - другой — за авторизацию (RBAC/ABAC), логирование попыток и метрики.
