# Prometheus как источник данных

## Общая модель: pull‑запросы из Grafana в Prometheus

Prometheus — это time‑series база для метрик, у которой есть HTTP API.  
Grafana при отрисовке панели формирует HTTP‑запрос к API Prometheus, передавая текстовый запрос на PromQL и временной диапазон.

Ключевые моменты:

- Grafana **ничего не пушит** в Prometheus, а только делает HTTP GET/POST запросы.  
- В запросе указываются: выражение PromQL, начало/конец интервала и шаг дискретизации.  
- Ответ Prometheus — JSON с массивом рядов, где каждый ряд описан набором меток и списком точек вида `[timestamp, value]`.

Упрощённый вид ответа Prometheus:

```json
{
  "status": "success",
  "data": {
    "resultType": "matrix",
    "result": [
      {
        "metric": {
          "__name__": "http_requests_total",
          "method": "GET",
          "handler": "/api",
          "instance": "app-1:8080",
          "job": "my_service"
        },
        "values": [
          [1710000000, "1234"],
          [1710000300, "1300"]
        ]
      }
    ]
  }
}
```

Grafana берёт такие ряды, приводит значения к числам и рисует их как time series.

## Настройка Prometheus datasource в Grafana

При добавлении Prometheus в Grafana указываются:

- URL HTTP‑API Prometheus (например, `http://prometheus:9090`);  
- дополнительные настройки:  
  - настройки HTTP (TLS, заголовки, auth);  
  - опциональный proxy (запрос идёт от бекенда Grafana, а не из браузера).

Grafana общается с Prometheus:

- по HTTP(S);  
- как правило, через backend Grafana (Server mode), чтобы не светить Prometheus наружу.

### HTTP / timeout / TTL

При настройке datasource можно задать:

- timeout запросов — сколько ждать ответа Prometheus, прежде чем считать запрос неудачным;  
- max idle connections, keep‑alive и другие параметры HTTP‑клиента (обычно скрыты за простыми полями);  
- кеширование результатов (на уровне Grafana или внешних прокси), чтобы одни и те же запросы (например, при автообновлении дашборда) не нагружали Prometheus лишний раз.

TTL как явная настройка есть не всегда, но идея такая:

- либо используется встроенный кеш Grafana/браузера;  
- либо ставится внешний reverse‑proxy (nginx, varnish) с кешем для одинаковых запросов, если нагрузка большая.

## Метки (labels) и выбор instance

Prometheus хранит каждый ряд метрик как набор:

- имени метрики: `__name__` (например, `http_requests_total`);  
- набора меток: `job`, `instance`, `method`, `status_code`, `pod`, `namespace` и т.д.;  
- списка точек `[timestamp, value]`.

PromQL‑запрос в Grafana всегда оперирует метками:

```
sum by (status_code) (
  rate(http_requests_total{job="my_service"}[5m])
)
```

Здесь:

- фильтрация по меткам `job="my_service"`;  
- агрегация по `status_code`;  
- `instance` можно явно указывать, чтобы выбрать конкретную ноду/под:

```c
rate(http_requests_total{job="my_service", instance="app-1:8080"}[5m])
```

Grafana даёт удобный UI для:

- подстановки переменных в фильтры (например, выбирать `instance` из дропаун‑списка);  
- построения шаблонов запросов (`$namespace`, `$pod`, `$app` и т.д.).

## Формат данных и преобразование внутри Grafana

Grafana принимает от Prometheus JSON, но внутри превращает это в:

- одну или несколько time series (рядов) для графиков;  
- таблицы для табличных панелей (при необходимости делает pivot/transform).

Особенности:

- временная метка — обычно Unix time в секундах или миллисекундах;  
- значение — строка в JSON, приводится к числу с плавающей точкой;  
- метки Prometheus становятся «полями» метаданных ряда (series name / legend / labels).

## Типовые проблемы: down, timeout, high cardinality

Частые проблемы:

- **Prometheus недоступен (down)**:  
  - Grafana не может достучаться до URL;  
  - панель показывает ошибку соединения.  

- **Timeout**:  
  - запрос слишком тяжёлый или Prometheus перегружен;  
  - нужно оптимизировать PromQL, уменьшать диапазон, настраивать recording rules.  

- **Слишком много рядов (high cardinality)**:  
  - если у метрик много разных комбинаций меток (например, добавили `user_id` или `request_id`), каждый запрос возвращает тысячи/миллионы рядов;  
  - Grafana начинает тормозить при рендере, Prometheus тяжело считать такие запросы;  
  - решается дизайном метрик и агрегацией (убираем лишние метки, считаем заранее).
