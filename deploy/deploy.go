package deploy

import (
	"time"

	"github.com/andrewslotin/slack-deploy-command/slack"
)

type Deploy struct {
	User      slack.User
	Subject   string
	StartedAt time.Time
}

func New(user slack.User, subject string) Deploy {
	return Deploy{User: user, Subject: subject}
}
