package main

import (
	"fmt"
	"unsafe"
)

func main() {
	s := make([]bool, 0, 3)
	fmt.Println(s, len(s), cap(s), unsafe.Sizeof(s))

	s2 := make([]bool, 0, 5)
	fmt.Println(s2, len(s2), cap(s2), unsafe.Sizeof(s2))

	s3 := make([]bool, 1, 5)
	s3 = append(s3, true, true, false)
	fmt.Println(s3, len(s3), cap(s3), unsafe.Sizeof(s3))

	s4 := make([]int16, 1, 5)
	fmt.Println(unsafe.Sizeof(s4))
}
