package router

import (
	"bookshelf/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func New() *echo.Echo {

	e := echo.New()
	//e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// simple handler
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	//TODO middlerware
	b := controller.NewBook()
	e.Use(b.Process)
	e.GET("/", b.Handle)

	return e
}
