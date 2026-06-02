/*
# 01. Заказы с именами пользователей

Вывести все заказы и имя пользователя, сделавшего заказ.

**Таблицы:** `orders`, `users`
**Колонки на выходе:** `o.id` (как `order_id`), `o.status`, `o.amount`, `u.name` (как `user_name`)
*/

SELECT o.id   AS order_id,
       o.status,
       o.amount,
       u.name AS user_name
FROM orders o
         JOIN users u ON o.user_id = u.id;

