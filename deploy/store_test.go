package deploy_test

import (
	"time"

	"github.com/andrewslotin/michael/deploy"
	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
	Setup func() (store deploy.Store, teardownFn func(), err error)
}

func (suite *StoreSuite) TestGetSet() {
	store, teardown, err := suite.Setup()
	if teardown != nil {
		defer teardown()
	}
	require.NoError(suite.T(), err)

	_, ok := store.Get("key1")
	assert.False(suite.T(), ok)

	// Store a value
	channel1Deploy := deploy.Deploy{
		User:        slack.User{ID: "1", Name: "Test User"},
		Subject:     "Deploy subject a/b#1 and c/d#2 for @user1 and @user2",
		StartedAt:   time.Now().Add(-5 * time.Minute),
		FinishedAt:  time.Now().Add(-1 * time.Minute),
		Aborted:     true,
		AbortReason: "something went wrong",
		PullRequests: []deploy.PullRequestReference{
			{ID: "1", Repository: "a/b"},
			{ID: "2", Repository: "c/d"},
		},
		Subscribers: []deploy.UserReference{
			{Name: "user1"},
			{Name: "user2"},
		},
	}
	store.Set("key1", channel1Deploy)
	if d, ok := store.Get("key1"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), channel1Deploy, d)
	}

	// Populate another key
	channel2Deploy := deploy.Deploy{
		User:      slack.User{ID: "2", Name: "Second User"},
		Subject:   "Another deploy c/d#2 for @another_user",
		StartedAt: time.Now().Add(-4 * time.Minute),
		PullRequests: []deploy.PullRequestReference{
			{ID: "2", Repository: "c/d"},
		},
		Subscribers: []deploy.UserReference{
			{Name: "another_user"},
		},
	}
	store.Set("key2", channel2Deploy)
	if d, ok := store.Get("key2"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), channel2Deploy, d)
	}

	// Check that another record wasn't changed
	if d, ok := store.Get("key1"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), channel1Deploy, d)
	}
}

func (suite *StoreSuite) TestSet_Update() {
	store, teardown, err := suite.Setup()
	if teardown != nil {
		defer teardown()
	}
	require.NoError(suite.T(), err)

	channel1Deploy := deploy.New(slack.User{ID: "1", Name: "First User"}, "Deploy subject")
	channel1Deploy.Start()
	store.Set("key1", channel1Deploy)

	d, ok := store.Get("key1")
	require.True(suite.T(), ok)

	d.Subject = "Updated subject"
	d.User = slack.User{ID: "2", Name: "Updated User"}
	store.Set("key1", d)

	if updated, ok := store.Get("key1"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), updated, d)
	}
}
