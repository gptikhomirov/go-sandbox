package main

import "fmt"

func sum(nums []int) int {
	total := 0
	// 1
	for _, v := range nums {
		total += v
	}

	// 2
	//for i := 0; i < len(nums); i++ {
	//	total += nums[i]
	//}

	return total
}

func main() {
	sl1 := []int{1, 7, 3, 4, 5}
	fmt.Println(sum(sl1), "\n")

	sl2 := make([]int, 0)
	fmt.Println(sum(sl2), "\n")

	var sl3 []int
	fmt.Println(sum(sl3), "\n")
}
