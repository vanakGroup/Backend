package exception

type AppError struct {
	Message    string
	StatusCode int
}

func (e AppError) Error() string {
	return e.Message
}
