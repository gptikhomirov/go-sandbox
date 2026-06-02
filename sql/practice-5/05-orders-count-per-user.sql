/*
# 05. Количество заказов у каждого пользователя

Для каждого пользователя вывести его имя и количество заказов. Пользователи без заказов должны быть в результате с `0`. Отсортировать по убыванию количества.

**Таблицы:** `users`, `orders`
**Колонки на выходе:** `u.id`, `u.name`, `orders_count`
*/

SELECT u.id        AS user_id,
       u.name      AS user_name,
       COUNT(o.id) AS orders_count
FROM users u
         LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name
ORDER BY orders_count DESC;
