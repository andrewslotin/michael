package deploy_test

import (
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
	channel1Deploy := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Deploy subject a/b#1 for @user")
	channel1Deploy.Start()
	store.Set("key1", channel1Deploy)
	if d, ok := store.Get("key1"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), channel1Deploy, d)
		assert.Equal(suite.T(), channel1Deploy.PullRequests, d.PullRequests)
		assert.Equal(suite.T(), channel1Deploy.Subscribers, d.Subscribers)
	}

	// Populate another key
	channel2Deploy := deploy.New(slack.User{ID: "2", Name: "Second User"}, "Another deploy c/d#2 for @another_user")
	channel2Deploy.Start()
	store.Set("key2", channel2Deploy)
	if d, ok := store.Get("key2"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), channel2Deploy, d)
		assert.Equal(suite.T(), channel2Deploy.PullRequests, d.PullRequests)
		assert.Equal(suite.T(), channel2Deploy.Subscribers, d.Subscribers)
	}

	// Check that another record wasn't changed
	if d, ok := store.Get("key1"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), channel1Deploy, d)
		assert.Equal(suite.T(), channel1Deploy.PullRequests, d.PullRequests)
		assert.Equal(suite.T(), channel1Deploy.Subscribers, d.Subscribers)
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
