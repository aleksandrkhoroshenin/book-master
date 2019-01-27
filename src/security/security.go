package security

import (
	"net/http"
	"time"
)

type Security interface {
	Login(w http.ResponseWriter, r *http.Request)
	CheckSession(r *http.Request) (*Session, error)
}

type service struct {
	Sm *SessionManager
}

func CreateInstance(sm *SessionManager) Security {
	return &service{
		Sm: sm,
	}
}

func (s *service) Login(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(10 * time.Hour)

	cookieUserName, err := r.Cookie("username")
	if err == http.ErrNoCookie {
		return
	} else if err != nil {
		return
	}

	cookiePassword, err := r.Cookie("password")
	if err == http.ErrNoCookie {
		return
	} else if err != nil {
		return
	}

	sessionId, err := s.Sm.Create(&Session{
		Login:    cookieUserName.Value,
		Password: cookiePassword.Value,
	})

	if err != nil {
		return
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sessionId.ID,
		Expires: expiration,
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("session_id", sessionId.ID)
	http.Redirect(w, r, "/book", http.StatusFound)
}

func (s *service) CheckSession(r *http.Request) (*Session, error) {
	cookieSessionID, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	session := s.Sm.Check(&SessionID{
		ID: cookieSessionID.Value,
	})
	return session, nil
}
