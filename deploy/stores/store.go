package stores

import (
	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/slack"
)

type Store interface {
	Get(channelID string) (deploy deploy.Deploy, ok bool)
	Set(channelID string, user slack.User, subject string)
	Del(channelID string) (deploy deploy.Deploy, ok bool)

	Archive(channelID string, deploy deploy.Deploy) bool
	FetchAllArchives(channelID string) (deploys []*deploy.Deploy, ok bool)
}
