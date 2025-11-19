# gRPC health check

> Health‑check в gRPC — это отдельный сервис/метод, по которому инфраструктура (Kubernetes, балансировщики, сервис‑мэш) может понять, «жив» ли наш сервис и готов ли он принимать трафик.

---

## Зачем нужен health‑check

Health‑check используется для:

- liveness: процесс жив, не завис, не ушёл в unrecoverable state;
- readiness: сервис готов принимать запросы (подключился к БД, поднял кэш, прогрелся);
- integration: Kubernetes, Envoy, Consul, service‑mesh и прочие инструменты могут автоматически выключать/включать инстансы из балансировки.

Важно разделять:

- **TCP‑живость** (порт открыт) ≠ **приложение работает корректно**;
- gRPC‑health даёт более точный сигнал на уровне бизнес‑логики и зависимостей.

---

## Стандартный gRPC Health Checking Protocol

gRPC определяет общий контракт для health‑сертиса:

```
syntax = "proto3";

package grpc.health.v1;

message HealthCheckRequest {
string service = 1;
}

message HealthCheckResponse {
enum ServingStatus {
UNKNOWN = 0;
SERVING = 1;
NOT_SERVING = 2;
SERVICE_UNKNOWN = 3; // для Watch
}
ServingStatus status = 1;
}

service Health {
rpc Check(HealthCheckRequest) returns (HealthCheckResponse);
rpc Watch(HealthCheckRequest) returns (stream HealthCheckResponse);
}
```

Идея:

- `Check` — разовый запрос: «здоров ли сервис X?»;
- `Watch` — стрим обновлений статуса.

В Go этот контракт уже сгенерирован и доступен как `google.golang.org/grpc/health/grpc_health_v1`.

---

## Быстрый старт: health‑сервер в Go

Базовая интеграция с gRPC‑health в сервере:

```
import (
"google.golang.org/grpc"
"google.golang.org/grpc/health"
healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
s := grpc.NewServer()

	// Создаём health-сервер
	healthServer := health.NewServer()

	// Регистрируем его в gRPC-сервере
	healthpb.RegisterHealthServer(s, healthServer)

	// Регистрация твоих gRPC-сервисов ниже:
	// usersv1.RegisterUserServiceServer(s, userServer)

	// Устанавливаем статус для всего приложения (пустая строка = "system")
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// ... запустить s.Serve(lis)
}
```

Что важно:

- health‑сервер регистрируется как обычный gRPC‑сервис;
- статус можно задавать отдельно для всего приложения и для конкретных сервисов.

---

## Пер‑сервисный статус

Если в инстансе несколько gRPC‑сервисов, можно отслеживать их по отдельности:

```
healthServer.SetServingStatus("users.v1.UserService", healthpb.HealthCheckResponse_SERVING)
healthServer.SetServingStatus("billing.v1.BillingService", healthpb.HealthCheckResponse_NOT_SERVING)
```

Здесь `"users.v1.UserService"` — это service name, который клиент укажет в `HealthCheckRequest.service`.

Паттерн:

- при старте:
  - устанавливаем `NOT_SERVING`,
  - поднимаем зависимости (БД, кэш, внешние API),
  - при успешной инициализации меняем на `SERVING`;
- при критической ошибке/потере зависимости:
  - переключаем статус обратно на `NOT_SERVING`.

---

## Readiness / liveness для Kubernetes

Частый вариант:

- liveness‑probe — простой TCP/HTTP‑пинг (процесс жив);
- readiness‑probe — gRPC‑health‑check.

Пример HTTP‑readiness через `grpc-health-probe` (готовый бинарь, который делает gRPC‑запрос `Health.Check`):

```
readinessProbe:
exec:
command:
- /bin/grpc-health-probe
- -addr=:50051
- -service=users.v1.UserService
initialDelaySeconds: 5
periodSeconds: 10
```

Плюсы:

- Kubernetes не начнёт слать трафик, пока health не в статусе `SERVING`;
- при падении БД можно менять статус на `NOT_SERVING` и инстанс выпадает из балансировки.

---

## Собственный health‑метод вместо стандарта

Иногда достаточно простого, «кастомного» health‑метода в своём сервисе:

```
service UserService {
rpc Health(google.protobuf.Empty) returns (HealthResponse);
}

message HealthResponse {
bool ok = 1;
string message = 2;
}
```

Плюсы:

- можно вернуть больше контекста (например, состояние зависимостей);
- не нужно подключать отдельный health‑сервис.

Минусы:

- «нестандартный» контракт, tooling вокруг gRPC‑health использовать сложнее (grpc‑health‑probe и т.д.);
- меньше совместимости с готовыми решениями.

Практика: для инфраструктуры — стандартный gRPC‑health, для людей/дебага — дополнительный кастомный `/debug/health` или RPC.

---

## Что проверять внутри health‑чека

Минимум:

- возможность сделать простой запрос к основной БД;
- состояние критичных внешних сервисов (например, очередь сообщений).

Варианты:

- быстрый health: только базовая самопроверка (все зависимости не проверяются на каждый запрос);
- «глубокий» health по отдельному флагу/методу — с полной проверкой зависимостей.

Важно: health‑чек не должен сам становиться «точкой DDOS» по зависимостям (частые запросы, тяжёлые проверки).

---

## Пример: простая обвязка поверх health‑сервера

Небольшой хелпер:

```
type HealthManager struct {
server *health.Server
}

func NewHealthManager() *HealthManager {
return &HealthManager{server: health.NewServer()}
}

func (h *HealthManager) Register(grpcServer *grpc.Server) {
healthpb.RegisterHealthServer(grpcServer, h.server)
}

func (h *HealthManager) SetReady() {
h.server.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
}

func (h *HealthManager) SetNotReady() {
h.server.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
}
```

Использование:

```
healthMgr := NewHealthManager()
healthMgr.Register(s)

// после инициализации зависимостей
healthMgr.SetReady()

// при фатальном сбое
// healthMgr.SetNotReady()
```

Такой слой удобно переиспользовать во всех сервисах.