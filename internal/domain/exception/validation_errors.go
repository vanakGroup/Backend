package exception

import (
	"bytes"
	"strings"
)

type ValidationError struct {
	FieldErros []FieldError
}

func NewValidationErrors() *ValidationError {
	return &ValidationError{}
}

func (ve ValidationError) Error() string {
	buff := bytes.NewBufferString("")
	for i := 0; i < len(ve.FieldErros); i++ {
		buff.WriteString(ve.FieldErros[i].Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

func (ve *ValidationError) Add(field, tag string) {
	ve.FieldErros = append(ve.FieldErros, *NewFieldError(field, tag))
}

func (ve *ValidationError) HasErrors() bool {
	return len(ve.FieldErros) > 0
}
