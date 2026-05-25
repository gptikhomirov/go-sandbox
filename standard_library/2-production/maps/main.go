/*
Пакет maps (уровень 2) — только новое поверх base.
В base уже: Clone, Copy, Equal, Keys, Values, DeleteFunc.
Здесь: EqualFunc, All, Insert, Collect — мост к итераторам (Go 1.23).

Зачем: сравнивать словари по своему правилу, перебирать пары и собирать словарь из итератора —
полезно при работе с потоками данных и преобразованиях.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	prices := map[string]int{"яблоко": 50, "банан": 30}

	// All — итератор пар (ключ, значение). Можно перебрать в for range.
	total := 0
	for _, v := range maps.All(prices) {
		total += v
	}
	fmt.Println(total) // => 80

	// EqualFunc — сравнить словари по СВОЕМУ правилу для значений (не строгим ==).
	a := map[string]int{"x": 10}
	b := map[string]int{"x": 12}
	closeEnough := maps.EqualFunc(a, b, func(v1, v2 int) bool {
		return v1/10 == v2/10 // считаем равными, если совпадает десяток
	})
	fmt.Println(closeEnough) // => true

	// Collect собирает словарь ИЗ итератора пар. Insert — добавляет пары из итератора в существующий.
	src := maps.All(prices)
	copyMap := maps.Collect(src)
	fmt.Println(len(copyMap)) // => 2

	dst := map[string]int{"вишня": 70}
	maps.Insert(dst, maps.All(prices)) // влить пары из prices в dst
	keys := slices.Sorted(maps.Keys(dst))
	fmt.Println(keys) // => [банан вишня яблоко]
}

/*
Что важно запомнить:
  • maps.Keys/Values/All возвращают ИТЕРАТОРЫ (iter.Seq), а не срезы. Чтобы получить срез — slices.Collect,
    чтобы словарь — maps.Collect. Это «мост» между map и итераторами Go 1.23.
  • EqualFunc vs Equal: Equal сравнивает значения строго (==), EqualFunc — по твоей функции
    (полезно для float с допуском, регистронезависимых строк и т.п.).
  • Insert(dst, seq) — влить пары из итератора в существующий словарь (как Copy, но из итератора).

Задача:
  1) Через maps.All и Collect сделай копию словаря, удвоив все значения (собери новый словарь).
*/
