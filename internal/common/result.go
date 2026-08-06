package common

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

const (
	CodeSuccess      = 0
	CodeError        = 1
	CodeParamError   = 400
	CodeUnauthorized = 401
	CodeNotFound     = 404
	CodeServerError  = 500
)

func Success(data interface{}) Result {
	return Result{Code: CodeSuccess, Data: data, Message: "ok"}
}

func SuccessWithMessage(message string, data interface{}) Result {
	return Result{Code: CodeSuccess, Data: data, Message: message}
}

func Error(code int, message string) Result {
	return Result{Code: code, Data: nil, Message: message}
}

func ErrorWithData(code int, message string, data interface{}) Result {
	return Result{Code: code, Data: data, Message: message}
}
