package security

import (
	"../errors"
	"../users"
	"../utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Security interface {
	Login(w http.ResponseWriter, r *http.Request)
	CheckSession(h http.HandlerFunc) http.HandlerFunc
}

type service struct {
	Sm    *SessionManager
	Users users.UserHandler
}

type LoginResponse struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func CreateInstance(sm *SessionManager, user users.UserHandler) Security {
	return &service{
		Sm:    sm,
		Users: user,
	}
}

func (s *service) Login(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(10 * time.Hour)

	cookieUserName, err := r.Cookie("username")
	if err == http.ErrNoCookie {
		errors.ErrorHandler(w, "Access denied, username is not found", http.StatusUnauthorized, err)
		return
	} else if err != nil {
		errors.ErrorHandler(w, "Get cookies error", http.StatusInternalServerError, err)
		return
	}

	cookiePassword, err := r.Cookie("password")
	if err == http.ErrNoCookie {
		errors.ErrorHandler(w, "Access denied, password is not found", http.StatusUnauthorized, err)
		return
	} else if err != nil {
		errors.ErrorHandler(w, "Get cookies error", http.StatusInternalServerError, err)
		return
	}

	user := &users.User{
		Id:       utils.GenerateUUID(),
		Login:    cookieUserName.Value,
		Password: cookiePassword.Value,
	}
	err = s.Users.PutUser(user)

	if err != nil {
		errors.ErrorHandler(w, "Put User error", http.StatusInternalServerError, err)
		return
	}

	sessionUser, err := s.Users.GetUser(user.Id)

	if err != nil {
		errors.ErrorHandler(w, "Get user error", http.StatusInternalServerError, err)
		return
	}

	sessionId, err := s.Sm.Create(sessionUser, &Session{
		Login:    cookieUserName.Value,
		Password: cookiePassword.Value,
	})

	if err != nil {
		errors.ErrorHandler(w, "Create error", http.StatusInternalServerError, err)
		return
	}

	// TODO:: add token in cookie and expire time for session_id
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sessionId.ID,
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	log.Println("log In")

	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(&LoginResponse{
		Message: "Log In",
		Error:   err,
	})
	w.Write([]byte(b))
	//http.Redirect(w, r, "/books", http.StatusFound)
}

func (s *service) CheckSession(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieSessionID, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			log.Println("No session_id", err)
			s.Login(w, r)
		} else if err != nil {
			log.Println("Error cookie", err)
			errors.ErrorHandler(w, "Error cookie", http.StatusUnauthorized, err)
			return
		}
		ok, err := s.Sm.Check(&SessionID{ID: cookieSessionID.Value})
		if err != nil {
			log.Println("Error check session", err)
			errors.ErrorHandler(w, "Error check session", http.StatusUnauthorized, err)
			return
		}
		if !ok {
			s.Login(w, r)
		}
		h.ServeHTTP(w, r)
	})
}
