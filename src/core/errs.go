package core

import (
	"fmt"
	"strings"
)

var (
	_err        Error
	CallerEmpty = Caller{}
)

type Error interface {
	Code() string
	Message() string
	Cause() error
	error
}

type Caller struct {
	File string
	Func string
	Line int
}

type ErrorWrapper struct {
	code    string
	message string
	cause   error
	caller  Caller
}

func (e ErrorWrapper) Code() string {
	return e.code
}

func (e ErrorWrapper) Message() string {
	return e.message
}

func (e ErrorWrapper) Cause() error {
	return e.cause
}

func (e ErrorWrapper) Error() string {
	stack := e.message
	// caller info
	if e.caller.File != "" {
		stack = strings.Join([]string{e.callerInfo(), stack}, " - ")
	}
	// append cause
	if e.cause == nil {
		return stack
	}
	return strings.Join([]string{stack, e.Cause().Error()}, "\nCause by: ")
}

func (e ErrorWrapper) callerInfo() string {
	return fmt.Sprintf("%s.%s:%d", e.caller.File, e.caller.Func, e.caller.Line)
}

func NewError(code string, message string, source error) Error {
	return ErrorWrapper{
		code:    code,
		message: message,
		cause:   source,
	}
}

func NewErrorWithCaller(code string, message string, source error, caller Caller) Error {
	return ErrorWrapper{
		code:    code,
		message: message,
		cause:   source,
		caller:  caller,
	}
}
