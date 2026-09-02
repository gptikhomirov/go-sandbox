package main

import "fmt"

func demoBasicDefer() {
	fmt.Println("start func")       // 2
	defer fmt.Println("defer func") // 4
	fmt.Println("end func")         // 3
}

func main() {
	fmt.Println("start main") // 1
	demoBasicDefer()
	defer fmt.Println("defer main") // 6
	fmt.Println("end main")         // 5
}
