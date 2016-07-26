package deploy

import (
	"fmt"
	"time"

	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/boltdb/bolt"
)

const (
	userIDKey     = "user.id"
	userNameKey   = "user.name"
	subjectKey    = "subject"
	startedAtKey  = "started_at"
	finishedAtKey = "finished_at"
)

type BoltDBStore struct {
	db *bolt.DB
}

func NewBoltDBStore(path string) (*BoltDBStore, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open db %s: %s", path, err)
	}

	return &BoltDBStore{db: db}, nil
}

func (s *BoltDBStore) Get(key string) (deploy Deploy, ok bool) {
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(key))
		if b == nil {
			return nil
		}

		var err error
		if deploy, err = s.readDeploy(b); err != nil {
			return err
		}

		ok = true

		return nil
	})

	return deploy, ok
}

func (s *BoltDBStore) Set(key string, d Deploy) {
	s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(key))
		if err != nil {
			return fmt.Errorf("failed to store deploy of %s by %s in channel %s: %s", d.Subject, d.User.Name, key, err)
		}

		s.writeDeploy(d, b)

		return nil
	})
}

func (*BoltDBStore) writeDeploy(deploy Deploy, b *bolt.Bucket) {
	b.Put([]byte(subjectKey), []byte(deploy.Subject))
	b.Put([]byte(userIDKey), []byte(deploy.User.ID))
	b.Put([]byte(userNameKey), []byte(deploy.User.Name))
	b.Put([]byte(startedAtKey), []byte(deploy.StartedAt.Format(time.RFC3339Nano)))

	if !deploy.FinishedAt.IsZero() {
		b.Put([]byte(finishedAtKey), []byte(deploy.FinishedAt.Format(time.RFC3339Nano)))
	}
}

func (*BoltDBStore) readDeploy(b *bolt.Bucket) (deploy Deploy, err error) {
	deploy.User = slack.User{
		ID:   string(b.Get([]byte(userIDKey))),
		Name: string(b.Get([]byte(userNameKey))),
	}
	deploy.Subject = string(b.Get([]byte(subjectKey)))

	if startedAt, err := time.Parse(time.RFC3339Nano, string(b.Get([]byte(startedAtKey)))); err != nil {
		return deploy, fmt.Errorf("malformed started_at time for deploy of %s by %s: %s", deploy.Subject, deploy.User.Name, err)
	} else {
		deploy.StartedAt = startedAt
	}

	if value := b.Get([]byte(finishedAtKey)); value != nil {
		if finishedAt, err := time.Parse(time.RFC3339Nano, string(value)); err != nil {
			return deploy, fmt.Errorf("malformed finished_at time for deploy of %s by %s: %s", deploy.Subject, deploy.User.Name, err)
		} else {
			deploy.FinishedAt = finishedAt
		}
	}

	return deploy, nil
}
