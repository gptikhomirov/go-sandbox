# Практика 3 — Агрегаты и группировка

Конструкции: `COUNT`, `SUM`, `AVG`, `MIN`, `MAX`, `GROUP BY`, `HAVING`, `FILTER`.

Таблицы: `users`, `orders`, `products`, `order_items`, `payments`.

Главные правила:
- `COUNT(*)` считает все строки, `COUNT(col)` — только не-NULL значения колонки.
- Если в `SELECT` обычная колонка + агрегат — обычная колонка **должна** быть в `GROUP BY`.
- `WHERE` фильтрует **строки до** группировки, `HAVING` — **группы после** группировки.
- В Postgres alias из `SELECT` нельзя использовать в `HAVING` — повторяй агрегат.
- `FILTER (WHERE ...)` — условный агрегат, Postgres-style.

| # | Файл                                                                  | Что нового                       |
|---|-----------------------------------------------------------------------|----------------------------------|
| 1 | [01-total-users.sql](./01-total-users.sql)                            | `COUNT(*)`                       |
| 2 | [02-count-with-email.sql](./02-count-with-email.sql)                  | `COUNT(col)` vs `COUNT(*)`       |
| 3 | [03-users-by-city.sql](./03-users-by-city.sql)                        | `GROUP BY`                       |
| 4 | [04-avg-age-by-city.sql](./04-avg-age-by-city.sql)                    | `AVG` + `GROUP BY`               |
| 5 | [05-order-stats-by-status.sql](./05-order-stats-by-status.sql)        | несколько агрегатов на группу    |
| 6 | [06-cities-with-multiple.sql](./06-cities-with-multiple.sql)          | `HAVING`                         |
| 7 | [07-big-spenders.sql](./07-big-spenders.sql)                          | `GROUP BY user_id` + `HAVING SUM`|
| 8 | [08-active-vs-inactive-by-city.sql](./08-active-vs-inactive-by-city.sql) | `FILTER`                       |
| 9 | [09-distinct-paid-users.sql](./09-distinct-paid-users.sql)            | `COUNT(DISTINCT col)`            |

---

## Подсказки и теория по задачам

### [01. Общее число пользователей](./01-total-users.sql)

**Что нового.** `COUNT(*)` — самый частый агрегат. Возвращает число строк. **Никогда не отбрасывает NULL**, потому что считает строки, а не значения колонки.

```sql
SELECT COUNT(*) FROM users;
```

Без `GROUP BY` агрегат возвращает **одну строку** на весь запрос — это нормально.

---

### [02. Количество всех пользователей и количество с email](./02-count-with-email.sql)

**Что нового.** Разница `COUNT(*)` vs `COUNT(column)`:
- `COUNT(*)` — все строки.
- `COUNT(email)` — только строки, где `email IS NOT NULL`.

Это самый быстрый способ узнать, сколько NULL в колонке: `COUNT(*) - COUNT(column)`.

**Алиасы.** Для именования колонок в выводе используй `AS`:

```sql
SELECT COUNT(*) AS total, COUNT(email) AS with_email
FROM users;
```

`AS` можно опустить (`COUNT(*) total`), но с `AS` читается явнее.

---

### [03. Пользователей по городам](./03-users-by-city.sql)

**Что нового.** `GROUP BY col` — собирает строки в группы по одинаковому значению `col`. К каждой группе можно применить агрегат.

```sql
SELECT city, COUNT(*) AS users_count
FROM users
GROUP BY city
ORDER BY users_count DESC;
```

**Правило.** Всё, что в `SELECT` и **не** агрегат, должно быть в `GROUP BY`. Иначе Postgres даст ошибку:
> column "city" must appear in the GROUP BY clause or be used in an aggregate function

**Что попробовать.** Что вернёт запрос, если убрать `GROUP BY city` (оставить только `SELECT city, COUNT(*)`)? Попробуй и прочти текст ошибки.

---

### [04. Средний возраст по городам](./04-avg-age-by-city.sql)

**Что нового.**
- `AVG(col)` — среднее. Возвращает дробное число (`numeric`).
- `WHERE` применяется **до** группировки: сначала отбрасываются неактивные, потом считаются группы.

**Production-нюанс.** Среднее по NULL: `AVG` **игнорирует** NULL-значения. Это удобно, но иногда нужно знать — например, если у части пользователей `age IS NULL`, среднее посчитается только по тем, у кого возраст есть.

