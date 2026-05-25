/*
Пакет slices — generic-операции над срезами (Go 1.21+). Современная замена части sort.

Зачем:   поиск/проверка наличия, сортировка, безопасное копирование, экстремумы — без циклов.
Когда:   почти везде, где раньше писали ручной for: Contains, Index, Sort, Max/Min.
Грабли:  срезы делят один массив — присваивание b := a НЕ копирует данные; чтобы менять копию,
         нужен slices.Clone. BinarySearch требует предварительно отсортированный срез.
*/
package main

import (
	"fmt"
	"slices"
)

type Product struct {
	Name  string
	Price int
}

func main() {
	nums := []int{5, 2, 9, 2, 1}

	// Contains/Index — «есть ли» и «где». Без ручного цикла.
	fmt.Println(slices.Contains(nums, 9)) // => true
	fmt.Println(slices.Index(nums, 9))    // => 2

	// Sort — сортировка ordered-типов на месте (короче, чем sort.Ints).
	slices.Sort(nums)
	fmt.Println(nums)                         // => [1 2 2 5 9]
	fmt.Println(slices.Max(nums), slices.Min(nums)) // => 9 1

	// BinarySearch — поиск в отсортированном срезе: индекс + флаг «найдено».
	idx, found := slices.BinarySearch(nums, 9)
	fmt.Println(idx, found) // => 4 true

	// Clone — независимая копия. Критично перед изменением, иначе испортите исходные данные.
	backup := slices.Clone(nums)
	backup[0] = -1
	fmt.Println(nums, backup) // => [1 2 2 5 9] [-1 2 2 5 9]

	// SortFunc — сортировка структур по полю. Компаратор возвращает -1/0/1.
	products := []Product{{"B", 30}, {"A", 10}, {"C", 20}}
	slices.SortFunc(products, func(a, b Product) int {
		return a.Price - b.Price // <0: a раньше b
	})
	fmt.Println(products) // => [{A 10} {C 20} {B 30}]

	// Equal — поэлементное сравнение (оператор == для срезов запрещён).
	fmt.Println(slices.Equal([]int{1, 2}, []int{1, 2})) // => true
}

/*
Что запомнить (что чаще и почему):
  • Sort vs SortFunc — главное различие, которое надо понимать:
      slices.Sort(s)            — для срезов ORDERED-типов (int, string, float). Простейший случай, чаще всего.
      slices.SortFunc(s, cmp)   — для СТРУКТУР или своего порядка; cmp возвращает -1/0/1.
    Правило: сортируете числа/строки → Sort; сортируете объекты по полю → SortFunc.
    Для устойчивой сортировки структур есть SortStableFunc (когда важен порядок равных).
  • Contains vs ContainsFunc: Contains(s, v) — точное значение; ContainsFunc — по условию
    (например, найти первого пользователя с Age>18). Для структур обычно нужен *Func-вариант.
  • Clone — обязателен перед мутацией чужого среза. Самый частый баг новичков: изменили
    переданный срез и неожиданно поменяли данные у вызывающего.
  • slices vs sort: для новых проектов предпочитают slices (типобезопасно, меньше кода).

Типичные сценарии:
  1) Проверка прав:   if slices.Contains(user.Roles, "admin") { ... }
  2) Сорт по полю:    slices.SortFunc(items, func(a, b Item) int { return cmp.Compare(a.Name, b.Name) })
  3) Защита данных:   safe := slices.Clone(input); mutate(safe)
*/
