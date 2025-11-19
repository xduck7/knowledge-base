## Зачем нужен Docker Compose

Одиночный контейнер можно поднять одной командой `docker run`.
Однако реальное backend‑приложение почти никогда не ограничивается одним процессом: ему нужны база данных, кеш, брокер сообщений, миграции, метрики и вспомогательные утилиты.

Docker Compose решает эту проблему.
Он позволяет описать многоконтейнерное приложение декларативно — в одном YAML‑файле — и управлять им как единой системой: запускать, останавливать, пересоздавать и обновлять.

---

## Концепция: проект и сервисы

Compose оперирует понятием «проекта».
Проект — это всё, что описано в одном `docker-compose.yml` (или `compose.yaml`): набор сервисов, сетей, томов и конфигураций.

Сервис — это логическая единица приложения, обычно один тип контейнера.
Например, в типичном backend‑проекте можно выделить сервисы `api`, `postgres`, `redis`, `mq`, `prometheus`, `migrations`.

Один сервис может порождать несколько контейнеров (при масштабировании), но в конфигурации он описывается один раз.

---

## Базовая структура compose‑файла

Минимальный compose‑файл версия 3 выглядит так:

```sh
version: "3.9"

services:
  app:
    image: my-app:latest
    ports:
      - "8080:8080"
```

Три ключевые части:

- `version` — версия синтаксиса Compose (в линейке 3.x используется для Swarm‑совместимого формата);
- `services` — определение сервисов (контейнеров);
- `app` — имя сервиса и его конфигурация.

---

## Сервисы

Сервис описывает, как запускать контейнер: из какого образа, с какими переменными окружения, томами и портами.

Пример сервиса `api`:

```sh
services:
  api:
    build: ./api
    environment:
      - ENV=local
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
```

Здесь:

- `build` указывает, что образ нужно собрать из директории `./api`;
- `environment` описывает переменные окружения;
- `ports` публикует порт 8080 наружу;
- `depends_on` задаёт зависимость от других сервисов проекта.

---

## Сети

Compose создаёт отдельную сеть для проекта по умолчанию.
Все сервисы подключаются к этой сети и получают имена, совпадающие с именами сервисов, — это встроенный DNS и service discovery.

В простом случае достаточно дефолтной сети.
При необходимости можно явно описывать сети:

```sh
services:
  api:
    networks:
      - backend

  worker:
    networks:
      - backend

networks:
  backend:
    driver: bridge
```

Все сервисы, подключённые к сети `backend`, могут обращаться друг к другу по имени сервиса: `api`, `worker` и т.д.

---

## Томa (volumes)

Тома позволяют сохранять данные вне жизненного цикла контейнера: базы данных, очереди, файлы, кеш и т.п.
Для backend‑приложений это особенно важно: потеря данных в контейнере после пересоздания неприемлема.

Compose даёт два уровня работы с томами:

- объявление именованных томов на верхнем уровне;
- подключение этих томов к сервисам.

Пример:

```sh
services:
  postgres:
    image: postgres:16
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

Том `postgres_data` будет создан и привязан к каталогу данных Postgres.
Пересоздание контейнера не повлияет на содержимое тома.

---

## Env‑файлы и переменные окружения

Переменные окружения — стандартный способ конфигурировать приложения в контейнерах.
Compose позволяет задавать их несколькими способами:

- напрямую в `environment`;
- через `env_file`;
- через подстановку переменных окружения хоста.

Пример использования `env_file`:

```sh
services:
  api:
    build: ./api
    env_file:
      - .env
