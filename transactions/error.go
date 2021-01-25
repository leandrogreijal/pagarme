package transactions

import "fmt"

type error interface {
	Error() string
}

type InternalError struct {
	Path string
}

type InvalidValueError struct {
	ValueParam string
	Value      string
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("%v is invalid. Value: %v", e.ValueParam, e.Value)
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("Mundipagg internal error. Path: %v", e.Path)
}
