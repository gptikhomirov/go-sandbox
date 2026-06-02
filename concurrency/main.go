package main

import (
	"concurrency/miner"
	"concurrency/postman"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var coal atomic.Int64

	mtx := sync.Mutex{}
	var mails []string

	minerCtx, cancelMinerCtx := context.WithCancel(context.Background())
	postmanCtx, cancelPostmanCtx := context.WithCancel(context.Background())

	go func() {
		time.Sleep(2 * time.Second)
		cancelMinerCtx()
	}()

	go func() {
		time.Sleep(4 * time.Second)
		cancelPostmanCtx()
	}()

	coalTransferPoint := miner.MinerPool(minerCtx, 20000)
	mailTransferPoint := postman.PostmanPool(postmanCtx, 20000)

	initTime := time.Now()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for c := range coalTransferPoint {
			coal.Add(int64(c))
			fmt.Println("Добыли ", c, " угля")
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()

		for mail := range mailTransferPoint {
			mtx.Lock()
			mails = append(mails, mail)
			fmt.Println("Получили письмо:", mail)
			mtx.Unlock()
		}
	}()

	wg.Wait()

	fmt.Println("Итого добыли угля:", coal.Load())

	mtx.Lock()
	fmt.Println("Итого получено писем:", len(mails))
	mtx.Unlock()

	fmt.Println("Всего ушло времени:", time.Since(initTime))
}
