package gosmartstring

func CreateSSError(code int, message string) *SSError {
	return &SSError{
		IObject: &SSObject{},
		Code:    code,
		Message: CreateString(message),
	}
}

func (err SSError) Error() string {

	return err.Message.Value
}

type SSError struct {
	IObject
	Code    int
	Message *SSString
}

func (err *SSError) GetType() string {
	return "sserror"
}
