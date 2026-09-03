package main

import "fmt"

func sender(ch chan<- int) {
	for _, value := range []int{1, 2, 3, 4, 5} {
		ch <- value
	}
	close(ch) // without close can be deadlock
}

func main() {
	ch := make(chan int)

	go sender(ch)

	for value := range ch {
		fmt.Println(value)
	}

	fmt.Println("+1 read, zero value", <-ch)
	value, closed := <-ch
	fmt.Println("+1 read, zero value", value, "", closed)

	close(ch) // getter DO NOT close channel => panic
}
