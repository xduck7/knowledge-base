## Ambassador паттерн

Ambassador (Амбассадор) — это архитектурный паттерн для микросервисов, который размещает вспомогательный прокси-сервис между клиентским приложением и внешними сервисами для управления коммуникацией. Паттерн является частным случаем Sidecar и специализируется на работе с внешними сервисами, доступность которых не является стабильной.

## Когда использовать Ambassador

Ambassador особенно полезен в следующих сценариях:

- **Работа с внешними сервисами** — когда микросервисы используют разные протоколы или обновляются в разное время
- **Нестабильная сеть** — для управления сетевыми сбоями, retry логикой и таймаутами при обращении к внешним ресурсам
- **Кросс-сетевые вызовы** — при необходимости вызовов внешних сервисов через ненадежные сети
- **Service Discovery** — для динамического обнаружения сервисов в децентрализованной среде
- **Унифицированный контроль доступа** — когда нужна централизованная аутентификация и авторизация для множества сервисов

## Функциональность Ambassador

Ambassador может выполнять следующие задачи подключения к удаленным сервисам:

- Мониторинг и логирование запросов/ответов
- Маршрутизация и балансировка нагрузки между экземплярами сервисов
- Обеспечение безопасности (TLS, проверка токенов, аутентификация)
- Трансляция протоколов между различными форматами
- Retry механизм и Circuit Breaker для устойчивости к сбоям
- Кэширование ответов и батчинг запросов
- Rate limiting и request throttling для API Management

## ASCII схема

```
┌────────────────────────────────────────────┐
│              Application Pod               │
│                                            │
│  ┌──────────────┐      ┌───────────────┐   │
│  │    Main      │      │  Ambassador   │   │
│  │  Application │─────►│    Proxy      │   │
│  │              │      │               │   │
│  │  (Go App)    │      │ - Monitoring  │   │
│  │              │      │ - Retry       │   │
│  │              │      │ - Security    │   │
│  └──────────────┘      └───────┬───────┘   │
│                                │           │
└────────────────────────────────┼───────────┘
                                 │
                                 │ Network
                                 ▼
                    ┌──────────────────────────┐
                    │   External Services      │
                    │                          │
                    │  - Legacy Systems        │
                    │  - Third-party APIs      │
                    │  - Remote Databases      │
                    └──────────────────────────┘
```

## Пример Docker Compose

```yaml
version: '3.8'

services:
  main-app:
    build:
      context: ./app
      dockerfile: Dockerfile
    container_name: go-main-app
    environment:
      - AMBASSADOR_URL=http://ambassador:8081
    networks:
      - app-network
    depends_on:
      - ambassador

  ambassador:
    build:
      context: ./ambassador
      dockerfile: Dockerfile
    container_name: ambassador-proxy
    ports:
      - "8081:8081"
    environment:
      - EXTERNAL_API_URL=https://api.external-service.com
      - RETRY_ATTEMPTS=3
      - TIMEOUT_SECONDS=10
    volumes:
      - ./ambassador/config.yaml:/etc/ambassador/config.yaml:ro
      - ./logs:/var/log/ambassador
    networks:
      - app-network
      - external-network

  external-service:
    image: httpbin/httpbin
    container_name: external-service-mock
    ports:
      - "8082:80"
    networks:
      - external-network

networks:
  app-network:
    driver: bridge
  external-network:
    driver: bridge
```

Пример простого Ambassador на Go:

```go
package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "time"
)

type Ambassador struct {
    targetURL *url.URL
    proxy     *httputil.ReverseProxy
}

func NewAmbassador(target string) *Ambassador {
    targetURL, _ := url.Parse(target)
    return &Ambassador{
        targetURL: targetURL,
        proxy:     httputil.NewSingleHostReverseProxy(targetURL),
    }
}

func (a *Ambassador) Handle(w http.ResponseWriter, r *http.Request) {
    log.Printf("Request to: %s %s", r.Method, r.URL.Path)
    
    r.Header.Set("X-Forwarded-Host", r.Host)
    
    a.proxy.ServeHTTP(w, r)
}

func main() {
    ambassador := NewAmbassador("http://external-service:80")
    
    server := &http.Server{
        Addr:         ":8081",
        Handler:      http.HandlerFunc(ambassador.Handle),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    
    log.Println("Ambassador starting on :8081")
    log.Fatal(server.ListenAndServe())
}
```

В этой конфигурации `main-app` обращается только к `ambassador`, который берет на себя всю сложность взаимодействия с внешними сервисами, включая retry, мониторинг и безопасность
