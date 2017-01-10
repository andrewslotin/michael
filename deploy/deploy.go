package deploy

import (
	"time"

	"github.com/andrewslotin/michael/slack"
)

type Deploy struct {
	User         slack.User
	Subject      string
	StartedAt    time.Time
	FinishedAt   time.Time
	PullRequests []PullRequestReference
}

func New(user slack.User, subject string) Deploy {
	return Deploy{
		User:         user,
		Subject:      subject,
		PullRequests: FindPullRequestReferences(subject),
	}
}

func (d *Deploy) Start() bool {
	if !d.StartedAt.IsZero() {
		return false
	}

	d.StartedAt = time.Now().UTC()
	return true
}

func (d *Deploy) Finish() {
	if !d.FinishedAt.IsZero() {
		return
	}

	d.FinishedAt = time.Now().UTC()
}

func (d1 Deploy) Equal(d2 Deploy) bool {
	return d1.User == d2.User &&
		d1.Subject == d2.Subject &&
		d1.StartedAt.Equal(d2.StartedAt) &&
		d1.FinishedAt.Equal(d2.FinishedAt)
}
