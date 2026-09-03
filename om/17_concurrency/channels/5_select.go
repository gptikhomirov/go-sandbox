package main

import (
	"fmt"
	"time"
)

// 1. select блокируется пока 1 из каналов не будет готов
// 2. несколько каналов готовы => выбирается рандомный
// 3. ни один канал не готов => БЕЗ default блокирует горутину | иначе сразу выполняет default

func main() {
	movies := make(chan string)

	go func() {
		time.Sleep(3 * time.Second)
		movies <- "Gladiator"
		fmt.Println("")
	}()

	// ограничить по времени
	select {
	case movie := <-movies:
		fmt.Println("Got from channel", movie)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout")
	}

	// неблокирующий селект
	ch := make(chan string)

	select {
	case value := <-ch:
		fmt.Println("Value from channel", value)
	default: // without default deadlock
		fmt.Println("No data, no block")
	}
}
