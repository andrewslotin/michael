package stores

import "github.com/andrewslotin/slack-deploy-command/deploy"

type Store interface {
	Get(key string) (deploy deploy.Deploy, ok bool)
	Set(key string, deploy deploy.Deploy)
	Del(key string) (deploy deploy.Deploy, ok bool)
}
