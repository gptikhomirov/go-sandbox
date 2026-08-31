package main

import (
	"fmt"
	"net/http"
)

type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string {
	return e.Msg
}

type NotFoundError struct {
	Msg string
}

func (e *NotFoundError) Error() string {
	return e.Msg
}

func errorToStatus(err error) int {
	switch err := err.(type) {
	case *ValidationError:
		fmt.Println("validation:", err.Msg)
		return http.StatusBadRequest
	case *NotFoundError:
		fmt.Println("not found:", err.Msg)
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func main() {
	err1 := &NotFoundError{
		Msg: "user not found",
	}
	err2 := &ValidationError{
		Msg: "user not found",
	}

	status1 := errorToStatus(err1)
	status2 := errorToStatus(err2)

	fmt.Println(status1) // 404
	fmt.Println(status2) // 400
}
