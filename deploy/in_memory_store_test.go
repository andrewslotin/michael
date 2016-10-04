package deploy_test

import (
	"testing"
	"time"

	"github.com/andrewslotin/michael/deploy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestInMemoryStore_AsStore(t *testing.T) {
	suite.Run(t, &StoreSuite{Setup: func() (store deploy.Store, teardownFn func(), err error) {
		return deploy.NewInMemoryStore(), nil, nil
	}})
}

func TestInMemoryStore_AsRepository(t *testing.T) {
	suite.Run(t, &RepositorySuite{Setup: func() (repo deploy.Repository, setFn func(string, deploy.Deploy), teardownFn func(), err error) {
		r := deploy.NewInMemoryStore()
		return r, r.Set, nil, nil
	}})
}

func TestInMemoryStore_Since(t *testing.T) {
	store := deploy.NewInMemoryStore()

	history := []deploy.Deploy{
		deploy.Deploy{
			StartedAt:  time.Now().Add(-60 * time.Minute),
			FinishedAt: time.Now().Add(-55 * time.Minute),
		},
		deploy.Deploy{
			StartedAt:  time.Now().Add(-40 * time.Minute),
			FinishedAt: time.Now().Add(-35 * time.Minute),
		},
		deploy.Deploy{
			StartedAt:  time.Now().Add(-20 * time.Minute),
			FinishedAt: time.Now().Add(-15 * time.Minute),
		},
		deploy.Deploy{
			StartedAt: time.Now(),
		},
	}

	for _, d := range history {
		store.Set("key1", d)
	}

	t.Run("Multiple", func(t *testing.T) {
		deploys := store.Since("key1", time.Now().Add(-58*time.Minute))
		assert.Len(t, deploys, 3)
	})
	t.Run("Missing key", func(t *testing.T) {
		deploys := store.Since("key2", time.Now().Add(-58*time.Minute))
		assert.Len(t, deploys, 0)
	})
	t.Run("No deploys since", func(t *testing.T) {
		deploys := store.Since("key1", time.Now().Add(time.Minute))
		assert.Len(t, deploys, 0)
	})
}
