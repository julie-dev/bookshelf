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

	e.POST("/books/update", controller.SaveBook)
	e.GET("/books/search", controller.SearchBook)
	e.GET("/books/:code", controller.GetBook)
	e.GET("/books", controller.GetBookList)

	return e
}
