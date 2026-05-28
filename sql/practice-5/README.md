# Практика 5 — Простые связи и транзакции

Конструкции: `JOIN`, `LEFT JOIN`, `ON`, `BEGIN`, `COMMIT`, `ROLLBACK`.

Таблицы: все 6.

Главные правила:
- Используй **алиасы таблиц**: `FROM users u JOIN orders o ON o.user_id = u.id`.
- В `SELECT`/`WHERE`/`ON`/`GROUP BY` префиксуй колонки алиасами: `u.id`, `o.amount` — иначе при росте запроса непонятно, откуда колонка, и легко поймать `ambiguous column`.
- `JOIN` (он же `INNER JOIN`) — только пересечение.
- `LEFT JOIN` — все строки слева, NULL справа при отсутствии связи.
- После `LEFT JOIN` для счётчика по правой таблице используй `COUNT(o.id)`, а не `COUNT(*)` — иначе посчитаешь строку-«пустышку».
- Фильтр по правой таблице, который должен **сохранить все строки слева**, кладётся в `ON`, а не в `WHERE`.
- `BEGIN; ... COMMIT;` — атомарная транзакция. `ROLLBACK` отменяет всё, что было сделано после `BEGIN`.

| # | Файл                                                                       | Что нового                              |
|---|----------------------------------------------------------------------------|-----------------------------------------|
| 1 | [01-orders-with-users.sql](./01-orders-with-users.sql)                     | `INNER JOIN`                            |
| 2 | [02-paid-orders-with-emails.sql](./02-paid-orders-with-emails.sql)         | `JOIN` + `WHERE`                        |
| 3 | [03-users-without-orders.sql](./03-users-without-orders.sql)               | `LEFT JOIN` + `IS NULL`                 |
| 4 | [04-users-without-paid-orders.sql](./04-users-without-paid-orders.sql)     | фильтр в `ON` (текущая задача из промпта) |
| 5 | [05-orders-count-per-user.sql](./05-orders-count-per-user.sql)             | `LEFT JOIN` + `GROUP BY` + `COUNT(o.id)`|
| 6 | [06-paid-total-per-user.sql](./06-paid-total-per-user.sql)                 | `COALESCE` + `SUM` после `LEFT JOIN`    |
| 7 | [07-products-never-sold.sql](./07-products-never-sold.sql)                 | `LEFT JOIN` через `order_items`         |
| 8 | [08-transfer-money-tx.sql](./08-transfer-money-tx.sql)                     | `BEGIN`/`COMMIT`                        |
| 9 | [09-transfer-rollback.sql](./09-transfer-rollback.sql)                     | `ROLLBACK`                              |

---

## Подсказки и теория по задачам

### [01. Заказы с именами пользователей](./01-orders-with-users.sql)

**Что нового — INNER JOIN.**

```sql
SELECT
    o.id     AS order_id,
    o.status,
    o.amount,
    u.name   AS user_name
FROM orders o
JOIN users u ON u.id = o.user_id;
```

- `JOIN` без префикса = `INNER JOIN`. Это просто короче.
- `ON` указывает условие связи: какая колонка слева равна какой колонке справа.
- Алиасы (`o`, `u`) — обязательная гигиена с момента, как в запросе ≥ 2 таблицы.

**INNER vs прочие JOIN.** `INNER JOIN` возвращает только строки, где связь нашлась с **обеих** сторон. Если в `orders` есть `user_id`, которому нет соответствия в `users` — такой заказ в результат не попадёт. (В нашем seed такого нет — `user_id` есть `FOREIGN KEY`, но в общем случае помни.)

---

### [02. Paid-заказы с email пользователей](./02-paid-orders-with-emails.sql)

**Что нового.** `JOIN` + дополнительные фильтры в `WHERE`. Когда `INNER JOIN` — нет разницы, кидать фильтр в `ON` или в `WHERE`, результат один и тот же. Но **читаемость** другая:

- В `ON` — то, что описывает **связь** таблиц.
- В `WHERE` — то, что фильтрует **строки**.

```sql
FROM orders o
JOIN users u ON u.id = o.user_id   -- связь
WHERE o.status = 'paid'             -- фильтр
  AND u.email IS NOT NULL
ORDER BY o.amount DESC;
```

Это конвенция, а не правило БД — но в production-коде так делают почти всегда.

