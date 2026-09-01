package main

import "fmt"

func removeDuplicatesWithSlices(s []string) []string {
	result := make([]string, 0, len(s))

	for _, v := range s {
		isAppended := false

		for _, appended := range result {
			if v == appended {
				isAppended = true
				break
			}
		}

		if !isAppended {
			result = append(result, v)
		}
	}

	return result
}

func removeDuplicatesWithMap(s []string) []string {
	seen := make(map[string]bool, len(s))
	result := make([]string, 0, len(s))

	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	return result
}

func main() {
	sl1 := []string{"abc", "a", "b", "a", "ab", "abs", "abc"}
	//sl2 := []int{2, 3, 2, 56, 34, 3, 1}

	fmt.Println(removeDuplicatesWithMap(sl1))
}
