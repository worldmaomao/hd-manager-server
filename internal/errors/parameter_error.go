package errors

type ParameterError struct {
	msg string
}

func NewParameterError(msg string) ParameterError {
	return ParameterError{
		msg: msg,
	}
}

func (err ParameterError) Error() string {
	return err.msg
}
