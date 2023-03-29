package main

import (
	"fmt"
	"sync"
)

type InMemorySessionStore struct {
	sessions map[string]*Session
	mu       sync.Mutex
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]*Session),
	}
}

func (s *InMemorySessionStore) getSession(id string) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}

	return session, nil
}

func (s *InMemorySessionStore) setSession(session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[session.ID] = session

	return nil
}

func (s *InMemorySessionStore) deleteSession(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[id]; !ok {
		return fmt.Errorf("session not found")
	}

	delete(s.sessions, id)

	return nil
}
