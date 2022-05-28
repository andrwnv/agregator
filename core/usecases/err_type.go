package usecases

type UsecaseError struct {
	Message string
}

func MakeUsecaseError(errText string) *UsecaseError {
	return &UsecaseError{Message: errText}
}

func (t *UsecaseError) Error() string {
	return t.Message
}
