package main

type JsonErrors map[string]string

type ValidationError interface {
	error

	JsonErrors() JsonErrors
}

type JsonValidationError struct {
	errMsg     string
	jsonErrors JsonErrors
}

func (v *JsonValidationError) Error() string {
	return v.errMsg
}

func (v *JsonValidationError) JsonErrors() JsonErrors {
	return v.jsonErrors
}

func NewValidationError(msg string, errors JsonErrors) ValidationError {
	return &JsonValidationError{msg, errors}
}
