/*
Пакет mime — определение MIME-типов (Content-Type) по расширению и разбор Content-Type.

MIME-тип — это «что за данные»: text/html, application/json, image/png. Браузер и сервер
по нему понимают, как обрабатывать содержимое.

Зачем: при отдаче файла поставить правильный Content-Type; при приёме — разобрать Content-Type
запроса (например, вытащить charset или boundary у multipart).

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"mime"
)

func main() {
	// TypeByExtension — какой Content-Type у файла с таким расширением (точка обязательна).
	fmt.Println(mime.TypeByExtension(".json")) // => application/json
	fmt.Println(mime.TypeByExtension(".html")) // => text/html; charset=utf-8
	fmt.Println(mime.TypeByExtension(".png"))  // => image/png

	// ParseMediaType — разобрать значение заголовка Content-Type на тип и ПАРАМЕТРЫ.
	mediaType, params, err := mime.ParseMediaType("text/html; charset=utf-8")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println(mediaType)            // => text/html
	fmt.Println(params["charset"])    // => utf-8

	// Частый случай: у multipart/form-data в параметрах лежит boundary (нужен для разбора формы).
	mt, p, _ := mime.ParseMediaType("multipart/form-data; boundary=xyz123")
	fmt.Println(mt, p["boundary"]) // => multipart/form-data xyz123
}

/*
Что важно запомнить:
  • TypeByExtension(".ext") — выбрать Content-Type при отдаче файла. Расширение с точкой!
    (net/http.ServeContent/ServeFile делают это сами — ручной mime нужен в своих обработчиках.)
  • ParseMediaType — разобрать входящий Content-Type на базовый тип + параметры (charset, boundary).
    Не парси Content-Type вручную через Split(";") — там бывают кавычки и пробелы, ParseMediaType корректен.
  • Для разбора multipart-формы boundary берут именно отсюда (см. пример mime/multipart).

Задача:
  1) Разбери "application/json; charset=utf-8" и выведи базовый тип и charset.
*/
