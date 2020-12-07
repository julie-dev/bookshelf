package controller

import "time"

func RequestOpenAPI(isbn string) (*Book, error) {
	return &Book{
		ISBN: isbn,
		Date: time.Now().Local(),
	}, nil
}
