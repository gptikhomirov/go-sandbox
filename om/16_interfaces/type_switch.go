package main

import "fmt"

type Movie struct {
	Name string
	Year int
}
type Person struct {
	Name  string
	Email string
}

func printByType(v interface{}) {
	switch value := v.(type) {
	case int:
		fmt.Println("its int")
	case string:
		fmt.Println("its string")
	case bool:
		fmt.Println("its bool")
	case Movie:
		fmt.Println("its movie", value.Name, value.Year)
	default:
		fmt.Println("unknown type")
	}
}

func main() {
	printByType(Movie{Name: "Kek", Year: 2000})
	printByType(Person{Name: "Ivan", Email: "kek@mail.com"})
	printByType(123)
}
