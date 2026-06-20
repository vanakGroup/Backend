package exception

type DisabledError struct {
	FieldsError []FieldError
}

func NewDisabledError() *DisabledError {
	return &DisabledError{
		make([]FieldError, 0),
	}
}

func (de *DisabledError) Add(field, tag string) {
	de.FieldsError = append(de.FieldsError, FieldError{
		Field: field,
		Tag:   tag,
	})
}

func (de *DisabledError) Error() string {
	return "field is Disabled"
}

func (de *DisabledError) HasErrors() bool {
	return len(de.FieldsError) > 0
}
