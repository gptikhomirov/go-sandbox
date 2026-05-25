/*
Пакет bufio (уровень 2) — только новое поверх base.
В base уже: NewScanner, Scan/Text/Split, NewReader/ReadString, NewWriter/WriteString/Flush.
Здесь: точное чтение (ReadBytes/ReadByte/ReadRune/Peek), размеры буферов, Scanner.Buffer,
готовые split-функции (ScanLines/ScanWords/ScanRunes), Scanner.Bytes/Err.

Зачем: разбирать поток точнее (по байту/руне, заглядывать вперёд), и снимать лимит длины строки
у Scanner (по умолчанию ~64 КБ — на длинных строках Scan молча упадёт без Buffer).

Как запустить:  go run main.go
*/
package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	r := bufio.NewReader(strings.NewReader("привет\nмир"))

	// Peek заглядывает вперёд на n байт, НЕ сдвигая позицию чтения.
	head, _ := r.Peek(4)
	fmt.Printf("%q\n", head) // => "пр" (4 байта = 2 кириллические буквы)

	// ReadBytes читает до разделителя включительно и возвращает []byte.
	line, _ := r.ReadBytes('\n')
	fmt.Printf("%q\n", string(line)) // => "привет\n"

	// ReadRune читает одну руну (символ), корректно для UTF-8.
	ru, _, _ := r.ReadRune()
	fmt.Printf("%c\n", ru) // => м

	// Scanner.Buffer увеличивает лимит длины токена — нужно для очень длинных строк.
	sc := bufio.NewScanner(strings.NewReader("очень длинная строка"))
	sc.Buffer(make([]byte, 1024), 1024*1024) // буфер 1 КБ, максимум 1 МБ
	sc.Split(bufio.ScanWords)                // готовая split-функция: по словам
	count := 0
	for sc.Scan() {
		count++
	}
	fmt.Println("слов:", count, "ошибка:", sc.Err()) // => слов: 3 ошибка: <nil>

	// NewReaderSize — reader с буфером заданного размера (когда дефолтных 4 КБ мало/много).
	br := bufio.NewReaderSize(strings.NewReader("data"), 64)
	_ = br

	// NewWriterSize — writer с буфером заданного размера; WriteByte/WriteRune — точная запись.
	var out strings.Builder
	w := bufio.NewWriterSize(&out, 64)
	w.WriteByte('>')
	w.WriteRune('я')
	w.Flush()
	fmt.Println(out.String()) // => >я
}

/*
Что важно запомнить:
  • Scanner удобен, но у него ЛИМИТ на длину токена (~64 КБ). Если строки длиннее — Scan вернёт false,
    а sc.Err() покажет "token too long". Лечится через sc.Buffer(buf, max). Частая прод-ловушка на больших логах/JSONL.
  • Готовые split-функции: ScanLines (по умолчанию), ScanWords (по словам), ScanRunes (по символам).
  • Reader.Peek — заглянуть вперёд без сдвига (например, определить формат по первым байтам).
  • ReadBytes/ReadString('\n') возвращают разделитель ВНУТРИ результата — помни про лишний \n.

Задача:
  1) Прочитай поток построчно, но для строк длиннее 64 КБ — подними лимит через Scanner.Buffer.
*/
