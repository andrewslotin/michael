package deploy

import (
	"time"

	"github.com/andrewslotin/slack-deploy-command/slack"
)

type Deploy struct {
	User      slack.User
	Subject   string
	StartedAt time.Time
	EndAt     time.Time
}
