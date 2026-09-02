package main

import "fmt"

func deferArgs() {
	someValue := 100
	defer fmt.Println("defer: someValue = ", someValue) // cals arg someValue as 100

	someValue = 200
	fmt.Println("someValue = ", someValue) // 200
}

func deferArgsWithClosure() {
	someValue := 100
	defer func() {
		fmt.Println("with closure: defer: someValue = ", someValue) // 200 after calc
	}()

	someValue = 200
	fmt.Println("with closure: someValue = ", someValue) // 200
}

func main() {
	deferArgs()

	fmt.Println("---")

	deferArgsWithClosure()
}
