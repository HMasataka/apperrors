package main

import (
	"errors"
	"fmt"
	"net/http"
)

type E struct {
	Code int
	Err  error
}

func (e *E) Error() string {
	return e.Err.Error()
}

func StatusCode(err error) int {
	if e, ok := err.(*E); ok {
		return e.Code
	}
	return http.StatusBadRequest
}

func main() {
	err := New()
	fmt.Println(err)
	fmt.Println(StatusCode(err))
}

func New() error {
	return &E{Code: http.StatusInternalServerError, Err: errors.New("error")}
}
