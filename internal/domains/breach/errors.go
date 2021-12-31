package breach

import "fmt"

type ErrorCode int

type Error struct {
	ErrCode ErrorCode
	Message string
}

const (
	DatasourceErr ErrorCode = iota
	BreachNotFoundErr
	BreachValidationErr
)

func NewError(code ErrorCode, msg string) *Error {
	return &Error{ErrCode: code, Message: msg}
}

func NewErrorf(code ErrorCode, msg string, params ...interface{}) *Error {
	nmsg := fmt.Sprintf(msg, params...)
	return NewError(code, nmsg)
}

func (e Error) Error() string { return e.Message }
