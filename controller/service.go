package controller

import (
	"bookshelf/config"
	e "bookshelf/error"
	"errors"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	ISBN_13_LEN = 13
)

type BookshelfService struct {
	config    *config.Config
	bookshelf map[string]Book //repo
}

func NewBookshelfService(config *config.Config) *BookshelfService {
	svc := new(BookshelfService)
	svc.bookshelf = make(map[string]Book)
	svc.config = config

	return svc
}

func (s *BookshelfService) UpdateBook(c echo.Context) error {
	var err error

	isbn := c.QueryParam("code")
	if len(isbn) != ISBN_13_LEN {
		return e.ErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("Wrong isbn - %v", isbn))
	}

	if _, exists := s.bookshelf[isbn]; exists {
		return e.ErrorResponse(c, http.StatusInternalServerError,
			errors.New("Book code is duplicated").Error())
	}

	book, err := RequestOpenAPI(s.config, isbn)
	if err != nil {
		return e.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	book.ID, err = uuid.GenerateUUID()
	if err != nil {
		return e.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	s.bookshelf[book.ISBN] = *book

	return c.JSON(http.StatusOK, book)
}

func (s *BookshelfService) GetBook(c echo.Context) error {

	isbn := c.Param("code")
	book, exists := s.bookshelf[isbn]
	if !exists {
		return e.ErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("Invalid book code - %v", isbn))
	}

	return c.JSON(http.StatusOK, book)
}

func (s *BookshelfService) GetBookList(c echo.Context) error {
	return c.JSON(http.StatusOK, s.bookshelf)
}

func (s *BookshelfService) SearchBook(c echo.Context) error {
	var result []*Book
	var err error

	target := c.QueryParam("target")
	query := c.QueryParam("query")
	switch target {
	case "title":
		result, err = searchBookUsingTitle(query)
		if err != nil {
			return e.ErrorResponse(c, http.StatusNotFound, "")
		}
	case "isbn":
		result, err = searchBookUsingISBN(query)
		if err != nil {
			return e.ErrorResponse(c, http.StatusNotFound, "")
		}
	default:
		return e.ErrorResponse(c, http.StatusBadRequest,
			fmt.Sprintf("Please check target string - %v", target))
	}

	return c.JSON(http.StatusOK, result)
}

func searchBookUsingTitle(title string) ([]*Book, error) {

	return nil, nil
}

func searchBookUsingISBN(title string) ([]*Book, error) {
	return nil, nil
}