```

Файл `.env`:

```sh
ENV=local
DB_DSN=postgres://postgres:secret@postgres:5432/app?sslmode=disable
REDIS_ADDR=redis:6379
```

Такое разделение позволяет хранить конфигурацию отдельно от кода и компоуз‑файла, менять секреты и параметры окружения без изменений в YAML.

---

## Управление жизненным циклом

Основные команды:

- `docker compose up` — поднять проект (создать сети, тома, контейнеры);
- `docker compose up -d` — то же, но в фоне;
- `docker compose down` — остановить и удалить контейнеры и сети проекта;
- `docker compose logs -f` — смотреть логи всех сервисов или выбранного;
- `docker compose ps` — список контейнеров проекта.

Compose превращает набор `docker run` в один декларативный описательный файл.
Это фундамент для последующих паттернов локальной разработки и staging/prod‑подобных окружений.


## Версия файла и формат

Compose‑файлы версии 3.x ориентированы на совместимость с Docker Swarm и широко используются для описания многоконтейнерных приложений.
Несмотря на наличие более новых особенностей, `version: "3.8"` или `version: "3.9"` остаются удобным и распространённым форматом.

Общий каркас:

```sh
version: "3.9"

services:
  ...

volumes:
  ...

networks:
  ...
```

Далее будут рассматриваться ключевые поля, необходимые backend‑разработчику для уверённой работы с Compose.

---

## services

Секция `services` описывает все сервисы (контейнеры), входящие в проект.

Пример:

```sh
services:
  api:
    build: ./api
    ports:
      - "8080:8080"
    environment:
      ENV: local
    depends_on:
      - postgres

  postgres:
    image: postgres:16
    volumes:
      - postgres_data:/var/lib/postgresql/data
```

Поля сервиса часто используемые в backend‑проектах:

- `image` / `build` — образ или инструкция по его сборке;
- `ports` — публикация портов;
- `environment` / `env_file` — переменные окружения;
- `volumes` — тома и bind‑mount’ы;
- `depends_on` — декларативные зависимости;
- `healthcheck` — проверка «здоровья» сервиса;
- `deploy` — параметры развёртывания (актуально для Swarm).

---

## depends_on

`depends_on` описывает, какие сервисы должны быть созданы и запущены раньше.

Простейшая форма:

```sh
services:
  api:
    depends_on:
      - postgres
      - redis
```

Этот вариант гарантирует порядок запуска: сначала `postgres` и `redis`, потом `api`.
Однако он не гарантирует готовность сервиса (например, Postgres может запускаться дольше).

В версиях Compose 3.9+ появилась расширенная форма для учёта `healthcheck`:

```sh
services:
  api:
    depends_on:
      postgres:
        condition: service_healthy
```

В этом случае `api` будет ждать, пока Postgres станет «здоровым» по `healthcheck`.

---

## healthcheck

`healthcheck` позволяет описать команду, по которой Docker будет регулярно проверять состояние контейнера.

Пример для базы данных:

```sh
services:
  postgres:
    image: postgres:16
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s
```

Поля:

- `test` — команда проверки;
- `interval` — интервал между проверками;
- `timeout` — максимальное время ожидания ответа команды;
- `retries` — количество неудачных попыток до признания контейнера «нездоровым»;
- `start_period` — время «разогрева» перед началом подсчёта неудачных проверок.

Статус healthcheck можно увидеть через `docker ps` или `docker inspect`.
Связка healthcheck + `depends_on: condition: service_healthy` особенно ценна для backend‑стеков с зависимостями.

---

## deploy

Секция `deploy` используется в основном для Swarm‑кластеров и не применяется локально в классическом `docker compose up`.
Тем не менее, понимание основных полей полезно для общего кругозора.

Пример:

```sh
services:
  api:
    image: my-api:1.0.0
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: "1.0"
          memory: "512M"
      restart_policy:
        condition: on-failure
```

Ключевые элементы:

- `replicas` — число копий сервиса;
- `resources` — лимиты ресурсов;
- `restart_policy` — политика рестартов в Swarm‑контексте.

В локальной разработке роль `deploy` часто переоценивают.
Для обычных `docker compose up` важнее грамотно настроить `healthcheck`, `depends_on`, сети и тома.

---

## volumes и networks на верхнем уровне

Топ‑уровневые секции `volumes` и `networks` позволяют переиспользовать их между сервисами и настраивать драйверы и параметры.

Пример томов:

```sh
volumes:
  postgres_data:
  redis_data:
