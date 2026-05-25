/*
Пакет sort (уровень 2) — только новое поверх base.
В base уже: Ints/Strings/Float64s, Slice/SliceStable, SearchInts/Search.
Здесь: сортировка СВОЕГО типа через интерфейс sort.Interface, Reverse, SearchFloat64s.

Зачем: иногда нужно сортировать не срез, а свой контейнер — для этого реализуют интерфейс
из трёх методов (Len/Less/Swap). А Reverse разворачивает любой порядок сортировки.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"sort"
)

// Свой тип, реализующий sort.Interface: три метода — и sort.Sort умеет его сортировать.
type ByLen []string

func (s ByLen) Len() int           { return len(s) }
func (s ByLen) Less(i, j int) bool { return len(s[i]) < len(s[j]) } // короче -> раньше
func (s ByLen) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func main() {
	words := ByLen{"банан", "ёж", "кот"}
	sort.Sort(words) // сортирует через методы Len/Less/Swap
	fmt.Println(words) // => [ёж кот банан]

	// Reverse оборачивает интерфейс и инвертирует Less -> сортировка по убыванию.
	sort.Sort(sort.Reverse(words))
	fmt.Println(words) // => [банан кот ёж]

	// Stable — стабильная версия Sort (сохраняет порядок равных).
	sort.Stable(ByLen{"аб", "вг", "д"})

	// SearchFloat64s — бинарный поиск в отсортированном []float64.
	fs := []float64{1.1, 2.2, 3.3}
	fmt.Println(sort.SearchFloat64s(fs, 2.2)) // => 1
}

/*
Что важно запомнить:
  • Три способа сортировать, от простого к гибкому:
      sort.Ints/Strings        — базовые срезы (см. base);
      sort.Slice(s, less)      — срез по своему правилу без объявления типа (см. base);
      sort.Sort(myType)        — СВОЙ контейнер: реализуешь Len/Less/Swap (sort.Interface).
    В новом коде для срезов чаще берут slices.Sort/SortFunc; sort.Interface нужен для нестандартных
    контейнеров и встречается в существующем коде.
  • sort.Reverse(data) — универсальный разворот порядка: оборачивает любой sort.Interface.

Задача:
  1) Сделай тип ByAge []Person и отсортируй людей по возрасту через sort.Sort.
*/
