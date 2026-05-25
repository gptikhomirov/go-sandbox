/*
Пакет fmt — форматированный ввод/вывод.

Зачем:   выводить данные, собирать строки из значений, оборачивать ошибки контекстом.
Когда:   отладка, логи, построение сообщений/ключей, запись в любой io.Writer (HTTP-ответ).
Грабли:  %v скрывает тип (легко спутать "1" и 1) — для отладки берите %#v или %T;
         Printf без \n не переносит строку; неверный глагол даёт в выводе %!d(...).
*/
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Println — пробелы между аргументами + перенос строки. Для быстрой отладки.
	fmt.Println("user:", "alice", "age:", 30) // => user: alice age: 30

	// Printf — точный контроль через глаголы (verbs):
	//   %d целое · %s строка · %.2f float c 2 знаками · %v «как есть»
	//   %+v структура с именами полей · %q строка в кавычках · %T тип
	user := struct {
		Name string
		Age  int
	}{"alice", 30}
	fmt.Printf("name=%s age=%d\n", user.Name, user.Age) // => name=alice age=30
	fmt.Printf("struct=%+v\n", user)                    // => struct={Name:alice Age:30}
	fmt.Printf("price=%.2f quoted=%q\n", 9.5, "hi")     // => price=9.50 quoted="hi"

	// Sprintf — собрать строку БЕЗ вывода. Ключи кэша, пути, сообщения.
	cacheKey := fmt.Sprintf("user:%d:profile", user.Age)
	fmt.Println(cacheKey) // => user:30:profile

	// Fprintf — пишет в io.Writer. В реальном коде это http.ResponseWriter, файл, буфер.
	var b strings.Builder
	fmt.Fprintf(&b, "report: %d items", 42)
	fmt.Println(b.String())                              // => report: 42 items
	fmt.Fprintln(os.Stderr, "warning: low disk space")  // предупреждения — в stderr

	// Errorf — ошибка с контекстом. %w оборачивает исходную, сохраняя её для errors.Is/As.
	err := fmt.Errorf("load config: %w", os.ErrNotExist)
	fmt.Println(err) // => load config: file does not exist
}

/*
Что запомнить (что чаще и почему):
  • Println — только для отладки/CLI. В сервисах для логов берут log или log/slog.
  • Printf vs Sprintf: одинаковые глаголы, но Printf печатает, Sprintf возвращает строку.
    Нужна строка для дальнейшего использования (ключ, путь) — Sprintf; нужен вывод — Printf.
  • Fprintf — обобщение: Printf это Fprintf(os.Stdout, ...). Пишите сразу в нужный Writer
    (ResponseWriter, буфер), а не собирайте строку и потом выводите.
  • Errorf с %w — стандарт оборачивания ошибок. Без %w (через %v) причина теряется
    и errors.Is/As её уже не найдут.

Типичные сценарии:
  1) HTTP-ответ:     fmt.Fprintf(w, "id=%d", id)
  2) Ключ кэша:      key := fmt.Sprintf("user:%d", id)
  3) Ошибка слоя:    return fmt.Errorf("save user: %w", err)
*/
