package controller

import (
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"time"
)

type Info struct {
	Code 		string 		`json:"code"`
}

type Book struct {
	ID			int 		`json:"id"`
	Name 		string 		`json:"name"`
	Auther 		string 		`json:"auther"`
	ISBN 		string 		`json:"isbn"`
	Date 		time.Time 	`json:"date"`
	Publisher 	string 		`json:"publisher"`
}

func NewBook() *Book {
	return &Book{
		Date: time.Now(),
	}
}

func Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

func TestHandler(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, string(body))
}

func SaveBook(c echo.Context) error {

	book, err := RequestOpenAPI(c.Param("code"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, book)
}
