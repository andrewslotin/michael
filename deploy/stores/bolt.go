package stores

import (
	"fmt"
	"time"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/boltdb/bolt"
)

const (
	userIDKey    = "user.id"
	userNameKey  = "user.name"
	subjectKey   = "subject"
	startedAtKey = "started_at"
)

type BoltDB struct {
	db *bolt.DB
}

func NewBoltDB(path string) (*BoltDB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open db %s: %s", path, err)
	}

	return &BoltDB{db: db}, nil
}

func (s *BoltDB) Get(key string) (deploy deploy.Deploy, ok bool) {
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

func (s *BoltDB) Set(key string, d deploy.Deploy) {
	s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(key))
		if err != nil {
			return fmt.Errorf("failed to store deploy of %s by %s in channel %s: %s", d.Subject, d.User.Name, key, err)
		}

		d.StartedAt = time.Now()
		s.writeDeploy(d, b)

		return nil
	})
}

func (s *BoltDB) Del(key string) (d deploy.Deploy, ok bool) {
	s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(key))
		if b == nil {
			return nil
		}

		var err error
		if d, err = s.readDeploy(b); err != nil {
			return err
		}

		if err = tx.DeleteBucket([]byte(key)); err != nil {
			return err
		}

		ok = true

		return nil
	})

	return d, ok
}

func (*BoltDB) writeDeploy(deploy deploy.Deploy, b *bolt.Bucket) {
	b.Put([]byte([]byte(subjectKey)), []byte(deploy.Subject))
	b.Put([]byte([]byte(userIDKey)), []byte(deploy.User.ID))
	b.Put([]byte([]byte(userNameKey)), []byte(deploy.User.Name))
	b.Put([]byte([]byte(startedAtKey)), []byte(deploy.StartedAt.Format(time.RFC1123Z)))
}

func (*BoltDB) readDeploy(b *bolt.Bucket) (deploy deploy.Deploy, err error) {
	deploy.User = slack.User{
		ID:   string(b.Get([]byte(userIDKey))),
		Name: string(b.Get([]byte(userNameKey))),
	}
	deploy.Subject = string(b.Get([]byte(subjectKey)))

	if startedAt, err := time.Parse(time.RFC1123Z, string(b.Get([]byte(startedAtKey)))); err != nil {
		return deploy, fmt.Errorf("malformed started_at time for deploy of %s by %s: %s", deploy.Subject, deploy.User.Name, err)
	} else {
		deploy.StartedAt = startedAt
	}

	return deploy, nil
}
