package endpoints

type EndpointError struct {
	Message string
}

func MakeEndpointError(errText string) *EndpointError {
	return &EndpointError{Message: errText}
}

func (t *EndpointError) Error() string {
	return t.Message
}
