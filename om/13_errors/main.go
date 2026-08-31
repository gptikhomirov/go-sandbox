package main

import (
	"errors"
	"fmt"
)

func simpleErrors() {
	errJustText := errors.New("my custom error")
	fmt.Println(errJustText)
	fmt.Println(errJustText.Error())

	err1 := fmt.Errorf("user %d not found", 123)
	fmt.Println(err1)

	err2 := fmt.Errorf("deep error: %v", err1) // loose err1
	err3 := fmt.Errorf("deep error: %w", err1) // save err1
	fmt.Println(err2)
	fmt.Println(err3)
	// errors.Is проверяет ошибку по значению
	fmt.Println("err2 is err1 =>", errors.Is(err2, err1)) // false
	fmt.Println("err3 is err1 =>", errors.Is(err3, err1)) // true
}

// custom error
type ValidationError struct {
	Field   string
	Value   any
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("custom error of field '%s': '%s'", e.Field, e.Message)
}

func someFuncWithErrReturn() error {
	return &ValidationError{
		Field:   "datepicker",
		Value:   "0000.00.00",
		Message: "year cant be empty",
	}
}

func main() {
	//simpleErrors()

	// custom error
	err := someFuncWithErrReturn()
	if err != nil {
		// errors.As проверяет ошибку по типу
		// если сходится с переданной кастомной ошибкой - записывает в нее значение
		var validationErr *ValidationError
		if errors.As(err, &validationErr) {
			fmt.Println("THIS IS CUSTOM ERR:", err.Error())
			fmt.Println(
				validationErr.Field,
				validationErr.Value,
				validationErr.Message,
			)
		} else {
			fmt.Println("unknown err")
		}
	}
}
