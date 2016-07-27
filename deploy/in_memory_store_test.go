package deploy_test

import (
	"testing"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/stretchr/testify/suite"
)

func TestInMemoryStore(t *testing.T) {
	suite.Run(t, &StoreSuite{Setup: func() (store deploy.Store, teardownFn func(), err error) {
		return deploy.NewInMemoryStore(), nil, nil
	}})
}
