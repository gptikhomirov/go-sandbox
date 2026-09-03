package main

import "fmt"

type NotFoundError struct {
	ID int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("user %d not found", e.ID)
}

func findUser(id int) error {
	var err *NotFoundError
	// true, тк обычная переменная-указатель с типом *NotFoundError и значением nil
	fmt.Println("in find user:", err == nil, err)

	if id == 0 {
		err = &NotFoundError{ID: id}
	}

	return err
}

func main() {
	err := findUser(42) // вернулся интерфейс с типом *NotFoundError и значением nil

	fmt.Println("in main:", err == nil, err) // false

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	fmt.Println("OK")
}
