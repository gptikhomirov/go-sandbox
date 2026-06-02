package miner

import (
	"context"
	"fmt"
	rand2 "math/rand/v2"
	"sync"
	"time"
)

func Miner(
	ctx context.Context,
	wg *sync.WaitGroup,
	transferPoint chan<- int,
	n int,
	power int,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Шахтер №", n, " закончил рабочий день")
			return
		default:
			fmt.Println("Шахтер №", n, " начал работу")
			time.Sleep(1 * time.Second)
			fmt.Println("Шахтер №", n, " добыл ", power, " угля")
			transferPoint <- power
			fmt.Println("Шахтер №", n, " отправил ", power, " угля")
		}
	}
}

func MinerPool(ctx context.Context, minerCount int) <-chan int {
	coalTransferPoint := make(chan int)

	wg := &sync.WaitGroup{}
	for i := 1; i <= minerCount; i++ {
		wg.Add(1)
		go Miner(ctx, wg, coalTransferPoint, i, 10*rand2.IntN(6))
	}

	go func() {
		wg.Wait()
		close(coalTransferPoint)
	}()

	return coalTransferPoint
}
