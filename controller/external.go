package controller

import (
	"bookshelf/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Document struct {
	Title     string    `json:"title"`
	Author    []string  `json:"authors"`
	PubDate   time.Time `json:"datetime"`
	Publisher string    `json:"publisher"`
	Thumbnail string    `json:"thumbnail"`
	ISBN      string    `json:"isbn"`
}

type OpenAPIResponse struct {
	Documents []Document `json:"documents"`
}

func (s *BookshelfService) RequestOpenAPI(isbn string) (*model.Book, error) {
	h := s.Handler

	switch h.CurrentStatus {
	case ConnectionStatusClose:
		//go straight
		req, err := NewRequest(s.config.OpenAPIUrl, s.config.OpenAPIKey, isbn)
		if err != nil {
			return nil, err
		}

		resp, err := h.Send(req)
		if err != nil {
			if s.Handler.RetryCount > s.Handler.MaxRetryCount {
				s.Handler.CurrentStatus = ConnectionStatusOpen
				s.Handler.Expire = time.Now().Add(time.Second * s.Handler.Timeout)
			} else {
				s.Handler.RetryCount++
			}
		}

		return ReturnBookInfo(resp)
	case ConnectionStatusHalfOpen:
		//randomly
	case ConnectionStatusOpen:
		//drop request
	default:
	}

	return nil, nil
}

func NewRequest(url string, key string, isbn string) (*http.Request, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("KakaoAK %v", key))

	q := req.URL.Query()
	q.Add("target", "isbn")
	q.Add("query", isbn)
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func DoRequest(req *http.Request) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	return resp, err
}

func ReturnBookInfo(resp *http.Response) (*model.Book, error) {

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var bookResponse OpenAPIResponse
	err := decoder.Decode(&bookResponse)
	if err != nil {
		return nil, err
	}

	if len(bookResponse.Documents) == 0 {
		return nil, errors.New("The book information could not be found")
	}

	bookInfo := bookResponse.Documents[0]
	return &model.Book{
		Title:     bookInfo.Title,
		Author:    bookInfo.Author,
		ISBN:      strings.Trim(bookInfo.ISBN, " "),
		PubDate:   bookInfo.PubDate.Local(),
		Publisher: bookInfo.Publisher,
		Thumbnail: bookInfo.Thumbnail,
	}, nil
}
