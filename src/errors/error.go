package errors

type PublicError struct {
	Code    int
	Message string
}

func (e PublicError) Error() string {
	return e.Message
}
