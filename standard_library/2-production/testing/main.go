/*
Пакет testing — тесты, бенчмарки и фаззинг (новый пакет уровня 2).

В Go тесты — часть языка и инструментов. Файлы с тестами называются *_test.go, функции —
TestXxx(t *testing.T). Запускаются командой `go test`, а не `go run`.

Здесь main.go — это КОД, который мы тестируем (функция Max). Сами тесты — в main_test.go рядом.

Как запустить:
  go test -v            — прогнать тесты
  go test -bench .      — прогнать бенчмарки
  go test -run TestMax  — только конкретный тест
*/
package main

import "fmt"

// Max — функция, которую будем тестировать в main_test.go.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	fmt.Println("Это пакет с тестами. Запусти:  go test -v")
	fmt.Println("Max(2, 5) =", Max(2, 5)) // => Max(2, 5) = 5
}
