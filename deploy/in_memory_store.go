package deploy

import "sync"

type InMemoryStore struct {
	mu sync.RWMutex
	m  map[string]Deploy
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		m: make(map[string]Deploy),
	}
}

func (s *InMemoryStore) Get(key string) (d Deploy, ok bool) {
	s.mu.RLock()
	d, ok = s.m[key]
	s.mu.RUnlock()

	return d, ok
}

func (s *InMemoryStore) Set(key string, d Deploy) {
	s.mu.Lock()
	s.m[key] = d
	s.mu.Unlock()
}

func (s *InMemoryStore) Del(key string) (d Deploy, ok bool) {
	s.mu.Lock()
	d, ok = s.m[key]
	delete(s.m, key)
	s.mu.Unlock()

	return d, ok
}
