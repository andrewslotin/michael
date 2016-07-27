package deploy_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/andrewslotin/slack-deploy-command/deploy"
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
