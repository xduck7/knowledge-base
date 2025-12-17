# Основные запросы к PostgreSQL

## 📖 Содержание

- [Основные запросы](#основные-запросы)
- [Операторы фильтрации](#операторы-фильтрации)
- [JOIN — объединение таблиц](#join--объединение-таблиц)
- [Подзапросы](#подзапросы)
- [CRUD операции](#crud-операции)
- [Работа с таблицами (DDL)](#работа-с-таблицами-ddl)

---

## Основные запросы

### `SELECT` — получение данных

Используется для выборки полей из таблицы: все (`*`) или только нужные.

```sql
-- Получить все поля из таблицы users
SELECT * FROM users;

-- Получить только id и email
SELECT id, email FROM users;
```

---

### `DISTINCT` — уникальные значения

Убирает дубликаты из результата.

```sql
-- Получить список уникальных городов
SELECT DISTINCT city FROM users;

-- Уникальные комбинации полей
SELECT DISTINCT city, country FROM users;
```

---

### `AS` — алиасы (псевдонимы)

Переименовывает колонки или таблицы в результате.

```sql
-- Переименовать колонку в результате
SELECT email AS user_email, name AS user_name FROM users;

-- Алиас для таблицы
SELECT u.name, o.total 
FROM users AS u 
JOIN orders AS o ON u.id = o.user_id;
```

---

### `ORDER BY` — сортировка

Сортирует результат по указанному полю.

```sql
-- Сортировка по убыванию
SELECT * FROM users ORDER BY created_at DESC;

-- Сортировка по возрастанию (по умолчанию)
SELECT * FROM users ORDER BY name ASC;

-- Сортировка по нескольким полям
SELECT * FROM users ORDER BY city, name DESC;
```

---

### `LIMIT` и `OFFSET` — пагинация

Ограничивает количество записей и пропускает первые N записей.

```sql
-- Получить первые 10 записей
SELECT * FROM users LIMIT 10;

-- Пропустить первые 20 и взять следующие 10 (страница 3)
SELECT * FROM users LIMIT 10 OFFSET 20;

-- Получить 10 последних пользователей
SELECT * FROM users ORDER BY created_at DESC LIMIT 10;
```

---

### `GROUP BY` — группировка

Группирует записи по значению поля.

```sql
-- Подсчитать количество заказов по каждому статусу
SELECT status, COUNT(*) FROM orders GROUP BY status;

-- Результат: status | count
--           pending | 5
--           done    | 12
```

---

### `HAVING` — фильтрация после группировки

Фильтрует сгруппированные данные (в отличие от `WHERE`, который работает до группировки).

```sql
-- Показать только те статусы, где количество заказов > 3
SELECT status, COUNT(*) 
FROM orders 
GROUP BY status 
HAVING COUNT(*) > 3;
```

---

## Операторы фильтрации

### `WHERE` — условие фильтрации

Аналог `if` — фильтрует записи по условию.

```sql
-- Найти пользователя с id = 10
SELECT * FROM users WHERE id = 10;

-- Найти пользователя по email
SELECT * FROM users WHERE email = 'test@mail.com';
```

---

### `AND` / `OR` / `NOT` — логические операторы

```sql
-- AND: оба условия должны быть true
SELECT * FROM users WHERE age > 18 AND city = 'Moscow';

-- OR: хотя бы одно условие true
SELECT * FROM users WHERE city = 'Moscow' OR city = 'SPb';

-- NOT: инвертирует условие
SELECT * FROM users WHERE NOT status = 'blocked';
```

---

### `IN` / `NOT IN` — проверка вхождения в список

```sql
-- Найти пользователей из списка городов
SELECT * FROM users WHERE city IN ('Moscow', 'SPb', 'Kazan');

-- Исключить определённые статусы
SELECT * FROM orders WHERE status NOT IN ('cancelled', 'refunded');
```

---

### `BETWEEN` — диапазон значений

```sql
-- Найти заказы в диапазоне сумм
SELECT * FROM orders WHERE total BETWEEN 100 AND 500;

-- Найти записи за период (даты)
SELECT * FROM orders 
WHERE created_at BETWEEN '2024-01-01' AND '2024-12-31';
```

---

### `LIKE` / `ILIKE` — поиск по паттерну

- `%` — любое количество символов
- `_` — ровно один символ
- `ILIKE` — регистронезависимый поиск (PostgreSQL)

```sql
-- Найти email, начинающийся с 'admin'
SELECT * FROM users WHERE email LIKE 'admin%';

-- Найти email, содержащий 'gmail'
SELECT * FROM users WHERE email LIKE '%gmail%';

-- Найти имена из 4 букв, начинающиеся на 'A'
SELECT * FROM users WHERE name LIKE 'A___';

-- Регистронезависимый поиск
SELECT * FROM users WHERE name ILIKE '%ivan%';
```

---

### `IS NULL` / `IS NOT NULL` — проверка на NULL

```sql
-- Найти пользователей без телефона
SELECT * FROM users WHERE phone IS NULL;

-- Найти пользователей с заполненным телефоном
SELECT * FROM users WHERE phone IS NOT NULL;
```

> ⚠️ **Важно:** Нельзя использовать `= NULL` или `!= NULL` — только `IS NULL` / `IS NOT NULL`!

---

## JOIN — объединение таблиц

Объединяет строки из нескольких таблиц по условию связи.

### `INNER JOIN` (или просто `JOIN`)

Возвращает только совпадающие записи из обеих таблиц.

```sql
-- Получить заказы с информацией о пользователях
SELECT * FROM orders 
JOIN users ON users.id = orders.user_id;
```

**Пример с алиасами:**

```sql
SELECT
    p.product_name,
    c.category_name
FROM products p
JOIN categories c ON c.category_id = p.category_id;
```

---

### `LEFT JOIN`

Возвращает **все записи из левой таблицы**, даже если нет совпадений в правой.

```sql
SELECT u.name, o.total
FROM users u
LEFT JOIN orders o ON u.id = o.user_id;
```

| name | total |
|------|-------|
| Ivan | 300   |
| Anna | 700   |
| Petr | NULL  |

---

### `RIGHT JOIN`

Возвращает **все записи из правой таблицы**, даже если нет совпадений в левой.

```sql
SELECT u.name, o.total
FROM users u
RIGHT JOIN orders o ON u.id = o.user_id;
```

| name | total |
|------|-------|
| Ivan | 300   |
| Anna | 700   |
| NULL | 100   |

---

### `FULL JOIN`

Возвращает **все записи из обеих таблиц**, независимо от совпадений.

```sql
SELECT u.name, o.total
FROM users u
FULL JOIN orders o ON u.id = o.user_id;
```

| name | total |
|------|-------|
| Ivan | 500   |
| Ivan | 300   |
| Anna | 700   |
| Petr | NULL  |
| NULL | 100   |

---

### `CROSS JOIN` — декартово произведение

Каждая запись первой таблицы соединяется с каждой записью второй.

```sql
-- Все комбинации размеров и цветов
SELECT s.size, c.color
FROM sizes s
CROSS JOIN colors c;
```

---

### `SELF JOIN` — соединение таблицы с самой собой

```sql
-- Найти сотрудников и их менеджеров
SELECT 
    e.name AS employee, 
    m.name AS manager
FROM employees e
LEFT JOIN employees m ON e.manager_id = m.id;
```

---

### ⚠️ Важно: `LEFT JOIN` + `WHERE`

Использование `WHERE` после `LEFT JOIN` превращает его в обычный `INNER JOIN`:

```sql
-- ❌ LEFT JOIN превратится в INNER JOIN
SELECT u.name, o.total
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.total > 500;
```

**Правильный способ** — добавить условие в `ON`:

```sql
-- ✅ LEFT JOIN сохранит своё поведение
SELECT u.name, o.total
FROM users u
LEFT JOIN orders o ON u.id = o.user_id AND o.total > 500;
```

---

## Подзапросы

Запрос внутри другого запроса.

### Подзапрос в `WHERE`

```sql
-- Найти пользователей, у которых есть заказы
SELECT * FROM users 
WHERE id IN (SELECT user_id FROM orders);

-- Найти товары дороже средней цены
SELECT * FROM products 
WHERE price > (SELECT AVG(price) FROM products);
```

---

### Подзапрос в `FROM`

```sql
-- Использовать результат подзапроса как таблицу
SELECT avg_by_city.city, avg_by_city.avg_salary
FROM (
    SELECT city, AVG(salary) AS avg_salary
    FROM employees
    GROUP BY city
) AS avg_by_city
WHERE avg_by_city.avg_salary > 50000;
```

---

### `EXISTS` / `NOT EXISTS`

Проверяет, существуют ли записи в подзапросе.

```sql
-- Найти пользователей, у которых есть хотя бы один заказ
SELECT * FROM users u
WHERE EXISTS (
    SELECT 1 FROM orders o WHERE o.user_id = u.id
);

-- Найти пользователей без заказов
SELECT * FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM orders o WHERE o.user_id = u.id
);
```

---

## CRUD операции

### `INSERT` — добавление записи

```sql
-- Добавить одну запись
INSERT INTO users (email, password)
VALUES ('a@mail.com', '123');

-- Добавить несколько записей
INSERT INTO users (email, password) VALUES 
    ('a@mail.com', '123'),
    ('b@mail.com', '456'),
    ('c@mail.com', '789');

-- Добавить и вернуть созданную запись
INSERT INTO users (email, password)
VALUES ('a@mail.com', '123')
RETURNING *;

-- Добавить и вернуть только id
INSERT INTO users (email, password)
VALUES ('a@mail.com', '123')
RETURNING id;
```

---

### `UPDATE` — обновление записи

```sql
-- Обновить одну запись
UPDATE users
SET email = 'new@mail.com'
WHERE id = 1;

-- Обновить несколько полей
UPDATE users
SET email = 'new@mail.com', name = 'New Name', updated_at = NOW()
WHERE id = 1;

-- Обновить и вернуть результат
UPDATE users
SET email = 'new@mail.com'
WHERE id = 1
RETURNING *;
```

---

### `DELETE` — удаление записи

```sql
-- Удалить запись
DELETE FROM users WHERE id = 1;

-- Удалить и вернуть удалённое
DELETE FROM users WHERE id = 1 RETURNING *;

-- Удалить все записи (осторожно!)
DELETE FROM users;
```

---

### `UPSERT` — вставка или обновление

```sql
-- Если запись существует — обновить, иначе — вставить
INSERT INTO users (id, email, name)
VALUES (1, 'a@mail.com', 'Ivan')
ON CONFLICT (id) 
DO UPDATE SET email = EXCLUDED.email, name = EXCLUDED.name;

-- Если конфликт — ничего не делать
INSERT INTO users (email, password)
VALUES ('a@mail.com', '123')
ON CONFLICT (email) DO NOTHING;
```

---

> 💡 **Совет:** Всегда используйте `WHERE` с `UPDATE` и `DELETE`, чтобы случайно не изменить/удалить все записи!

---

## Работа с таблицами (DDL)

### `CREATE TABLE` — создание таблицы

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100),
    age INTEGER CHECK (age >= 0),
    created_at TIMESTAMP DEFAULT NOW()
);

-- С внешним ключом
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    total DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);
```




