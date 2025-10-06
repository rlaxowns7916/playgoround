package utils

import (
	"fmt"
)

type NotFoundError struct {
	queryInfo QueryInfo
	err       error
}

func NewNotFoundError(query string, args []interface{}, err error) *NotFoundError {
	return &NotFoundError{
		queryInfo: NewQueryInfo(query, args),
		err:       err,
	}
}

func (e *NotFoundError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("not found: %s (cause: %v)", e.queryInfo.String(), e.err)
	}
	return fmt.Sprintf("not found: %s", e.queryInfo.String())
}

func (e *NotFoundError) Unwrap() error {
	return e.err
}

type SystemError struct {
	queryInfo QueryInfo
	err       error
}

func NewSystemError(query string, args []interface{}, err error) *SystemError {
	return &SystemError{
		queryInfo: NewQueryInfo(query, args),
		err:       err,
	}
}

func (e *SystemError) Error() string {
	return fmt.Sprintf("system error while executing query '%s': %v", e.queryInfo.String(), e.err)
}

func (e *SystemError) Unwrap() error {
	return e.err
}
