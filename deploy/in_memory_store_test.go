package deploy_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestInMemoryStore(t *testing.T) {
	suite.Run(t, &StoreSuite{Setup: func() (store deploy.Store, teardownFn func(), err error) {
		return deploy.NewInMemoryStore(), nil, nil
	}})
}

func TestBoltDBStore_All(t *testing.T) {
	store := deploy.NewInMemoryStore()

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
