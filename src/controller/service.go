package controller

import "database/sql"

type DictHandler interface {
	IsValid() bool
	TableName() string
}

type Session struct {
	SessionID string
}

type Context struct {
	Db *sql.DB
	Session
}

func (c *Context) GetSessionID() string {
	return c.SessionID
}

