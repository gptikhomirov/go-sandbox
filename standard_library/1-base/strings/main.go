/*
Пакет strings — действия со строками (текстом).

Важно про строки в Go: их НЕЛЬЗЯ изменить «на месте». Любая операция возвращает НОВУЮ строку,
а старая остаётся прежней. Поэтому функции strings.* всегда что-то возвращают — результат нужно сохранить.

Зачем нужно: привести текст к единому виду, найти/заменить кусок, разрезать и склеить.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"strings"
)

func main() {
	// TrimSpace убирает пробелы и переводы строк по краям. ToLower делает все буквы маленькими.
	// Частый приём: привести email/логин к единому виду перед сравнением.
	raw := "  Alice@Example.COM  "
	email := strings.ToLower(strings.TrimSpace(raw))
	fmt.Println(email) // => alice@example.com

	// Contains — есть ли подстрока внутри. HasPrefix/HasSuffix — начинается/заканчивается ли.
	fmt.Println(strings.Contains("hello world", "world")) // => true
	fmt.Println(strings.HasPrefix("report.txt", "report")) // => true
	fmt.Println(strings.HasSuffix("report.txt", ".txt"))   // => true

	// Split разрезает строку по разделителю и возвращает срез (список) кусков.
	parts := strings.Split("a,b,c", ",")
	fmt.Println(parts)      // => [a b c]
	fmt.Println(parts[0])   // => a  (обращение к элементу по номеру, нумерация с 0)

	// Fields разрезает по пробелам (любым) и выкидывает пустые куски — удобно для слов.
	fmt.Println(strings.Fields("один   два  три")) // => [один два три]

	// Join — обратное действие: склеить список строк в одну через разделитель.
	fmt.Println(strings.Join([]string{"2024", "03", "15"}, "-")) // => 2024-03-15

	// ReplaceAll заменяет ВСЕ вхождения. Count считает, сколько раз встретилось.
	fmt.Println(strings.ReplaceAll("ля-ля-ля", "ля", "ла")) // => ла-ла-ла
	fmt.Println(strings.Count("ля-ля-ля", "ля"))            // => 3

	// Builder — способ собирать длинную строку по кусочкам ЭФФЕКТИВНО.
	// (Если в цикле писать s = s + "...", Go каждый раз создаёт новую строку — это медленно.)
	var b strings.Builder
	for i := 1; i <= 3; i++ {
		// & означает «дай саму переменную b», чтобы Fprintf писал прямо в неё.
		fmt.Fprintf(&b, "строка%d;", i)
	}
	fmt.Println(b.String()) // => строка1;строка2;строка3;

	// ── Ещё функции пакета ──
	fmt.Println(strings.ToUpper("hello"))               // => HELLO
	fmt.Println(strings.TrimPrefix("ID-42", "ID-"))     // => 42 (убрать начало, если есть)
	fmt.Println(strings.TrimSuffix("file.txt", ".txt")) // => file (убрать конец, если есть)
	fmt.Println(strings.EqualFold("GET", "get"))        // => true (сравнение без учёта регистра)
	fmt.Println(strings.Index("abcdef", "cd"))          // => 2 (позиция, с 0; если нет -> -1)
	// NewReader превращает строку в источник (io.Reader) — нужно там, где требуется Reader.
	var word string
	fmt.Fscan(strings.NewReader("привет"), &word)
	fmt.Println(word) // => привет
}

/*
Что важно запомнить:
  • Строку нельзя менять «на месте» — операции возвращают новую строку, сохраняй результат.
  • Split — резать по разделителю (a,b,c). Fields — резать по пробелам на слова.
  • Join — склеить список строк обратно в одну.
  • Для сборки длинного текста в цикле бери strings.Builder, а не склейку через +.

Маленькие задачи:
  1) Из строки "  Hello  " получи "hello" (убери пробелы и сделай маленькими буквами).
  2) Разрежь "1;2;3" по ";" и напечатай второй элемент.
*/
