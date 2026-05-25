/*
Пакет slices (уровень 2) — только новое поверх base.
В base уже: Contains, Index, Sort, SortFunc, Max/Min, Clone, Equal, BinarySearch, Insert, Delete, IsSorted.
Здесь: *Func-варианты (по условию), Compact (дедуп), Reverse, Concat, Compare,
и мост к итераторам (All/Values/Collect/Sorted) из Go 1.23.

Как запустить:  go run main.go
*/
package main

import (
	"cmp"
	"fmt"
	"slices"
)

type User struct {
	Name string
	Age  int
}

func main() {
	// *Func-варианты: искать/сравнивать по СВОЕМУ условию (когда нужен не сам элемент, а признак).
	users := []User{{"Bob", 30}, {"Ann", 17}}
	fmt.Println(slices.ContainsFunc(users, func(u User) bool { return u.Age < 18 })) // => true (есть несовершеннолетний)
	i := slices.IndexFunc(users, func(u User) bool { return u.Name == "Ann" })
	fmt.Println(i) // => 1

	// Compact удаляет ПОДРЯД идущие дубликаты (обычно после Sort -> получаем уникальные).
	nums := []int{1, 1, 2, 3, 3, 3}
	fmt.Println(slices.Compact(nums)) // => [1 2 3]

	// Reverse разворачивает на месте. Concat склеивает несколько срезов в новый.
	r := []int{1, 2, 3}
	slices.Reverse(r)
	fmt.Println(r)                              // => [3 2 1]
	fmt.Println(slices.Concat([]int{1}, []int{2, 3})) // => [1 2 3]

	// Compare — лексикографическое сравнение двух срезов: -1 / 0 / 1.
	fmt.Println(slices.Compare([]int{1, 2}, []int{1, 3})) // => -1

	// SortStableFunc — стабильная сортировка структур по правилу (порядок равных сохраняется).
	slices.SortStableFunc(users, func(a, b User) int { return cmp.Compare(a.Age, b.Age) })
	fmt.Println(users) // => [{Ann 17} {Bob 30}]

	// Мост к итераторам (Go 1.23):
	//   Values(s) — итератор значений; Sorted(seq) — собрать отсортированный срез из итератора;
	//   Collect(seq) — собрать срез; All(s) — итератор пар (индекс, значение).
	sortedUnique := slices.Sorted(slices.Values([]int{3, 1, 2, 1}))
	fmt.Println(sortedUnique) // => [1 1 2 3]
	for idx, u := range slices.All(users) {
		_ = idx
		_ = u
	}

	// Ещё операции (для полноты):
	xs := []int{5, 1, 3, 1}
	xs = slices.DeleteFunc(xs, func(v int) bool { return v == 1 }) // удалить по условию
	fmt.Println(xs)                                                // => [5 3]
	xs = slices.Replace(xs, 0, 1, 9, 9)                            // заменить диапазон [0,1) на 9,9
	fmt.Println(xs)                                                // => [9 9 3]
	xs = slices.Grow(xs, 10)                                       // зарезервировать место (capacity)
	xs = slices.Clip(xs)                                           // убрать лишнюю capacity
	fmt.Println(slices.IsSortedFunc(xs, func(a, b int) int { return a - b })) // проверка сортировки по правилу
	pos, ok := slices.BinarySearchFunc([]int{1, 3, 5}, 3, func(a, b int) int { return a - b })
	fmt.Println(pos, ok) // => 1 true
	fmt.Println(slices.MaxFunc(users, func(a, b User) int { return a.Age - b.Age }).Name) // самый старший
	fmt.Println(slices.MinFunc(users, func(a, b User) int { return a.Age - b.Age }).Name) // самый младший
}

/*
Что важно запомнить:
  • Правило *Func: если ищешь/сравниваешь по ПОЛЮ или УСЛОВИЮ — бери ContainsFunc/IndexFunc/EqualFunc,
    а не Contains/Index (те ищут конкретное значение целиком).
  • Дедупликация: slices.Sort + slices.Compact = уникальные значения. Compact убирает только СОСЕДНИЕ дубли,
    поэтому сначала сортируют.
  • Reverse/Concat/Compare закрывают частые операции, ради которых раньше писали циклы.
  • Sorted/Collect/Values/All — мост между срезами и итераторами (iter.Seq) Go 1.23; пригодится с maps.Keys и т.п.

Задача:
  1) Получи уникальные отсортированные значения из []int{5,3,5,1,3} (Sort + Compact или Sorted+Compact).
*/
