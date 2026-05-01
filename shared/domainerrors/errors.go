package domainerrors

import "fmt"

type DomainError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

func (e *DomainError) ErrorCode() string    { return e.Code }
func (e *DomainError) ErrorMessage() string { return e.Message }

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func New(code, message string, statusCode int, err error) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}
