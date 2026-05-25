/*
Пакет maps — удобные действия над словарями (map).

map (словарь) — это набор пар «ключ -> значение», например map[string]int{"возраст": 30}.
По ключу быстро достаёшь значение: m["возраст"] вернёт 30.

Две важные особенности словарей в Go:
  1) Порядок перебора СЛУЧАЙНЫЙ — нельзя рассчитывать, что ключи пойдут по порядку.
  2) Как и срезы, словарь «делится»: b := a — это один и тот же словарь. Нужна копия — maps.Clone.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"slices"

	"maps"
)

func main() {
	// Создаём словарь настроек: ключ — строка, значение — число.
	settings := map[string]int{"timeout": 30, "retries": 3}
	fmt.Println(settings["timeout"]) // => 30  (достаём значение по ключу)

	// Clone — НЕЗАВИСИМАЯ копия словаря. Меняем копию — оригинал не трогается.
	cloned := maps.Clone(settings)
	cloned["timeout"] = 60
	fmt.Println(settings["timeout"], cloned["timeout"]) // => 30 60

	// Copy переносит пары из второго словаря в первый (с заменой совпадающих ключей).
	// Частый приём: взять настройки по умолчанию и «наложить» пользовательские.
	maps.Copy(cloned, map[string]int{"retries": 5, "workers": 8})
	fmt.Println(cloned["retries"], cloned["workers"]) // => 5 8

	// Из-за случайного порядка перебирать словарь «как есть» нельзя, если важен порядок.
	// Решение: достать все ключи, отсортировать их, и идти по отсортированным ключам.
	// maps.Keys даёт ключи, slices.Collect собирает их в список, slices.Sort сортирует.
	keys := slices.Collect(maps.Keys(cloned))
	slices.Sort(keys)
	fmt.Println(keys) // => [retries timeout workers]

	for _, k := range keys { // range по списку: k — очередной ключ
		fmt.Printf("%s=%d\n", k, cloned[k])
	}
	// => retries=5 / timeout=60 / workers=8

	// ── Ещё функции пакета ──
	fmt.Println(maps.Equal(map[string]int{"a": 1}, map[string]int{"a": 1})) // => true
	// Values перебирает ЗНАЧЕНИЯ словаря (порядок случайный). Сложим их.
	sum := 0
	for v := range maps.Values(cloned) {
		sum += v
	}
	fmt.Println(sum > 0) // => true
	// DeleteFunc удаляет пары по условию. Удалим все, где значение меньше 10.
	maps.DeleteFunc(cloned, func(k string, v int) bool { return v < 10 })
	fmt.Println(len(cloned) > 0) // => true
}

/*
Что важно запомнить:
  • map (словарь) — пары «ключ -> значение». Значение берут по ключу: m[key].
  • Порядок перебора словаря СЛУЧАЙНЫЙ. Нужен порядок — собери ключи и отсортируй (как выше).
  • maps.Clone — независимая копия (b := m копией НЕ является).
  • maps.Copy(dst, src) — наложить один словарь на другой (совпадающие ключи перезапишутся).

Маленькие задачи:
  1) Сделай словарь оценок map[string]int{"Аня": 5, "Боря": 4} и напечатай оценку Ани.
  2) Добавь в словарь нового ученика и выведи всех в порядке имён (через сортировку ключей).
*/
