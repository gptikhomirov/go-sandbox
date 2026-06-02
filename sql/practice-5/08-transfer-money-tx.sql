/*
# 08. Перевод денег — атомарная транзакция

Перевести 100 у.е. с аккаунта `id = 1` на аккаунт `id = 2` и записать в `payments` строку про этот перевод.

Условия:
- списать 100 с `accounts.id = 1`;
- зачислить 100 на `accounts.id = 2`;
- вставить новую запись в `payments` (`status = 'completed'`, `created_at = NOW()`, `id = 100`);
- всё должно произойти **атомарно** — либо всё, либо ничего.

**Таблицы:** `accounts`, `payments`
*/

BEGIN;

UPDATE accounts
SET balance = balance - 100
WHERE id = 1;

UPDATE accounts
SET balance = balance + 100
WHERE id = 2;

INSERT INTO payments (id, from_account_id, to_account_id, amount, status, created_at)
VALUES (11, 1, 2, 100, 'completed', CURRENT_TIMESTAMP);

COMMIT;