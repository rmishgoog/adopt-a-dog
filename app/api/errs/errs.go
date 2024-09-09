package errs

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

// Return a new Error.
func New(code ErrCode, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}

// Return a new Error with a formatted message.
func Newf(code ErrCode, format string, args ...any) Error {
	return Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Implements the error interface.
func (e Error) Error() string {
	return e.Message
}

// Check if the error is an Error type defined here.
func IsError(err error) bool {
	var er Error
	return errors.As(err, &er)
}

// Get the copy of the Error type from the error.
func GetError(err error) Error {
	var er Error
	if !errors.As(err, &er) {
		return Error{}
	}
	return er
}
