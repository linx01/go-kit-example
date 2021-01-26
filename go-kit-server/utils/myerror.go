package utils

// MyError ...
type MyError struct {
	Code    int
	Message string
}

// Error 实现error接口
func (m *MyError) Error() string {
	return m.Message
}

// GenError ...
func GenError(code int, message string) error {
	return &MyError{Code: code, Message: message}
}
