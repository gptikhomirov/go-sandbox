/*
# 04. Пользователи без paid-заказов

Вывести пользователей, у которых **нет** ни одного `paid`-заказа. Сюда должны попасть:
- те, кто вообще без заказов;
- те, у кого есть только `pending` / `cancelled` / `refunded`.

**Таблицы:** `users`, `orders`
**Колонки на выходе:** `u.id` (как `user_id`), `u.name` (как `user_name`)
*/

-- Вариант 1 - не правильно
-- SELECT u.id AS user_id, u.name AS user_name
-- FROM users u
--          LEFT JOIN orders o ON u.id = o.user_id
-- WHERE status != 'paid'
--    OR status ISNULL;

-- Вариант 2 - правильно
SELECT u.id AS user_id, u.name AS user_name
FROM users u
         LEFT JOIN orders o ON u.id = o.user_id AND o.status = 'paid'
WHERE status ISNULL;
