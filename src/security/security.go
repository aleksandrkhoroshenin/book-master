package security

import (
	"net/http"
	"time"
)

type Security interface {
	Login(w http.ResponseWriter, r *http.Request)
	//CheckSession(r *http.Request) (*Session, error)
	CheckSession(h http.Handler) http.Handler
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
		http.Error(w, "Access denied, username is not found", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Access denied, username is not found", http.StatusUnauthorized)
		return
	}

	cookiePassword, err := r.Cookie("password")
	if err == http.ErrNoCookie {
		http.Error(w, "Access denied, password is not found", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Access denied, password is not found", http.StatusUnauthorized)
		return
	}

	sessionId, err := s.Sm.Create(&Session{
		Login:    cookieUserName.Value,
		Password: cookiePassword.Value,
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

//func (s *service) CheckSession(r *http.Request) (*Session, error) {
//	cookieSessionID, err := r.Cookie("session_id")
//	if err == http.ErrNoCookie {
//		return nil, nil
//	} else if err != nil {
//		return nil, err
//	}
//
//	session := s.Sm.Check(&SessionID{
//		ID: cookieSessionID.Value,
//	})
//	return session, nil
//}

func (s *service) CheckSession(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieSessionID, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		} else if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if s.Sm.Check(&SessionID{ID: cookieSessionID.Value}) == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}
