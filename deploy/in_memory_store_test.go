package deploy_test

import (
	"testing"

	"github.com/andrewslotin/slack-deploy-command/deploy"
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
