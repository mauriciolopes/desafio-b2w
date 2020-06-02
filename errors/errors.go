package errors

import (
	"log"
	"runtime"
)

// InternalError represents internal core errors.
type InternalError string

func (ie InternalError) Error() string {
	return string(ie)
}

// InputError represents input errors.
type InputError string

func (ie InputError) Error() string {
	return string(ie)
}

// NotFoundError represents not found errors.
type NotFoundError string

func (nfe NotFoundError) Error() string {
	return string(nfe)
}

// ConflictError represents conflict error.
type ConflictError string

func (ce ConflictError) Error() string {
	return string(ce)
}

// UnexpectedError returns a unexpected internal error.
func UnexpectedError(cause error) InternalError {
	buf := make([]byte, 1<<10)
	runtime.Stack(buf, false)
	var errMsg string
	if cause != nil {
		errMsg = cause.Error() + "\n"
	}
	log.Println(errMsg + string(buf))
	return InternalError("Ocorreu um erro inesperado.")
}
