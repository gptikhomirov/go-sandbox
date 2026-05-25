/*
Пакет net/http (уровень 2) — только новое поверх base.
В base уже: сервер/роутер/хендлеры, http.Get, базовый Client, статусы.
Здесь: production-вещи — graceful shutdown, статика (FileServer/StripPrefix), cookies, basic-auth,
лимит тела (MaxBytesReader), запрос с context, настройка клиента (Transport), HTTPS, константы методов.

Как запустить:  go run main.go
*/
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Подготовим папку со статикой для FileServer.
	dir, _ := os.MkdirTemp("", "static")
	defer os.RemoveAll(dir)
	os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("я статика"), 0o644)

	mux := http.NewServeMux()

	// FileServer отдаёт файлы из папки. StripPrefix убирает префикс пути перед поиском файла:
	// запрос /files/hello.txt -> ищется hello.txt в dir.
	mux.Handle("GET /files/", http.StripPrefix("/files/", http.FileServer(http.Dir(dir))))

	// Хендлер с cookie и basic-auth.
	mux.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth() // разобрать заголовок Authorization: Basic ...
		if !ok || user != "admin" || pass != "secret" {
			w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(w, "нужна авторизация", http.StatusUnauthorized) // 401
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "xyz", HttpOnly: true}) // выдать cookie
		fmt.Fprintf(w, "привет, %s", user)
	})

	// Хендлер с лимитом размера тела (защита от гигантских запросов).
	mux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // максимум 1 МБ
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "слишком большое тело", http.StatusRequestEntityTooLarge)
			return
		}
		fmt.Fprintf(w, "принято %d байт", len(body))
	})

	// Загрузка файла через форму: ParseMultipartForm + FormFile (см. также пакет mime/multipart).
	mux.HandleFunc("POST /avatar", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil { // до 10 МБ в памяти
			http.Error(w, "плохая форма", http.StatusBadRequest)
			return
		}
		file, header, err := r.FormFile("avatar") // поле формы с файлом
		if err != nil {
			http.Error(w, "нет файла", http.StatusBadRequest)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "загружен %s", header.Filename)
	})

	server := &http.Server{Addr: "127.0.0.1:8098", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("сервер:", err)
		}
	}()
	time.Sleep(50 * time.Millisecond)

	// Клиент с настройками: Transport (пул соединений) и таймаут.
	client := &http.Client{
		Timeout:   2 * time.Second,
		Transport: &http.Transport{MaxIdleConns: 10},
	}
	defer client.CloseIdleConnections()

	// Запрос с context (можно отменить/ограничить по времени). Константа метода — http.MethodGet.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:8098/me", nil)
	req.SetBasicAuth("admin", "secret") // добавить basic-auth к запросу
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("ответ: %s\n", body) // => ответ: привет, admin
	if len(resp.Cookies()) > 0 {
		fmt.Println("получена cookie:", resp.Cookies()[0].Name) // => получена cookie: session
	}

	// Статика:
	resp2, _ := http.Get("http://127.0.0.1:8098/files/hello.txt")
	st, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()
	fmt.Printf("статика: %s\n", st) // => статика: я статика

	// Graceful shutdown: даём серверу аккуратно завершить активные запросы.
	shutCtx, shutCancel := context.WithTimeout(context.Background(), time.Second)
	defer shutCancel()
	server.Shutdown(shutCtx)
	fmt.Println("сервер остановлен")

	// HTTPS-сервер запускают так (нужны файлы сертификата и ключа):
	//   server.ListenAndServeTLS("cert.pem", "key.pem")
}

/*
Что важно запомнить:
  • Server.Shutdown(ctx) — корректная остановка: дожидается активных запросов вместо резкого обрыва.
    Обязателен в проде (по сигналу SIGTERM). Голый процесс-kill рвёт соединения.
  • FileServer + StripPrefix — отдача статики из папки; StripPrefix срезает URL-префикс перед поиском файла.
  • MaxBytesReader — лимит тела запроса: без него злоумышленник пришлёт гигабайты и положит сервис.
  • Cookies: http.SetCookie(w, &http.Cookie{...}) на сервере, resp.Cookies()/r.Cookie(name) для чтения.
    HttpOnly/Secure — важные флаги безопасности.
  • BasicAuth: r.BasicAuth() читает, req.SetBasicAuth() ставит.
  • Клиент в проде: свой http.Client с Timeout и Transport (пул), а не http.Get/DefaultClient без таймаута.
  • NewRequestWithContext — всегда предпочтительнее NewRequest: даёт отмену/таймаут на запрос.
  • Константы методов http.MethodGet/Post/... — вместо строк "GET" (меньше опечаток).

Задача:
  1) Добавь хендлер "GET /logout", который перезаписывает cookie session с MaxAge=-1 (удаление), и проверь клиентом.
*/
