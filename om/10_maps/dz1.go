package main

import "fmt"

func countChars(s string) map[rune]int {
	m := make(map[rune]int)

	for _, ch := range s {
		m[ch] += 1
		fmt.Println(string(ch), fmt.Sprintf("%c\n", ch))
	}

	return m
}

func main() {
	res := countChars("hello world")

	fmt.Println(res)
	fmt.Println(res['h'])
	fmt.Println(res['l'])
}
