package model

import "time"

type Book struct {
	ISBN      string    `json:"isbn,required" xorm:"varchar(13) notnull unique 'isbn' pk"`
	Title     string    `json:"title" xorm:"notnull"`
	Author    []string  `json:"authors" xorm:""`
	RegDate   time.Time `json:"reg_date" xorm:"created"`
	PubDate   time.Time `json:"pub_date" xorm:""`
	Publisher string    `json:"publisher" xorm:""`
	Thumbnail string    `json:"thumbnail" xorm:""`
}
