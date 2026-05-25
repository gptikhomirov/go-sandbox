/*
Пакет database/sql (уровень 2) — только новое поверх base.
В base уже: Open, Ping, QueryRow/Query/Exec, BeginTx, Scan, Rows.*, Commit/Rollback, ErrNoRows, Null-типы.
Здесь: настройка ПУЛА соединений, prepared statements, метрики пула, named-параметры, метаданные результата.

Зачем: в проде важно не «как сделать запрос», а как настроить пул под нагрузку, переиспользовать
подготовленные запросы и снимать метрики. Это и спрашивают на собеседовании.

ВАЖНО: как и в base, нужен драйвер и DSN. Файл компилируется на stdlib и показывает паттерны.

Как запустить:  go run main.go   (без DSN напечатает подсказку)
*/
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

func main() {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		fmt.Println("Нет DSN. Код ниже — production-настройки пула и prepared statements.")
		return
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Println("open:", err)
		return
	}
	defer db.Close()

	// Настройка ПУЛА — ключевое для прода. Без неё под нагрузкой кончатся соединения или БД ляжет.
	db.SetMaxOpenConns(25)                 // максимум одновременных соединений к БД
	db.SetMaxIdleConns(25)                 // сколько держать «про запас» (idle)
	db.SetConnMaxLifetime(5 * time.Minute) // как долго живёт соединение (защита от «протухших»)
	db.SetConnMaxIdleTime(time.Minute)     // как долго простаивает перед закрытием

	// Stats — метрики пула (для мониторинга: сколько соединений занято/ждут).
	stats := db.Stats()
	fmt.Println("открытых соединений:", stats.OpenConnections)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// PrepareContext — «подготовленный» запрос: БД один раз разбирает SQL, дальше выполняем много раз
	// с разными аргументами (быстрее и безопаснее для повторяющихся запросов, например в цикле).
	stmt, err := db.PrepareContext(ctx, "SELECT name FROM users WHERE id = $1")
	if err != nil {
		fmt.Println("prepare:", err)
		return
	}
	defer stmt.Close()
	for _, id := range []int{1, 2, 3} {
		var name string
		if err := stmt.QueryRowContext(ctx, id).Scan(&name); err == nil {
			fmt.Println(id, name)
		}
	}

	// Named-параметры — обращение к плейсхолдеру по ИМЕНИ (читабельнее, чем $1/$2 при многих аргументах).
	db.QueryRowContext(ctx,
		"SELECT name FROM users WHERE id = @id",
		sql.Named("id", 1))

	// Метаданные результата: Columns — имена колонок, ColumnTypes — их типы (для generic-обработки).
	rows, err := db.QueryContext(ctx, "SELECT * FROM users LIMIT 1")
	if err == nil {
		cols, _ := rows.Columns()
		fmt.Println("колонки:", cols)
		rows.Close()
	}
}

/*
Что важно запомнить:
  • Настройка пула — обязательна в проде:
      SetMaxOpenConns  — потолок одновременных соединений (защищает БД от перегруза);
      SetMaxIdleConns  — сколько держать наготове (обычно = MaxOpenConns, чтобы не пересоздавать);
      SetConnMaxLifetime — закрывать старые соединения (балансировщики/файрволы рвут «висящие»).
    Дефолты неудачные (MaxIdle=2), поэтому их почти всегда переопределяют.
  • PrepareContext выгоден для ПОВТОРЯЮЩИХСЯ запросов (в цикле): SQL разбирается один раз. Не забудь Close.
  • db.Stats() — отдавай в метрики (Prometheus): рост WaitCount/занятых соединений = пул мал или запросы медленные.
  • sql.Named — именованные параметры, когда $1/$2/$3 становится трудно читать.

Задача:
  1) Подбери MaxOpenConns под свою БД (например, 25) и снимай db.Stats().InUse в фоне — посмотри поведение под нагрузкой.
*/
