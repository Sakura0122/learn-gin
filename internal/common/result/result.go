package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

const (
	CodeSuccess      = 0
	CodeError        = 1
	CodeParamError   = 400
	CodeUnauthorized = 401
	CodeNotFound     = 404
	CodeServerError  = 500
)

func Result(code int, data any, message string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Data:    data,
		Message: message,
	})
}

func Success(data any, c *gin.Context) {
	Result(CodeSuccess, data, "ok", c)
}

func SuccessWithMessage(message string, data any, c *gin.Context) {
	Result(CodeSuccess, data, message, c)
}

func Error(code int, message string, c *gin.Context) {
	Result(code, nil, message, c)
}

func ErrorWithData(code int, message string, data any, c *gin.Context) {
	Result(code, data, message, c)
}
