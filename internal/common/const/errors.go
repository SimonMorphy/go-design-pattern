package consts

import (
	"errors"
	"fmt"
)

type Error struct {
	code int
	msg  string
	err  error
}

func (e *Error) Error() string {
	var msg string
	if e.msg != "" {
		msg = e.msg
	}
	msg = ErrMsg[e.code]
	return msg + " -> " + e.err.Error()
}

func New(code int) error {
	return &Error{
		code: code,
	}
}

func NewWithError(code int, err error) error {
	if err == nil {
		return New(code)
	}
	return &Error{
		code: code,
		err:  err,
	}
}

func NewWithMsgf(code int, format string, args ...any) error {
	return &Error{
		code: code,
		msg:  fmt.Sprintf(format, args...),
	}
}

func Errno(err error) int {
	if err == nil {
		return ErrnoSuccess
	}
	targetError := &Error{}
	if errors.As(err, &targetError) {
		return targetError.code
	}
	return -1
}

func Output(err error) (int, string) {
	if err == nil {
		return ErrnoSuccess, ErrMsg[ErrnoSuccess]
	}
	errno := Errno(err)
	if errno == -1 {
		return ErrnoUnknownError, err.Error()
	}
	return errno, err.Error()
}
