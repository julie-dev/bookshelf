package error

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorMessage struct {
	Code    int    `json:"code,required"`
	Status  string `json:"status,required"`
	Message string `json:"message"`
}

func ErrorResponse(c echo.Context, code int, msg string) error {
	return c.JSON(code, NewErrorMessage(code, msg))
}

func NewErrorMessage(code int, msg string) *ErrorMessage {
	return &ErrorMessage{
		Code:    code,
		Status:  http.StatusText(code),
		Message: msg,
	}
}
