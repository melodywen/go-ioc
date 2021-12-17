package exceptions

import "fmt"

// BuildException build error
type BuildException struct {
	concrete interface{}
	BaseException
}

// NewBuildException build exception construct method
func NewBuildException(concrete interface{}) *BuildException {
	return &BuildException{concrete: concrete}
}

// GetErrorCode defined error code name
func (_ *BuildException) GetErrorCode() string {
	return "build-exception"
}

//Error print error message
func (e *BuildException) Error() string {
	return fmt.Sprintf("Target [%T] is not instantiable. please see upper logger", e.concrete)
}
