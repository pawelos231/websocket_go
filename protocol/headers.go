package protocol

import "fmt"

type CustomHeader struct {
	value   string
	toThrow error
	status  int
}

type HeaderError struct {
	HeaderName string
	Message    string
	Code       int
}

func (e *HeaderError) Error() string {
	return fmt.Sprintf("error in header %s: %s", e.HeaderName, e.Message)
}
