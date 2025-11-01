package error

import (
	"runtime/debug"
)

type MyError struct {
	err        error
	stacktrace string
	message    string
}

func Wrap(err error, message string) *MyError {
	return &MyError{
		err:        err,
		stacktrace: string(debug.Stack()),
		message:    message}
}

func (e *MyError) Error() string {
	return e.message
}
