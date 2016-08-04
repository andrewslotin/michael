package deploy

import "sync"

type InMemoryStore struct {
	mu sync.RWMutex
	m  map[string][]Deploy
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		m: make(map[string][]Deploy),
	}
}

func (s *InMemoryStore) Get(key string) (d Deploy, ok bool) {
	s.mu.RLock()
	history, ok := s.m[key]
	ok = ok && len(history) > 0
	if ok {
		d = history[len(history)-1]
	}
	s.mu.RUnlock()

	return d, ok
}

func (s *InMemoryStore) Set(key string, d Deploy) {
	s.mu.Lock()
	history, ok := s.m[key]
	if ok && len(history) > 0 && history[len(history)-1].StartedAt == d.StartedAt { // Update last deploy
		history[len(history)-1] = d
	} else { // Add new deploy
		s.m[key] = append(s.m[key], d)
	}
	s.mu.Unlock()
}

func (s *InMemoryStore) All(key string) []Deploy {
	deploys := make([]Deploy, len(s.m[key]))
	s.mu.RLock()
	copy(deploys, s.m[key])
	s.mu.RUnlock()

	return deploys
}
