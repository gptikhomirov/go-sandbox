/*
# 03. Пользователи без заказов

Вывести пользователей, у которых **ни одного** заказа.

**Таблицы:** `users`, `orders`
**Колонки на выходе:** `u.id`, `u.name`
*/

SELECT u.id, u.name
FROM users u
         LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;
