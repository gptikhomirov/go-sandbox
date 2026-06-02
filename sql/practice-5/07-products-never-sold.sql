/*
# 07. Товары, которые никогда не продавались

Вывести товары, которых **нет** ни в одной строке `order_items`.

**Таблицы:** `products`, `order_items`
**Колонки на выходе:** `p.id`, `p.name`, `p.category`
*/

SELECT p.id,
       p.name,
       p.category
FROM products p
         LEFT JOIN order_items oi ON p.id = oi.product_id
WHERE oi.id IS NULL;
