package exception

type AlreadyExistsErrorFields struct {
	Field    string
	Tag      string
	TransKey string
}

type AlreadyExistsError struct {
	FieldErrors []AlreadyExistsErrorFields
}

func NewAlreadyExistsError() *AlreadyExistsError {
	return &AlreadyExistsError{
		make([]AlreadyExistsErrorFields, 0),
	}
}

func (aee *AlreadyExistsError) Add(field string, tag string, transKey string) {
	aee.FieldErrors = append(aee.FieldErrors, AlreadyExistsErrorFields{
		Field:    field,
		Tag:      tag,
		TransKey: transKey,
	})
}

func (aee *AlreadyExistsError) Error() string {
	return "already exists"
}

func (aee *AlreadyExistsError) HasErrors() bool {
	return len(aee.FieldErrors) > 0
}
