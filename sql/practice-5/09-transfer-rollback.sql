/*
# 09. Транзакция с ROLLBACK

Открыть транзакцию, **временно** изменить балансы и вставить запись в payments — а потом всё откатить.
Убедиться, что после `ROLLBACK` базе ничего не сделалось.

**Таблицы:** `accounts`, `payments`
*/

BEGIN;

UPDATE accounts
SET balance = balance - 999
WHERE id = 1;
UPDATE accounts
SET balance = balance + 999
WHERE id = 2;
INSERT INTO payments (id, from_account_id, to_account_id, amount, status, created_at)
VALUES (200, 1, 2, 999, 'completed', NOW());

-- внутри транзакции эти изменения видны
SELECT id, balance
FROM accounts
WHERE id IN (1, 2);
SELECT *
FROM payments
WHERE id = 200;

ROLLBACK;

-- после ROLLBACK — всё откатилось, базы как будто не трогали
SELECT id, balance
FROM accounts
WHERE id IN (1, 2);
SELECT *
FROM payments
WHERE id = 200;