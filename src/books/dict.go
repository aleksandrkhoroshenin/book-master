package books

import (
	"time"
)

type Books struct {
	Id string `json:"id" sql:"id"`
	Name string `json:"name" sql:"name"`
	PageCount int `json:"pageCount" sql:"page_count"`
	ReleaseYear time.Time `json:"releaseYear" sql:"release_year"`
}

func (b * Books) IsValid() bool {
	return len(b.Name) < 50 && b.PageCount < 1000
}

func (b * Books) TableName() string {
	return "books"
}

