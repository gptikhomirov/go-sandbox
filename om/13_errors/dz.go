package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID   string
	Name string
}

type FieldValidationError struct {
	Field   string
	Value   any
	Message string
}

func (e *FieldValidationError) Error() string {
	return fmt.Sprintf("поле %s со значением %v не валидно", e.Field, e.Value)
}

func parseDate(date string) error {
	dateError := &FieldValidationError{
		Field: "date",
		Value: date,
	}

	if len(date) == 0 {
		dateError.Message = "поле дата обязательно к заполнению"
		return dateError
	}

	parts := strings.Split(date, "-")
	if len(parts) != 3 ||
		len(parts[0]) != 4 ||
		len(parts[1]) != 2 ||
		len(parts[2]) != 2 {

		dateError.Message = fmt.Sprintf("неверный формат даты: %s", date)
		return dateError
	}

	if _, err := strconv.Atoi(parts[0]); err != nil {
		dateError.Message = fmt.Sprintf("некорректный год: %s", parts[0])
		return dateError
	}

	if _, err := strconv.Atoi(parts[1]); err != nil {
		dateError.Message = fmt.Sprintf("некорректный месяц: %s", parts[1])
		return dateError
	}

	if _, err := strconv.Atoi(parts[2]); err != nil {
		dateError.Message = fmt.Sprintf("некорректный день: %s", parts[2])
		return dateError
	}

	return nil
}

func activateSubscription(userID string, date string) error {
	users := map[string]User{
		"1":   User{Name: "Ivan"},
		"243": User{Name: "Max"},
	}

	if _, ok := users[userID]; !ok {
		return ErrUserNotFound
	}

	if err := parseDate(date); err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}

	return nil
}

func main() {
	err := activateSubscription("1", "1234-12-22") // 2001-01-01
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			fmt.Println("нет такого чела")
		}

		var validationErr *FieldValidationError
		if errors.As(err, &validationErr) {
			fmt.Println("ошибка валидации:", validationErr.Message)
		}
	} else {
		fmt.Println("Подписка активна")
	}
}
