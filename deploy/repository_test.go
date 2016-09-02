package deploy_test

import (
	"fmt"
	"time"

	"github.com/andrewslotin/michael/deploy"
	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RepositorySuite struct {
	suite.Suite
	Setup func() (deploy.Repository, func(string, deploy.Deploy), func(), error)
}

func (suite *RepositorySuite) TestAll() {
	repo, storeSet, teardown, err := suite.Setup()
	if teardown != nil {
		defer teardown()
	}
	require.NoError(suite.T(), err)

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

		storeSet("key1", d)
		deploys = append(deploys, d)
	}

	allDeploys := repo.All(key)
	if assert.Len(suite.T(), allDeploys, len(deploys)) {
		for i, d := range allDeploys {
			assert.True(suite.T(), d.Equal(deploys[i]), "expected %+v, got %+v", d, deploys[i])
		}
	}
}

func (suite *RepositorySuite) TestSince_Multiple() {
	repo, storeSet, teardown, err := suite.Setup()
	if teardown != nil {
		defer teardown()
	}
	require.NoError(suite.T(), err)

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
		storeSet("key1", d)
	}

	deploys := repo.Since("key1", time.Now().Add(-58*time.Minute))
	assert.Len(suite.T(), deploys, 3)
	assert.True(suite.T(), history[1].Equal(deploys[0]))
	assert.True(suite.T(), history[2].Equal(deploys[1]))
	assert.True(suite.T(), history[3].Equal(deploys[2]))
}

func (suite *RepositorySuite) TestSince_EmptyHistory() {
	repo, storeSet, teardown, err := suite.Setup()
	if teardown != nil {
		defer teardown()
	}
	require.NoError(suite.T(), err)

	storeSet("key1", deploy.Deploy{
		StartedAt: time.Now(),
	})

	deploys := repo.Since("key2", time.Now().Add(-10*time.Minute))
	assert.Len(suite.T(), deploys, 0)
}

func (suite *RepositorySuite) TestSince_NoMatchingDeploys() {
	repo, storeSet, teardown, err := suite.Setup()
	if teardown != nil {
		defer teardown()
	}
	require.NoError(suite.T(), err)

	storeSet("key1", deploy.Deploy{
		StartedAt:  time.Now().Add(-20 * time.Minute),
		FinishedAt: time.Now().Add(-15 * time.Minute),
	})

	deploys := repo.Since("key1", time.Now().Add(-17*time.Minute))
	assert.Len(suite.T(), deploys, 0)
}
