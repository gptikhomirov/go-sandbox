/*
Пакет net/http — HTTP-сервер и клиент. Ядро любого REST API на Go.

Зачем:   регистрировать маршруты/хендлеры, читать запрос, писать JSON-ответ, ходить в другие сервисы.
Когда:   веб-сервисы и API, интеграции с внешними HTTP-эндпоинтами.
Грабли:  заголовки (w.Header().Set) ставьте ДО WriteHeader/Write — после уже поздно. Тело ответа
         клиента ВСЕГДА закрывайте (defer resp.Body.Close), иначе утекут соединения. http.Get без
         таймаута может висеть вечно — в проде используйте http.Client с Timeout (см. production).

Пример поднимает сервер в фоне на 127.0.0.1:8099, делает к нему запрос клиентом и завершается.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// ServeMux — стандартный роутер. С Go 1.22 понимает метод и path-параметры:
	// "GET /users/{id}" — раньше так не умел, нужны были сторонние роутеры (chi, gorilla).
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")           // path-параметр {id}
		verbose := r.URL.Query().Get("v") // query-параметр ?v=1

		// Порядок важен: сначала заголовки, потом статус, потом тело.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200

		// ResponseWriter — это io.Writer, поэтому пишем JSON прямо в него.
		json.NewEncoder(w).Encode(map[string]any{"id": id, "verbose": verbose})
	})

	// Server с таймаутами — так запускают в проде (а не голый http.ListenAndServe).
	// Запуск в горутине, чтобы main мог продолжить и сделать клиентский запрос.
	server := &http.Server{Addr: "127.0.0.1:8099", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("server error:", err)
		}
	}()
	time.Sleep(50 * time.Millisecond) // дать серверу подняться

	// Клиент: http.Get — простой GET. Тело обязательно закрываем.
	resp, err := http.Get("http://127.0.0.1:8099/users/42?v=1")
	if err != nil {
		log.Fatal("request failed:", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("status:", resp.StatusCode) // => status: 200
	fmt.Printf("body: %s", body)            // => body: {"id":"42","verbose":"1"}
}

/*
Что запомнить (что чаще и почему):
  • ServeMux (Go 1.22+) vs сторонние роутеры: теперь stdlib сам умеет "METHOD /path/{param}"
    и r.PathValue — для большинства API внешний роутер уже не нужен. Это главное изменение последних версий.
  • http.Server{...} vs http.ListenAndServe(addr, h): голый ListenAndServe не задаёт таймауты
    (Read/Write/Idle) — это уязвимость к медленным клиентам. В проде создавайте Server явно.
  • http.Get/Post (пакетные) vs http.Client: пакетные функции используют DefaultClient БЕЗ таймаута —
    годятся для скриптов/примеров. В сервисе заводите &http.Client{Timeout: ...} и NewRequestWithContext.
  • Порядок записи ответа: Header().Set → WriteHeader(code) → Write(body). После первого Write
    статус и заголовки уже отправлены.
  • defer resp.Body.Close() — обязателен на КАЖДЫЙ успешный ответ, иначе соединения не переиспользуются.

Типичные сценарии:
  1) Ручка REST:    mux.HandleFunc("POST /users", createUser)
  2) JSON-ответ:    w.Header().Set("Content-Type","application/json"); json.NewEncoder(w).Encode(v)
  3) Вызов сервиса: req,_ := http.NewRequestWithContext(ctx, "GET", url, nil); resp,_ := client.Do(req)
*/
