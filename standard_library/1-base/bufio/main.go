/*
Пакет bufio — удобное и быстрое чтение текста по строкам/словам ("buf" = buffer, буфер).

Зачем нужно: когда надо прочитать ввод ПОСТРОЧНО (из файла, с клавиатуры, из любого источника текста).
bufio.Scanner делает это просто: «пока есть строки — давай мне их по одной».

Что такое «источник текста»: это может быть файл, ввод с клавиатуры (os.Stdin) или, как здесь,
строка, обёрнутая в strings.NewReader. Scanner работает с любым из них одинаково.

Как запустить:  go run main.go
*/
package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	// Готовим «источник»: три строки текста. \n — это перевод строки.
	text := "Аня 30\nБоря 25\nВера 41\n"

	// NewScanner делает сканер поверх источника. Дальше — цикл «пока есть строки».
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() { // Scan() читает следующую строку и возвращает true, пока строки есть
		line := scanner.Text() // Text() — текущая строка (без перевода строки в конце)
		fmt.Println("строка:", line)
	}
	// => строка: Аня 30 / строка: Боря 25 / строка: Вера 41

	// Можно читать не по строкам, а по СЛОВАМ. Для этого переключаем режим: Split(bufio.ScanWords).
	// Удобно, когда числа/слова записаны через пробел.
	numbers := bufio.NewScanner(strings.NewReader("10 20 30"))
	numbers.Split(bufio.ScanWords)

	sum := 0
	for numbers.Scan() {
		word := numbers.Text() // очередное слово, например "10"
		var n int
		fmt.Sscan(word, &n) // превращаем слово в число (& — «пиши результат в n»)
		sum += n
	}
	fmt.Println("сумма:", sum) // => сумма: 60

	// ── Ещё функции пакета ──
	// NewReader + ReadString читают до указанного символа (включая его). Здесь — до запятой.
	reader := bufio.NewReader(strings.NewReader("первое,второе"))
	part, _ := reader.ReadString(',')
	fmt.Printf("%q\n", part) // => "первое,"
	// Writer буферизует запись; в конце нужен Flush. Пишем в strings.Builder.
	var out strings.Builder
	bw := bufio.NewWriter(&out)
	bw.WriteString("буфер ок")
	bw.Flush() // без Flush данные не попадут в out
	fmt.Println(out.String()) // => буфер ок
}

/*
Что важно запомнить:
  • bufio.Scanner — стандартный способ читать текст ПО СТРОКАМ. Шаблон всегда один:
        for scanner.Scan() { line := scanner.Text(); ...использовать line... }
  • Split(bufio.ScanWords) переключает на чтение ПО СЛОВАМ (по пробелам) — удобно для чисел через пробел.
  • Для большого ввода Scanner намного быстрее, чем читать по одному значению через fmt.Scan.
  • Чтобы читать с клавиатуры, вместо strings.NewReader(...) подставляют os.Stdin.

Маленькие задачи:
  1) Замени strings.NewReader(text) на os.Stdin (добавь import "os"), запусти и введи пару строк руками.
  2) Посчитай, сколько ВСЕГО слов в строке "раз два три четыре" (режим ScanWords, считай в цикле).
*/
