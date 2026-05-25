/*
Пакет database/sql — универсальный интерфейс к SQL-базам (через драйвер).

Зачем:   слой работы с БД в REST API — выборки, вставки, транзакции, пул соединений.
Когда:   любой сервис с реляционной БД (PostgreSQL, MySQL, SQLite).
Грабли:  значения подставляйте ТОЛЬКО плейсхолдерами ($1 / ?), никогда не клейте в строку запроса —
         это SQL-инъекция. sql.Open НЕ открывает соединение (это делает первый запрос/Ping).
         rows нужно закрывать (defer rows.Close) и проверять rows.Err() после цикла.

ВАЖНО: database/sql сам не общается с БД — нужен драйвер «пустым» импортом, например:
  _ "github.com/jackc/pgx/v5/stdlib"   // PostgreSQL
  _ "modernc.org/sqlite"               // SQLite (чистый Go, без cgo)
Файл компилируется на чистой stdlib; чтобы реально выполнить запросы, добавьте драйвер и DSN.
*/
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"
)

type User struct {
	ID    int
	Name  string
	Email sql.NullString // колонка может быть NULL → Null-тип, а не обычная строка
}

func main() {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		// => без драйвера и DSN выполнить запросы нельзя — печатаем подсказку и выходим.
		fmt.Println("DSN не задан. Добавьте драйвер (import _ \"...\") и DSN, чтобы запустить.")
		fmt.Println("Ниже — канонические паттерны для реального сервиса.")
		return
	}

	// Open создаёт ПУЛ, а не соединение. "pgx"/"sqlite" — имя, под которым драйвер себя зарегистрировал.
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Println("open:", err)
		return
	}
	defer db.Close()

	// Настройка пула обязательна в проде, иначе соединения исчерпаются.
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil { // реально проверить доступность БД
		fmt.Println("ping:", err)
		return
	}

	queryOne(ctx, db, 1)
	queryMany(ctx, db)
	createUser(ctx, db, "Dave", "dave@x.io")
}

// queryOne — ОДНА строка: QueryRowContext + Scan. ErrNoRows = «не найдено» → 404.
func queryOne(ctx context.Context, db *sql.DB, id int) {
	var u User
	err := db.QueryRowContext(ctx,
		`SELECT id, name, email FROM users WHERE id = $1`, id). // $1 — плейсхолдер, защита от инъекций
		Scan(&u.ID, &u.Name, &u.Email)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		fmt.Println("not found → 404")
	case err != nil:
		fmt.Println("query error:", err)
	default:
		fmt.Printf("user: %+v\n", u)
	}
}

// queryMany — МНОГО строк: цикл Next/Scan, обязательны Close и проверка Err после цикла.
func queryMany(ctx context.Context, db *sql.DB) {
	rows, err := db.QueryContext(ctx, `SELECT id, name, email FROM users ORDER BY id`)
	if err != nil {
		fmt.Println("query error:", err)
		return
	}
	defer rows.Close() // вернуть соединение в пул

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			fmt.Println("scan error:", err)
			return
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil { // ошибка могла случиться В ПРОЦЕССЕ итерации
		fmt.Println("rows error:", err)
		return
	}
	fmt.Printf("loaded %d users\n", len(users))
}

// createUser — запись в транзакции. Паттерн: Begin → defer Rollback → Commit.
func createUser(ctx context.Context, db *sql.DB, name, email string) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("begin:", err)
		return
	}
	// defer Rollback безопасен: после успешного Commit станет no-op. Гарантирует откат при ошибке/панике.
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2)`, name, email); err != nil {
		fmt.Println("insert:", err)
		return
	}
	if err := tx.Commit(); err != nil {
		fmt.Println("commit:", err)
		return
	}
	fmt.Println("user created")
}

/*
Что запомнить (что чаще и почему):
  • QueryRow vs Query:
      QueryRowContext — ровно ОДНА строка (поиск по id, COUNT, EXISTS). Scan сразу, без цикла.
                        Если строки нет — Scan вернёт sql.ErrNoRows (это и есть ваш 404).
      QueryContext    — НЕСКОЛЬКО строк: цикл rows.Next()/rows.Scan(), потом Close + Err.
  • Exec vs Query: Exec* — для INSERT/UPDATE/DELETE (не возвращают строк, дают кол-во затронутых через
    Result.RowsAffected). Query* — для SELECT. Перепутать — частая ошибка.
  • Всегда *Context-версии (ExecContext/QueryContext): пробрасывают таймаут и отмену из ctx запроса.
    Версии без Context (Exec/Query) — только для скриптов без контекста.
  • Плейсхолдеры различаются по драйверу: $1,$2 (Postgres/pgx) и ?,? (MySQL/SQLite). Значения —
    только через них, НИКОГДА через fmt.Sprintf в текст запроса.
  • Транзакция: Begin → defer tx.Rollback() → ... → tx.Commit(). defer Rollback после Commit безвреден.
  • NULL-колонки читайте в sql.NullString/NullInt64/... — обычная строка не примет NULL и Scan упадёт.

Типичные сценарии:
  1) Получить по id:  db.QueryRowContext(ctx, "...WHERE id=$1", id).Scan(&u.ID, &u.Name)
  2) Список:          rows, _ := db.QueryContext(ctx, "SELECT ..."); defer rows.Close(); for rows.Next() {...}
  3) Запись:          db.ExecContext(ctx, "INSERT INTO ... VALUES ($1,$2)", a, b)
*/
