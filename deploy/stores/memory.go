package stores

import (
	"sync"
	"time"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/slack"
)

type Memory struct {
	mu sync.RWMutex
	m  map[string]deploy.Deploy
	a  map[string][]deploy.Deploy
}

func NewMemory() *Memory {
	return &Memory{
		m: make(map[string]deploy.Deploy),
		a: make(map[string][]deploy.Deploy),
	}
}

func (s *Memory) Get(channelID string) (deploy deploy.Deploy, ok bool) {
	s.mu.RLock()
	deploy, ok = s.m[channelID]
	s.mu.RUnlock()

	return deploy, ok
}

func (s *Memory) Set(channelID string, user slack.User, subject string) {
	s.mu.Lock()
	s.m[channelID] = deploy.Deploy{
		User:      user,
		Subject:   subject,
		StartedAt: time.Now(),
	}
	s.mu.Unlock()
}

func (s *Memory) Del(channelID string) (deploy deploy.Deploy, ok bool) {
	s.mu.Lock()
	deploy, ok = s.m[channelID]
	delete(s.m, channelID)
	s.mu.Unlock()

	return deploy, ok
}

func (s *Memory) Archive(channelID string, deploy deploy.Deploy) (id uint64, ok bool) {
	s.mu.Lock()
	s.a[channelID] = append(s.a[channelID], deploy)
	s.mu.Unlock()

	return uint64(len(s.a[channelID])), true
}

func (s *Memory) FetchAllArchives(channelID string) (deploys []*deploy.Deploy, ok bool) {
	deploys = []*deploy.Deploy{}
	s.mu.Lock()
	for _, d := range s.a[channelID] {
		newDeploy := d
		deploys = append(deploys, &newDeploy)
	}
	s.mu.Unlock()

	return deploys, true
}
