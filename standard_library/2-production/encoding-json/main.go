/*
Пакет encoding/json (уровень 2) — только новое поверх base.
В base уже: Marshal, MarshalIndent, Unmarshal, NewEncoder/Encode, NewDecoder/Decode, DisallowUnknownFields.
Здесь: Valid, Compact/Indent, Encoder.SetIndent/SetEscapeHTML, Decoder.UseNumber, Number,
потоковый разбор (Decoder.More/Token).

Зачем: проверять/переформатировать готовый JSON, точно работать с числами (Number — без потери точности),
читать ПОТОК объектов (JSONL/большие массивы) без загрузки всего в память.

Как запустить:  go run main.go
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func main() {
	raw := []byte(`{"a": 1,   "b":  2}`)

	// Valid — быстро проверить, что байты вообще являются корректным JSON (без разбора в структуру).
	fmt.Println(json.Valid(raw))             // => true
	fmt.Println(json.Valid([]byte(`{bad}`))) // => false

	// Compact убирает лишние пробелы; Indent — наоборот, делает с отступами.
	var buf bytes.Buffer
	json.Compact(&buf, raw)
	fmt.Println(buf.String()) // => {"a":1,"b":2}

	// SetEscapeHTML(false) отключает экранирование <, >, & (по умолчанию включено для безопасности в HTML).
	// Нужно, когда JSON НЕ идёт в HTML и важна читаемость (URL, формулы).
	var out bytes.Buffer
	enc := json.NewEncoder(&out)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(map[string]string{"url": "a?b=1&c=2"})
	fmt.Print(out.String()) // => { "url": "a?b=1&c=2" }  (& не превратился в &)

	// UseNumber + json.Number — хранить числа как есть, без конвертации в float64.
	// Важно для больших целых (id), где float теряет точность.
	dec := json.NewDecoder(strings.NewReader(`{"id": 123456789012345678}`))
	dec.UseNumber()
	var m map[string]json.Number
	dec.Decode(&m)
	id, _ := m["id"].Int64()
	fmt.Println(id) // => 123456789012345678 (точно, без округления)

	// Потоковый разбор массива: Token читает по одному элементу, More — «есть ли ещё».
	stream := json.NewDecoder(strings.NewReader(`[1, 2, 3]`))
	stream.Token() // прочитать открывающую '['
	sum := 0
	for stream.More() { // пока в массиве есть элементы
		var n int
		stream.Decode(&n)
		sum += n
	}
	fmt.Println(sum) // => 6
}

/*
Что важно запомнить:
  • json.Valid — дешёвая проверка «это вообще JSON?» без разбора в структуру.
  • SetEscapeHTML(false): по умолчанию json экранирует <, >, & в \uXXXX (безопасно для вставки в HTML).
    Отключай, только если точно знаешь, что вывод не попадёт в HTML.
  • UseNumber + json.Number — спасение для БОЛЬШИХ целых (id из БД): обычный Unmarshal кладёт числа в float64
    и теряет точность после ~2^53. Number хранит как строку, конвертируешь явно (.Int64()).
  • Decoder.More()/Token() — читать большой массив/поток объектов (JSONL) по одному, не держа всё в памяти.

Задача:
  1) Разбери `{"big": 9007199254740993}` обычным Unmarshal в map[string]interface{} и через UseNumber —
     сравни значения и убедись, что float теряет точность.
*/
