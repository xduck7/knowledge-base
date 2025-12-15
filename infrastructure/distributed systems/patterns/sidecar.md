## Sidecar паттерн

Sidecar паттерн — это архитектурный паттерн проектирования, при котором вспомогательный контейнер (sidecar) развертывается вместе с основным приложением для предоставления дополнительной функциональности. Как прицеп едет за основным транспортом, так и sidecar-контейнер работает бок о бок с главным контейнером, разделяя с ним жизненный цикл и сетевое пространство. Особенностью паттерна является независимость, потому что главный контейнер не знает о существовании своего "прицепа", который дополняет его функционал.

## Когда использовать Sidecar

Паттерн особенно полезен в следующих сценариях:

- **Кросс-функциональные задачи** — когда нужно централизовать логирование, мониторинг, безопасность или конфигурацию для нескольких микросервисов
- **Изоляция ответственности** — для отделения вспомогательной функциональности от основной бизнес-логики приложения
- **Service mesh и прокси** — для перехвата и маршрутизации трафика, балансировки нагрузки, service discovery
- **Динамическая конфигурация** — для управления настройками без перезапуска основного приложения
- **Canary deployment** — для постепенного переключения трафика между версиями

Не стоит использовать sidecar в performance-критичных системах, где важна минимальная латентность, и в ресурсо-ограниченных средах.

## Функциональность Sidecar

Sidecar-контейнер может выполнять различные функции:

- Логирование и агрегация логов
- Мониторинг метрик и трейсинг
- Прокси для HTTP/gRPC трафика
- Управление TLS сертификатами и шифрованием
- Service discovery и регистрация сервисов
- Circuit breaker и retry логика
- Конвертация протоколов (например, HTTP в gRPC)

## ASCII схема

```
┌─────────────────────────────────────┐
│           Pod / Host                │
│                                     │
│  ┌─────────────┐   ┌─────────────┐  │
│  │   Main      │   │   Sidecar   │  │
│  │ Application │◄─►│  Container  │  │
│  │             │   │             │  │
│  │    (App)    │   │ (Logging,   │  │
│  │             │   │  Proxy,     │  │
│  │ Port: 8080  │   │  Monitor)   │  │
│  └─────────────┘   └─────────────┘  │
│         ▲                  ▲        │
│         │                  │        │
│         └──shared network──┘        │
│             (localhost)             │
└─────────────────────────────────────┘
                 ▲
                 │
              External
              Requests
```

## Пример Docker Compose

```yaml
version: '3.8'

services:
  # Основное приложение на Go
  main-app:
    build:
      context: ./main-app
      dockerfile: Dockerfile
    container_name: main-app
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
    networks:
      - app-network
    depends_on:
      - sidecar-proxy

  # Sidecar контейнер (например, nginx proxy)
  sidecar-proxy:
    image: nginx:alpine
    container_name: sidecar-proxy
    ports:
      - "8081:8081"
    volumes:
      - ./sidecar/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./sidecar/logs:/var/log/nginx
    networks:
      - app-network

  # Sidecar для логирования (Fluent Bit)
  sidecar-logger:
    image: fluent/fluent-bit:latest
    container_name: sidecar-logger
    volumes:
      - ./sidecar/fluent-bit.conf:/fluent-bit/etc/fluent-bit.conf
      - ./sidecar/logs:/logs
    networks:
      - app-network
    depends_on:
      - main-app

networks:
  app-network:
    driver: bridge
```

Для Go-приложения можно использовать следующий Dockerfile:

```dockerfile
FROM golang:1.22-alpine as builder
WORKDIR /go/app
COPY . .
RUN go build -v -o app cmd/main.go

FROM alpine
COPY --from=builder /go/app/app .
EXPOSE 8080
CMD ["./app"]
```

В этой конфигурации `main-app` — ваше основное приложение, а `sidecar-proxy` и `sidecar-logger` выполняют вспомогательные функции, разделяя общую сеть и могут общаться через localhost
