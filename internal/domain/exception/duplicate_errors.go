package exception

import "fmt"

type DuplicateError struct {
	Item    string
	Message string
}

func (e DuplicateError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("duplicate %s", e.Item)
}

func NewDuplicateError(item string) *DuplicateError {
	return &DuplicateError{Item: item}
}
