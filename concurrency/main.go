package main

import (
	"concurrency/miner"
	"concurrency/postman"
	"context"
	"fmt"
	"time"
)

func main() {
	var coal int
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

	coalTransferPoint := miner.MinerPool(minerCtx, 2)
	mailTransferPoint := postman.PostmanPool(postmanCtx, 4)

	isCoalClosed := false
	isMailClosed := false
	for !isCoalClosed || !isMailClosed {
		select {
		case c, ok := <-coalTransferPoint:
			if !ok {
				isCoalClosed = true
				continue
			}

			coal += c
			fmt.Println("Добыли ", c, " угля")

		case mail, ok := <-mailTransferPoint:
			if !ok {
				isMailClosed = true
				continue
			}

			mails = append(mails, mail)
			fmt.Println("Получили письмо:", mail)
		}
	}

	fmt.Println("Итого добыли угля:", coal)
	fmt.Println("Итого получено писем:", len(mails))

}
