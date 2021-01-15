package errors

type DuplicationError struct {
	msg string
}

func NewDuplicationError(msg string) DuplicationError {
	return DuplicationError{
		msg: msg,
	}
}

func (err DuplicationError) Error() string {
	return err.msg
}
