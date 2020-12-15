package error

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorMessage struct {
	Code    int    `json:"code,required"`
	Status  string `json:"status,required"`
	Message string `json:"message"`
} //@name ErrorResponseFormat

func ErrorResponse(c echo.Context, code int, msg string) error {
	return c.JSON(code, newErrorMessage(code, msg))
}

func BadRequest(c echo.Context, msg string) error {
	return c.JSON(http.StatusBadRequest, newErrorMessage(http.StatusBadRequest, msg))
}

func InternalError(c echo.Context, msg string) error {
	return c.JSON(http.StatusInternalServerError, newErrorMessage(http.StatusInternalServerError, msg))
}

func newErrorMessage(code int, msg string) *ErrorMessage {
	return &ErrorMessage{
		Code:    code,
		Status:  http.StatusText(code),
		Message: msg,
	}
}
