package main

import (
	"fmt"
	"sync"
)

type SaveCounter struct {
	mu    sync.Mutex
	count int
}

func (sc *SaveCounter) Increment() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.count++
}

func (sc *SaveCounter) Value() int {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	return sc.count
}

func main() {
	counter := SaveCounter{}
	fmt.Println("before:", counter.count)

	var wg sync.WaitGroup
	wg.Add(1000)
	for range 1000 {
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()

	fmt.Println("after:", counter.count)
}
