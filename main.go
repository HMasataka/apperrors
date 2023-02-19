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

func (e *E) Unwrap() error {
	return e.Err
}

type Frame struct {
	frames [3]uintptr
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

func caller(skip int) Frame {
	var s Frame
	runtime.Callers(skip+1, s.frames[:])
	return s
}

func StatusCode(err error) int {
	if e, ok := err.(*E); ok {
		fmt.Println(e.frame.location())
		return e.Code
	}
	return http.StatusBadRequest
}

func Wrap(code int, err error) error {
	return &E{
		Code:  code,
		Err:   err,
		frame: caller(1),
	}
}

func New(code int, msg string) error {
	return &E{
		Code:  code,
		Err:   errors.New(msg),
		frame: caller(1),
	}
}

func main() {
	err := Wrap(http.StatusAccepted, errors.New("err"))
	fmt.Println(err)
	fmt.Println(StatusCode(err))
}
