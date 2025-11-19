# gRPC error handling

> В gRPC ошибка — это **не поле в ответе**, а отдельный канал: RPC либо успешно возвращает ответ, либо завершается с ошибкой статуса.

---

## Базовая модель ошибок в gRPC

В gRPC каждый RPC‑вызов заканчивается:

- **успешно** — статус `OK` и непустой ответ;
- **с ошибкой** — статус с кодом (`codes.*`) и сообщением, а тело ответа при этом отсутствует/игнорируется.

Ключевые моменты:

- ошибки кодируются в статусе (metadata/трейлеры), а не в `message` тела;
- в Go все ошибки с клиента приходят как `error`, который можно преобразовать в `status.Status`.

---

## Статусы и коды (`codes.*`)

Основные коды из `google.golang.org/grpc/codes`:

- `codes.InvalidArgument` — невалидный ввод (валидация, формат данных);
- `codes.NotFound` — ресурс не найден;
- `codes.AlreadyExists` — конфликт при создании (уже существует);
- `codes.PermissionDenied` / `codes.Unauthenticated` — нет прав / не авторизован;
- `codes.FailedPrecondition` — неверное состояние (нельзя выполнить в текущем контексте);
- `codes.Unavailable` — временная недоступность (сетевые проблемы, перезапуск);
- `codes.DeadlineExceeded` — превышен таймаут;
- `codes.Internal` — внутренняя ошибка сервера (panic, баг, неожиданные состояния).

Рекомендуется:

- использовать максимально **точный** код, а не всегда `Internal`;
- заранее договориться о том, какие коды считаются **retryable** (обычно `Unavailable`, `DeadlineExceeded`).

---

## Сервер: возврат ошибок в Go

На сервере методы имеют сигнатуру вида:

```go
func (s *userServer) GetUser(ctx context.Context, req *usersv1.GetUserRequest) (*usersv1.GetUserResponse, error) {
// ...
}
```

Правило: **либо успешный ответ, либо ошибка**, но не оба сразу.

Примеры:

```go
import (
"fmt"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func (s *userServer) GetUser(ctx context.Context, req *usersv1.GetUserRequest) (*usersv1.GetUserResponse, error) {
if req.GetId() <= 0 {
return nil, status.Errorf(codes.InvalidArgument, "id must be positive")
}

    user, err := s.repo.GetByID(ctx, req.GetId())
    if errors.Is(err, repository.ErrNotFound) {
        return nil, status.Errorf(codes.NotFound, "user %d not found", req.GetId())
    }
    if err != nil {
        return nil, status.Errorf(codes.Internal, "get user: %v", err)
    }

    return &usersv1.GetUserResponse{User: toProtoUser(user)}, nil
}
```

Анти‑паттерн:

- возвращать «ошибку» внутри поля ответа (например, `error_code`/`error_message`) и одновременно статус `OK`.

---

## Клиент: разбор ошибок

На клиенте, после RPC‑вызова:

```go
resp, err := client.GetUser(ctx, &usersv1.GetUserRequest{Id: 1})
if err != nil {
st, ok := status.FromError(err)
if !ok {
// не gRPC-ошибка (что-то совсем низкоуровневое)
log.Printf("unknown error: %v", err)
return
}

    switch st.Code() {
    case codes.NotFound:
        log.Printf("user not found: %v", st.Message())
    case codes.InvalidArgument:
        log.Printf("invalid argument: %v", st.Message())
    case codes.DeadlineExceeded:
        log.Printf("timeout: %v", st.Message())
    default:
        log.Printf("rpc error (%v): %v", st.Code(), st.Message())
    }
    return
}

log.Printf("user: %+v", resp.GetUser())
```

Лучше:

- централизовать обработку ошибок (например, через обёртку над клиентом или интерсептор);
- разделять ошибки пользователя (валидация, права) и ошибки инфраструктуры (таймаут, недоступность).

---

## Расширенная модель: `ErrorDetails`

Иногда одного кода/сообщения мало — нужны подробности (валидация, поля, локации). Для этого можно использовать `status.WithDetails` и стандартные типы `google.rpc.*` (или свои).

Пример сервера с деталями:

```go
import (
"fmt"

    "google.golang.org/genproto/googleapis/rpc/errdetails"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func validateCreateUser(req *usersv1.CreateUserRequest) error {
var violations []*errdetails.BadRequest_FieldViolation

    if req.GetEmail() == "" {
        violations = append(violations, &errdetails.BadRequest_FieldViolation{
            Field:       "email",
            Description: "must not be empty",
        })
    }

    if len(violations) == 0 {
        return nil
    }

    st := status.New(codes.InvalidArgument, "validation failed")
    br := &errdetails.BadRequest{FieldViolations: violations}

    stWithDetails, err := st.WithDetails(br)
    if err != nil {
        // если не получилось добавить детали, возвращаем базовый статус
        return st.Err()
    }

    return stWithDetails.Err()
}
```

