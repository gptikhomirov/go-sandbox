package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("stop by context")
				return // <- return! иначе бесконечный цикл
			default:
				fmt.Println("default, continue wait...")
				time.Sleep(1 * time.Second)
			}
		}
	}()

	time.Sleep(3 * time.Second)
	cancel()
	wg.Wait() // <- wait after cancel!

}
