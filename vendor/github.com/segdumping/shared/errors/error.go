package errors

import (
	"bytes"
	"errors"
	"strconv"
)

type Error struct {
	error
	errCode   int
}

func New(s string) *Error {
	return &Error{error: errors.New(s)}
}

func (e *Error) WithCode(code int) *Error {
	e.errCode = code
	return e
}

func (e *Error) ErrCode() int {
	return e.errCode
}

func (e *Error) Error() string {
	var buf bytes.Buffer
	buf.WriteString(e.error.Error())
	buf.WriteString(", errCode: ")
	buf.WriteString(strconv.FormatInt(int64(e.errCode), 10))

	return buf.String()
}
