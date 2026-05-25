/*
Пакет math (уровень 2) — только новое поверх base.
В base уже: Abs, Ceil/Floor/Round/Trunc, Sqrt, Pow, Log/Log2/Log10, MaxInt/MinInt, Pi/E.
Здесь: Pow10, Mod, тригонометрия, проверки спецзначений (NaN/Inf), границы float и целых типов.

Зачем: расчёты с дробными (углы, остатки), безопасная проверка результатов на NaN/бесконечность,
и точные границы числовых типов для валидации/переполнений.

Как запустить:  go run main.go
*/
package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(math.Pow10(3))   // => 1000 (10 в степени 3)
	fmt.Println(math.Mod(10, 3)) // => 1 (остаток от деления для float64)

	// Тригонометрия (углы в радианах).
	fmt.Printf("%.2f\n", math.Sin(math.Pi/2)) // => 1.00
	fmt.Printf("%.2f\n", math.Cos(0))         // => 1.00
	fmt.Printf("%.2f\n", math.Tan(0))         // => 0.00
	fmt.Printf("%.2f\n", math.Atan2(1, 1))    // => 0.79 (угол по координатам, рад)

	// NaN («не число») и Inf (бесконечность) появляются при делении 0/0, log отрицательного и т.п.
	// Их НЕЛЬЗЯ сравнивать через == (NaN != NaN!). Проверяй только через IsNaN/IsInf.
	bad := math.Sqrt(-1)
	fmt.Println(math.IsNaN(bad))         // => true
	fmt.Println(math.IsInf(math.Inf(1), 1)) // => true
	fmt.Println(math.NaN() == math.NaN())   // => false (важная ловушка!)

	// Границы типов — для валидации диапазонов и защиты от переполнения.
	fmt.Println(int64(math.MaxInt64) > 0) // => true
	fmt.Println(math.MaxFloat64 > 1e300)  // => true
}

/*
Что важно запомнить:
  • NaN/Inf — это валидные значения float64, а не паника. Любая «странная» математика (0/0, sqrt(-1))
    даёт их молча. Проверяй результат через math.IsNaN / math.IsInf, НИКОГДА через == (NaN != NaN).
  • Mod — остаток для float64 (для int используй оператор %).
  • MaxInt64/MinInt64/MaxInt32 — точные границы; пригодятся для проверки «влезет ли» при конвертациях.
  • Углы тригонометрии — в радианах (градусы * Pi / 180).

Задача:
  1) Проверь, что 0.0/0.0 даёт NaN (через переменную, не литералом), и поймай это через math.IsNaN.
*/
