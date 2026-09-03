package main

import (
	"fmt"
	"sync"
	"time"
)

func writer2(ch chan<- int, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("writer: end")
		wg.Done()
	}()

	ch <- 1
	fmt.Println("writer:", 1)
	ch <- 2
	fmt.Println("writer:", 2)
	ch <- 3
	fmt.Println("writer:", 3)

	fmt.Println("writer: sleep for 5s")
	time.Sleep(5 * time.Second)
	ch <- 4
}

func reader2(ch <-chan int, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("reader: end")
		wg.Done()
	}()

	fmt.Println("reader: messages", len(ch), "value =", <-ch)
	fmt.Println("reader: messages", len(ch), "value =", <-ch)
	fmt.Println("reader: messages", len(ch), "value =", <-ch)
	fmt.Println("reader: messages", len(ch), "value =", <-ch)
}

func main() {
	ch := make(chan int, 3) // Блокируется если отправок больше 3

	var wg sync.WaitGroup
	wg.Add(2)
	go writer2(ch, &wg)
	go reader2(ch, &wg)

	wg.Wait()
	fmt.Println("main: end")
}
