package stores

import (
	"sync"
	"time"

	"github.com/andrewslotin/slack-deploy-command/deploy"
)

type Memory struct {
	mu sync.RWMutex
	m  map[string]deploy.Deploy
}

func NewMemory() *Memory {
	return &Memory{
		m: make(map[string]deploy.Deploy),
	}
}

func (s *Memory) Get(key string) (deploy deploy.Deploy, ok bool) {
	s.mu.RLock()
	deploy, ok = s.m[key]
	s.mu.RUnlock()

	return deploy, ok
}

func (s *Memory) Set(key string, deploy deploy.Deploy) {
	deploy.StartedAt = time.Now()
	s.mu.Lock()
	s.m[key] = deploy
	s.mu.Unlock()
}

func (s *Memory) Del(key string) (deploy deploy.Deploy, ok bool) {
	s.mu.Lock()
	deploy, ok = s.m[key]
	delete(s.m, key)
	s.mu.Unlock()

	return deploy, ok
}
