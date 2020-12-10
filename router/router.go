package router

import (
	"bookshelf/controller"
	_ "bookshelf/docs"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func New(service *controller.BookshelfService, ver string) *echo.Echo {

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.Use(service.Repository.Transaction)

	books := e.Group(fmt.Sprintf("/api/%v/books", ver))
	books.POST("/update", service.UpdateBook)
	books.GET("/search", service.SearchBook)
	books.GET("/:code", service.GetBook)
	books.GET("/list", service.GetBookList)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}
