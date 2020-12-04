package controller

import "time"

func RequestOpenAPI(s string) (*Book, error) {
	return &Book{
		ISBN: s,
		Date: time.Now(),
	}, nil
}