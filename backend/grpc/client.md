# gRPC client

## Общая идея

gRPC‑клиент в Go — это:

- один или несколько `grpc.ClientConn` (каналов) к серверу;
- сгенерированные client‑stub’ы поверх этих соединений;
- использование context с таймаутом/дедлайном на каждый вызов.

Клиентский код должен максимально просто использовать контракты из `.proto`.

---

## Создание соединения и клиента

Минимальный пример:

```go
package main

import (
"context"
"log"
"time"

	usersv1 "github.com/example/project/gen/proto/users/v1"
	"google.golang.org/grpc"
)

func main() {
conn, err := grpc.Dial(
"localhost:50051",
grpc.WithInsecure(), // для продакшена — заменить на TLS
grpc.WithBlock(),
)
if err != nil {
log.Fatalf("failed to connect: %v", err)
}
defer conn.Close()

	client := usersv1.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &usersv1.GetUserRequest{Id: 1})
	if err != nil {
		log.Fatalf("GetUser error: %v", err)
	}

	log.Printf("user: %+v\n", resp.GetUser())
}
```

Ключевые моменты:

- `grpc.Dial` создаёт connection; обычно его **переиспользуют** и не создают на каждый запрос;
- `NewUserServiceClient` — сгенерированный конструктор клиента;
- каждый вызов — с `context.Context` и таймаутом.

---

## Работа с streaming‑методами

Пример server‑streaming (клиент читает поток):

```go
stream, err := client.ListUsersStream(ctx, &usersv1.ListUsersRequest{})
if err != nil {
    log.Fatalf("ListUsersStream error: %v", err)
}

for {
    user, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatalf("stream recv error: %v", err)
        }
    log.Printf("user: %+v\n", user)	
}
```

Для client‑ и bidi‑streaming логика похожая, только добавляется отправка в поток через `Send` / `SendMsg`.

---

## Настройки клиента (опции, балансировка, retry)

При подключении можно передавать дополнительные опции:

```go
conn, err := grpc.Dial(
    "localhost:50051",
    grpc.WithTransportCredentials(creds), // TLS/mTLS
    grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
    grpc.WithUnaryInterceptor(clientLoggingInterceptor),
    )
```

Чаще всего на клиенте настраивают:

- TLS / mTLS;
- client‑side load balancing (при работе с несколькими адресами);
- retry через интерсепторы;
- метрики/трейсинг.

---

## Рекомендации по использованию клиента

Кратко:

- не создавать `grpc.ClientConn` на каждый запрос — это тяжёлая операция; держать один connection на процесс/инстанс (или малое фиксированное число);
- обязательно использовать context с таймаутами/дедлайнами;
- логировать и метрики снимать по кодам ошибок (`codes.Unavailable`, `DeadlineExceeded` и т.д.);
- выносить создание клиента в отдельный конструктор, а не размазывать `grpc.Dial` по коду.

Пример конструктора:

```go
func NewUserClient(addr string, creds credentials.TransportCredentials) (usersv1.UserServiceClient, func() error, error) {
conn, err := grpc.Dial(
addr,
grpc.WithTransportCredentials(creds),
)
if err != nil {
return nil, nil, err
}

	closeFn := func() error {
		return conn.Close()
	}

	return usersv1.NewUserServiceClient(conn), closeFn, nil
}
```