package main

import (
	"fmt"
	"os"
)

func writeToTempFile() {
	file, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		fmt.Println("error create file:", err)
		return
	}

	defer func() {
		if err := os.Remove(file.Name()); err != nil {
			fmt.Println("error remove file:", err)
		}
		fmt.Println("file removed") // 2
	}()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("error close file:", err)
		}

		fmt.Println("file closed") // 1
	}()

	fmt.Println("file created")

	if _, err := file.WriteString("Some data"); err != nil {
		fmt.Println("error write", err)
		return
	}

	fmt.Println("data is written")
}

func main() {
	writeToTempFile()
}
