package model

import "time"

type Book struct {
	ID        string    `json:"uuid" xorm:"varchar(36) notnull unique 'id' pk"`
	Title     string    `json:"title" xorm:"notnull"`
	Author    []string  `json:"authors" xorm:""`
	ISBN      string    `json:"isbn,required" xorm:"varchar(13) notnull unique 'isbn'"`
	RegDate   time.Time `json:"reg_date" xorm:"created"`
	PubDate   time.Time `json:"pub_date" xorm:""`
	Publisher string    `json:"publisher" xorm:""`
	Thumbnail string    `json:"thumbnail" xorm:""`
}
