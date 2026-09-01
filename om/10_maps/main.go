package main

import "fmt"

func main() {
	m := map[int]int{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
	}

	for i := 0; i < len(m); i++ {
		fmt.Println(i, m[i+1])
	}
}
