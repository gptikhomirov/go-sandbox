# SQL Sandbox

Учебная среда для практики PostgreSQL под Go-backend.
Поднимает локальный Postgres в Docker и автоматически заливает схему + seed-данные, покрывающие 5 практик базового уровня (от `SELECT` до `JOIN` + транзакций).

## Что нужно поставить

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (или совместимая среда: Colima, OrbStack)
- SQL-клиент: **GoLand / DataGrip / DBeaver** или `psql`

## Запуск с нуля

```bash
cd sql
docker compose up -d
docker compose ps                                                  # status: healthy
docker compose exec postgres psql -U postgres -d sandbox -c '\dt'  # должно быть 6 таблиц
```

Креды (локальные, учебные):

| param    | value     |
|----------|-----------|
| host     | localhost |
| port     | 5432      |
| database | sandbox   |
| user     | postgres  |
| password | postgres  |

## Reset базы между практиками

Volume `sandbox-pgdata` хранит данные между рестартами. Чтобы вернуть базу к исходному seed-состоянию:

```bash
docker compose down -v
docker compose up -d
```

Флаг `-v` удаляет volume — при следующем `up` Postgres снова прогонит `init/01_schema.sql`.

## Остановить / удалить

```bash
docker compose stop      # пауза, данные остаются
docker compose down      # удалить контейнер, данные остаются
docker compose down -v   # удалить контейнер + данные
```

## Подключение из GoLand / DataGrip

1. Database tool window → `+` → **Data Source → PostgreSQL**.
2. Введите параметры из таблицы выше.
3. При первом запуске GoLand предложит **Download** JDBC-драйвер → согласиться.
4. **Test Connection** → Apply.
5. Запросы — правый клик на источнике → **New → Query Console**.

Из терминала:

```bash
docker compose exec postgres psql -U postgres -d sandbox
```

## Как идёт обучение

Учим SQL слоями. Базовый уровень разбит на 5 практик:

| # | Тема                          | Ключевые конструкции                                                  | Таблицы                                          |
|---|-------------------------------|-----------------------------------------------------------------------|--------------------------------------------------|
| 1 | Чтение и фильтрация           | `SELECT`, `WHERE`, `AND/OR/NOT`, `IN`, `BETWEEN`, `ORDER BY`, `LIMIT` | `users`, `orders`, `products`                    |
| 2 | NULL и текстовые фильтры      | `IS NULL`, `IS NOT NULL`, `LIKE`, `ILIKE`                             | `users`, `orders`                                |
| 3 | Агрегаты и группировка        | `COUNT`, `SUM`, `AVG`, `MIN`, `MAX`, `GROUP BY`, `HAVING`, `FILTER`   | `users`, `orders`, `products`, `payments`, `order_items` |
| 4 | Изменение данных              | `INSERT`, `UPDATE`, `DELETE`, `RETURNING`, `ON CONFLICT`              | `users`, `orders`, `products`, `accounts`        |
| 5 | Простые связи и транзакции    | `JOIN`, `LEFT JOIN`, `ON`, `BEGIN/COMMIT/ROLLBACK`                    | все                                              |

Полный детальный план для LLM-наставника (с правилами проверки, production-стилем, типами задач) — в [PROMPT.md](./PROMPT.md). Скопируйте его в системный промпт любого LLM-агента (Claude / ChatGPT / Cursor / etc.) и попросите вести вас по практикам.

## Структура

```
sql/
├── docker-compose.yml      # postgres:16-alpine + healthcheck + volume
├── init/
│   └── 01_schema.sql       # DROP + CREATE + INSERT для всех таблиц
├── PROMPT.md               # системный промпт для SQL-наставника
└── README.md
```

## Таблицы

- `users` — пользователи (email/без email, активные/неактивные, soft-deleted)
- `products` — товары разных категорий, часть неактивна
- `orders` — заказы разных статусов (`paid`, `pending`, `cancelled`, `refunded`)
- `order_items` — позиции заказов (M:N между orders и products)
- `accounts` — счета пользователей (с блокировками)
- `payments` — переводы между счетами (`completed`, `pending`, `failed`)

Полные определения и seed смотри в [init/01_schema.sql](./init/01_schema.sql).
