/*
# 10. Что купила Anna

Вывести список товаров, которые покупала пользователь с именем `Anna`, с количеством и ценой за единицу в момент покупки.

**Таблицы:** `users`, `orders`, `order_items`, `products`
**Колонки на выходе:** `o.id` (как `order_id`), `o.status`, `p.name` (как `product_name`), `oi.quantity`, `oi.unit_price`

## Условия

- только заказы пользователя `Anna`
- учитывать заказы любого статуса
- одна строка на каждую позицию в заказе (один заказ — несколько строк, если в нём несколько товаров)
- сортировка: `o.id ASC`, затем `p.name ASC`
*/

SELECT o.id AS order_id, o.status, p.name AS product_name, oi.quantity, oi.unit_price
FROM orders o
         JOIN users u ON o.user_id = u.id
         JOIN order_items oi ON o.id = oi.order_id
         JOIN products p ON oi.product_id = p.id
WHERE u.name = 'Anna'
ORDER BY o.id, p.name;
