package store

import (
	"sync"

	"github.com/andrewslotin/slack-deploy-command/slack"
)

type UserStore struct {
	mu sync.RWMutex
	m  map[string]slack.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		m: make(map[string]slack.User),
	}
}

func (s *UserStore) Get(key string) (user slack.User, ok bool) {
	s.mu.RLock()
	user, ok = s.m[key]
	s.mu.RUnlock()

	return user, ok
}

func (s *UserStore) Set(key string, user slack.User) {
	s.mu.Lock()
	s.m[key] = user
	s.mu.Unlock()
}
