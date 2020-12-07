package controller

import (
	e "bookshelf/error"
	"errors"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	ISBN_13_LEN = 13
)

type Book struct {
	ID        string    `json:"uuid"`
	Name      string    `json:"name"`
	Author    string    `json:"author"`
	ISBN      string    `json:"isbn,required"`
	Date      time.Time `json:"date"`
	Publisher string    `json:"publisher"`
}

var (
	bookshelf map[string]Book
)

func init() {
	bookshelf = make(map[string]Book)
}

func SaveBook(c echo.Context) error {
	var err error

	isbn := c.QueryParam("code")
	if len(isbn) != ISBN_13_LEN {
		return e.ErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("Wrong isbn - %v", isbn))
	}

	if _, exists := bookshelf[isbn]; exists {
		return e.ErrorResponse(c, http.StatusInternalServerError,
			errors.New("Book code is duplicated").Error())
	}

	book, err := RequestOpenAPI(isbn)
	if err != nil {
		return e.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	book.ID, err = uuid.GenerateUUID()
	if err != nil {
		return e.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	bookshelf[book.ISBN] = *book

	return c.JSON(http.StatusOK, book)
}

func GetBook(c echo.Context) error {

	isbn := c.Param("code")
	book, exists := bookshelf[isbn]
	if !exists {
		return e.ErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("Invalid book code - %v", isbn))
	}

	return c.JSON(http.StatusOK, book)
}

func GetBookList(c echo.Context) error {
	return c.JSON(http.StatusOK, bookshelf)
}

func SearchBook(c echo.Context) error {
	var result []Book

	target := c.QueryParam("target")
	query := c.QueryParam("query")
	switch target {
	case "title":
		searchBookTitle(query)
	case "isbn":
		searchBookISBN(query)
	default:
		return e.ErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("Please check target string - %v", target))
	}

	return c.JSON(http.StatusOK, result)
}

func searchBookTitle(title string) {
}

func searchBookISBN(title string) {
}
