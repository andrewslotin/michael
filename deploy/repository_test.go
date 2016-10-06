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
