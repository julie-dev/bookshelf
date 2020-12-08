package controller

type Book struct {
	ID        string   `json:"uuid"`
	Title     string   `json:"title"`
	Author    []string `json:"authors"`
	ISBN      string   `json:"isbn,required"`
	RegDate   string   `json:"reg_date"`
	PubDate   string   `json:"pub_date"`
	Publisher string   `json:"publisher"`
	Thumbnail string   `json:"thumbnail"`
}
