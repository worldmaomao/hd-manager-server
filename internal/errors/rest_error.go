package errors

type RestError struct {
	HttpResponseCode int
	ErrorCode        int
	ErrorMsg         string
}

func (e RestError) Error() string {
	return e.ErrorMsg
}

func NewRestError(httpResponseCode int, errorCode int, msg string) RestError {
	return RestError{
		HttpResponseCode: httpResponseCode,
		ErrorCode:        errorCode,
		ErrorMsg:         msg,
	}
}
