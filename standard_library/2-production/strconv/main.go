/*
Пакет strconv (уровень 2) — только новое поверх base.
В base уже: Atoi, Itoa, ParseInt, ParseFloat, ParseBool, FormatInt, FormatFloat.
Здесь: ParseUint, FormatUint, FormatBool, Quote/Unquote, QuoteRune, AppendInt/AppendFloat.

Зачем: беззнаковые числа (ParseUint), экранирование строк в кавычки (Quote — для логов/генерации кода),
и добавление чисел в []byte без лишних аллокаций (Append*).

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"strconv"
)

func main() {
	// ParseUint — для БЕЗ знаковых (счётчики, размеры, id, которые не бывают отрицательными).
	u, _ := strconv.ParseUint("4096", 10, 64)
	fmt.Println(u) // => 4096

	// FormatUint / FormatBool — обратно в строку.
	fmt.Println(strconv.FormatUint(255, 16)) // => ff (в 16-ричной)
	fmt.Println(strconv.FormatBool(true))    // => true

	// Quote оборачивает строку в кавычки и экранирует спецсимволы (\n, \t, кавычки).
	// Удобно для логов и сообщений: видно «настоящее» содержимое строки.
	fmt.Println(strconv.Quote("строка\tс\tтабами\n")) // => "строка\tс\tтабами\n"
	// Unquote — обратно: снять кавычки и раскрыть экранирование.
	unq, _ := strconv.Unquote(`"hello\nworld"`)
	fmt.Printf("%q\n", unq) // => "hello\nworld"
	// QuoteRune — то же для одного символа.
	fmt.Println(strconv.QuoteRune('я')) // => 'я'

	// AppendInt дописывает число в УЖЕ существующий []byte, не создавая новую строку.
	// Полезно в горячем коде (сборка большого вывода без аллокаций).
	buf := []byte("id=")
	buf = strconv.AppendInt(buf, 42, 10)
	fmt.Println(string(buf)) // => id=42
}

/*
Что важно запомнить:
  • ParseUint/FormatUint — когда число заведомо неотрицательное (uint64). Иначе — Parse/FormatInt.
  • Quote vs обычный вывод: Quote показывает спецсимволы явно (\n, \t) и экранирует — незаменим в логах,
    когда нужно понять, что РЕАЛЬНО в строке (есть ли пробелы/переводы). Тот же эффект даёт %q в Printf.
  • Append* (AppendInt/AppendFloat) — оптимизация: дописать в []byte без промежуточной строки.
    В обычном коде проще Itoa/FormatInt; Append* берут, когда важна производительность сборки.

Задача:
  1) Выведи строку "a\tb" через strconv.Quote и сравни с обычным fmt.Println той же строки.
*/