```

Пример сетей:

```sh
networks:
  backend:
    driver: bridge
```

Использование в сервисах:

```sh
services:
  postgres:
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - backend

  redis:
    volumes:
      - redis_data:/data
    networks:
      - backend
```

Такое разделение делает конфигурацию явной и переносимой.

---

## Чтение и валидация конфигурации

Compose умеет показывать итоговую конфигурацию с учётом всех подстановок и override‑файлов:

```sh
docker compose config
```

Эта команда:

- валидирует YAML;
- показывает развёрнутую конфигурацию;
- помогает диагностировать ошибки до запуска.

Для серьёзных проектов полезно приучаться смотреть на результат `docker compose config` так же внимательно, как на вывод `go test`.

## Цель: комфортная локальная разработка

Compose удобно использовать не только для «мини‑прода», но и для повседневной разработки.
Он позволяет запускать всю инфраструктуру разработчика одной командой, при этом оставляя сам код на машине, чтобы работали hot reload, IDE и отладка.

Локальный dev‑стенд обычно включает:

- приложение (API);
- базу данных;
- кеш;
- брокер сообщений (опционально);
- вспомогательные утилиты (миграции, админки).

---

## Маппинг кода (bind mounts)

Главная идея локального стенда: контейнеры отвечают за окружение, а код при этом остаётся в рабочей директории разработчика.
Для этого используется bind‑mount исходников.

Пример:

```sh
services:
  api:
    build: ./api
    volumes:
      - ./api:/app
    working_dir: /app
    command: air
    ports:
      - "8080:8080"
```

Здесь:

- `./api` монтируется внутрь контейнера в `/app`;
- релодер (`air`, `reflex`, `nodemon` и т.п.) отслеживает изменения файлов;
- IDE работает с кодом напрямую на хосте, без потерь по DX.

---

## Hot reload

Для языков вроде Go, Node.js, Python и многих других нет нужды собирать образ при каждом изменении кода.

Паттерн:

1. Образ содержит только рантайм и инструменты (Go, air, gcc и т.п.).
2. Код монтируется через volume.
3. Точка входа — процесс‑релодер.

Грубый пример Dockerfile для Go‑dev:

```sh
FROM golang:1.23

RUN go install github.com/cosmtrek/air@latest

WORKDIR /app
CMD ["air"]
```

В компоуз‑файле:

```sh
services:
  api:
    build:
      context: ./api
      dockerfile: Dockerfile.dev
    volumes:
      - ./api:/app
    ports:
      - "8080:8080"
    environment:
      - ENV=local
```

Такой подход даёт удобство локальной разработки при наличии полного окружения в контейнере.

---

## env‑файлы для разработки

Конфигурация dev‑стенда сосредотачивается в одном или нескольких env‑файлах.

Например:

```sh
services:
  api:
    env_file:
      - ./config/dev/api.env

  postgres:
    env_file:
      - ./config/dev/postgres.env
```

Файлы `.env` могут содержать:

- креды для локальной БД;
- адреса сервисов внутри сети Compose;
- настройки логирования, дебага и т.п.

Иногда удобно иметь общий `.env` в корне проекта с базовыми значениями и отдельные файлы для сервис‑специфичных настроек.

---

## Override‑файлы

Compose поддерживает несколько файлов, которые накладываются друг на друга.
Частый паттерн:

- базовый `docker-compose.yml` описывает «близкое к прод» окружение;
- `docker-compose.override.yml` содержит dev‑специфику: bind mounts, debug‑флаги, дополнительные порты.

При запуске `docker compose up` оба файла учитываются автоматически.

Пример базового файла:

```sh
# docker-compose.yml
version: "3.9"

services:
  api:
    image: my-api:1.0.0
    ports:
      - "8080:8080"
    depends_on:
      - postgres

  postgres:
    image: postgres:16
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

Пример override‑файла:

```sh
# docker-compose.override.yml
services:
  api:
    build: ./api
    volumes:
      - ./api:/app
    environment:
      - ENV=local
```

В результате:

