package exception

type FormatErrorFields struct {
	Field    string
	Tag      string
	TransKey string
}

type FormatError struct {
	FieldErrors []FormatErrorFields
}

func NewFormatError() *FormatError {
	return &FormatError{
		make([]FormatErrorFields, 0),
	}
}

func (fe *FormatError) Add(field string, tag string, transkey string) {
	fe.FieldErrors = append(fe.FieldErrors, FormatErrorFields{
		Field:    field,
		Tag:      tag,
		TransKey: transkey,
	})
}

func (fe *FormatError) HasErrors() bool {
	return len(fe.FieldErrors) > 0
}

func (fe *FormatError) Error() string {
	return "wrong format"
}
