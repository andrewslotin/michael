package deploy_test

import (
	"time"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/slack"
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
	store.Set("key1", deploy.New(slack.User{ID: "1", Name: "Test User"}, "Deploy subject"))
	if d, ok := store.Get("key1"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), "1", d.User.ID)
		assert.Equal(suite.T(), "Test User", d.User.Name)
		assert.Equal(suite.T(), "Deploy subject", d.Subject)
		assert.WithinDuration(suite.T(), time.Now(), d.StartedAt, time.Second)
		assert.True(suite.T(), d.FinishedAt.IsZero())
	}

	// Override existing value
	store.Set("key1", deploy.New(slack.User{ID: "2", Name: "First User"}, "Updated deploy subject"))
	if d, ok := store.Get("key1"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), "2", d.User.ID)
		assert.Equal(suite.T(), "First User", d.User.Name)
		assert.Equal(suite.T(), "Updated deploy subject", d.Subject)
		assert.WithinDuration(suite.T(), time.Now(), d.StartedAt, time.Second)
		assert.True(suite.T(), d.FinishedAt.IsZero())
	}

	// Populate another key
	store.Set("key2", deploy.New(slack.User{ID: "3", Name: "Second User"}, "Another deploy"))
	if d, ok := store.Get("key2"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), "3", d.User.ID)
		assert.Equal(suite.T(), "Second User", d.User.Name)
		assert.Equal(suite.T(), "Another deploy", d.Subject)
		assert.WithinDuration(suite.T(), time.Now(), d.StartedAt, time.Second)
		assert.True(suite.T(), d.FinishedAt.IsZero())
	}

	// Check that another record wasn't changed
	if d, ok := store.Get("key1"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), "2", d.User.ID)
		assert.Equal(suite.T(), "First User", d.User.Name)
		assert.Equal(suite.T(), "Updated deploy subject", d.Subject)
		assert.WithinDuration(suite.T(), time.Now(), d.StartedAt, time.Second)
		assert.True(suite.T(), d.FinishedAt.IsZero())
	}
}

func (suite *StoreSuite) TestDel() {
	store, teardown, err := suite.Setup()
	if teardown != nil {
		defer teardown()
	}
	require.NoError(suite.T(), err)

	_, ok := store.Del("key1")
	assert.False(suite.T(), ok)

	store.Set("key1", deploy.New(slack.User{ID: "1", Name: "First User"}, "Deploy subject"))
	store.Set("key2", deploy.New(slack.User{ID: "2", Name: "Second User"}, "Another deploy"))

	_, ok = store.Get("key1")
	require.True(suite.T(), ok)
	_, ok = store.Get("key2")
	require.True(suite.T(), ok)

	if d, ok := store.Del("key1"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), "1", d.User.ID)
		assert.Equal(suite.T(), "First User", d.User.Name)
		assert.Equal(suite.T(), "Deploy subject", d.Subject)
		assert.WithinDuration(suite.T(), time.Now(), d.StartedAt, time.Second)
	}

	_, ok = store.Get("key1")
	assert.False(suite.T(), ok)
	_, ok = store.Get("key2")
	assert.True(suite.T(), ok)
}
