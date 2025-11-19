# Protocol Buffers и `.proto` файлы

> Protocol Buffers (protobuf) — это язык описания интерфейсов (IDL) и бинарный формат сериализации данных, который чаще всего используется вместе с gRPC.  
> В `.proto` файлах описываются сообщения и сервисы, а по ним уже генерируется код для клиента и сервера.

---

## Зачем нужен `.proto`

Основные задачи `.proto` файла:

- описать структуру данных (сообщения);
- описать контракты сервисов (RPC‑методы, их запросы и ответы);
- служить единым источником правды для всех языков и сервисов.

Плюсы использования protobuf:

- компактный бинарный формат (меньше трафика, быстрее сериализация);
- строгая типизация и схема;
- удобная эволюция контрактов (добавление полей без ломания старых клиентов).

---

## Базовая структура `.proto`

Минимальный пример gRPC‑контракта:

```protobuf
syntax = "proto3";

package example.users.v1;

option go_package = "github.com/yourorg/yourrepo/gen/proto/users/v1;usersv1";

service UserService {
rpc GetUser(GetUserRequest) returns (GetUserResponse);
}

message GetUserRequest {
int64 id = 1;
}

message GetUserResponse {
int64 id = 1;
string name = 2;
}
```

Ключевые элементы:

- `syntax` — версия синтаксиса (`proto3` или `proto2`, для gRPC почти всегда `proto3`);
- `package` — логическое пространство имён для сообщений и сервисов;
- `option go_package` — маппинг proto‑пакета на Go‑пакет;
- `service` — описание gRPC‑сервиса и его методов;
- `message` — описание структур данных (запросы, ответы и т.д.).

---

## syntax = "proto3"

### Почему proto3

`proto3` — актуальная версия синтаксиса protobuf:

- проще, чем `proto2` (нет `required`, меньше опций);
- лучше поддерживается tooling’ом;
- дефолт для gRPC в разных языках.

### Особенности proto3

- все поля по умолчанию «optional» (либо присутствуют, либо нет, с дефолтными значениями);
- есть тип `oneof` для взаимоисключающих полей;
- поддерживаются well‑known types (`google.protobuf.Timestamp`, `Duration` и т.д.).

---

## package и go_package

### package

`package` задаёт логическое пространство имён внутри protobuf‑миров:

```
package example.users.v1;
```

Рекомендации:

- включать домен/организацию: `yourorg.users.v1`, `example.billing.v1`;
- версионировать через суффикс (`v1`, `v2`, `v1beta1`);
- не менять package задним числом: это ломает совместимость.

### option go_package

Для Go почти всегда добавляют:

```
option go_package = "github.com/yourorg/yourrepo/gen/proto/users/v1;usersv1";
```

Здесь:

- до `;` — путь модуля/директории в Go;
- после `;` — имя пакета (`package usersv1` в Go‑коде).

Это избавляет от проблем с импортами и именами пакетов после генерации.

---

## service и RPC‑методы

### Описание сервиса

Сервис описывает набор RPC‑методов:

```protobuf
service UserService {
rpc GetUser(GetUserRequest) returns (GetUserResponse);
rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}
```

Обычно:

- имя сервиса — существительное + `Service` ([translate:UserService], [translate:AuthService], [translate:BillingService]);
- каждый `rpc` — отдельный метод с чётко определёнными типами запросов и ответов.

### Типы RPC через `.proto`

Тип RPC задаётся сигнатурой:

```protobuf
// Unary
rpc GetUser(GetUserRequest) returns (GetUserResponse);

// Server-streaming
rpc ListUsersStream(ListUsersRequest) returns (stream User);

// Client-streaming
rpc ImportUsers(stream ImportUserRequest) returns (ImportUsersResponse);

// Bidirectional streaming
rpc Chat(stream ChatMessage) returns (stream ChatMessage);
```

Ключевое слово `stream` обозначает поток сообщений в соответствующем направлении.

---

## message: поля и нумерация

### Общий вид сообщения

```protobuf
message User {
int64 id = 1;
string name = 2;
string email = 3;
}
```

Правила:

- каждое поле имеет:
  - тип (`int64`, `string`, `bool`, `User`, `google.protobuf.Timestamp` и т.д.);
  - имя (lower_snake_case);
  - уникальный номер (`= 1`, `= 2`, ...).

### Номера полей

Номер поля используется во внутреннем бинарном формате, а не имя:

- диапазоны:
  - `1`–`15` — компактнее (используются чаще для важных полей);
  - `16`–`2047` — обычные поля;
- нельзя переиспользовать номер под другое поле в той же `message`;
- при удалении поля — лучше пометить `reserved`.

Пример:

```protobuf
message User {
reserved 4, 6;
reserved "old_name", "old_status";

int64 id = 1;
string name = 2;
string email = 3;
bool active = 5;
}
```

---

