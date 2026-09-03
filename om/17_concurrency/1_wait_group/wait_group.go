package main

import (
	"fmt"
	rand2 "math/rand/v2"
	"sync"
	"time"
)

func main() {
	// Используется, чтобы дождаться пока все параллельно запущенные горутины завершатся
	// Обычно такая инициализация
	var wg sync.WaitGroup

	n := 10
	wg.Add(n)
	for i := range n {
		go func() {
			defer wg.Done()

			delay := time.Duration(rand2.IntN(3)+1) * time.Second
			time.Sleep(delay)
			fmt.Println(fmt.Sprintf("goroutine %d with delay %ds ended", i, delay/time.Second))
		}()
	}
	wg.Wait()

	fmt.Println("end main")
}
