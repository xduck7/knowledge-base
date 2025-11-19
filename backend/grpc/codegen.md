# Генерация кода из `.proto`

> `.proto` — это источник правды, а codegen — мост между контрактом и исполняемым кодом.  
> В этом файле — обзор основных подходов к генерации кода для gRPC‑сервисов на Go: через `protoc`, `buf`, `easyp` и дополнительные инструменты поверх них.

---

## Зачем вообще нужна генерация

В gRPC мы **не пишем руками** интерфейсы клиента/сервера и структуры сообщений — всё это генерируется из `.proto` файлов.

Что именно генерируется:

- структуры Go для сообщений (`message`);
- интерфейсы и реализации [translate:stub]’ов для gRPC‑клиентов;
- интерфейсы для gRPC‑серверов (которые реализует твой код);
- вспомогательные функции (marshal/unmarshal, проверки и т.д.).

Плюсы:

- единый контракт для всех языков;
- меньше ручного кода и ошибок;
- проще эволюционировать API (добавление полей, новых методов и сервисов).

---

## Генерация через `protoc` + плагины

Самый «каноничный» путь — использовать официальный компилятор protobuf `protoc` + плагины для Go.

### Необходимые инструменты

Минимальный набор:

- `protoc` — бинарник компилятора protobuf;
- `protoc-gen-go` — генерация Go‑типов для сообщений;
- `protoc-gen-go-grpc` — генерация gRPC‑клиентов и серверных интерфейсов.

Установка плагинов (Go ≥ 1.20, модульный режим):

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Убедись, что `$GOPATH/bin` или `$GOBIN` находится в `$PATH`.

### Простейшая команда генерации

Допустим, `.proto` лежат в `proto/`, а генерировать хотим в `gen/proto`:

```sh
protoc \
--proto_path=proto \
--go_out=gen/proto \
--go_opt=paths=source_relative \
--go-grpc_out=gen/proto \
--go-grpc_opt=paths=source_relative \
proto/example/users/v1/users.proto
```

Ключевые опции:

- `--proto_path` — корень, откуда `protoc` ищет файлы;
- `--go_out` / `--go-grpc_out` — директории для вывода;
- `paths=source_relative` — генерировать файлы рядом с исходным `.proto`‑путём (удобно для Go‑модулей).

### Типичная интеграция через `Makefile`

```makefile
PROTO_DIR := proto
OUT_DIR := gen/proto

PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

.PHONY: proto
proto:
protoc \
--proto_path=$(PROTO_DIR) \
--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
$(PROTO_FILES)
```

Дальше:

```
make proto
```

и всё сгенерируется в `gen/proto`.

---

## Генерация через `buf`

[translate:Buf] — это более «современный» слой поверх `protoc`, который решает сразу несколько задач:

- управление зависимостями `.proto` файлов;
- проверка стиля и линтинг схем;
- единый интерфейс для генерации кода.

### Базовая конфигурация

Минимальный `buf.yaml` в корне:

```yaml
version: v1
deps: []
build:
roots:
- proto
```

Конфигурация генерации `buf.gen.yaml`:

```yaml
version: v1
plugins:
- name: go
  out: gen/proto
  opt:
    - paths=source_relative
- name: go-grpc
  out: gen/proto
  opt:
    - paths=source_relative
```

### Команды

Проверка схем:

```sh
buf lint
```

Генерация кода:

```sh
buf generate
```

Плюсы `buf`:

- централизованный конфиг генерации для разных языков;
- встроенный линтинг и breaking‑change‑чекер;
- удобная работа с реиспользуемыми proto‑пакетами (registry).

---

## Генерация через `easyp`

[translate:easyp] - это пример инструмента, который поверх protobuf/gRPC добавляет:

- шаблоны для типичных слоёв (handler/service/repository);
- автогенерацию «обвязки» вокруг gRPC;
- интеграцию с конкретным стеком (например, Go + gRPC + HTTP‑gateway и т.п.).

Схема работы обычно такая:

1. Пишешь `.proto` (как обычно).
2. Запускаешь `easyp` (или аналог) с указанием пути к `.proto`.
3. На выходе — не только gRPC‑клиент/сервер, но и:
   - HTTP‑обёртки;
   - базовые реализации интерфейсов;
   - заготовки для тестов.

Пример (условный):

```
easyp \
--proto_path=proto \
--out=internal/app \
proto/example/users/v1/users.proto
```

Такой подход полезен, когда:

- в команде есть стандартизированный паттерн архитектуры;
- хочется сократить шаблонный код вокруг gRPC (конвертеры, DTO, врапперы).

---

## Дополнительные варианты и плагины

Помимо базового `protoc + go + go-grpc`, часто используют дополнительные плагины.

### HTTP / gRPC‑gateway

Генерация HTTP‑gateway поверх gRPC:

- `protoc-gen-grpc-gateway` — HTTP‑прокси → gRPC;
- `protoc-gen-openapiv2` — генерация OpenAPI‑спеки.

Пример интеграции в команду:

```
protoc \
--proto_path=proto \
--go_out=gen/proto --go_opt=paths=source_relative \
--go-grpc_out=gen/proto --go-grpc_opt=paths=source_relative \
--grpc-gateway_out=gen/proto --grpc-gateway_opt=paths=source_relative \
--openapiv2_out=gen/openapi \
proto/example/users/v1/users.proto
```

### Валидация, mock’и и прочий «сахар»

Популярные направления:

- генерация валидаторов по proto‑аннотациям;
- генерация mock‑объектов для юнит‑тестов;
- генерация client SDK под фронтенд (TypeScript и т.д.).

Общий паттерн: любой новый инструмент подключается как дополнительный `protoc-gen-*` плагин + опции в команде/конфиге.

---

## Организация сгенерированного кода в репозитории

Типичный layout для Go‑проекта:

```
.
├── proto/                # исходные .proto
│   └── example/
│       └── users/v1/users.proto
├── gen/
│   └── proto/            # сгенерированный Go-код
│       └── example/
│           └── users/v1/
│               ├── users.pb.go
│               └── users_grpc.pb.go
└── internal/ or pkg/     # бизнес-логика, реализующая интерфейсы
```

Рекомендации:

- сгенерированный код — **не** править руками;
- при изменении `.proto` всегда запускать генерацию;
- удобно повесить генерацию на:
  - `Makefile` таргет;
  - `go generate` комментарии;
  - CI (валидация, что сгенерированный код актуален).

---

## Пример: `go generate` поверх `protoc`

В `.proto` рядом или в отдельном `.go` файле можно добавить:

```go
//go:generate protoc \
//  --proto_path=../../proto \
//  --go_out=../../gen/proto --go_opt=paths=source_relative \
//  --go-grpc_out=../../gen/proto --go-grpc_opt=paths=source_relative \
//  ../../proto/example/users/v1/users.proto
```

Тогда достаточно:

```
go generate ./...
```

и вся генерация подтянется автоматически.

---

## Когда что выбирать

Кратко:

- **`protoc` + плагины**  
  Минимальный обязательный набор; подойдёт всегда, когда нужен контроль и простота.

- **`buf`**  
  Когда много `.proto`, несколько команд/репозиториев, нужен линтинг и контроль совместимости.

- **`easyp` / похожие генераторы «обвязки»**  
  Когда хочется стандартизировать архитектуру и сократить шаблонный код поверх gRPC.