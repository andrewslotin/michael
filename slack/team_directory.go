package slack

import (
	"fmt"
	"sync"
)

type userLister interface {
	ListUsers() ([]User, error)
}

type TeamDirectory struct {
	api userLister

	mu    sync.RWMutex
	cache map[string]User
}

func NewTeamDirectory(api userLister) *TeamDirectory {
	return &TeamDirectory{api: api, cache: make(map[string]User)}
}

func (dir *TeamDirectory) Fetch(username string) (User, error) {
	dir.mu.RLock()
	user, ok := dir.cache[username]
	dir.mu.RUnlock()

	if ok {
		return user, nil
	}

	dir.mu.Lock()
	defer dir.mu.Unlock()

	user, ok = dir.cache[username]
	if !ok {
		users, err := dir.api.ListUsers()
		if err != nil {
			return User{}, fmt.Errorf("failed to fetch team users: %s", err)
		}

		dir.cache = make(map[string]User, len(users))
		for _, user := range users {
			dir.cache[user.Name] = user
		}

		user, ok = dir.cache[username]
	}

	if !ok {
		return User{}, NoSuchUserError{Name: username}
	}

	return user, nil
}
