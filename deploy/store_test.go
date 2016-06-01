package deploy_test

import (
	"testing"
	"time"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_GetSet(t *testing.T) {
	store := deploy.NewStore()

	d, ok := store.Get("key1")
	assert.False(t, ok)

	// Store a value
	store.Set("key1", slack.User{ID: "1", Name: "Test User"}, "Deploy subject")
	d, ok = store.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "1", d.User.ID)
	assert.Equal(t, "Test User", d.User.Name)
	assert.Equal(t, "Deploy subject", d.Subject)
	assert.WithinDuration(t, time.Now(), d.StartedAt, time.Second)

	// Override existing value
	store.Set("key1", slack.User{ID: "2", Name: "First User"}, "Updated deploy subject")
	d, ok = store.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "2", d.User.ID)
	assert.Equal(t, "First User", d.User.Name)
	assert.Equal(t, "Updated deploy subject", d.Subject)
	assert.WithinDuration(t, time.Now(), d.StartedAt, time.Second)

	// Populate another key
	store.Set("key2", slack.User{ID: "3", Name: "Second User"}, "Another deploy")
	d, ok = store.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "3", d.User.ID)
	assert.Equal(t, "Second User", d.User.Name)
	assert.Equal(t, "Another deploy", d.Subject)
	assert.WithinDuration(t, time.Now(), d.StartedAt, time.Second)

	d, ok = store.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "2", d.User.ID)
	assert.Equal(t, "First User", d.User.Name)
	assert.Equal(t, "Updated deploy subject", d.Subject)
	assert.WithinDuration(t, time.Now(), d.StartedAt, time.Second)
}

func TestStore_Del(t *testing.T) {
	store := deploy.NewStore()

	_, ok := store.Del("key1")
	assert.False(t, ok)

	store.Set("key1", slack.User{ID: "1", Name: "First User"}, "Deploy subject")
	store.Set("key2", slack.User{ID: "2", Name: "Second User"}, "Another deploy")

	_, ok = store.Get("key1")
	require.True(t, ok)
	_, ok = store.Get("key2")
	require.True(t, ok)

	d, ok := store.Del("key1")
	assert.True(t, ok)
	assert.Equal(t, "1", d.User.ID)
	assert.Equal(t, "First User", d.User.Name)
	assert.Equal(t, "Deploy subject", d.Subject)
	assert.WithinDuration(t, time.Now(), d.StartedAt, time.Second)

	_, ok = store.Get("key1")
	assert.False(t, ok)
	_, ok = store.Get("key2")
	assert.True(t, ok)
}