- в прод‑подобной среде можно использовать только базовый файл;
- локально — автоматически применяется override, и код монтируется с хоста.

---

## Раздельные профили

Для более сложных сценариев существуют профили (`profiles`).
Они позволяют включать и выключать наборы сервисов для разных задач:

```sh
services:
  api:
    profiles:
      - app

  postgres:
    profiles:
      - infra

  redis:
    profiles:
      - infra
```

Запуск только приложения:

```sh
docker compose --profile app up
```

Запуск приложения с инфраструктурой:

```sh
docker compose --profile app --profile infra up
```

Профили превращают compose‑файл в гибкий инструмент для разных сценариев разработки.

---

## Повседневный дев‑флоу

Для разработчика daily‑routine с Compose выглядит так:

1. Один раз поднять инфраструктуру:

   ```sh
   docker compose up -d postgres redis
   ```

2. Запускать и перезапускать `api` по необходимости (с hot reload).
3. Сбрасывать данные при необходимости, удаляя тома или пересоздавая их.
4. Выключать стенд в конце рабочего дня:

   ```
   docker compose down
   ```

## Зачем нужен «prod‑like» стенд на Compose

Полноценный прод обычно живёт в Kubernetes, Swarm или другой оркестрации.
Тем не менее, staging и окружения для интеграционных тестов часто удобно собирать на Compose: достаточно «похожего» поведения без усложнения инфраструктуры.

Цель prod‑like окружения на Compose:

- использовать те же образы, что и прод;
- минимально приближенные конфиги (переменные окружения, порты, таймауты, лимиты);
- возможность запускать его локально или на выделенной машине.

---

## Разделение файлов по окружениям

Здравый подход — разделять файлы:

- `docker-compose.yml` — базовая конфигурация;
- `docker-compose.dev.yml` — dev‑настройки;
- `docker-compose.staging.yml` — staging‑настройки;
- `docker-compose.prod.yml` — приближенные к бою настройки.

Запуск staging:

```sh
docker compose -f docker-compose.yml -f docker-compose.staging.yml up -d
```

Запуск «продового» варианта:

