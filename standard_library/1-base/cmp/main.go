/*
Пакет cmp — сравнение упорядоченных значений (Go 1.21+).

Зачем:   безопасные компараторы для slices.SortFunc и компактная многоуровневая сортировка.
Когда:   сортировка структур по нескольким полям, выбор первого непустого значения (дефолты).
Грабли:  не путайте с пакетом github.com/google/go-cmp (это другой, для тестов). Здесь — stdlib cmp.
         cmp.Compare безопаснее, чем a-b: разность int может переполниться на крайних значениях.
*/
package main

import (
	"cmp"
	"fmt"
	"slices"
)

type Employee struct {
	Dept string
	Name string
	Age  int
}

func main() {
	// Compare возвращает -1/0/1. Безопасно для любых int (в отличие от a-b).
	fmt.Println(cmp.Compare(3, 7), cmp.Compare(7, 7), cmp.Compare(9, 7)) // => -1 0 1

	// Or — первое НЕнулевое значение. Удобно для дефолтов.
	fmt.Println(cmp.Or("", "", "8080")) // => 8080

	// Главный приём: МНОГОУРОВНЕВАЯ сортировка. cmp.Or склеивает компараторы по приоритету:
	// сначала по отделу, при равенстве — по возрасту, затем по имени.
	staff := []Employee{
		{"eng", "Bob", 30},
		{"eng", "Ann", 30},
		{"sales", "Carl", 25},
		{"eng", "Ann", 22},
	}
	slices.SortFunc(staff, func(a, b Employee) int {
		return cmp.Or(
			cmp.Compare(a.Dept, b.Dept),
			cmp.Compare(a.Age, b.Age),
			cmp.Compare(a.Name, b.Name),
		)
	})
	for _, e := range staff {
		fmt.Printf("%s %s %d\n", e.Dept, e.Name, e.Age)
	}
	// => eng Ann 22 / eng Ann 30 / eng Bob 30 / sales Carl 25
}

/*
Что запомнить (что чаще и почему):
  • cmp.Compare vs a-b: для сортировки структур пишите cmp.Compare(a.X, b.X), а не a.X-b.X.
    Разность может переполнить int (например, MinInt - MaxInt) и сломать порядок. Compare всегда корректен.
  • cmp.Or для сортировки — ключевая идиома: вместо вложенных if-ов «если равны, сравни дальше»
    просто перечисляете компараторы по приоритету. Это то, ради чего cmp чаще всего и берут.
  • cmp.Or для значений: cmp.Or(req.Limit, cfg.Limit, 50) — первый ненулевой лимит. Замена цепочке if.
  • Пара cmp + slices.SortFunc — стандартный современный способ сортировать объекты.

Типичные сценарии:
  1) Мультисортировка: cmp.Or(cmp.Compare(a.LastName,b.LastName), cmp.Compare(a.FirstName,b.FirstName))
  2) Дефолт значения:  port := cmp.Or(flagPort, envPort, "8080")
  3) Простой порядок:  slices.SortFunc(xs, func(a,b T) int { return cmp.Compare(a.Key, b.Key) })
*/
