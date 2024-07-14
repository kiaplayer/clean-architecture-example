package errors

import (
	"fmt"
)

type AppError struct {
	reason string
	cause  error
}

func (e *AppError) Error() string {
	if e.cause == nil {
		return e.reason
	}

	return fmt.Sprintf("%s (cause: %s)", e.reason, e.cause.Error())
}

func NewAppError(reason string, cause error) AppError {
	return AppError{
		reason: reason,
		cause:  cause,
	}
}
