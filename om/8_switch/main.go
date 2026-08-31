package main

import "fmt"

// SWITCH FALLTHROUGH
//func main() {
//	value := 1
//
//	switch {
//	case value == 0:
//		fmt.Println(0)
//	case value == 1:
//		fmt.Println(1)
//		fallthrough
//	case value == 2: // FALSE but executed after 1
//		fmt.Println(2)
//		fallthrough
//	case value == 3: // FALSE but executed after 2
//		fmt.Println(3)
//		fallthrough
//	default: // FALSE but executed after 3
//		fmt.Println("default")
//	}
//}

// TYPE SWITCH
func identifyType(v interface{}) {
	switch v.(type) {
	case int32:
		fmt.Println("ITS int32")
	case int64:
		fmt.Println("ITS int64")
	case int:
		fmt.Println("ITS int")
	case bool:
		fmt.Println("ITS bool")
	case string:
		fmt.Println("ITS string")
	default:
		fmt.Println("unknown")
	}
}
func main() {
	value := 123

	identifyType(value)
}
