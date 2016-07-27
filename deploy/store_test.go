package deploy_test

import (
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
	channel1Deploy := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Deploy subject")
	channel1Deploy.Start()
	store.Set("key1", channel1Deploy)
	if d, ok := store.Get("key1"); assert.True(suite.T(), ok) {
		assert.Equal(suite.T(), channel1Deploy, d)
	}

	// Populate another key
	channel2Deploy := deploy.New(slack.User{ID: "2", Name: "Second User"}, "Another deploy")
	channel2Deploy.Start()
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