Клиент может читать эти детали:

```go
err := client.CreateUser(ctx, req)
if err != nil {
st := status.Convert(err)

    for _, d := range st.Details() {
        switch info := d.(type) {
        case *errdetails.BadRequest:
            for _, v := range info.FieldViolations {
                log.Printf("validation error: field=%s desc=%s", v.Field, v.Description)
            }
        default:
            // другие типы деталей
        }
    }
}
```

Это позволяет:

- не выдумывать свои «обёртки» в полях сообщений;
- передавать структурированные ошибки между сервисами и до клиентов.

---

## Ошибки и retry‑логика

Часть кодов подразумевает, что **retry бессмысленен**, часть — что он допустим или ожидаем.

Чаще всего:

- **не retry**:
  - `InvalidArgument`, `NotFound`, `AlreadyExists`, `PermissionDenied`, `Unauthenticated` — проблема на стороне запроса/прав;
- **можно retry**:
  - `Unavailable` — временная недоступность, сетевые глюки, рестарт;
  - `DeadlineExceeded` — иногда (если это был внешний вызов, который мог успеть со второй попытки);
- **зависит от контекста**:
  - `Aborted`, `ResourceExhausted`, `FailedPrecondition`.

Хорошая практика:

- оформить политику retry‑код в одном месте (конфиг/хелпер);
- не делать retry бесконечным;
- для небезопасных операций (не‑идемпотентных) быть особенно осторожным (см. idempotency в `best_practices.md`).

---

## Mapping бизнес‑ошибок на gRPC‑коды

Типичный слой преобразования:

```go
func toStatusError(err error) error {
if err == nil {
return nil
}

    if errors.Is(err, repository.ErrNotFound) {
        return status.Error(codes.NotFound, "resource not found")
    }

    var vErr *ValidationError
    if errors.As(err, &vErr) {
        st := status.New(codes.InvalidArgument, "validation failed")
        // можно добавить details по полям
        return st.Err()
    }

    // по умолчанию
    return status.Error(codes.Internal, "internal error")
}
```

На сервере:

```go
if err != nil {
return nil, toStatusError(err)
}
```

Это:

- отделяет бизнес‑ошибки от транспорта;
- позволяет централизованно менять политику кодов и сообщений.

---

## Ошибки, deadlines и cancellations

Контекст в gRPC важен не только для таймаутов, но и для обработки отмен:

- клиент устанавливает дедлайн (`context.WithTimeout` / `WithDeadline`);
- при истечении времени вызов на клиенте заканчивается с `DeadlineExceeded`;
- сервер получает контекст с `ctx.Err() == context.DeadlineExceeded` или `context.Canceled`.

На сервере стоит:

```go
if err := ctx.Err(); err != nil {
// клиент уже ушёл / истёк таймаут — дальше работать бессмысленно
return nil, status.Error(codes.Canceled, "request canceled")
}
```

или просто прекращать работу и возвращать ошибку, если это уместно.

---

## Interceptors и логирование ошибок

Часто ошибки логируют централизованно через интерсепторы.

Пример unary‑интерсептора:

```go
func loggingUnaryInterceptor(
    ctx context.Context,
    req any,
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (any, error) {
    resp, err := handler(ctx, req)
    if err != nil {
        st := status.Convert(err)
        // логируем метод, код, сообщение
        log.Printf("grpc error: method=%s code=%s msg=%s", info.FullMethod, st.Code(), st.Message())
	}
        return resp, err
}
```

Так:

- бизнес‑код не захламлён логированием;
- все ошибки проходят через один фильтр и могут быть дополнены контекстом (метод, метаданные, user‑id и т.п.).

---

## Краткий чек‑лист по error handling

- Не шьём «ошибки» внутрь message‑ответов — используем статус.
- На сервере:
  - возвращаем либо `(response, nil)`, либо `(nil, error)`;
  - маппим доменные ошибки на корректные `codes.*`;
  - по возможности используем `ErrorDetails` для структурированных ошибок.
- На клиенте:
  - всегда проверяем `err` и разбираем его через `status.FromError` / `status.Convert`;
  - различаем клиентские ошибки (валидация, права) и инфраструктурные (таймаут, недоступность);
  - реализуем осмысленную retry‑политику.
- Логирование и метрики по ошибкам — через интерсепторы, а не точечно в каждом методе.
