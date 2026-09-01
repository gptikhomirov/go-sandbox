package main

import "fmt"

func hasDuplicates(nums []int) bool {
	m := make(map[int]struct{}, len(nums))

	for _, v := range nums {
		if _, ok := m[v]; ok {
			return true
		}
		m[v] = struct{}{}
	}

	return false
}

func main() {
	fmt.Println(hasDuplicates([]int{1, 233, 3, 4, 233, 22, 21}))
}
