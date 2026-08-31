package main

import "fmt"

func fmtMainFuncs() {
	fmt.Print("Hello")
	fmt.Print(" world\n")

	fmt.Println("Hello world")

	fmt.Printf("Hello %v", "world")
	
	str := fmt.Sprintf("Hello %d %s", 4, "world")
	fmt.Println(str, "\n\n")
}

var outter = 1 // NOT outter := 1 (only in funcs)

func main() {
	fmtMainFuncs()

	//fmt.Printf("type of 1 is %T, type of 2 is %T", "some str", outter)

	//
	const (
		base_price  = 999.00
		tariff_name = "Premium"
	)
	var month int = 12
	var discount = 15.0
	total := base_price * float64(month) * (100 - discount) / 100

	fmt.Println("Taриф:", tariff_name)
	fmt.Println("Месяцев:", month)
	fmt.Printf("Базовая цена: %.2f руб/мес\n", base_price)
	fmt.Printf("Скидка: %.f%%\n", discount)
	fmt.Printf("Итого: %.2f руб.\n", total)
}
