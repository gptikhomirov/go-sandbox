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
)

// struct — своя структура (набор полей). Текст в кавычках после поля — это «тег»:
// он говорит, КАК поле будет называться в JSON. Без тега имя берётся как есть (с большой буквы).
// Важно: в JSON попадают ТОЛЬКО поля с большой буквы (экспортируемые).
type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"` // omitempty = «не выводить, если пусто»
}

func main() {
	// Создаём значение и превращаем его в JSON. Marshal возвращает байты и ошибку.
	u := User{Name: "Alice", Age: 30}
	data, err := json.Marshal(u)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	// Email пустой -> из-за omitempty его в JSON нет.
	fmt.Println(string(data)) // => {"name":"Alice","age":30}

	// MarshalIndent делает «красивый» JSON с отступами — удобно читать глазами.
	pretty, _ := json.MarshalIndent(u, "", "  ")
	fmt.Println(string(pretty))
	// => {
	//      "name": "Alice",
	//      "age": 30
	//    }

	// Обратное действие: текст JSON -> структура Go.
	// & перед переменной означает «заполни саму переменную parsed результатом».
	input := `{"name": "Bob", "age": 25, "email": "bob@mail.com"}`
	var parsed User
	if err := json.Unmarshal([]byte(input), &parsed); err != nil {
		fmt.Println("не смог разобрать JSON:", err)
		return
	}
	fmt.Println(parsed.Name, parsed.Age, parsed.Email) // => Bob 25 bob@mail.com
}

/*
Что важно запомнить:
  • JSON — текстовый формат обмена данными: {"ключ": значение}.
  • Marshal: данные Go -> JSON. Unmarshal: JSON -> данные Go (передавай &переменную — её и заполнят).
  • Теги `json:"имя"` задают, как поле зовётся в JSON. omitempty прячет пустые поля.
  • В JSON попадают только поля с БОЛЬШОЙ буквы. Поле с маленькой буквы json просто не увидит.
  • Всегда проверяй ошибку Unmarshal: если пришёл «битый» JSON, она подскажет, что не так.

Маленькие задачи:
  1) Добавь в User поле City с тегом `json:"city"`, заполни и посмотри на вывод Marshal.
  2) Разбери JSON `{"name":"Ann","age":19}` в структуру и напечатай имя и возраст.
*/
