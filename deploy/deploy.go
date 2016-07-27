package deploy

import (
	"time"

	"github.com/andrewslotin/slack-deploy-command/slack"
)

type Deploy struct {
	User       slack.User
	Subject    string
	StartedAt  time.Time
	FinishedAt time.Time
}

func New(user slack.User, subject string) Deploy {
	return Deploy{User: user, Subject: subject}
}

func (d *Deploy) Start() bool {
	if !d.StartedAt.IsZero() {
		return false
	}

	d.StartedAt = time.Now()
	return true
}

func (d *Deploy) Finish() {
	if !d.FinishedAt.IsZero() {
		return
	}

	d.FinishedAt = time.Now()
}
