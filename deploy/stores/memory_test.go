package stores_test

import (
	"testing"
	"time"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/deploy/stores"
	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemory_GetSet(t *testing.T) {
	store := stores.NewMemory()

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

func TestMemory_Del(t *testing.T) {
	store := stores.NewMemory()

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

func TestMemory_Archive(t *testing.T) {
	store := stores.NewMemory()

	now := time.Now().Truncate(time.Second).UTC()
	deploy1 := deploy.Deploy{
		Subject: "sub1",
		User: slack.User{
			ID:   "1",
			Name: "a",
		},
		StartedAt: now,
		EndAt:     now.Add(time.Hour),
	}

	deploy2 := deploy.Deploy{
		Subject: "sub2",
		User: slack.User{
			ID:   "2",
			Name: "b",
		},
		StartedAt: now.Add(time.Hour),
		EndAt:     now.Add(time.Hour * 2),
	}

	deploy3 := deploy.Deploy{
		Subject: "sub3",
		User: slack.User{
			ID:   "3",
			Name: "c",
		},
		StartedAt: now.Add(time.Hour),
		EndAt:     now.Add(time.Hour * 2),
	}

	id, ok := store.Archive("channel1", deploy1)
	assert.True(t, ok)
	assert.Equal(t, id, uint64(1))

	id, ok = store.Archive("channel1", deploy2)
	assert.True(t, ok)
	assert.Equal(t, id, uint64(2))

	id, ok = store.Archive("channel2", deploy3)
	assert.True(t, ok)
	assert.Equal(t, id, uint64(1))

	deploys1, ok := store.FetchAllArchives("channel1")
	assert.True(t, ok)
	assert.Len(t, deploys1, 2)

	assert.Equal(t, "sub1", deploys1[0].Subject)
	assert.Equal(t, "1", deploys1[0].User.ID)
	assert.Equal(t, "a", deploys1[0].User.Name)
	assert.Equal(t, now, deploys1[0].StartedAt)
	assert.Equal(t, now.Add(time.Hour), deploys1[0].EndAt)

	assert.Equal(t, "sub2", deploys1[1].Subject)
	assert.Equal(t, "2", deploys1[1].User.ID)
	assert.Equal(t, "b", deploys1[1].User.Name)
	assert.Equal(t, now.Add(time.Hour), deploys1[1].StartedAt)
	assert.Equal(t, now.Add(time.Hour*2), deploys1[1].EndAt)

	deploys2, ok := store.FetchAllArchives("channel2")
	assert.True(t, ok)
	assert.Len(t, deploys2, 1)
}