---

### [05. Статистика заказов по статусам](./05-order-stats-by-status.sql)

**Что нового.** Несколько агрегатов в одном `SELECT` — обычный паттерн для дашбордов.

```sql
SELECT
    status,
    COUNT(*)    AS count,
    SUM(amount) AS total,
    AVG(amount) AS avg,
    MIN(amount) AS min,
    MAX(amount) AS max
FROM orders
GROUP BY status;
```

Каждый агрегат считается независимо по той же группе. Один проход по таблице — несколько метрик.

**Style.** Когда колонок много — каждая на новой строке, выравнивание `AS` помогает читать.

---

### [06. Города, где больше одного пользователя](./06-cities-with-multiple.sql)

**Что нового.** `HAVING` — фильтр **для групп**. Применяется **после** `GROUP BY`. `WHERE` так не умеет — `WHERE` работает на строках до группировки.

```sql
SELECT city, COUNT(*) AS users_count
FROM users
GROUP BY city
HAVING COUNT(*) > 1;
```

**Postgres-нюанс.** В `HAVING` нельзя использовать алиас из `SELECT` — повтори агрегат полностью:

```sql
HAVING COUNT(*) > 1   -- ✓ ok
HAVING users_count > 1 -- ✗ ошибка
```

**Что попробовать.** Что произойдёт, если поставить `WHERE COUNT(*) > 1` вместо `HAVING`? Прочти текст ошибки — это полезно для понимания, чем эти секции отличаются.

---

### [07. Крупные покупатели](./07-big-spenders.sql)

**Что нового.** Композиция: `WHERE` режет строки до группировки (только `paid`), `GROUP BY user_id` собирает по покупателю, `HAVING SUM(amount) > 200` фильтрует получившиеся группы.

**Порядок секций:**

```sql
SELECT ...
FROM ...
WHERE ...           -- до группировки
GROUP BY ...
HAVING ...          -- после группировки
ORDER BY ...
LIMIT ...
```

Это фиксированный порядок в SQL. Можно запомнить как мнемонику: **S**elect, **F**rom, **W**here, **G**roup, **H**aving, **O**rder, **L**imit.

---

### [08. Активные vs неактивные по городам](./08-active-vs-inactive-by-city.sql)

**Что нового.** `FILTER (WHERE условие)` — Postgres-style условный агрегат. Считает агрегат только по строкам, удовлетворяющим условию:

```sql
SELECT
    city,
    COUNT(*)                              AS total,
    COUNT(*) FILTER (WHERE is_active)     AS active,
    COUNT(*) FILTER (WHERE NOT is_active) AS inactive
FROM users
GROUP BY city;
```

Это **production-стандарт в Postgres** для нескольких метрик из одной выборки. До `FILTER` писали через `COUNT(CASE WHEN ... THEN 1 END)` — теперь так делают редко.

**Преимущество перед двумя запросами.** Один проход по таблице вместо двух — быстрее и меньше нагрузка на БД.

---

### [09. Уникальные пользователи с paid-заказами](./09-distinct-paid-users.sql)

**Что нового — `COUNT(DISTINCT col)`.** Считает количество **уникальных** значений колонки (NULL отбрасывается, как у обычного `COUNT(col)`).

```sql
SELECT COUNT(DISTINCT user_id) AS unique_paid_users
FROM orders
WHERE status = 'paid';
```

Если у одного `user_id` пять `paid`-заказов — он посчитается **один раз**. Без `DISTINCT` получишь общее число `paid`-заказов, а не покупателей.

**Где используется в реальности.** «Сколько уникальных посетителей», «сколько разных городов в базе», «MAU/DAU». Везде, где важно «штуки сущностей», а не «штуки строк».

**Production-нюанс.** `COUNT(DISTINCT)` дороже обычного `COUNT(*)` — Postgres должен дедуплицировать значения. На больших таблицах иногда заменяют на approximation (`approx_count_distinct` в расширениях, HyperLogLog). Это уже уровень 3 — пока просто запомни.

**Что попробовать.** Сравни результат `COUNT(*)` vs `COUNT(DISTINCT user_id)` на этом же фильтре — разница покажет, у скольких пользователей >1 paid-заказа.
