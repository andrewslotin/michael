package deploy_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestBoltDBStore(t *testing.T) {
	suite.Run(t, &StoreSuite{Setup: func() (store deploy.Store, teardownFn func(), err error) {
		fd, err := ioutil.TempFile(os.TempDir(), "doppelganger")
		if err != nil {
			return nil, nil, err
		}
		fd.Close()

		store, err = deploy.NewBoltDBStore(fd.Name())
		if err != nil {
			return nil, func() { os.Remove(fd.Name()) }, err
		}

		return store, func() { os.Remove(fd.Name()) }, nil
	}})
}

func TestBoltDBStore_All(t *testing.T) {
	fd, err := ioutil.TempFile(os.TempDir(), "doppelganger")
	require.NoError(t, err)
	fd.Close()
	defer os.Remove(fd.Name())

	store, err := deploy.NewBoltDBStore(fd.Name())
	require.NoError(t, err)

	now := time.Now()
	key := "key1"
	user := slack.User{ID: "1", Name: "User 1"}

	var deploys []deploy.Deploy
	for delta := -10 * time.Minute; delta < 0; delta += time.Minute {
		d := deploy.New(user, fmt.Sprintf("Deploy from %s ago", delta))
		d.StartedAt = now.Add(delta)
		if delta+time.Minute < 0 {
			d.FinishedAt = now.Add(delta + time.Minute)
		}

		store.Set("key1", d)
		deploys = append(deploys, d)
	}

	allDeploys := store.All(key)
	if assert.Len(t, allDeploys, len(deploys)) {
		for i, d := range allDeploys {
			assert.Equal(t, deploys[i].User, d.User)
			assert.Equal(t, deploys[i].Subject, d.Subject)
			assert.True(t, deploys[i].StartedAt.Equal(d.StartedAt), "expected %s, got %s", deploys[i].StartedAt, d.StartedAt)
			assert.True(t, deploys[i].FinishedAt.Equal(d.FinishedAt), "expected %s, got %s", deploys[i].FinishedAt, d.FinishedAt)
		}
	}
}