**Что попробовать.** Попробуй перенести `o.status = 'paid'` в `ON`. Результат тот же — но фильтр «потерялся» среди условий связи, читается хуже.

---

### [03. Пользователи без заказов](./03-users-without-orders.sql)

**Что нового — LEFT JOIN.** `LEFT JOIN` (= `LEFT OUTER JOIN`) возвращает **все** строки левой таблицы, даже если справа связи нет. В таких строках поля справа = NULL.

```sql
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
```

Если у пользователя нет заказов — он попадёт в результат, и `o.id, o.status, o.amount` будут NULL.

**Паттерн «найти строки без связи»:**

```sql
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
WHERE o.id IS NULL
```

Логика: сначала `LEFT JOIN` склеивает всех пользователей со всеми их заказами (или NULL, если заказов нет). Потом `WHERE o.id IS NULL` оставляет только тех, у кого связи **не нашлось**.

**Production-альтернатива.** То же можно через `NOT EXISTS` — на уровне 2 ты увидишь, что в больших таблицах `NOT EXISTS` часто быстрее. Но `LEFT JOIN + IS NULL` — классика, его обязательно нужно уметь.

---

### [04. Пользователи без paid-заказов](./04-users-without-paid-orders.sql)

**Что нового — фильтр в `ON`, а не в `WHERE`.**

«Хочется» написать так:

```sql
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
WHERE o.status = 'paid'    -- ← НЕПРАВИЛЬНО
```

Это сломает `LEFT JOIN`. Условие в `WHERE` отбросит строки, где `o.status IS NULL` (пользователи без заказов), и `LEFT JOIN` превратится в `INNER JOIN`.

**Правильно** — кинуть фильтр в `ON`:

```sql
FROM users u
LEFT JOIN orders o
  ON o.user_id = u.id
 AND o.status = 'paid'      -- ← часть условия связи
WHERE o.id IS NULL;
```

Логика:
1. `LEFT JOIN ... ON o.user_id = u.id AND o.status = 'paid'` — каждому пользователю «приклеиваем» только его paid-заказы (если есть). У кого только pending — справа NULL.
2. `WHERE o.id IS NULL` — оставляем тех, у кого слева не нашлось ни одного paid.

**Правило.** Если хочешь сохранить все строки **левой** таблицы — фильтр по правой таблице кидай в `ON`, а не в `WHERE`. В `WHERE` остаются только фильтры по самой левой таблице.

**Почему это важно.** Это самая частая ошибка при работе с `LEFT JOIN` — её ловят на собеседованиях и в проде.

---

### [05. Количество заказов у каждого пользователя](./05-orders-count-per-user.sql)

**Что нового — `COUNT(o.id)` vs `COUNT(*)` после LEFT JOIN.**

После `LEFT JOIN` пользователь без заказов всё равно даёт **одну** строку в результате (с NULL справа). Если применить `COUNT(*)` — для такого пользователя получится `1`, что неверно.

Решение — `COUNT(колонка_правой_таблицы)`, потому что `COUNT(col)` не считает NULL:

```sql
SELECT
    u.id,
    u.name,
    COUNT(o.id) AS orders_count
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
GROUP BY u.id, u.name
ORDER BY orders_count DESC;
```

**GROUP BY на нескольких колонках.** `u.name` — не агрегат и не уникальный ключ, поэтому добавляем в `GROUP BY`. Можно `GROUP BY u.id, u.name` или `GROUP BY u.id` если в Postgres-конфигурации это разрешено (когда `u.id` — `PRIMARY KEY`, Postgres знает, что остальные колонки `u` однозначно определены). На уровне обучения — пиши обе.

**Что попробовать.** Замени `COUNT(o.id)` на `COUNT(*)` — увидишь, что пользователь без заказов получит 1 вместо 0. Это и есть классический баг.

---

### [06. Сумма paid-заказов по пользователям](./06-paid-total-per-user.sql)

**Что нового — `COALESCE` для замены NULL.**

`SUM` без строк возвращает `NULL`, а не `0`. Когда хочется ноль (например, для отчёта в API), оборачивай в `COALESCE`:

```sql
COALESCE(SUM(o.amount), 0) AS paid_total
```

`COALESCE(a, b, c)` возвращает первый не-NULL аргумент. Самое частое использование — `COALESCE(value, default)`.

**Запрос целиком.**

