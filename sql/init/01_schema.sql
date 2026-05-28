DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    age INTEGER,
    city TEXT NOT NULL,
    email TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE products (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    status TEXT NOT NULL,
    amount INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    paid_at TIMESTAMP,
    cancelled_at TIMESTAMP
);

CREATE TABLE order_items (
    id INTEGER PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    unit_price INTEGER NOT NULL
);

CREATE TABLE accounts (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    balance INTEGER NOT NULL,
    is_blocked BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE payments (
    id INTEGER PRIMARY KEY,
    from_account_id INTEGER NOT NULL REFERENCES accounts(id),
    to_account_id INTEGER NOT NULL REFERENCES accounts(id),
    amount INTEGER NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

INSERT INTO users (id, name, age, city, email, is_active, created_at, deleted_at) VALUES
(1,  'Alex',  25, 'Moscow', 'alex@mail.com', true,  '2025-01-01 10:00:00', NULL),
(2,  'Maria', 31, 'Berlin', 'maria@gmail.com', true, '2025-01-02 11:00:00', NULL),
(3,  'Ivan',  17, 'Moscow', NULL, false, '2025-01-03 12:00:00', NULL),
(4,  'Olga',  22, 'Paris', 'olga@mail.com', true, '2025-01-04 13:00:00', NULL),
(5,  'Kate',  19, 'London', 'kate@gmail.com', true, '2025-01-05 14:00:00', NULL),
(6,  'John',  40, 'Berlin', NULL, false, '2025-01-06 15:00:00', NULL),
(7,  'Bob',   28, 'Rome', 'bob@company.com', true, '2025-01-07 16:00:00', NULL),
(8,  'Anna',  21, 'Paris', 'anna@gmail.com', true, '2025-01-08 17:00:00', NULL),
(9,  'Alex',  35, 'Madrid', 'alex.work@company.com', true, '2025-01-09 18:00:00', NULL),
(10, 'Tom',   45, 'London', NULL, true, '2025-01-10 19:00:00', '2025-02-01 10:00:00');

INSERT INTO products (id, name, category, price, is_active) VALUES
(1, 'Keyboard', 'electronics', 100, true),
(2, 'Mouse', 'electronics', 50, true),
(3, 'Monitor', 'electronics', 300, true),
(4, 'Desk', 'furniture', 250, true),
(5, 'Chair', 'furniture', 150, true),
(6, 'Notebook', 'stationery', 10, true),
(7, 'Pen', 'stationery', 5, false),
(8, 'Webcam', 'electronics', 120, true);

INSERT INTO orders (id, user_id, status, amount, created_at, paid_at, cancelled_at) VALUES
(1,  1, 'paid',      100, '2025-02-01 10:00:00', '2025-02-01 10:05:00', NULL),
(2,  1, 'pending',   250, '2025-02-02 10:00:00', NULL, NULL),
(3,  2, 'paid',      300, '2025-02-03 10:00:00', '2025-02-03 10:10:00', NULL),
(4,  4, 'cancelled', 150, '2025-02-04 10:00:00', NULL, '2025-02-04 10:20:00'),
(5,  6, 'pending',   500, '2025-02-05 10:00:00', NULL, NULL),
(6,  7, 'refunded',  200, '2025-02-06 10:00:00', '2025-02-06 10:05:00', NULL),
(7,  8, 'paid',      450, '2025-02-07 10:00:00', '2025-02-07 10:05:00', NULL),
(8,  8, 'paid',       60, '2025-02-08 10:00:00', '2025-02-08 10:05:00', NULL),
(9,  9, 'cancelled', 700, '2025-02-09 10:00:00', NULL, '2025-02-09 10:30:00');

INSERT INTO order_items (id, order_id, product_id, quantity, unit_price) VALUES
(1, 1, 1, 1, 100),
(2, 2, 4, 1, 250),
(3, 3, 3, 1, 300),
(4, 4, 5, 1, 150),
(5, 5, 3, 1, 300),
(6, 5, 4, 1, 200),
(7, 6, 8, 1, 120),
(8, 6, 2, 2, 40),
(9, 7, 3, 1, 300),
(10, 7, 5, 1, 150),
(11, 8, 6, 6, 10),
(12, 9, 3, 2, 300),
(13, 9, 2, 2, 50);

INSERT INTO accounts (id, user_id, balance, is_blocked) VALUES
(1, 1, 1000, false),
(2, 2, 500, false),
(3, 3, 200, false),
(4, 4, 0, false),
(5, 5, 1500, false),
(6, 6, 50, true),
(7, 7, 700, false),
(8, 8, 900, false),
(9, 9, 300, false),
(10, 10, 100, true);

INSERT INTO payments (id, from_account_id, to_account_id, amount, status, created_at) VALUES
(1, 1, 2, 100, 'completed', '2025-03-01 10:00:00'),
(2, 2, 1, 50, 'completed', '2025-03-02 10:00:00'),
(3, 5, 3, 200, 'pending', '2025-03-03 10:00:00'),
(4, 7, 8, 150, 'failed', '2025-03-04 10:00:00'),
(5, 8, 1, 300, 'completed', '2025-03-05 10:00:00');
