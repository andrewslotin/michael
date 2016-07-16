package stores

import (
	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/slack"
)

type Store interface {
	Get(key string) (deploy deploy.Deploy, ok bool)
	Set(key string, user slack.User, subject string)
	Del(key string) (deploy deploy.Deploy, ok bool)
}
