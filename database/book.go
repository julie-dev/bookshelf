package database

import (
	"bookshelf/model"
	"github.com/go-xorm/xorm"
)

func (r *Repository) SaveBook(session *xorm.Session, book *model.Book) error {
	if _, err := session.Insert(book); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetBook(session *xorm.Session, isbn string) (*model.Book, error) {
	book := new(model.Book)
	has, err := session.Where("isbn = ?", isbn).Get(book)
	if err != nil {
		return nil, err
	}

	if has == false {
		book = nil
	}

	return book, nil
}

func (r *Repository) GetBookList(session *xorm.Session) (*[]model.Book, error) {

	books := new([]model.Book)
	if err := session.Find(books); err != nil {
		return nil, err
	}

	return books, nil
}
