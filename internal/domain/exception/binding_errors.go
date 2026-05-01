package exception

type BindingError struct {
	Err error
}

func NewBinding(err error) *BindingError {
	return &BindingError{Err: err}
}

func (be BindingError) Error() string {
	return be.Err.Error()
}
