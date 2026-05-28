# Практика 4 — Изменение данных

Конструкции: `INSERT`, `VALUES`, `UPDATE`, `SET`, `DELETE`, `RETURNING`.

Таблицы: `users`, `products`, `orders`, `accounts`.

Главные правила:
- `UPDATE` и `DELETE` почти всегда должны иметь `WHERE` — иначе изменишь/удалишь всё.
- `RETURNING` (Postgres) — сразу возвращает изменённые строки. Идиоматично для backend: один RTT, не нужно дёргать `SELECT` после.
- State transition: `UPDATE ... WHERE id = $1 AND status = 'pending'` — проверяй текущее состояние, чтобы избежать гонок.
- Во время обучения оборачивай мутации в `BEGIN; ... ROLLBACK;` — изменения откатятся, можно повторять задачу.

| # | Файл                                                       | Что нового                       |
|---|------------------------------------------------------------|----------------------------------|
| 1 | [01-insert-user.sql](./01-insert-user.sql)                 | базовый `INSERT`                 |
| 2 | [02-insert-returning.sql](./02-insert-returning.sql)       | `RETURNING`                      |
| 3 | [03-update-email.sql](./03-update-email.sql)               | `UPDATE` + `WHERE id`            |
| 4 | [04-soft-delete-user.sql](./04-soft-delete-user.sql)       | soft-delete через `UPDATE`       |
| 5 | [05-mark-order-paid.sql](./05-mark-order-paid.sql)         | state transition `UPDATE`        |
| 6 | [06-bulk-discount.sql](./06-bulk-discount.sql)             | `UPDATE` по группе строк         |
| 7 | [07-delete-test-product.sql](./07-delete-test-product.sql) | `DELETE`                         |
| 8 | [08-deactivate-blocked.sql](./08-deactivate-blocked.sql)   | `UPDATE` + подзапрос `IN`        |
| 9 | [09-batch-insert-products.sql](./09-batch-insert-products.sql) | batch `INSERT` (`VALUES` с несколькими строками) |

## Перед практикой и после

Чтобы данные не мутировали между задачами — оборачивай в транзакцию:

```sql
BEGIN;

-- твой запрос

ROLLBACK;  -- откатить, чтобы не задеть seed
```

Или после серии — сделай `docker compose down -v && docker compose up -d` в `sql/` для полного reset.

---

## Подсказки и теория по задачам

### [01. Создать пользователя](./01-insert-user.sql)

**Что нового.** Базовая форма:

```sql
INSERT INTO table (col1, col2, ...)
VALUES (val1, val2, ...);
```

Колонки в `INSERT` указываем **явно** — не полагайся на порядок колонок в таблице, это хрупко.

**Подсказка.** `CURRENT_TIMESTAMP` или `NOW()` — текущее время в Postgres. Не пиши `'2025-01-01'` руками.

**Что НЕ указывать.**
- `deleted_at` — он nullable, по умолчанию будет NULL.
- `is_active` — есть `DEFAULT true`, можно не указывать. Но для тренировки укажи явно.

---

### [02. Insert + RETURNING](./02-insert-returning.sql)

**Что нового.** `RETURNING` — Postgres-specific: после `INSERT` сразу возвращает вставленные данные. **Зачем нужно в backend:**

- Получить сгенерированный `id` (когда он `SERIAL` / `IDENTITY`), без отдельного `SELECT`.
- Получить значения `DEFAULT`-полей (например, `created_at = NOW()`).
- Один round-trip к БД вместо двух.

```sql
INSERT INTO products (...)
VALUES (...)
RETURNING id, name;
```

Возвращается так же, как `SELECT` — таблицей.

**Production-паттерн в Go.**

```go
var id int
err := db.QueryRow(
    "INSERT INTO products (...) VALUES (...) RETURNING id",
    ...,
).Scan(&id)
```

То есть `QueryRow` вместо `Exec` — потому что есть результат.

---

### [03. Обновить email](./03-update-email.sql)

**Что нового.** Базовая форма:

```sql
UPDATE users
SET email = 'ivan@mail.com'
WHERE id = 3;
```

**Главное правило.** **Всегда** пиши `WHERE` в `UPDATE`. Без него обновятся **все строки таблицы** — это самая частая авария новичков (и не только).

В реальных проектах часто настраивают safe-update режим в IDE / в Postgres, который блокирует `UPDATE` без `WHERE`. Но рассчитывать на это нельзя — пиши `WHERE` сам.

**Проверка.** После выполнения сделай `SELECT id, name, email FROM users WHERE id = 3` — убедись, что email изменился.

Не забудь обернуть в `BEGIN; ... ROLLBACK;` если не хочешь оставить изменение.

---

### [04. Soft-delete пользователя](./04-soft-delete-user.sql)

**Что нового.** **Soft-delete** — production-паттерн: вместо физического `DELETE` помечаем строку как «удалённую». Зачем:
- сохраняется история и аудит;
- внешние ссылки (`FOREIGN KEY`) не ломаются;
- можно «восстановить» — просто проставить `deleted_at = NULL`.

```sql
UPDATE users
SET deleted_at = NOW(),
    is_active = false
WHERE id = 5
RETURNING id, deleted_at;
```

