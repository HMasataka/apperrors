package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
)

type E struct {
	Code  int
	Err   error
	frame Frame
}

func (e *E) Error() string {
	return e.Err.Error()
}

type Frame struct {
	frames [3]uintptr
}

func caller(skip int) Frame {
	var s Frame
	runtime.Callers(skip+1, s.frames[:])
	return s
}

func (f Frame) location() (function, file string, line int) {
	frames := runtime.CallersFrames(f.frames[:])
	if _, ok := frames.Next(); !ok {
		return "", "", 0
	}
	fr, ok := frames.Next()
	if !ok {
		return "", "", 0
	}
	return fr.Function, fr.File, fr.Line
}

func StatusCode(err error) int {
	if e, ok := err.(*E); ok {
		fmt.Println(e.frame.location())
		return e.Code
	}
	return http.StatusBadRequest
}

func New() error {
	return &E{
		Code:  http.StatusInternalServerError,
		Err:   errors.New("error"),
		frame: caller(1),
	}
}

func main() {
	err := New()
	fmt.Println(err)
	fmt.Println(StatusCode(err))
}
