package misc

type JwtError struct{}

func (e *JwtError) Error() string {
	return "Failed claims extract."
}
