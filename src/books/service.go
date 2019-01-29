package books

import (
	"database/sql"
)

// Обработчик бизнес-логики
type HandlerBooks interface {
	AddBooks(book *Books) error
	EditBooks(book *Books) error
	DeleteBooks(id string) error
	GetBook(id string) (*Books, int, error)
	GetBooks() ([]*Books, int, error)
}

type service struct {
	db *sql.DB
}

// Создание инстанса
func CreateInstance(db *sql.DB) *service {
	return &service{
		db: db,
	}
}

// Добавление книги
func (s *service) AddBooks(book *Books) error {
	_, err := s.db.Exec("insert into books(id, name, page_count, release_year) values($1, $2, $3, $4)", book.Id, book.Name, book.PageCount, book.ReleaseYear)
	return err
}

// Редактирование книги
func (s *service) EditBooks(book *Books) error {
	_, err := s.db.Exec("update books set name = $1, page_count = $2, release_year=$3 where id = $4", book.Name, book.PageCount, book.ReleaseYear, book.Id)
	return err
}

// Удаление книги
func (s *service) DeleteBooks(id string) error {
	_, err := s.db.Exec("delete from books where id = $1", id)
	return err
}

// Получить книгу
func (s *service) GetBook(id string) (*Books, int, error) {
	var count int
	book := &Books{}
	err := s.db.QueryRow("select id, name, page_count, release_year, count(*) over()  from books where id = $1", id).Scan(&book.Id, &book.Name, &book.PageCount, &book.ReleaseYear, &count)
	return book, count, err
}

// Получить список книг
func (s *service) GetBooks() ([]*Books, int, error) {
	var count int
	books := make([]*Books, 0)
	rows, err := s.db.Query("select id, name, page_count, release_year, count(*) over()  from books")
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		book := &Books{}
		if err := rows.Scan(&book.Id, &book.Name, &book.PageCount, &book.ReleaseYear, &count); err != nil {
			return nil, 0, err
		}
		books = append(books, book)
	}

	return books, count, nil
}
