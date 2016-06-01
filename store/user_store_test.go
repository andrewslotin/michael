package store_test

import (
	"testing"

	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/andrewslotin/slack-deploy-command/store"
	"github.com/stretchr/testify/assert"
)

func TestStore_GetSet(t *testing.T) {
	store := store.NewUserStore()

	user, ok := store.Get("key1")
	assert.False(t, ok)

	// Store an empty value
	store.Set("key1", slack.User{})
	user, ok = store.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, slack.User{}, user)

	// Store non-empty value
	store.Set("key1", slack.User{ID: "1", Name: "Test User"})
	user, ok = store.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "Test User", user.Name)

	// Override existing value
	store.Set("key1", slack.User{ID: "2", Name: "First User"})
	user, ok = store.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "2", user.ID)
	assert.Equal(t, "First User", user.Name)

	// Populate another key
	store.Set("key2", slack.User{ID: "3", Name: "Second User"})
	user, ok = store.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "3", user.ID)
	assert.Equal(t, "Second User", user.Name)

	user, ok = store.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "2", user.ID)
	assert.Equal(t, "First User", user.Name)
}
