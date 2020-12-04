package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Book struct {
	ID			int 		`json:"id"`
	Name 		string 		`json:"name"`
	Auther 		string 		`json:"auther"`
	ISBN 		string 		`json:"isbn"`
	Date 		time.Time 	`json:"date"`
	Publisher 	string 		`json:"publisher"`
}

var (
	bookList map[string]Book
)

func NewBook() *Book {
	return &Book{
		Date: time.Now(),
	}
}

func (b *Book) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Debug("pre-process")
		if err := next(c); err != nil {
			c.Error(err)
		}
		c.Logger().Debug("post process")
		return nil
	}
}

func (b *Book) Handle(c echo.Context) error {
	c.Logger().Debug("handle")
	return c.JSON(http.StatusOK, b)
}