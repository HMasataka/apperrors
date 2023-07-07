package apperrors

import (
	"fmt"
	"net/http"
	"runtime"
)

type ErrorCode struct {
	Code int
	Name string
}

type Error struct {
	Code  ErrorCode
	Err   error
	frame Frame
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Unwrap() error {
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
	if e, ok := err.(*Error); ok {
		return e.Code.Code
	}
	return http.StatusBadRequest
}

func New(code ErrorCode, msg error) error {
	return &Error{
		Code:  code,
		Err:   msg,
		frame: caller(1),
	}
}

func Print(args ...any) {
	f := make([]any, len(args))
	for i, v := range args {
		if e, ok := v.(*Error); ok {
			fn, file, line := e.frame.location()
			f[i] = fmt.Sprintf("func: %v, file: %v, line: %v, err: %v", fn, file, line, e.Error())
			continue
		}

		f[i] = v
	}

	fmt.Print(f...)
}

func Printf(format string, args ...any) {
	f := make([]any, len(args))
	for i, v := range args {
		if e, ok := v.(*Error); ok {
			fn, file, line := e.frame.location()
			f[i] = fmt.Sprintf("func: %v, file: %v, line: %v, err: %v", fn, file, line, e.Error())
			continue
		}

		f[i] = v
	}

	fmt.Printf(format, f...)
}
