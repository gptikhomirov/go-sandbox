/*
Пакет sort — сортировка и бинарный поиск (классический API, до дженериков).

Зачем:   упорядочить данные перед выдачей, найти элемент в отсортированном срезе.
Когда:   сортировка структур по полю (sort.Slice) — всё ещё частый и читаемый приём.
Грабли:  бинарный поиск работает ТОЛЬКО на отсортированных данных. less-функция должна быть
         строгой (a<b), а не a<=b, иначе порядок будет нестабильным/неверным.

NB: для простых срезов (int/string) сейчас короче пакет slices — см. соседний пример slices.
*/
package main

import (
	"fmt"
	"sort"
)

type User struct {
	Name string
	Age  int
}

func main() {
	// Ints/Strings/Float64s — сортировка базовых срезов на месте (по возрастанию).
	nums := []int{5, 2, 9, 1}
	sort.Ints(nums)
	fmt.Println(nums) // => [1 2 5 9]

	// sort.Slice — сортировка по произвольному правилу. Самый частый случай в вебе:
	// упорядочить список сущностей по полю перед отдачей в API.
	users := []User{{"Bob", 30}, {"Alice", 25}, {"Carol", 30}}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Age < users[j].Age // по возрасту, по возрастанию
	})
	fmt.Println(users) // => [{Alice 25} {Bob 30} {Carol 30}]

	// SliceStable — сохраняет относительный порядок равных. Нужен для многоуровневой сортировки.
	sort.SliceStable(users, func(i, j int) bool {
		return users[i].Name < users[j].Name
	})
	fmt.Println(users) // => [{Alice 25} {Bob 30} {Carol 30}]

	// SearchInts — бинарный поиск в ОТСОРТИРОВАННОМ срезе. Возвращает индекс вставки.
	i := sort.SearchInts(nums, 9)
	fmt.Println(i, i < len(nums) && nums[i] == 9) // => 3 true
}

/*
Что запомнить (что чаще и почему):
  • sort.Slice vs slices.SortFunc (новый пакет): делают одно и то же, но компаратор разный:
      sort.Slice(s, func(i, j int) bool { ... })   — возвращает bool (a<b), работает с индексами.
      slices.SortFunc(s, func(a, b T) int { ... })  — возвращает -1/0/1, работает со значениями.
    Новый код чаще пишут на slices.SortFunc (типобезопаснее, удобно с cmp.Or для мультисортировки).
    sort.Slice остаётся валидным и встречается повсеместно в существующих проектах.
  • sort.Ints/Strings vs slices.Sort: для базовых срезов slices.Sort короче — обычно берут его.
  • Slice vs SliceStable: Stable нужен, только когда важен порядок равных элементов
    (двухуровневая сортировка). Он чуть медленнее, поэтому по умолчанию — обычный Slice.
  • Бинарный поиск — только после сортировки; иначе результат бессмысленный.

Типичные сценарии:
  1) Выдача API:        sort.Slice(items, func(i,j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
  2) Базовый срез:      slices.Sort(ids)          // вместо sort.Ints
  3) Поиск в словаре:   i := sort.SearchStrings(sortedWords, w)
*/
