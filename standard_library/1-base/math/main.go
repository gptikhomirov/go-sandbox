/*
Пакет math — математические функции и константы (всё работает с float64).

Зачем:   геометрия и расстояния, округления, логарифмы, границы для алгоритмов.
Когда:   расчёты цен/метрик, оценка сложности (Log2), «бесконечность» при поиске min/max.
Грабли:  math.Max/Min/Abs принимают и возвращают float64, НЕ int. Для целых в Go 1.21+
         есть встроенные min()/max() — их и используют для int, а math.* для дробных.
*/
package main

import (
	"fmt"
	"math"
)

func main() {
	// Abs/Max/Min — модуль и экстремумы float64.
	fmt.Println(math.Abs(-3.5))      // => 3.5
	fmt.Println(math.Max(2, 9))      // => 9
	fmt.Println(min(2, 9))           // => 2  (встроенный min для int — без math)

	// Округления: Ceil вверх, Floor вниз, Round до ближайшего, Trunc отбрасывает дробь.
	fmt.Println(math.Ceil(2.1), math.Floor(2.9), math.Round(2.5), math.Trunc(2.9))
	// => 3 2 3 2

	// Sqrt/Pow — корень и степень. Расстояние между точками (3,0)-(0,4).
	dist := math.Sqrt(math.Pow(3, 2) + math.Pow(4, 2))
	fmt.Println(dist) // => 5

	// Log2 — оценка глубины дерева / сложности O(log n).
	fmt.Println(math.Log2(1024)) // => 10

	// MaxInt/MinInt — типичная инициализация при поиске минимума в алгоритмах.
	best := math.MaxInt
	for _, v := range []int{7, 3, 9, 1} {
		best = min(best, v) // встроенный min
	}
	fmt.Println(best) // => 1
}

/*
Что запомнить (что чаще и почему):
  • Для int — встроенные min()/max() (Go 1.21+), а НЕ math.Max/Min:
      math.Max(2, 9) требует float64 и вернёт 9.0; для целых это лишние конвертации.
      min(a, b)/max(a, b) работают с любыми ordered-типами и читаются проще.
  • math.* нужен именно для float64-математики: Sqrt, Pow, Ceil/Floor/Round, тригонометрия, логарифмы.
  • math.MaxInt/MinInt — удобные границы-«бесконечности» для алгоритмов поиска.
  • Округление денег: math.Round(x*100)/100 даёт 2 знака, но для финансов точнее хранить
    суммы в копейках (int), а не во float (float64 неточен: 0.1+0.2 != 0.3).

Типичные сценарии:
  1) Расстояние:     d := math.Hypot(dx, dy)   // = Sqrt(dx²+dy²), без ручного Pow
  2) Поиск минимума: best := math.MaxInt; best = min(best, v)
  3) Округление %:   pct := math.Round(part / total * 100)
*/
