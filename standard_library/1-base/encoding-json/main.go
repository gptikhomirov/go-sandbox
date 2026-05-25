/*
Пакет encoding/json — превращение данных Go в формат JSON и обратно.

JSON — это текстовый формат для обмена данными, его понимают почти все программы и сайты.
Выглядит так:  {"name": "Alice", "age": 30}

Два главных действия:
  - Marshal:   данные Go -> текст JSON   (например, чтобы отправить или сохранить);
  - Unmarshal: текст JSON -> данные Go   (например, чтобы прочитать пришедшие данные).

Как запустить:  go run main.go
*/
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// struct — своя структура (набор полей). Текст в кавычках после поля — это «тег»:
// он говорит, КАК поле будет называться в JSON. Важно: в JSON попадают ТОЛЬКО поля с большой буквы.
type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"` // omitempty = «не выводить, если пусто»
}

func main() {
	// Marshal: значение Go -> JSON (в виде байтов).
	u := User{Name: "Alice", Age: 30}
	data, err := json.Marshal(u)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println(string(data)) // => {"name":"Alice","age":30}  (Email пустой -> опущен)

	// MarshalIndent: «красивый» JSON с отступами — удобно читать глазами.
	pretty, _ := json.MarshalIndent(u, "", "  ")
	fmt.Println(string(pretty))
	// => {
	//      "name": "Alice",
	//      "age": 30
	//    }

	// Unmarshal: текст JSON -> структура Go. & означает «заполни саму переменную parsed».
	input := `{"name": "Bob", "age": 25, "email": "bob@mail.com"}`
	var parsed User
	if err := json.Unmarshal([]byte(input), &parsed); err != nil {
		fmt.Println("не смог разобрать JSON:", err)
		return
	}
	fmt.Println(parsed.Name, parsed.Age, parsed.Email) // => Bob 25 bob@mail.com

	// ── Ещё функции пакета ──
	// Encoder пишет JSON сразу в «приёмник» (io.Writer). Удобно в вебе: писать ответ прямо клиенту.
	// Здесь приёмник — os.Stdout (экран).
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(User{Name: "Carol", Age: 41}) // => {"name":"Carol","age":41}

	// Decoder читает JSON из «источника» (io.Reader). В вебе источник — тело запроса.
	// Здесь источник — строка через strings.NewReader.
	dec := json.NewDecoder(strings.NewReader(`{"name":"Dan","age":50,"extra":true}`))
	// DisallowUnknownFields включает строгую проверку: лишнее поле "extra" -> ошибка.
	dec.DisallowUnknownFields()
	var strict User
	err = dec.Decode(&strict)
	fmt.Println(err != nil) // => true (отклонено из-за лишнего поля "extra")
}

/*
Что важно запомнить:
  • JSON — текстовый формат обмена данными: {"ключ": значение}.
  • Marshal: данные Go -> JSON. Unmarshal: JSON -> данные Go (передавай &переменную — её заполнят).
  • Теги `json:"имя"` задают имя поля в JSON. omitempty прячет пустые поля.
  • В JSON попадают только поля с БОЛЬШОЙ буквы.
  • Marshal/Unmarshal работают с готовым текстом; Encoder/Decoder — сразу с потоком (приёмник/источник),
    это удобно для веба: Encoder пишет ответ, Decoder читает тело запроса.
  • Всегда проверяй ошибку: «битый» JSON подскажет, что не так.

Маленькие задачи:
  1) Добавь в User поле City с тегом `json:"city"`, заполни и посмотри вывод Marshal.
  2) Разбери JSON `{"name":"Ann","age":19}` через Unmarshal и напечатай имя и возраст.
*/
