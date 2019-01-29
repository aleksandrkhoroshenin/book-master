package users

import (
	"database/sql"
)

type UserHandler interface {
	PutUser(user *User) error
	GetUser(id string) (*User, error)
}

type service struct {
	db *sql.DB
}

func CreateInstance(db *sql.DB) UserHandler {
	return &service{
		db: db,
	}
}

func (s *service) PutUser(user *User) error {
	_, err := s.db.Exec("insert into users(id, login, password) values($1, $2, $3)", user.Id, user.Login, user.Password)
	return err
}

func (s *service) GetUser(id string) (*User, error) {
	user := &User{}
	err := s.db.QueryRow("select id, login, password from users where id = $1", id).Scan(&user.Id, &user.Login, &user.Password)
	return user, err
}
