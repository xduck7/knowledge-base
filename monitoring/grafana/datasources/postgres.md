# PostgreSQL как источник метрик и данных

## Два режима: прямой SQL и metрик через postgres-exporter

С PostgreSQL в Grafana обычно работают в двух разных ролях:

1. **Прямой SQL‑datasource**  
   Grafana подключается к PostgreSQL как к обычной БД, посылает SQL‑запросы и визуализирует результат.

2. **Инфраструктурные метрики через `postgres-exporter`**  
   Отдельный процесс `postgres-exporter` снимает метрики с PostgreSQL и отдаёт их в Prometheus, а уже Prometheus является datasource для Grafana.

Во втором режиме цепочка такая:

1. `postgres-exporter` подключается к PostgreSQL как клиент и периодически выполняет SQL к системным вьюхам (`pg_stat_database`, `pg_stat_activity`, `pg_locks` и т.д.).  
2. Exporter экспонирует результаты в формате Prometheus‑метрик на HTTP‑эндпоинте `/metrics`.  
3. Prometheus скрейпит этот эндпоинт по pull‑модели и сохраняет метрики как time series.  
4. Grafana подключается к Prometheus и визуализирует эти метрики: QPS, связи, коннекты, deadlocks, cache hit ratio, репликацию и т.п.

То есть для мониторинга состояния самой базы (health/perf) Grafana обычно использует не прямой PostgreSQL‑datasource, а связку `postgres-exporter → Prometheus → Grafana`.

### PostgreSQL как прямой datasource: как это работает

При добавлении PostgreSQL‑datasource указываются:

- `Host`, `Port` (например, `postgres:5432`);  
- `Database`, `User`, `Password`;  
- TLS‑параметры, схема, дополнительные опции.

Grafana:

- устанавливает соединение с БД;  
- выполняет SQL‑запросы, прописанные в панели;  
- получает результат в табличном виде и мапит колонки на:

  - время (`time`);  
  - значения (`value1`, `value2`, …);  
  - измерения (`status`, `country`, `service` и т.д.).

PostgreSQL всегда возвращает «чистую» таблицу.  
Grafana решает, как превратить её:

- в time series (когда есть колонка времени);  
- в таблицу (когда нужна просто табличная форма);  
- в другие типы визуализаций.

### SQL‑запросы как источник временных рядов

Grafana ожидает от SQL:

- колонку времени (`time` или помеченную как time);  
- числовую колонку как значение;  
- дополнительные колонки как измерения/лейблы.

Пример:

```c
SELECT
  time_bucket('5 minutes', created_at) AS time,
  status,
  count(*) AS value
FROM http_requests
WHERE $__timeFilter(created_at)
GROUP BY 1, 2
ORDER BY 1, 2;
```

Здесь:

- `time_bucket`/`date_trunc` — агрегация по интервалам;  
- `status` — измерение (аналог label’а);  
- `value` — числовая метрика;  
- макрос `$__timeFilter(created_at)` подставляет интервал времени из панели.

Grafana получит таблицу и построит:

- отдельный ряд для каждого `status`;  
- точки `(time, value)` для каждого ряда.

### Как «летят» данные от PostgreSQL до панели

Поток данных для прямого datasource:

1. Пользователь открывает дашборд.  
2. Панель с PostgreSQL‑datasource содержит SQL с макросами (`$__timeFilter`, `$__interval`, переменные).  
3. Grafana подставляет реальные значения (таймдиапазон, значения переменных).  
4. Grafana отправляет итоговый SQL в PostgreSQL.  
5. PostgreSQL выполняет запрос и возвращает набор строк/колонок.  
6. Grafana интерпретирует:

   - колонку времени → ось X;  
   - числовые колонки → значения;  
   - остальные → лейблы/легенды/категории;

   и рисует график или таблицу.

Для технических метрик PostgreSQL (QPS, коннекты, кеш, репликация) обычно используют `postgres-exporter` и Prometheus.  
Для бизнес‑метрик и отчётов (заказы, пользователи, деньги) — прямой PostgreSQL‑datasource с SQL.

### Два мира: бизнес‑SQL и метрики через exporter

Итого:

- **Бизнес‑данные / аналитика**  
  - Grafana → PostgreSQL (SQL);  
  - результат — time series/таблицы по данным приложения.

- **Технические метрики PostgreSQL**  
  - PostgreSQL → `postgres-exporter` (SQL к системным вьюхам);  
  - `postgres-exporter` → Prometheus (HTTP `/metrics` в формате Prometheus);  
  - Prometheus → Grafana (PromQL и JSON);  
  - результат — классические DB‑дашборды.

На одном дашборде можно комбинировать:

- графики QPS/latency БД через `postgres-exporter` (Prometheus);  
- графики количества заказов/выручки напрямую из PostgreSQL;  
- логи запросов или ошибок (Loki/Elasticsearch);  
- метрики приложения (Prometheus).
