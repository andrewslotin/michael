package deploy_test

import (
	"testing"
	"time"

	"github.com/andrewslotin/michael/deploy"
	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
   Test objects
*/
type StoreMock struct {
	mock.Mock
}

func (m *StoreMock) Get(key string) (deploy.Deploy, bool) {
	args := m.Called(key)
	return args.Get(0).(deploy.Deploy), args.Bool(1)
}

func (m *StoreMock) Set(key string, d deploy.Deploy) {
	m.Called(key, d)
}

func (m *StoreMock) Del(key string) (d deploy.Deploy, ok bool) {
	args := m.Called(key)
	return args.Get(0).(deploy.Deploy), args.Bool(1)
}

/*
   Tests
*/

func TestChannelDeploys_Current(t *testing.T) {
	current := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Test subject")
	current.StartedAt = time.Now().Add(-5 * time.Minute)

	store := new(StoreMock)
	store.
		On("Get", "key1").Return(current, true).
		On("Get", "key2").Return(deploy.Deploy{}, false)

	repo := deploy.NewChannelDeploys(store)

	if d, ok := repo.Current("key1"); assert.True(t, ok) {
		assert.Equal(t, current, d)
	}

	_, ok := repo.Current("key2")
	assert.False(t, ok)

	store.AssertExpectations(t)
}

func TestChannelDeploys_Start(t *testing.T) {
	current := deploy.New(slack.User{ID: "2", Name: "Another User"}, "Active deploy")
	current.StartedAt = time.Now().Add(-2 * time.Minute)

	store := new(StoreMock)
	store.
		On("Get", "key1").Return(deploy.Deploy{}, false).
		On("Get", "key2").Return(current, true)
	store.
		On("Set", "key1", mock.AnythingOfType("deploy.Deploy")).Return()

	repo := deploy.NewChannelDeploys(store)

	d := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Test subject")
	if started, ok := repo.Start("key1", d); assert.True(t, ok) {
		assert.Equal(t, d.User, started.User)
		assert.Equal(t, d.Subject, started.Subject)
		assert.WithinDuration(t, time.Now(), started.StartedAt, time.Second)
	}

	d = deploy.New(slack.User{ID: "1", Name: "Test User"}, "Test subject")
	if started, ok := repo.Start("key2", d); assert.False(t, ok) {
		assert.Equal(t, current, started)
	}

	store.AssertExpectations(t)
}

func TestChannelDeploys_Start_UpdateCurrent(t *testing.T) {
	current := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Active deploy")
	current.StartedAt = time.Now().Add(-2 * time.Minute)

	store := new(StoreMock)
	store.
		On("Get", "key1").Return(current, true).Once(). // return running deploy
		On("Set", "key1", mock.AnythingOfType("deploy.Deploy")).Return(current, true).Once().
		On("Get", "key1").Return(deploy.Deploy{}, false). // current deploy has already been finished
		On("Set", "key1", mock.AnythingOfType("deploy.Deploy")).Return()

	repo := deploy.NewChannelDeploys(store)

	d := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Test subject")
	if started, ok := repo.Start("key1", d); assert.True(t, ok) {
		assert.Equal(t, d.User, started.User)
		assert.Equal(t, d.Subject, started.Subject)
		assert.WithinDuration(t, time.Now(), started.StartedAt, time.Second)
	}

	store.AssertExpectations(t)
}

func TestChannelDeploys_Finish(t *testing.T) {
	current := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Test subject")
	current.StartedAt = time.Now().Add(-2 * time.Second)

	store := new(StoreMock)
	store.
		On("Get", "key1").Return(current, true).
		On("Get", "key2").Return(deploy.Deploy{}, false).
		On("Set", "key1", mock.AnythingOfType("deploy.Deploy")).Return()

	repo := deploy.NewChannelDeploys(store)

	if d, ok := repo.Finish("key1"); assert.True(t, ok) {
		assert.Equal(t, current.User, d.User)
		assert.Equal(t, current.Subject, d.Subject)
		assert.WithinDuration(t, time.Now(), d.FinishedAt, time.Second)
	}

	_, ok := repo.Finish("key2")
	assert.False(t, ok)
}
