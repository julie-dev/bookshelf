package router

import (
	"bookshelf/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(service *controller.BookshelfService) *echo.Echo {

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.Use(service.Repository.Transaction)

	books := e.Group("/books")
	books.POST("/update", service.UpdateBook)
	books.GET("/search", service.SearchBook)
	books.GET("/:code", service.GetBook)
	books.GET("/", service.GetBookList)

	return e
}
