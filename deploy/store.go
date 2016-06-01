package deploy

import (
	"sync"
	"time"

	"github.com/andrewslotin/slack-deploy-command/slack"
)

type Store struct {
	mu sync.RWMutex
	m  map[string]Deploy
}

func NewStore() *Store {
	return &Store{
		m: make(map[string]Deploy),
	}
}

func (s *Store) Get(key string) (deploy Deploy, ok bool) {
	s.mu.RLock()
	deploy, ok = s.m[key]
	s.mu.RUnlock()

	return deploy, ok
}

func (s *Store) Set(key string, user slack.User, subject string) {
	s.mu.Lock()
	s.m[key] = Deploy{
		User:      user,
		Subject:   subject,
		StartedAt: time.Now(),
	}
	s.mu.Unlock()
}

func (s *Store) Del(key string) (deploy Deploy, ok bool) {
	s.mu.Lock()
	deploy, ok = s.m[key]
	delete(s.m, key)
	s.mu.Unlock()

	return deploy, ok
}
