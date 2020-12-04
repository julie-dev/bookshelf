package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Book struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Author    string    `json:"author"`
	ISBN      string    `json:"isbn"`
	Date      time.Time `json:"date"`
	Publisher string    `json:"publisher"`
}

var (
	bookshelf map[string]Book
)

func SaveBook(c echo.Context) error {

	isbn := c.Param("code")
	if _, exists := bookshelf[isbn]; exists {
		return errors.New("Book code is duplicated")
	}

	book, err := RequestOpenAPI(isbn)
	if err != nil {
		return err
	}

	bookshelf[book.ISBN] = *book

	return c.JSON(http.StatusOK, book)
}

func GetBook(c echo.Context) error {

	return nil
}

func GetBookList(c echo.Context) error {

	return nil
}

func init() {
	bookshelf = make(map[string]Book)
}
