package exception

import "fmt"

type ConflictError struct {
	Item string
}

func (c ConflictError) Error() string {
	return fmt.Sprintf("Conflict %s", c.Item)
}

func NewConflictError(item string) *ConflictError {
	return &ConflictError{Item: item}
}
