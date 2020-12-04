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
	e.Use(controller.Process)
	e.GET("/", controller.TestHandler)
	e.PUT("/save/:code", controller.SaveBook)

	return e
}
