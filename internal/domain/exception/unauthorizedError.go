package exception

import "net/http"

type UnauthorizedError struct {
	Message    string
	StatusCode int
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

func NewUnauthorizedError(message string) UnauthorizedError {
	return UnauthorizedError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}
