package types

type ServiceError struct {
	StatusCode int
	ErrorCode  string
}

func (s ServiceError) Error() string {
	return ""
}
