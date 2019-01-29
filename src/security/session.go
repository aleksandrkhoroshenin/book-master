package security

import (
	"../users"
	"../utils"
	"database/sql"
)

// Обработчик сессий
type SessionHandler interface {
	Create(user *users.User, in *Session) (*SessionID, error)
	Check(in *SessionID) (bool, error)
	Delete(in *SessionID)
}

// Описание сессии
type Session struct {
	Login    string
	Password string
}

// id сессии
type SessionID struct {
	ID string `json:"session_id"`
}

// TODO::Redis
type SessionManager struct {
	db *sql.DB
}

func NewSessionManager(db *sql.DB) *SessionManager {
	return &SessionManager{
		db: db,
	}
}

func (sm *SessionManager) Create(user *users.User, in *Session) (*SessionID, error) {
	session := SessionID{utils.GenerateUUID()}
	err := sm.putSession(session.ID, user.Id)
	return &session, err
}

func (sm *SessionManager) putSession(id, user_id string) error {
	_, err := sm.db.Exec("insert into session(id, user_id) values($1, $2)", id, user_id)
	return err
}

func (sm *SessionManager) getSession(idIn string) (int, error) {
	var count int
	err := sm.db.QueryRow("select count(*) from session where id = $1", idIn).Scan(&count)
	return count, err
}

func (sm *SessionManager) Check(sessionId *SessionID) (bool, error) {
	count, err := sm.getSession(sessionId.ID)
	if err != nil {
		return false, err
	}
	return count > 0, err
}

// TODO:: удаление сессий
func (sm *SessionManager) Delete(in *SessionID) {

}
