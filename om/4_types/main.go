package main

import (
	"fmt"
	"reflect"
)

func main() {
	someBool := true
	someInt := 1
	someFloat := 3.1415926
	someStr := "some string"

	fmt.Printf("type of 1 is %T, type of 2 is %T\n", someStr, someInt)

	fmt.Println(reflect.TypeOf(someBool))
	fmt.Println(reflect.TypeOf(someFloat))
}