## Базовые типы и коллекции

### Скаляры

Часто используемые типы:

- `int32`, `int64` — целые (signed);
- `uint32`, `uint64` — целые без знака;
- `bool` — логический тип;
- `string` — UTF‑8 строки;
- `bytes` — произвольный бинарный буфер;
- `float`, `double` — числа с плавающей запятой.

### Повторяющиеся поля

Списки описываются через `repeated`:

```protobuf
message ListUsersResponse {
repeated User users = 1;
}
```

`repeated` — это именно список, порядок элементов сохраняется.

---

## enum: перечисления

Перечисления удобно использовать для статусов/типов:

```protobuf
enum UserStatus {
USER_STATUS_UNSPECIFIED = 0;
USER_STATUS_ACTIVE = 1;
USER_STATUS_BLOCKED = 2;
}

message User {
int64 id = 1;
string name = 2;
UserStatus status = 3;
}
```

Рекомендации:

- первое значение обязательно с номером `0`;
- использовать префикс с именем enum: `USER_STATUS_*`;
- не менять значения задним числом (иначе сломается совместимость).

---

## oneof: взаимоисключающие поля

`oneof` позволяет описать, что установлено только одно из нескольких полей:

```protobuf
message SearchRequest {
string query = 1;

oneof filter {
int64 user_id = 2;
string email = 3;
}
}
```

Особенности:

- в конкретном сообщении может быть установлено только одно поле из `oneof`;
- удобно моделирует sum‑типы / tagged unions.

---

## Well-known types

Вместо самодельных структур стоит использовать стандартные типы из `google.protobuf`:

```protobuf
import "google/protobuf/timestamp.proto";
import "google/protobuf/duration.proto";
import "google/protobuf/wrappers.proto";

message Session {
int64 id = 1;
google.protobuf.Timestamp created_at = 2;
google.protobuf.Duration ttl = 3;
google.protobuf.StringValue description = 4;
}
```

Часто встречающиеся:

- `google.protobuf.Timestamp` — момент времени (аналог `time.Time`);
- `google.protobuf.Duration` — длительность;
- `google.protobuf.*Value` — nullable‑обёртки (`StringValue`, `Int64Value` и т.д.).

---

## Организация `.proto` файлов

### Разделение по доменам

Обычно:

- один домен/подсистема — один или несколько `.proto` файлов;
- пример:
  - `users/v1/users.proto`
  - `auth/v1/auth.proto`
  - `billing/v1/invoices.proto`

Это помогает:

- чётко отделять контракты разных подсистем;
- версионировать каждый домен независимо.

### Имена файлов и пакетов

Рекомендации:

- путь: `proto/<domain>/v1/<name>.proto`;
- `package`: `<org>.<domain>.v1`;
- `go_package`: `github.com/org/project/gen/proto/<domain>/v1;<short_name>`.

Пример:

```
syntax = "proto3";

package example.auth.v1;

option go_package = "github.com/example/project/gen/proto/auth/v1;authv1";
```

---

## Эволюция контрактов и совместимость

### Что можно делать безопасно

- добавлять новые поля с новыми номерами;
- помечать поля как `reserved` и больше их не использовать;
- добавлять новые значения в конец `enum`.

### Чего лучше избегать

- менять тип существующего поля;
- переиспользовать номера удалённых полей под другой смысл;
- удалять или переименовывать значения `enum`, меняя их номера.

Общая идея: старый код должен продолжать корректно работать с новыми сообщениями, игнорируя незнакомые поля.

---

## Примеры хорошего стиля

Краткий чек‑лист:

- имена сообщений — в [translate:PascalCase] (`GetUserRequest`, `User`, `InvoiceItem`);
- имена полей — в [translate:snake_case] (`user_id`, `created_at`);
- у каждого RPC свой `*Request` / `*Response` (не переиспользовать одно сообщение в десятке методов);
- версии в пакете (`example.users.v1`) и в путях;
- `google.protobuf.Empty` использовать, когда нет полезного payload.

Пример:

```protobuf
syntax = "proto3";

package example.users.v1;

option go_package = "github.com/example/project/gen/proto/users/v1;usersv1";

import "google/protobuf/empty.proto";

service UserService {
rpc GetUser(GetUserRequest) returns (GetUserResponse);
rpc DeleteUser(DeleteUserRequest) returns (google.protobuf.Empty);
}

message GetUserRequest {
int64 id = 1;
}

message GetUserResponse {
User user = 1;
}

message DeleteUserRequest {
int64 id = 1;
}

message User {
int64 id = 1;
string name = 2;
string email = 3;
}
```

---

## Связь `.proto` с генерацией кода

`.proto` файл — входная точка для генерации:

- серверные интерфейсы (которые реализует твой сервис);
- клиентские stub’ы (которые используют другие сервисы);
- вспомогательные типы и функции для работы с сообщениями.