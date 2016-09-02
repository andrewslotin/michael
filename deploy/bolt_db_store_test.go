package deploy_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/andrewslotin/michael/deploy"
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

func tempDBFilePath() (string, error) {
	fd, err := ioutil.TempFile(os.TempDir(), "doppelganger")
	if err != nil {
		return "", err
	}
	fd.Close()

	return fd.Name(), nil
}
