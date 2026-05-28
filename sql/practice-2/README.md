# Практика 2 — NULL и текстовые фильтры

Конструкции: `IS NULL`, `IS NOT NULL`, `LIKE`, `ILIKE`, `%`.

Таблицы: `users`, `orders`.

Главные правила:
- `NULL` нельзя проверять через `=` или `<>`. Только `IS NULL` / `IS NOT NULL`.
- `LIKE` — регистрозависимый.
- `ILIKE` (Postgres-specific) — регистронезависимый.
- `%` — wildcard для любого числа символов. `_` — для одного.

| # | Файл                                                              | Что нового                  |
|---|-------------------------------------------------------------------|-----------------------------|
| 1 | [01-without-email.sql](./01-without-email.sql)                    | `IS NULL`                   |
| 2 | [02-with-email.sql](./02-with-email.sql)                          | `IS NOT NULL`               |
| 3 | [03-names-starting-with-a.sql](./03-names-starting-with-a.sql)    | `LIKE 'A%'`                 |
| 4 | [04-case-insensitive-search.sql](./04-case-insensitive-search.sql)| `ILIKE`                     |
| 5 | [05-gmail-users.sql](./05-gmail-users.sql)                        | `LIKE '%@gmail.com'`        |
| 6 | [06-unpaid-orders.sql](./06-unpaid-orders.sql)                    | `IS NULL` на timestamp      |
| 7 | [07-cancelled-orders.sql](./07-cancelled-orders.sql)              | `IS NOT NULL` + сортировка  |
| 8 | [08-active-with-domain.sql](./08-active-with-domain.sql)          | композит NULL + ILIKE + AND |

---

## Подсказки и теория по задачам

### [01. Пользователи без email](./01-without-email.sql)

**Что нового.** `NULL` нельзя сравнивать через `=`. Запись `WHERE email = NULL` всегда вернёт **пустой результат** (это специфика NULL-логики в SQL). Используй только `IS NULL`.

---

### [02. Пользователи с email](./02-with-email.sql)

**Что нового.** `IS NOT NULL` — отрицание `IS NULL`. Тоже не сравнение, а специальный оператор.

---

### [03. Имена на букву A](./03-names-starting-with-a.sql)

**Что нового.**
- `LIKE 'pattern'` — паттерн-матчинг по строке.
- `%` — wildcard для любого количества символов (включая ноль).
- `_` (подчёркивание) — wildcard для **ровно одного** символа.

Примеры паттернов:
- `'A%'` — начинается на `A`
- `'%a'` — заканчивается на `a`
- `'%al%'` — содержит `al` где-то внутри

`LIKE` чувствителен к регистру: `'a%'` и `'A%'` дадут разный результат.

---

### [04. Поиск без учёта регистра](./04-case-insensitive-search.sql)

**Что нового.** `ILIKE` — Postgres-specific версия `LIKE`, нечувствительная к регистру. Это то, что чаще всего нужно в поисках в реальных API (search by name, filter by email substring).

Альтернатива в стандартном SQL: `WHERE LOWER(name) LIKE '%al%'`. Но в Postgres `ILIKE` идиоматичнее.

**Production-нюанс.** Запрос вида `ILIKE '%al%'` не использует обычный B-tree индекс (паттерн начинается с `%`). Для production-поиска применяют `pg_trgm` + GIN-индекс. Сейчас это не важно, но запомни — пригодится на уровне 2.

---

### [05. Пользователи на gmail.com](./05-gmail-users.sql)

**Подсказка.** Паттерн: `'%@gmail.com'`.

**Думай о NULL.** Если у пользователя `email = NULL`, `email LIKE '%@gmail.com'` вернёт NULL (а не FALSE), и эта строка не попадёт в выдачу — то есть в данном случае получится без явной проверки. Но привыкай явно фильтровать NULL, когда это важно: добавь `email IS NOT NULL` — это и читаемее, и безопаснее.

---

### [06. Неоплаченные заказы](./06-unpaid-orders.sql)

**Что нового.** `IS NULL` работает не только со строками — но и с любым nullable-полем (timestamp, integer, boolean). В `orders` поле `paid_at TIMESTAMP` опционально.

**Что попробовать.** Сравни количество строк, которые вернёт:
- `WHERE paid_at IS NULL`
- `WHERE status <> 'paid'`

Это **не одно и то же**. Подумай, чем именно.

---

### [07. Отменённые заказы — последние сверху](./07-cancelled-orders.sql)

**Что нового.** `ORDER BY` на timestamp работает так же, как на числах: `DESC` — от свежих к старым.

---

### [08. Активные пользователи с корпоративным email](./08-active-with-domain.sql)

**Что нового.** Комбинация трёх независимых условий через `AND`. Production-style — каждое условие на новой строке, выравнено:

```sql
WHERE is_active
  AND email IS NOT NULL
  AND email ILIKE '%@company.com'
```

Так читается быстрее, чем «склейка в одну строку».

**Что попробовать.** Что произойдёт, если убрать `email IS NOT NULL`? Останется ли результат тот же? Почему?
