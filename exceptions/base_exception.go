package exceptions

// BaseExceptionInterface this package error interface
type BaseExceptionInterface interface {
	GetErrorModule() string
	GetErrorCode() string
}

// BaseException base exception error
type BaseException struct {
}

// GetErrorModule defined error model name
func (_ *BaseException) GetErrorModule() string {
	return "github.com/melodywen/go-ioc"
}
