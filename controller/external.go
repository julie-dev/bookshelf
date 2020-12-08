package controller

import (
	"bookshelf/config"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Document struct {
	Title     string    `json:"title"`
	Author    []string  `json:"authors"`
	PubDate   time.Time `json:"datetime"`
	Publisher string    `json:"publisher"`
	Thumbnail string    `json:"thumbnail"`
}

type OpenAPIResponse struct {
	Documents []Document `json:"documents"`
}

func RequestOpenAPI(config *config.Config, isbn string) (*Book, error) {

	client := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest("GET", config.RESTAPI_URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("KakaoAK %v", config.RESTAPI_KEY))

	q := req.URL.Query()
	q.Add("target", "isbn")
	q.Add("query", isbn)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var bookResponse OpenAPIResponse
	err = decoder.Decode(&bookResponse)
	if err != nil {
		return nil, err
	}

	if len(bookResponse.Documents) == 0 {
		return nil, errors.New("The book information could not be found")
	}

	bookInfo := bookResponse.Documents[0]
	return &Book{
		Title:     bookInfo.Title,
		Author:    bookInfo.Author,
		ISBN:      isbn,
		RegDate:   time.Now().Local().Format(time.RFC3339),
		PubDate:   bookInfo.PubDate.Local().Format(time.RFC3339),
		Publisher: bookInfo.Publisher,
		Thumbnail: bookInfo.Thumbnail,
	}, nil
}
