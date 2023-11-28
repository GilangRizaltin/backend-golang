package helpers

type Response struct {
	Message string
	Error   string
	Data    interface{}
	Meta    interface{}
}

func InitializeResponse(message, error string, data, meta interface{}) Response {
	return Response{
		Message: message,
		Error:   error,
		Data:    data,
		Meta:    meta,
	}
}