```sh
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

Такое перекладывание позволяет использовать один и тот же набор сервисов с разными параметрами.

---

## Отдельные env‑файлы для окружений

Файлы окружения — ключевой элемент разделения конфигурации:

- `config/dev/*.env`;
- `config/staging/*.env`;
- `config/prod/*.env` (жёстко ограниченный доступ).

Пример:

```sh
services:
  api:
    env_file:
      - ./config/staging/api.env

  postgres:
    env_file:
      - ./config/staging/postgres.env
```

Структура env‑файлов может быть согласована между окружениями: одинаковые ключи, разные значения.
Это упрощает переносимость и уменьшает риск ошибки при деплое.

---

## Настройки близкие к бою

Прод‑подобный стенд должен учитывать ограничения, характерные для реальной среды:

- лимиты ресурсов контейнеров (пусть и мягкие);
- настройки логирования;
- количество воркеров, пулов подключений;
- уровни логов, timeouts, retries.

На уровне Compose это частично отражается в:

- `deploy.resources` (если используется Swarm‑контекст);
- переменных окружения, задающих лимиты в самом приложении.

Пример:

```sh
services:
  api:
    image: my-api:1.2.0
    environment:
      - GOMAXPROCS=2
      - HTTP_SERVER_TIMEOUT=15s
      - DB_MAX_OPEN_CONNS=20
    depends_on:
      postgres:
        condition: service_healthy
```

---

## Логи и тома

Для prod‑подобного стенда важно, чтобы:

- логи были доступны снаружи контейнера;
- данные сервисов (БД, кеш, брокер) жили в именованных томах.

Пример:

```sh
services:
  api:
    image: my-api:1.2.0
    volumes:
      - api_logs:/var/log/api

  postgres:
    image: postgres:16
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  api_logs:
  postgres_data:
```

Далее эти тома можно резервировать, анализировать и переносить.
Логи — подключать к внешним системам мониторинга и анализа.

---

## Безопасность и экспонирование портов

Даже в staging‑окружении порты не должны быть открыты без необходимости.

Рекомендации:

- наружу публиковать только API и нужные UI (Prometheus, Grafana, RabbitMQ UI и т.п.);
- внутренние сервисы (БД, кеш, брокер) оставлять только внутри сети Compose;
- при необходимости использовать reverse proxy с авторизацией.

Разумный минимум:

```sh
services:
  api:
    ports:
      - "80:8080"

  prometheus:
    ports:
      - "9090:9090"
```

Остальные сервисы остаются доступны только внутри сети проекта.

---

## Автоматизация и CI/CD

Prod‑like стенд на Compose хорошо сочетается с CI:

- сборка образов;
- прогон интеграционных тестов на стеке `api + db + cache + mq`;
- запись артефактов (логи, дампы БД, метрики).

Схематично:

1. CI собирает образы и пушит их в registry.
2. CI поднимает Compose‑стенд нужной версии образов.
3. Запускаются интеграционные/контрактные тесты.
4. По завершении стенд выключается командой `docker compose down -v`.

Compose‑файлы в репозитории становятся явной спецификацией окружения, а не «магией» на серверах.

## Цель главы

В этой главе рассматривается полный пример Go‑стека, описанного на Docker Compose.
Стек включает:

- `api` — HTTP‑API на Go;
- `postgres` — основная база данных;
- `redis` — кеш/брокер для простых задач;
- `migrations` — сервис для запуска миграций;
- `prometheus` — сбор метрик.

Этот пример можно использовать как отправную точку для реального проекта, постепенно усложняя его.

---

## Обзор архитектуры

Внутри Compose‑проекта предполагается одна сеть `backend`.
Все сервисы подключены к этой сети и общаются по DNS‑именам:

- `api` подключается к БД `postgres:5432` и кешу `redis:6379`;
- `migrations` использует тот же DSN, что и `api`;
- `prometheus` опрашивает метрики `api` по HTTP.

Данные БД и Redis сохраняются в именованных томах.
Метрики Prometheus сохраняются в отдельном томе.

---

## Compose‑файл (базовая верся. Пример)

Ниже приведён минимальный, но законченный пример `docker-compose.yml`:

```sh
version: "3.9"

services:
  api:
    build:
      context: ./api
      dockerfile: Dockerfile
    image: my-api:local
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    environment:
      - ENV=local
      - DB_DSN=postgres://postgres:secret@postgres:5432/app?sslmode=disable
      - REDIS_ADDR=redis:6379
      - PROMETHEUS_PUSH_ENABLED=false
    ports:
      - "8080:8080"
    networks:
      - backend

  postgres:
    image: postgres:16
    environment:
      - POSTGRES_PASSWORD=secret
      - POSTGRES_DB=app
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - backend
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s

  redis:
    image: redis:7
    command: ["redis-server", "--save", "", "--appendonly", "no"]
    networks:
      - backend

  migrations:
    build:
      context: ./migrations
      dockerfile: Dockerfile
    image: my-migrations:local
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_DSN=postgres://postgres:secret@postgres:5432/app?sslmode=disable
    networks:
      - backend
    # одноразовый запуск миграций и выход
    command: ["./migrate", "up"]

  prometheus:
    image: prom/prometheus:v2.54.0
    volumes:
      - ./infra/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus_data:/prometheus
    networks:
      - backend
    ports:
      - "9090:9090"

networks:
  backend:
    driver: bridge

volumes:
  postgres_data:
  prometheus_data:
```

Этот файл описывает целостный стек, в котором каждый сервис играет свою роль, а Compose связывает их в единое приложение.

---

## Сервис api (Go)

Сервис `api` — центральная часть стека.
Он использует:

- HTTP‑порт 8080;
- подключение к Postgres;
- подключение к Redis;
- опциональную интеграцию с Prometheus (в данном примере — только как target, а не push).

Важные моменты конфигурации:

- `depends_on` с `service_healthy` для Postgres гарантирует, что приложение не стартует раньше «здоровой» БД;
- переменные окружения задают DSN и адрес Redis;
- `ports: "8080:8080"` публикует API наружу для тестирования.

Dockerfile для Go‑api в прод‑стиле может быть multi‑stage: билд бинарника, затем минимальный runtime с только бинарником и сертификатами.

---

## Сервис postgres

База данных описана минимально, но с упором на сохранение данных:

- указаны пароль и имя БД через `environment`;
- каталог данных вынесен в именованный том `postgres_data`;
- настроен `healthcheck` для корректной работы `depends_on`.

Простой конфиг:

```sh
postgres:
  image: postgres:16
  environment:
    - POSTGRES_PASSWORD=secret
    - POSTGRES_DB=app
  volumes:
    - postgres_data:/var/lib/postgresql/data
  networks:
    - backend
  healthcheck:
    test: ["CMD-SHELL", "pg_isready -U postgres"]
    interval: 10s
    timeout: 5s
    retries: 5
    start_period: 10s
```

Такой уровень деталей достаточен для локального и staging‑стендов.
Для более продвинутых сценариев добавляются backup‑механизмы, настройки WAL, репликация и т.д.

---

## Сервис redis

Redis используется как кеш и/или простейший брокер задач.
В данном примере он запускается без сохранения на диск, чтобы не создавать лишних томов и не мешать разработке.

Конфигурация:

```sh
redis:
  image: redis:7
  command: ["redis-server", "--save", "", "--appendonly", "no"]
  networks:
    - backend
```

Сервис `api` подключается к нему по адресу `redis:6379`.
При необходимости настройки Redis легко вынести в конфигурационный файл и смонтировать его через volume.

---

## Сервис migrations

Сервис `migrations` иллюстрирует подход «миграции как отдельный процесс/джоба».
Он:

- использует тот же DSN, что и `api`;
- зависит от здорового Postgres;
- запускает бинарник мигратора и завершает работу.

Такой подход удобен тем, что:

- миграции можно запускать руками: `docker compose run --rm migrations`;
- миграции можно включить в CI‑pipeline;
- не нужно встраивать миграции в старт приложения (что усложняет логику и повышает риск ошибок).

---

## Сервис prometheus

Prometheus собирает метрики со всех интересующих целей (в данном примере — с `api`).
Файл конфигурации `prometheus.yml` монтируется через volume:

```sh
prometheus:
  image: prom/prometheus:v2.54.0
  volumes:
    - ./infra/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:ro
    - prometheus_data:/prometheus
  networks:
    - backend
  ports:
    - "9090:9090"
```

Пример `prometheus.yml`:

```sh
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: "api"
    static_configs:
      - targets: ["api:9090"]
```

Предполагается, что Go‑сервис `api` экспонирует метрики на `:9090/metrics`.
Таким образом, локальный стек сразу включает полноценное наблюдение.

---

## Dev‑вариант с override‑файлом

Для разработки удобно добавить `docker-compose.override.yml`, который:

- пересобирает `api` из исходников;
- монтирует код внутрь контейнера;
- включает hot reload.

Пример:

```sh
services:
  api:
    build:
      context: ./api
      dockerfile: Dockerfile.dev
    volumes:
      - ./api:/app
    environment:
      - ENV=local
```

В результате:

- базовый `docker-compose.yml` остаётся ближе к прод;
- override добавляет только dev‑специфику;
- запуск простой командой `docker compose up` поднимает стек в режиме, удобном разработчику.

---

## Итог

Этот пример демонстрирует типичный «микростенд» для Go‑backend’а, описанный целиком в Docker Compose.
Он сочетает в себе:

- единый конфигурационный файл;
- воспроизводимое окружение;
- разделение ответственности между сервисами;
- расширяемость за счёт override‑файлов и дополнительных сервисов.

Такой compose‑стек легко адаптировать под реальные проекты: заменить Go‑API на другой сервис, добавить worker‑ы, подключить Grafana и дописать миграции.
Главное — сохранять декларативность и избегать «магии» за пределами репозитория.
