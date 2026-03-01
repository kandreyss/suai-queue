package sessions

import "sync"

type SessionStore struct {
	mu sync.RWMutex
	m  map[int64]*UserSession
}

func NewSessionStore() *SessionStore {
	return &SessionStore{m: make(map[int64]*UserSession)}
}

func (s *SessionStore) Get(userID int64) (*UserSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	us, ok := s.m[userID]
	return us, ok
}

func (s *SessionStore) Set(userID int64, us *UserSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[userID] = us
}

func (s *SessionStore) Delete(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, userID)
}
