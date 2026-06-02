/*
# 02. Paid-заказы с email пользователей

Вывести только оплаченные заказы и email пользователей. Только тех, у кого email указан. Отсортировать по убыванию `amount`.

**Таблицы:** `orders`, `users`
**Колонки на выходе:** `o.id`, `o.amount`, `u.name`, `u.email`
*/

SELECT o.id, o.amount, u.name, u.email
FROM orders o
         JOIN users u ON o.user_id = u.id
WHERE u.email IS NOT NULL
  AND o.status = 'paid'
ORDER BY o.amount DESC;