```sql
SELECT
    u.id,
    u.name,
    COALESCE(SUM(o.amount), 0) AS paid_total
FROM users u
LEFT JOIN orders o
  ON o.user_id = u.id
 AND o.status = 'paid'        -- фильтр в ON, чтобы сохранить всех users
GROUP BY u.id, u.name
ORDER BY paid_total DESC;
```

**Закрепление.** Та же тема, что в задаче 04: фильтр на правую таблицу — в `ON`. Только теперь мы не ищем тех, у кого 0, а считаем сумму, и хотим показать всех. Логика та же.

---

### [07. Товары, которые никогда не продавались](./07-products-never-sold.sql)

**Что нового.** Тот же паттерн «нет связи», что в задаче 03 — но теперь связь через таблицу-связку M:N.

```sql
FROM products p
LEFT JOIN order_items oi ON oi.product_id = p.id
WHERE oi.id IS NULL;
```

Логика:
1. `LEFT JOIN` склеивает каждый товар со **всеми** строками `order_items`, где он встречается. Если товар не покупали — справа NULL.
2. `WHERE oi.id IS NULL` — оставляем только товары без покупок.

**Production-нюанс.** Здесь мы могли бы проверить `oi.product_id IS NULL` — но `oi.id` нагляднее: смотрит на «существование строки в order_items», а не «совпадение по product_id».

Альтернатива через `NOT EXISTS` (уровень 2):

```sql
WHERE NOT EXISTS (
    SELECT 1 FROM order_items oi WHERE oi.product_id = p.id
)
```

Часто быстрее на больших таблицах, потому что Postgres может остановиться на первом найденном.

---

### [08. Перевод денег — атомарная транзакция](./08-transfer-money-tx.sql)

**Что нового — транзакция.**

```sql
BEGIN;

UPDATE accounts
SET balance = balance - 100
WHERE id = 1;

UPDATE accounts
SET balance = balance + 100
WHERE id = 2;

INSERT INTO payments (id, from_account_id, to_account_id, amount, status, created_at)
VALUES (100, 1, 2, 100, 'completed', NOW());

COMMIT;
```

Между `BEGIN` и `COMMIT` — все изменения **не видны** другим сессиям. Если что-то падает посередине, можно `ROLLBACK` — и базу как будто не трогали.

**Зачем именно тут транзакция.** Без транзакции возможна ситуация:
- `UPDATE` списания прошёл;
- сервер упал;
- зачисление не случилось → деньги «потерялись».

С транзакцией это невозможно: либо обе операции применятся, либо ни одна.

**Проверка.** После `COMMIT`:

```sql
SELECT id, balance FROM accounts WHERE id IN (1, 2);
SELECT * FROM payments WHERE id = 100;
```

Если хочешь повторить задачу — сначала сделай reset базы (`docker compose down -v && up -d` в `sql/`), потому что мы реально закоммитили изменения.

---

### [09. Транзакция с ROLLBACK](./09-transfer-rollback.sql)

**Что нового — ROLLBACK.**

```sql
BEGIN;

UPDATE accounts SET balance = balance - 999 WHERE id = 1;
UPDATE accounts SET balance = balance + 999 WHERE id = 2;
INSERT INTO payments (id, from_account_id, to_account_id, amount, status, created_at)
VALUES (200, 1, 2, 999, 'completed', NOW());

-- внутри транзакции эти изменения видны
SELECT id, balance FROM accounts WHERE id IN (1, 2);
SELECT * FROM payments WHERE id = 200;

ROLLBACK;

-- после ROLLBACK — всё откатилось, базы как будто не трогали
SELECT id, balance FROM accounts WHERE id IN (1, 2);
SELECT * FROM payments WHERE id = 200;
```

**Зачем это нужно.**
- **Тестирование запросов**: можно проверить эффект `UPDATE`/`DELETE`, не боясь испортить данные.
- **Откат при ошибке**: в backend-коде — try/catch вокруг транзакции, в catch — `ROLLBACK`.
- **Учебная среда**: оборачивай в `BEGIN; ... ROLLBACK;` все мутации, пока учишься, чтобы не приходилось ресетить базу после каждой задачи.

**Что произойдёт если забыть COMMIT/ROLLBACK.** Открытая транзакция остаётся «висеть» в сессии. Другие сессии её не видят, но: блокировки могут удерживаться, длинная транзакция мешает `VACUUM` (это уже уровень 3). В psql / GoLand при закрытии соединения незакоммиченное автоматически откатится — но в production-коде так делать нельзя.
