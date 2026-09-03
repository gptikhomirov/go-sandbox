package main

import (
	"fmt"
	"time"
)

func worker(ch chan<- struct{}) {
	fmt.Println("worker: doing...")
	time.Sleep(3 * time.Second)

	ch <- struct{}{}
	fmt.Println("worker: end")

}

func main() {
	done := make(chan struct{}) // as signal

	go worker(done)

	fmt.Println("main: blocks on read")
	<-done
	fmt.Println("main: end")
}
