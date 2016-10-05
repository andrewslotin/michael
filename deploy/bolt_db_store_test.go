package deploy_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/andrewslotin/michael/deploy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestBoltDBStore_AsStore(t *testing.T) {
	suite.Run(t, &StoreSuite{Setup: func() (store deploy.Store, teardownFn func(), err error) {
		path, err := tempDBFilePath()
		if err != nil {
			return nil, nil, err
		}

		teardownFn = func() { os.Remove(path) }

		store, err = deploy.NewBoltDBStore(path)
		if err != nil {
			return nil, teardownFn, err
		}

		return store, teardownFn, nil
	}})
}

func TestBoltDBStore_AsRepository(t *testing.T) {
	suite.Run(t, &RepositorySuite{Setup: func() (repo deploy.Repository, setFn func(string, deploy.Deploy), teardownFn func(), err error) {
		path, err := tempDBFilePath()
		if err != nil {
			return nil, nil, nil, err
		}

		teardownFn = func() { os.Remove(path) }

		r, err := deploy.NewBoltDBStore(path)
		if err != nil {
			return nil, nil, teardownFn, err
		}

		return r, r.Set, teardownFn, nil
	}})
}

func TestBoltDBStore_Since(t *testing.T) {
	path, err := tempDBFilePath()
	require.NoError(t, err)

	defer func() { os.Remove(path) }()

	store, err := deploy.NewBoltDBStore(path)
	require.NoError(t, err)

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

func tempDBFilePath() (string, error) {
	fd, err := ioutil.TempFile(os.TempDir(), "doppelganger")
	if err != nil {
		return "", err
	}
	fd.Close()

	return fd.Name(), nil
}
