package exception

import "fmt"

type BadRequestError struct {
	Message string
	Field   string
}

func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{
		Message: message,
	}
}

func NewBadRequestErrorWithField(field string) *BadRequestError {
	return &BadRequestError{
		Field: field,
	}
}

func (e *BadRequestError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}
