package model

import (
	"time"
)

type Article struct {
	Id    int
	Title string
	Date  time.Time
	Body  string
}

type Tag struct {
	Id   int
	Name string
}

type ArticleTag struct {
	ArticleId int
	TagId     int
}
