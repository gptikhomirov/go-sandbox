/*
Пакет mime/multipart — формы с файлами (multipart/form-data) и составные тела запросов.

Когда форма содержит файл (загрузка аватара), тело запроса разбивается на ЧАСТИ, разделённые
«boundary». multipart умеет и СОБИРАТЬ такое тело (клиент), и РАЗБИРАТЬ его (сервер).

На сервере чаще используют готовое r.ParseMultipartForm / r.FormFile (см. net/http), но под капотом — этот пакет.

Как запустить:  go run main.go
*/
package main

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
)

func main() {
	// --- Сторона КЛИЕНТА: собрать multipart-тело (текстовое поле + файл) ---
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	w.WriteField("title", "мой документ") // обычное поле формы

	// CreateFormFile создаёт часть для файла; пишем в неё содержимое как в io.Writer.
	fw, _ := w.CreateFormFile("upload", "hello.txt")
	io.WriteString(fw, "содержимое файла")

	contentType := w.FormDataContentType() // "multipart/form-data; boundary=..." — это пойдёт в заголовок
	w.Close()                              // ОБЯЗАТЕЛЬНО: дописывает завершающий boundary

	fmt.Println(len(body.Bytes()) > 0)            // => true (тело собрано)
	fmt.Println(contentType[:19])                 // => multipart/form-data

	// --- Сторона СЕРВЕРА: разобрать это тело обратно по частям ---
	// boundary достаём из Content-Type через mime.ParseMediaType (в реальном запросе — из заголовка).
	_, params, _ := mime.ParseMediaType(contentType)
	r := multipart.NewReader(&body, params["boundary"])
	for {
		part, err := r.NextPart()
		if err == io.EOF {
			break
		}
		data, _ := io.ReadAll(part)
		fmt.Printf("часть %q (файл %q): %q\n", part.FormName(), part.FileName(), string(data))
	}
	// => часть "title" (файл ""): "мой документ"
	//    часть "upload" (файл "hello.txt"): "содержимое файла"
}

/*
Что важно запомнить:
  • multipart/form-data — формат тела, когда в форме есть ФАЙЛЫ. Состоит из частей, разделённых boundary.
  • Клиент: multipart.NewWriter -> WriteField (текст) + CreateFormFile (файл) -> Close() (важно!) ->
    заголовок Content-Type берут из w.FormDataContentType().
  • Сервер: обычно НЕ парсят вручную, а зовут r.ParseMultipartForm(maxMemory) + r.FormFile("upload")
    (см. пример net/http). Прямой multipart.Reader нужен для нестандартных случаев/стриминга.
  • Не забудь w.Close() у писателя — без него тело без завершающего boundary и не распарсится.

Задача:
  1) Собери форму с двумя текстовыми полями (WriteField) и распечатай их, разобрав через multipart.Reader.
*/
