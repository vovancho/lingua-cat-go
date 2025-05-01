package error

import "errors"

var (
	InvalidDecodeJsonError = errors.New("INVALID_JSON")
	InvalidPathParamError  = errors.New("INVALID_PATH_PARAM")
)

type AppErrorInterface interface {
	error
	StatusCode() int
	Error() string
	Unwrap() error
}

type AppError struct {
	status  int
	message string
	err     error
}

func (e *AppError) StatusCode() int {
	return e.status
}

func (e *AppError) Error() string {
	return e.message
}

func (e *AppError) Unwrap() error {
	return e.err
}

func NewAppError(status int, message string, err error) *AppError {
	return &AppError{
		status:  status,
		message: message,
		err:     err,
	}
}
