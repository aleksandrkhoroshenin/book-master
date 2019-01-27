package security

import (
	"fmt"
	"math/rand"
	"sync"
)

type SessionHandler interface {
	Create(in *Session) (*SessionID, error)
	Check(in *SessionID) *Session
	Delete(in *SessionID)
}

type Session struct {
	Login     string
	Password string
}

type SessionID struct {
	ID string
}

type SessionManager struct {
	mu sync.RWMutex
	sessions map[SessionID]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mu: sync.RWMutex{},
		sessions: map[SessionID]*Session{},
	}
}

func (sm *SessionManager) Create(in *Session) (*SessionID, error) {
	sm.mu.Lock()
	id := SessionID{GenerateUUID()}
	sm.mu.Unlock()
	sm.sessions[id] = in
	return &id, nil
}

func (sm *SessionManager) Check(in *SessionID) *Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if sess, ok := sm.sessions[*in]; ok {
		return sess
	}
	return nil
}

func (sm *SessionManager) Delete(in *SessionID) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, *in)
}


func GenerateUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}