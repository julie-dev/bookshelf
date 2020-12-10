package controller

import (
	"bookshelf/config"
	"bookshelf/database"
	e "bookshelf/error"
	"bookshelf/model"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
)

const (
	ISBN_13_LEN = 13
)

type BookshelfService struct {
	config     *config.Config
	bookshelf  map[string]model.Book //repo
	Repository *database.Repository
}

func NewBookshelfService(config *config.Config, repository *database.Repository) *BookshelfService {
	svc := new(BookshelfService)
	svc.bookshelf = make(map[string]model.Book)
	svc.config = config
	svc.Repository = repository

	return svc
}

func (s *BookshelfService) UpdateBook(c echo.Context) error {
	var err error

	isbn := c.QueryParam("code")
	if len(isbn) != ISBN_13_LEN {
		return e.BadRequest(c, fmt.Sprintf("Wrong isbn - %v", isbn))
	}

	session := s.Repository.GetDBConn(c.Request().Context())
	if session == nil {
		return e.InternalError(c, errors.New("DB is not exists").Error())
	}

	temp, err := s.Repository.GetBook(session, isbn)
	if err != nil {
		return e.InternalError(c, err.Error())
	}

	if temp != nil {
		return e.BadRequest(c, "Book code already exists")
	}

	book, err := RequestOpenAPI(s.config, isbn)
	if err != nil {
		return e.InternalError(c, err.Error())
	}

	if err := s.Repository.SaveBook(session, book); err != nil {
		return e.InternalError(c, err.Error())
	}

	return c.JSON(http.StatusOK, book)
}

func (s *BookshelfService) GetBook(c echo.Context) error {

	isbn := c.Param("code")
	if len(isbn) != ISBN_13_LEN {
		return e.BadRequest(c, fmt.Sprintf("Wrong isbn - %v", isbn))
	}

	session := s.Repository.GetDBConn(c.Request().Context())
	if session == nil {
		return e.InternalError(c, errors.New("DB is not exists").Error())
	}

	book, err := s.Repository.GetBook(session, isbn)
	if err != nil {
		return e.InternalError(c, err.Error())
	}

	if book == nil {
		return e.BadRequest(c, fmt.Sprintf("Book code doesn't exists - %v", isbn))
	}

	return c.JSON(http.StatusOK, book)
}

func (s *BookshelfService) GetBookList(c echo.Context) error {

	session := s.Repository.GetDBConn(c.Request().Context())
	if session == nil {
		return e.InternalError(c, errors.New("DB is not exists").Error())
	}

	books, err := s.Repository.GetBookList(session)
	if err != nil {
		return e.InternalError(c, err.Error())
	}

	return c.JSON(http.StatusOK, books)
}

func (s *BookshelfService) SearchBook(c echo.Context) error {

	target := c.QueryParam("target")
	query := c.QueryParam("query")

	session := s.Repository.GetDBConn(c.Request().Context())
	if session == nil {
		return e.InternalError(c, errors.New("DB is not exists").Error())
	}

	var prep *xorm.Session
	switch target {
	case "title":
		sql, args, _ := builder.ToSQL(builder.Like{target, query})
		prep = session.Where(sql, args...)
	case "isbn":
		prep = session.Where("isbn = ?", query)
	default:
		return e.BadRequest(c, fmt.Sprintf("Unsupported target - %v", target))
	}

	books, err := s.Repository.GetBookList(prep)
	if err != nil {
		return e.InternalError(c, err.Error())
	}

	if books == nil {
		return e.BadRequest(c, "Book information not found")
	}

	return c.JSON(http.StatusOK, books)
}
