/*
Пакет net/http/httptest — тестирование HTTP-хендлеров и клиентов БЕЗ реального порта/сети.

Зачем: проверить, что хендлер возвращает правильный код и тело, не поднимая настоящий сервер.
Два инструмента:
  - NewRecorder — «фейковый» ResponseWriter: вызываешь хендлер напрямую и смотришь, что он записал.
  - NewServer — настоящий локальный сервер на случайном порту (для тестов клиента/интеграции).

Здесь показано в main для наглядности; в реальности это код внутри *_test.go (см. пакет testing).

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

// Хендлер, который мы хотим протестировать.
func greet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name обязателен", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Привет, %s", name)
}

func main() {
	// --- Способ 1: NewRecorder — вызвать хендлер напрямую, без сети. Самый частый в unit-тестах. ---
	req := httptest.NewRequest("GET", "/greet?name=Аня", nil) // фейковый запрос
	rec := httptest.NewRecorder()                              // фейковый ResponseWriter
	greet(rec, req)                                            // вызываем хендлер напрямую

	res := rec.Result()
	body, _ := io.ReadAll(res.Body)
	fmt.Println(rec.Code)        // => 200 (rec.Code — записанный статус)
	fmt.Printf("%s\n", body)     // => Привет, Аня

	// Проверим ветку ошибки (без name -> 400).
	rec2 := httptest.NewRecorder()
	greet(rec2, httptest.NewRequest("GET", "/greet", nil))
	fmt.Println(rec2.Code) // => 400

	// --- Способ 2: NewServer — настоящий локальный сервер (для тестов КЛИЕНТА/интеграции). ---
	srv := httptest.NewServer(http.HandlerFunc(greet))
	defer srv.Close() // важно закрыть

	resp, _ := http.Get(srv.URL + "/greet?name=Боб") // srv.URL — реальный адрес вида http://127.0.0.1:PORT
	got, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("через сервер: %s\n", got) // => через сервер: Привет, Боб
}

/*
Что важно запомнить:
  • httptest.NewRecorder + httptest.NewRequest — основной способ unit-тестировать хендлер: вызвал напрямую,
    проверил rec.Code, rec.Body, rec.Header(). Быстро, без портов и сети.
  • httptest.NewServer — поднимает РЕАЛЬНЫЙ сервер на случайном порту (srv.URL): нужен, когда тестируешь
    КЛИЕНТА или поведение по сети (редиректы, таймауты). Не забудь srv.Close().
  • Это код для *_test.go (функции TestXxx), здесь — в main лишь для запуска. В реальном тесте проверки
    делают через t.Errorf (см. пакет testing).
  • Пара testing + httptest — стандарт тестирования HTTP-сервисов на Go.

Задача:
  1) Перенеси этот код в greet_test.go: TestGreet с проверкой кода 200 и тела через httptest.NewRecorder.
*/
