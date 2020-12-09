package controller

import (
	"bookshelf/config"
	"bookshelf/model"
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

func RequestOpenAPI(config *config.Config, isbn string) (*model.Book, error) {

	client := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest("GET", config.OpenAPIUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("KakaoAK %v", config.OpenAPIKey))

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
	return &model.Book{
		Title:  bookInfo.Title,
		Author: bookInfo.Author,
		ISBN:   isbn,
		//RegDate:   time.Now().Local().Format(time.RFC3339),
		PubDate:   bookInfo.PubDate.Local(),
		Publisher: bookInfo.Publisher,
		Thumbnail: bookInfo.Thumbnail,
	}, nil
}
