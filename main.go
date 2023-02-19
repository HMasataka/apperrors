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

func (e E) Error() string {
	return e.Err.Error()
}

func main() {
	fmt.Println(New())
}

func New() error {
	return E{Code: http.StatusBadRequest, Err: errors.New("error")}
}
