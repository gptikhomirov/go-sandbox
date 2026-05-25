/*
Пакет database/sql — работа с базой данных (SQL).

⚠️ ТЕМА ПОСЛОЖНЕЕ И ТРЕБУЕТ ОТДЕЛЬНОЙ НАСТРОЙКИ. Если только начал — сперва освой остальное.
Сюда вернись, когда захочешь хранить данные «по-взрослому» (в базе, а не в файле/памяти).

Что такое база данных простыми словами: программа-хранилище, где данные лежат в таблицах
(строки и столбцы), а достают их запросами на языке SQL, например:
    SELECT name FROM users WHERE id = 1   — «дай имя пользователя с номером 1».

Важные мысли:
  - database/sql — ОБЩИЙ интерфейс. Чтобы говорить с конкретной базой, нужен «драйвер»
    (отдельная библиотека под PostgreSQL, SQLite и т.п.), подключаемый «пустым» импортом:
        import _ "modernc.org/sqlite"
  - Значения в запрос вставляй ТОЛЬКО через «плейсхолдеры» ($1, ?), НИКОГДА не склеивай строку —
    это защита от взлома (SQL-инъекций).

Файл КОМПИЛИРУЕТСЯ без базы; запросы выполнятся, только если задать переменную DSN и драйвер.
Он показывает, КАК выглядит правильный код — его копируют в реальный проект.

Как запустить:  go run main.go   (без базы просто напечатает подсказку)
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

// Структура под одну строку таблицы users.
// sql.NullString — для столбца, который МОЖЕТ быть пустым (NULL). Обычная строка NULL не примет.
type User struct {
	ID    int
	Name  string
	Email sql.NullString
}

func main() {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		fmt.Println("База не настроена (нет переменной DSN).")
		fmt.Println("Код ниже показывает, как правильно работать с базой — изучи его как образец.")
		return
	}

	// Open готовит ПОДКЛЮЧЕНИЕ (точнее, пул соединений). "pgx" — имя драйвера.
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Println("не смог открыть базу:", err)
		return
	}
	defer db.Close() // закрыть подключение в конце

	// Контекст с лимитом времени, чтобы запрос не висел вечно (см. пример context).
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// PingContext реально проверяет, что база доступна (Open этого ещё не делает).
	if err := db.PingContext(ctx); err != nil {
		fmt.Println("база недоступна:", err)
		return
	}

	getUser(ctx, db, 1)
	listUsers(ctx, db)
	createUser(ctx, db, "Дима", "dima@mail.com")
}

// Одна строка: QueryRowContext + Scan. Нет строки -> ошибка ErrNoRows.
func getUser(ctx context.Context, db *sql.DB, id int) {
	var u User
	// $1 — плейсхолдер, вместо него БЕЗОПАСНО подставится id. Scan кладёт столбцы в переменные (& — «сюда»).
	err := db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Name, &u.Email)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("пользователь не найден")
		return
	}
	if err != nil {
		fmt.Println("ошибка запроса:", err)
		return
	}
	fmt.Printf("нашли: %+v\n", u)
}

// Много строк: QueryContext + цикл Next/Scan. Обязательны Close и проверка Err после цикла.
func listUsers(ctx context.Context, db *sql.DB) {
	rows, err := db.QueryContext(ctx, "SELECT id, name FROM users ORDER BY id")
	if err != nil {
		fmt.Println("ошибка запроса:", err)
		return
	}
	defer rows.Close() // вернуть соединение в пул

	var users []User
	for rows.Next() { // переходим к следующей строке, пока они есть
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			fmt.Println("ошибка чтения строки:", err)
			return
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil { // ошибка могла случиться В ПРОЦЕССЕ перебора
		fmt.Println("ошибка после перебора:", err)
		return
	}
	fmt.Printf("всего пользователей: %d\n", len(users))
}

// Изменение данных в транзакции (несколько действий «всё или ничего»): Begin -> Rollback/Commit.
func createUser(ctx context.Context, db *sql.DB, name, email string) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("не начать транзакцию:", err)
		return
	}
	// defer Rollback безопасен: после успешного Commit он ничего не делает.
	// Гарантирует откат, если до Commit случится ошибка.
	defer tx.Rollback()

	// ExecContext — для запросов БЕЗ результата-строк (INSERT/UPDATE/DELETE).
	_, err = tx.ExecContext(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
	if err != nil {
		fmt.Println("ошибка вставки:", err)
		return
	}

	if err := tx.Commit(); err != nil { // подтвердить изменения
		fmt.Println("ошибка commit:", err)
		return
	}
	fmt.Println("пользователь добавлен")
}

/*
Что важно запомнить:
  • База хранит данные в таблицах; запросы пишут на SQL (SELECT/INSERT/UPDATE/DELETE).
  • database/sql сам не работает с базой — нужен ДРАЙВЕР (отдельная библиотека, «пустой» импорт).
  • Open готовит подключение, PingContext проверяет доступность.
  • Получить данные:
      QueryRowContext + Scan — ОДНА строка (нет -> sql.ErrNoRows, это твой «не найдено»).
      QueryContext + Next/Scan — МНОГО строк (не забудь rows.Close() и rows.Err()).
  • Изменить данные: ExecContext (INSERT/UPDATE/DELETE).
  • Несколько изменений «всё или ничего» — транзакция: BeginTx -> defer Rollback -> Commit.
  • Значения ТОЛЬКО через плейсхолдеры ($1/?), иначе SQL-инъекция.
  • Столбцы, которые могут быть NULL, читай в sql.NullString/NullInt64/NullBool/NullTime.

Чтобы реально запустить:
  1) Подключи драйвер, например SQLite (чистый Go): import _ "modernc.org/sqlite", sql.Open("sqlite", ...).
  2) Создай таблицу users (id, name, email) и задай переменную DSN с адресом базы.
*/
