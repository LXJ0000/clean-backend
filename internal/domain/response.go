package domain

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   error       `json:"error"`
}

func ErrorResp(msg string, err error) Response {
	return Response{
		Code:    1,
		Message: msg,
		Error:   err,
	}
}

func SuccesResp(data interface{}) Response {
	return Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}
