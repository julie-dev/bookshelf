package controller

import (
	"bookshelf/config"
	"bookshelf/database"
	e "bookshelf/error"
	"errors"
	"fmt"
	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	ISBN_13_LEN = 13
)

type BookshelfService struct {
	config     *config.Config
	Repository *database.Repository
	Handler    *ConnectionHandler
}

func NewBookshelfService(config *config.Config, repository *database.Repository) *BookshelfService {
	svc := new(BookshelfService)
	svc.config = config
	svc.Repository = repository
	svc.Handler = &ConnectionHandler{
		Send:          DoRequest,
		CurrentStatus: ConnectionStatusClose,
		RetryCount:    0,
		MaxRetryCount: 100,
	}

	return svc
}

// UpdateBook godoc
// @Description search book info using book code and store it in DB
// @Produce json
// @Param code query string true "Book Code"
// @Success 200 {array} model.Book
// @Failure 400 {object} e.ErrorMessage
// @Failure 500 {object} e.ErrorMessage
// @Router /update [post]
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

	book, err := s.RequestOpenAPI(isbn)
	if err != nil {
		return e.InternalError(c, err.Error())
	}

	if err := s.Repository.SaveBook(session, book); err != nil {
		return e.InternalError(c, err.Error())
	}

	return c.JSON(http.StatusOK, book)
}

// GetBook godoc
// @Description get book by Code
// @Produce json
// @Param code path string true "Book Code"
// @Success 200 {object} model.Book
// @Failure 400 {object} e.ErrorMessage
// @Failure 500 {object} e.ErrorMessage
// @Router /{code} [get]
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

// GetBookList godoc
// @Description get total book list
// @Produce json
// @Success 200 {array} model.Book
// @Failure 400 {object} e.ErrorMessage
// @Failure 500 {object} e.ErrorMessage
// @Router /list [get]
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

// SearchBook godoc
// @Description query book info stored in DB
// @Produce json
// @Param target query string true "Search Type"
// @Param query query string true "Search String"
// @Success 200 {array} model.Book
// @Failure 400 {object} e.ErrorMessage
// @Failure 500 {object} e.ErrorMessage
// @Router /search [get]
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
