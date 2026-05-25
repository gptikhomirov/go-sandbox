/*
Пакет fmt (уровень 2) — только новое поверх base.
В base уже: Println, Printf, Sprintf, Fprintf, Errorf, Fscan, Scan.
Здесь: Print/Sprint/Sprintln, Fprint/Fprintln, Scanf, Sscan/Sscanf.

Зачем: варианты без формата (Print/Sprint) и разбор значений ПРЯМО ИЗ СТРОКИ (Sscan/Sscanf) —
частый приём, когда данные уже в строке и не хочется тянуть strconv по одному полю.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"os"
)

func main() {
	// Print/Sprint — как Println/Sprintln, но БЕЗ переноса строки и без пробелов между не-строками.
	fmt.Print("a", "b", "\n")        // => ab
	s := fmt.Sprint("id=", 42)       // собрать строку без формата
	fmt.Println(s)                   // => id=42
	fmt.Println(fmt.Sprintln("x"))   // Sprintln добавляет пробелы и \n -> "x\n" + ещё \n от Println

	// Fprint/Fprintln — то же, но в произвольный io.Writer (здесь stderr).
	fmt.Fprintln(os.Stderr, "лог в stderr без формата")

	// Sscanf — разобрать строку по шаблону в переменные. Удобно для "ключ:значение", дат, координат.
	var name string
	var age int
	fmt.Sscanf("Аня 30", "%s %d", &name, &age)
	fmt.Println(name, age) // => Аня 30

	// Sscan — то же, но просто по пробелам, без шаблона.
	var x, y int
	fmt.Sscan("3 4", &x, &y)
	fmt.Println(x + y) // => 7

	// Scanf — чтение по шаблону с КЛАВИАТУРЫ (stdin). Раскомментируй и введи "имя возраст":
	//   fmt.Scanf("%s %d", &name, &age)
}

/*
Что важно запомнить:
  • Семейство по приёмнику: Print* -> stdout, Sprint* -> строка, Fprint* -> любой io.Writer.
    По формату: без суффикса — без формата; f — по шаблону; ln — с переносом строки.
  • Sscanf/Sscan читают из СТРОКИ, Scanf/Scan — из stdin, Fscanf/Fscan — из io.Reader. Одинаковые глаголы.
  • Для одного-двух чисел из строки Sscan короче, чем strconv; для надёжного разбора с проверкой
    ошибок по полю — лучше strconv (Sscanf «прощает» лишнее и труднее ловить ошибки).

Задача:
  1) Разбери строку "12:30" в часы и минуты через Sscanf("%d:%d", ...).
*/
