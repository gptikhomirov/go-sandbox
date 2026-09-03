package main

import (
	"fmt"
	"reflect"
)

func main() {
	var i interface{} = 123
	fmt.Println(reflect.TypeOf(i))

	str, ok := i.(string)
	if ok {
		fmt.Println("ok: ", reflect.TypeOf(str))
	}

	integer, ok := i.(int)
	if ok {
		fmt.Println("ok: ", reflect.TypeOf(integer))
	}

	str2 := i.(string)
	fmt.Println(reflect.TypeOf(str2)) // panic

}
