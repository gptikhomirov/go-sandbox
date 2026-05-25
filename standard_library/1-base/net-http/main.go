/*
Пакет net/http — веб-сервер и веб-клиент (то, на чём делают сайты и API).

⚠️ ТЕМА ПОСЛОЖНЕЕ. Если только начал — сначала освой fmt, strings, slices, encoding/json.
Сюда вернись потом. Но даже сейчас полезно увидеть, как мало кода нужно для веб-сервера на Go.

Простыми словами:
  - СЕРВЕР — программа, которая ждёт запросы по сети и отвечает на них.
  - ХЕНДЛЕР (обработчик) — функция, которая решает, что ответить на конкретный запрос.
  - w (ResponseWriter) — «куда писать ОТВЕТ» клиенту.
  - r (*Request) — пришедший ЗАПРОС (адрес, параметры, данные).
  - КЛИЕНТ — тот, кто шлёт запрос серверу (браузер или другая программа).

Этот пример: запускает маленький сервер, сам отправляет ему запросы как клиент и печатает ответы.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// mux («роутер») решает, какой хендлер вызвать для какого адреса.
	mux := http.NewServeMux()

	// Хендлер 1: отвечает на "GET /hello/{name}". {name} — часть пути (достаём через PathValue).
	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")

		// Формируем ответ по шагам:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // заголовок (ДО тела!)
		w.WriteHeader(http.StatusOK)                                // код 200 (успех)
		w.Write([]byte("Привет, " + name))                          // тело ответа (байты)
	})

	// Хендлер 2: показывает FormValue (параметр после "?") и типовые ответы Error/NotFound/Redirect.
	mux.HandleFunc("GET /find", func(w http.ResponseWriter, r *http.Request) {
		q := r.FormValue("q") // значение из ?q=...
		switch q {
		case "":
			http.Error(w, "параметр q обязателен", http.StatusBadRequest) // ответ с кодом 400
		case "old":
			http.Redirect(w, r, "/hello/new", http.StatusFound) // перенаправить на другой адрес
		case "missing":
			http.NotFound(w, r) // ответ 404 «не найдено»
		default:
			fmt.Fprintf(w, "ищем: %s", q)
		}
	})

	// Эти строки просто показывают типы Handler и HandlerFunc (для справки):
	var _ http.Handler = mux // ServeMux реализует интерфейс Handler
	var _ http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {} // функция -> Handler

	// Server — настраиваемый сервер (адрес, таймауты). Запускаем в отдельной задаче (горутина),
	// чтобы программа продолжила работу и смогла ниже сама сделать запросы.
	server := &http.Server{Addr: "127.0.0.1:8099", Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("сервер упал:", err)
		}
	}()
	time.Sleep(50 * time.Millisecond) // подождём, пока сервер поднимется

	// Клиент 1: http.Get — простой GET-запрос.
	resp, err := http.Get("http://127.0.0.1:8099/hello/Аня")
	if err != nil {
		log.Fatal("запрос не удался:", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close() // тело ответа надо закрыть
	fmt.Println("код:", resp.StatusCode)   // => код: 200
	fmt.Printf("ответ: %s\n", body)         // => ответ: Привет, Аня

	// Клиент 2: настраиваемый http.Client с таймаутом + ручной запрос через Do.
	client := &http.Client{Timeout: 2 * time.Second}
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8099/find?q=книга", nil)
	resp2, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()
	fmt.Printf("ответ2: %s\n", body2) // => ответ2: ищем: книга
}

/*
Что важно запомнить:
  • Сервер слушает запросы, хендлер на них отвечает. Хендлер — это func(w http.ResponseWriter, r *http.Request).
  • Ответ формируют по порядку: Header().Set(...) -> WriteHeader(код) -> Write(данные). После Write менять поздно.
  • Из запроса читают: r.PathValue("name") — часть пути; r.FormValue("q") — параметр после "?".
  • Типовые ответы: http.Error (код+текст), http.NotFound (404), http.Redirect (перенаправление).
  • Коды статусов: StatusOK 200, StatusCreated 201, StatusBadRequest 400, StatusNotFound 404,
    StatusInternalServerError 500 (и другие константы Status...).
  • Клиент: http.Get — простой запрос; http.Client{} + Client.Do — когда нужны таймаут и настройки.
    Тело ответа всегда закрывай: resp.Body.Close().

Маленькая задача (когда будешь готов):
  1) Добавь хендлер "GET /sum/{a}/{b}": достань a и b через PathValue, преврати в числа (strconv.Atoi),
     и ответь их суммой.
*/
