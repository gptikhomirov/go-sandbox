/*
Пакет bufio — буферизованный ввод/вывод поверх io.Reader/io.Writer.

Зачем:   читать большие потоки построчно/по токенам и писать пачками, не дёргая ОС на каждый байт.
Когда:   построчное чтение файлов и stdin, быстрый ввод в алгоритмических задачах.
Грабли:  у Scanner есть лимит на длину строки (~64 КБ) — для длинных строк нужен Buffer();
         у Writer данные лежат в буфере, пока не вызовешь Flush() — забыли Flush → потеряли вывод.
*/
package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	// Scanner — самый частый инструмент: читать вход построчно.
	// Источник — strings.NewReader, но в реальности это os.Stdin или открытый файл.
	scanner := bufio.NewScanner(strings.NewReader("alice 30\nbob 25\ncarol 41\n"))
	for scanner.Scan() {
		fmt.Println("line:", scanner.Text()) // строка без \n
	}
	// => line: alice 30 / line: bob 25 / line: carol 41
	if err := scanner.Err(); err != nil { // ошибку проверять ОБЯЗАТЕЛЬНО после цикла
		fmt.Println("read error:", err)
	}

	// Split задаёт способ разбиения. ScanWords — читать по словам/токенам.
	// Идеально для задач «считать N чисел через пробел».
	ws := bufio.NewScanner(strings.NewReader("10 20 30"))
	ws.Split(bufio.ScanWords)
	sum := 0
	for ws.Scan() {
		var n int
		fmt.Sscan(ws.Text(), &n)
		sum += n
	}
	fmt.Println("sum:", sum) // => sum: 60

	// Writer — буферизует запись. Критично при большом выводе. В конце ОБЯЗАТЕЛЬНО Flush.
	var sb strings.Builder
	w := bufio.NewWriter(&sb)
	for i := 0; i < 3; i++ {
		fmt.Fprintf(w, "item-%d ", i)
	}
	w.Flush() // без этого sb останется пустым
	fmt.Println(sb.String()) // => item-0 item-1 item-2
}

/*
Что запомнить (что чаще и почему):
  • bufio.Scanner vs fmt.Scan: для чтения многих строк/чисел Scanner в разы быстрее
    (fmt.Scan медленный из-за рефлексии). В алго-задачах с большим вводом — всегда Scanner.
  • Split-режимы: по умолчанию ScanLines (по строкам) — самый частый. ScanWords — по словам,
    когда числа/токены разделены любыми пробелами. ScanRunes — посимвольно (редко).
  • Scanner vs Reader.ReadString('\n'): Scanner удобнее (сам убирает \n, простой цикл).
    ReadString берут, когда нужен разделитель, отличный от строки, или сам символ \n в данных.
  • Writer.Flush() — главный источник «почему ничего не вывелось». Ставьте defer w.Flush()
    сразу после создания писателя.

Типичные сценарии:
  1) Чтение файла:   sc := bufio.NewScanner(file); for sc.Scan() { process(sc.Text()) }
  2) Ввод чисел:     sc.Split(bufio.ScanWords); for sc.Scan() { n,_ := strconv.Atoi(sc.Text()) }
  3) Быстрый вывод:  w := bufio.NewWriter(os.Stdout); defer w.Flush(); fmt.Fprintln(w, ...)
*/
