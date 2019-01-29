package books

// Справочник книг
type Books struct {
	Id          string `json:"id" sql:"id"`
	Name        string `json:"name,omitempty" sql:"name"`
	PageCount   int    `json:"pageCount,omitempty" sql:"page_count"`
	ReleaseYear int    `json:"releaseYear,omitempty" sql:"release_year"`
}

// проверка на валидность
func (b *Books) IsValid() bool {
	return len(b.Name) < 50 && b.PageCount < 1000
}

// Имя таблицы
func (b *Books) TableName() string {
	return "books"
}
