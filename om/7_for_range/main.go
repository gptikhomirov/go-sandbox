package main

import "fmt"

func main() {
	word := "Фильм"
	for pos, ch := range word {
		fmt.Printf("Позиция %d: '%c' (код %d, тип %T)\n", pos, ch, ch, ch)
	}

	ratings := map[string]int{
		"first":  1,
		"second": 2,
		"third":  3,
	}

	for key, value := range ratings {
		fmt.Println(key, value)
	}
}
