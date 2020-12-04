package router

import (
	"bookshelf/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.POST("/book/:code", controller.SaveBook)
	e.GET("/book", controller.GetBookList)
	e.GET("/book/:code", controller.GetBook)

	return e
}
