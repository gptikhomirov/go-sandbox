package main

import (
	"fmt"
	"time"
)

func writer(ch chan<- string) {
	time.Sleep(3 * time.Second)
	ch <- "value for writer"
	fmt.Println("writer: sent value to channel")
}

func reader(ch <-chan string) {
	fmt.Println("reader: blocks until get value")
	value := <-ch
	fmt.Println("reader: got value", value)
	fmt.Println("reader: end")
}

func main() {
	ch := make(chan string)

	go writer(ch)
	fmt.Println("main: blocks")
	fmt.Println(<-ch)
	fmt.Println("main: unblock")

	go reader(ch)
	time.Sleep(4 * time.Second)
	ch <- "value for reader"
	fmt.Println("main: end")

}
