# gRPC server

## Общая идея

gRPC‑сервер в Go — это:

- реализация интерфейсов, сгенерированных из `.proto`;
- экземпляр `grpc.Server` с зарегистрированными сервисами;
- сетевой `Listener`, на котором сервер принимает соединения (обычно TCP).

Ниже минимальный пример и несколько практических советов.

---

## Базовый пример сервера

Предположим, у тебя есть сгенерированный код для `UserService` в пакете `usersv1`.

```go
package main

import (
"log"
"net"

	usersv1 "github.com/example/project/gen/proto/users/v1"
	"google.golang.org/grpc"
)

type userServer struct {
usersv1.UnimplementedUserServiceServer
// зависимости: repo, logger, cfg и т.п.
}

func main() {
lis, err := net.Listen("tcp", ":50051")
if err != nil {
log.Fatalf("failed to listen: %v", err)
}

	s := grpc.NewServer(
		// сюда потом добавятся interceptors, options и т.д.
	)

	usersv1.RegisterUserServiceServer(s, &userServer{})

	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

Ключевые шаги:

- создать `net.Listener` на нужном порту;
- создать `grpc.Server`;
- зарегистрировать реализацию сервиса;
- запустить `Serve`.

---

## Реализация методов сервиса

Сигнатуры берутся из сгенерированного интерфейса.

```go
func (s *userServer) GetUser(ctx context.Context, req *usersv1.GetUserRequest) (*usersv1.GetUserResponse, error) {
// 1. Достать данные из req
// 2. Вызвать бизнес-логику / репозиторий
// 3. Вернуть ответ в виде protobuf-сообщения

	user := &usersv1.User{
		Id:    req.GetId(),
		Name:  "John Doe",
		Email: "john@example.com",
	}

	return &usersv1.GetUserResponse{User: user}, nil
}
```

Рекомендация: не класть в серверную структуру всё подряд, а инжектить зависимости (сервисы/репозитории) в конструктор `userServer`.

---

## Настройка сервера (options, TLS, interceptors)

Чаще всего сервер создаётся не просто через `grpc.NewServer()`, а с опциями:

```go
s := grpc.NewServer(
grpc.ChainUnaryInterceptor(
authUnaryInterceptor,
loggingUnaryInterceptor,
),
// grpc.Creds(credentials.NewTLS(tlsConfig)),
// другие опции…
)
```

Типичные настройки:

- TLS / mTLS для транспорта;
- unary/stream interceptors (логирование, метрики, трейсинг, auth);
- лимиты по размеру сообщений, keepalive и т.д.

Детали лучше раскрывать уже в отдельных файлах `auth.md`, `interceptors.md`, `observability.md`.

---

## Graceful shutdown

Для прод‑проекта важно корректно останавливать сервер:

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

go func() {
    if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
    }
}()

<-ctx.Done()
    stop()
    log.Println("shutting down gRPC server...")
    s.GracefulStop()
```

Идея: дождаться завершения активных запросов и не обрывать их посередине.