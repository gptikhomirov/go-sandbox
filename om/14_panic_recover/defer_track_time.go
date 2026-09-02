package main

import (
	"fmt"
	"time"
)

func trackTime(start time.Time, endText string) {
	end := time.Since(start).Round(time.Millisecond)
	fmt.Printf("%s: %v\n", endText, end)
}

func someLoadFunc() {
	defer trackTime(time.Now(), "it took") // time.Now() calculated immediately

	fmt.Println("Data loading...")
	time.Sleep(3 * time.Second)
	fmt.Println("Data loaded")
}

func main() {
	someLoadFunc()
}
