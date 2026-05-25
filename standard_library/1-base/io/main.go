/*
Пакет io — базовые интерфейсы потоков: Reader, Writer, Closer.

Зачем:   единый способ работать с файлами, сетью, HTTP-телами и буферами.
Когда:   везде, где есть «поток данных»: тело запроса (r.Body), файл, соединение, буфер.
Грабли:  io.ReadAll грузит ВЕСЬ поток в память — на больших/недоверенных данных это DoS-риск;
         ограничивайте через io.LimitReader или http.MaxBytesReader. io.EOF — это не ошибка,
         а нормальный сигнал «поток закончился».
*/
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// io.Reader — абстракция «откуда читать». strings.NewReader, файл, r.Body, соединение —
	// всё это io.Reader, и код ниже работает с любым из них без изменений.
	var src io.Reader = strings.NewReader("payload data 12345")

	// ReadAll — прочитать весь поток в []byte. Частый случай: тело HTTP-запроса/ответа.
	data, err := io.ReadAll(src)
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	fmt.Printf("%d bytes: %q\n", len(data), data) // => 18 bytes: "payload data 12345"

	// Copy — перелить из Reader в Writer без ручного буфера и без загрузки всего в память.
	// Так копируют файлы, проксируют ответы, сохраняют загрузки на диск.
	// os.Stdout — это io.Writer.
	fmt.Print("copied -> ")
	io.Copy(os.Stdout, strings.NewReader("streamed line\n")) // => copied -> streamed line

	// io.EOF — сигнал конца потока, возвращается при чтении до конца. Это не ошибка.
	r := strings.NewReader("ab")
	buf := make([]byte, 1)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			fmt.Printf("%q ", buf[:n]) // => "a" "b"
		}
		if err == io.EOF {
			fmt.Println("\nEOF")
			break
		}
	}
}

/*
Что запомнить (что чаще и почему):
  • io.ReadAll vs io.Copy — ключевой выбор по памяти:
      ReadAll — когда нужны ВСЕ данные как []byte и они заведомо небольшие (JSON-тело запроса).
      Copy    — когда данные надо просто ПЕРЕЛИТЬ (файл→диск, ответ→клиент): работает потоково,
                фиксированным буфером, не держит весь объём в памяти. Для больших данных — только Copy.
  • Ограничивайте размер: на публичных эндпоинтах оборачивайте тело в http.MaxBytesReader,
    иначе io.ReadAll на гигабайтном теле положит сервис.
  • Reader/Writer — это «утиная типизация» Go: пишите функции, принимающие io.Reader/io.Writer,
    и они заработают с файлом, сетью, буфером и тестовым strings.Reader без изменений.
  • io.EOF — проверяйте через err == io.EOF (или errors.Is). Не считайте его сбоем.

Типичные сценарии:
  1) Тело запроса:   body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
  2) Сохранить файл: io.Copy(dstFile, r.Body)
  3) Тестируемость:  func parse(r io.Reader) {...}  // в тесте: parse(strings.NewReader("..."))
*/
