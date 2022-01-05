package exceptions

import "fmt"

// AbstractToStringException transform package info error
type AbstractToStringException struct {
	abstract interface{}
	BaseException
}

// NewAbstractToStringException construct method
func NewAbstractToStringException(abstract interface{}) *AbstractToStringException {
	return &AbstractToStringException{abstract: abstract}
}

// GetErrorCode defined error code name
func (e *AbstractToStringException) GetErrorCode() string {
	return "-exception"
}

//Error print error message
func (e *AbstractToStringException) Error() string {
	return fmt.Sprintf("this object can not get abstract index info [%T=>%#v]", e.abstract, e.abstract)
}
