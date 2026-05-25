/*
Пакет encoding/json — (де)сериализация JSON.

Зачем:   тела HTTP-запросов и ответов, конфиги, обмен с внешними сервисами.
Когда:   ядро REST API — почти каждый хендлер читает JSON-тело и пишет JSON-ответ.
Грабли:  Marshal видит ТОЛЬКО экспортируемые (с большой буквы) поля — поле name с маленькой
         не попадёт в JSON. Имена и поведение задаются тегами `json:"..."`. Без проверки ошибки
         Unmarshal невалидный ввод «молча» оставит нулевые значения.
*/
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Теги управляют JSON: имя поля, omitempty (опустить пустое), "-" (исключить полностью).
type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email,omitempty"`
	Password string   `json:"-"`       // никогда не попадёт в JSON
	Roles    []string `json:"roles"`
}

func main() {
	u := User{ID: 1, Name: "Alice", Password: "secret", Roles: []string{"admin"}}

	// Marshal — struct → JSON []byte. Email пустой → опущен; Password исключён тегом.
	data, _ := json.Marshal(u)
	fmt.Println(string(data)) // => {"id":1,"name":"Alice","roles":["admin"]}

	// MarshalIndent — человекочитаемый JSON (логи, файлы конфигурации, отладка).
	pretty, _ := json.MarshalIndent(u, "", "  ")
	fmt.Println(pretty != nil) // => true (вывод многострочный, опущен для краткости)

	// Unmarshal — JSON → struct. Неизвестные поля игнорируются, отсутствующие остаются нулевыми.
	var parsed User
	if err := json.Unmarshal([]byte(`{"id":2,"name":"Bob","roles":["user"]}`), &parsed); err != nil {
		fmt.Println("unmarshal error:", err) // на практике → 400 Bad Request
		return
	}
	fmt.Printf("%+v\n", parsed) // => {ID:2 Name:Bob Email: Password: Roles:[user]}

	// Decoder — потоковое чтение из io.Reader (тело запроса, файл).
	// DisallowUnknownFields включает строгую валидацию: лишнее поле → ошибка.
	dec := json.NewDecoder(strings.NewReader(`{"id":3,"name":"Carol","extra":true}`))
	dec.DisallowUnknownFields()
	var strict User
	fmt.Println(dec.Decode(&strict) != nil) // => true (отвергнут из-за "extra")
}

/*
Что запомнить (что чаще и почему):
  • Marshal/Unmarshal ([]byte) vs NewEncoder/NewDecoder (io.Reader/Writer):
      Marshal/Unmarshal — когда JSON уже в []byte/строке (или нужен []byte).
      Encoder/Decoder   — когда работаете с ПОТОКОМ: в HTTP это идеально, т.к. ResponseWriter и
                          r.Body — это Writer/Reader. json.NewEncoder(w).Encode(v) — стандарт для ответа,
                          json.NewDecoder(r.Body).Decode(&v) — для запроса. Меньше аллокаций.
  • Теги обязательны: без них имена полей пойдут как есть (Name, а не name). omitempty для опциональных,
    "-" — чтобы скрыть пароли/внутренние поля.
  • Всегда проверяйте ошибку Unmarshal/Decode и возвращайте 400 — иначе битый ввод даст «тихие» нули.
  • DisallowUnknownFields — для строгих API, чтобы клиент не слал лишнее (опечатки в полях видны сразу).

Типичные сценарии:
  1) Ответ хендлера:  json.NewEncoder(w).Encode(resp)
  2) Тело запроса:    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w,"",400) }
  3) Лог объекта:     b, _ := json.MarshalIndent(obj, "", "  "); log.Println(string(b))
*/
