package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	ch := make(chan string)

	go func() {
		time.Sleep(time.Second * 3)
		ch <- "goroutine done"
	}()

	fmt.Println("start waiting...")
	select {
	case v := <-ch:
		fmt.Println("read from channel:", v)
	case <-ctx.Done():
		fmt.Println("timeout")
	}
}