Несколько колонок в `SET` — через запятую.

**Следствие.** Если в проекте принят soft-delete — все `SELECT` должны добавлять `WHERE deleted_at IS NULL`, иначе будут показывать «удалённых». Это легко забыть, поэтому такие фильтры часто выносят в общий слой / `VIEW`.

---

### [05. Перевести заказ в paid (state transition)](./05-mark-order-paid.sql)

**Что нового — production-style state transition.**

```sql
UPDATE orders
SET status = 'paid',
    paid_at = NOW()
WHERE id = 2
  AND status = 'pending'
RETURNING id, status, paid_at;
```

**Зачем `AND status = 'pending'`** — защита от гонок:
- Два процесса одновременно пытаются оплатить один и тот же заказ.
- Первый успел: заказ стал `paid`.
- Второй приходит с тем же `id` — но `WHERE id = $1 AND status = 'pending'` не сматчится, `UPDATE` затронет 0 строк, `RETURNING` вернёт пусто.
- Бэкенд видит, что строк 0 → знает, что оплата уже произошла кем-то другим (либо заказ был `cancelled`, и оплачивать его нельзя).

Без проверки `status` мы можем перезатереть `cancelled` или `refunded` заказ, что хуже потерянных денег.

**Что попробовать.** После выполнения попробуй ещё раз — `RETURNING` вернёт **пусто**, потому что заказ уже не `pending`. Это правильное поведение.

---

### [06. Bulk discount](./06-bulk-discount.sql)

**Что нового.** В `SET` можно использовать **выражения**, в том числе со ссылкой на текущее значение колонки:

```sql
SET price = price * 0.8
```

Postgres вычислит новое значение на основе старого.

```sql
UPDATE products
SET price = price * 0.8
WHERE category = 'furniture'
  AND is_active = true
RETURNING id, name, price;
```

**Production-нюанс.** При работе с деньгами `INTEGER` в копейках/центах — это норм. Когда поле `numeric(10,2)` — будь внимательнее с округлением; `INTEGER * 0.8` в Postgres вернёт `double precision`, и при записи обратно в INT — округление.

`RETURNING` тут особенно полезно — сразу видишь, какие именно строки задело, для аудита.

---

### [07. Удалить товар](./07-delete-test-product.sql)

**Что нового.**

```sql
DELETE FROM products
WHERE id = 7
RETURNING id, name;
```

`DELETE` без `WHERE` сотрёт всю таблицу. Те же правила, что у `UPDATE`: **всегда** пиши `WHERE`.

`RETURNING` на `DELETE` — удобно для аудита: можно записать удалённое в лог.

**Production-нюанс: FOREIGN KEY.** Если на эту строку есть ссылки из других таблиц (например, `order_items.product_id`), `DELETE` упадёт с ошибкой `violates foreign key constraint`. В нашем случае на product 7 нет ссылок в `order_items` — пройдёт.

В production обычно избегают физического `DELETE` для сущностей с историей (см. soft-delete в задаче 04). `DELETE` остаётся для технических чисток.

---

### [08. Деактивировать пользователей с заблокированными счетами](./08-deactivate-blocked.sql)

**Что нового — `IN` с подзапросом.** Можно подставить в `IN` не список, а результат другого `SELECT`:

```sql
UPDATE users
SET is_active = false
WHERE id IN (
    SELECT user_id
    FROM accounts
    WHERE is_blocked = true
)
RETURNING id, name;
```

Это даёт связать таблицы без `JOIN` — пока для апдейтов это норм. Для `SELECT` в `JOIN`-е стиль будет другой (это практика 5).

**Production-нюанс.** Если внутренний `SELECT` может вернуть очень много строк — могут быть проблемы с производительностью. Альтернатива в Postgres — `UPDATE ... FROM`, но это уже уровень 2.

`RETURNING` тут особенно ценен — backend увидит точный список тех, кого затронули, и сможет, например, отправить им уведомления.

---

### [09. Batch INSERT — три товара одним запросом](./09-batch-insert-products.sql)

**Что нового — `VALUES` с несколькими строками.** Можно вставить сразу несколько строк, перечислив их через запятую после `VALUES`:

```sql
INSERT INTO products (id, name, category, price, is_active)
VALUES
    (101, 'Mousepad', 'electronics', 25, true),
    (102, 'Lamp',     'furniture',   70, true),
    (103, 'Sticker',  'stationery',   2, true)
RETURNING id, name;
```

Один запрос — одна транзакция — один round-trip. Гораздо быстрее, чем 3 отдельных `INSERT`.

**Production-нюанс.** Это **самый частый способ массовой вставки** в backend. Если вставляешь >1000 строк за раз — есть ещё `COPY` (это уже уровень 2), но до тысяч `INSERT ... VALUES (...), (...), ...` отлично справляется.

**Атомарность.** Если хоть одна строка нарушает constraint (например, `id = 101` уже есть) — **весь batch** откатывается. Все три не вставятся. Это безопасно по умолчанию.

**Что попробовать.** Попробуй вставить дважды подряд — второй запуск упадёт на `duplicate key value violates unique constraint "products_pkey"`, потому что `id = 101` уже занят. Это естественное поведение `PRIMARY KEY`.
