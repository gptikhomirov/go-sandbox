/*
Пакет bytes — то же, что strings, но для []byte (новый пакет уровня 2).

Зачем отдельный пакет: данные из сети/файлов приходят как []byte. Конвертировать в string туда-сюда —
лишние копии и аллокации. bytes работает с байтами напрямую. API почти 1-в-1 как у strings.

Главный инструмент — bytes.Buffer: удобный «растущий» буфер, который умеет и читать, и писать
(реализует io.Reader и io.Writer).

Как запустить:  go run main.go
*/
package main

import (
	"bytes"
	"fmt"
)

func main() {
	data := []byte("привет,мир,go")

	// Те же операции, что в strings, но для []byte (без конвертации в string).
	fmt.Println(bytes.Contains(data, []byte("мир"))) // => true
	fmt.Println(bytes.HasPrefix(data, []byte("при"))) // => true
	parts := bytes.Split(data, []byte(","))
	fmt.Println(len(parts)) // => 3
	fmt.Println(bytes.Equal([]byte("ab"), []byte("ab"))) // => true

	// Cut — как strings.Cut: разрезать по первому разделителю.
	before, after, _ := bytes.Cut(data, []byte(","))
	fmt.Printf("%s | %s\n", before, after) // => привет | мир,go

	// bytes.Buffer — растущий буфер. Пишем в него как в io.Writer, читаем как из io.Reader.
	// Частый приём: собрать тело HTTP-запроса или ответ в память.
	var buf bytes.Buffer
	buf.WriteString("часть1;")
	buf.Write([]byte("часть2"))
	fmt.Println(buf.String()) // => часть1;часть2
	fmt.Println(buf.Len())    // => число байт в буфере

	// Buffer реализует io.Reader — из него можно читать (например, отдать в json.NewDecoder).
	line, _ := buf.ReadString(';')
	fmt.Printf("%q\n", line) // => "часть1;"

	buf.Reset() // очистить и переиспользовать
	fmt.Println(buf.Len()) // => 0

	// NewBuffer создаёт буфер с НАЧАЛЬНЫМ содержимым (в отличие от пустого var buf bytes.Buffer).
	nb := bytes.NewBuffer([]byte("старт"))
	nb.WriteString("-ещё")
	fmt.Println(nb.String()) // => старт-ещё
}

/*
Что важно запомнить:
  • bytes vs strings — одинаковый API, но bytes для []byte. Выбор простой: данные уже []byte (сеть, файл) -> bytes,
    чтобы не делать лишних string([]byte) конвертаций (это копирование памяти).
  • bytes.Buffer — швейцарский нож: и Writer (собрать данные), и Reader (отдать их). Идеален как тело запроса
    (bytes.NewReader(data)) или для накопления вывода. Аналог strings.Builder, но двусторонний.
  • Reset() переиспользует буфер без новой аллокации — полезно в циклах/пулах.
  • bytes.NewReader(b) — обернуть []byte в io.Reader (для http-запроса, json-декодера и т.п.).

Задача:
  1) Собери JSON-тело в bytes.Buffer (через json.NewEncoder(&buf).Encode(...)) и передай buf как io.Reader в запрос.
*/
