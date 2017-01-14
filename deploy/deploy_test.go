package deploy_test

import (
	"testing"
	"time"

	"github.com/andrewslotin/michael/deploy"
	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
)

func TestNewDeploy(t *testing.T) {
	var (
		user    = slack.User{ID: "1", Name: "Test User"}
		subject = "Test deploy https://github.com/a/b/pull/123 and x/y#4 for @user1"
	)

	d := deploy.New(user, subject)
	assert.Equal(t, user, d.User)
	assert.Equal(t, subject, d.Subject)
	assert.Zero(t, d.StartedAt)
	assert.Zero(t, d.FinishedAt)

	if assert.Len(t, d.PullRequests, 2) {
		assert.Contains(t, d.PullRequests, deploy.PullRequestReference{ID: "123", Repository: "a/b"})
		assert.Contains(t, d.PullRequests, deploy.PullRequestReference{ID: "4", Repository: "x/y"})
	}

	if assert.Len(t, d.Subscribers, 1) {
		assert.Contains(t, d.Subscribers, deploy.UserReference{Name: "user1"})
	}
}

func TestDeploy_Start(t *testing.T) {
	d := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Test deploy")
	assert.True(t, d.Start())
	assert.WithinDuration(t, time.Now(), d.StartedAt, time.Second)

	startTime := d.StartedAt
	assert.False(t, d.Start())
	assert.Equal(t, startTime, d.StartedAt)
}

func TestDeploy_Finish(t *testing.T) {
	d := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Test deploy")
	d.Finish()
	assert.WithinDuration(t, time.Now(), d.FinishedAt, time.Second)
	assert.False(t, d.Aborted)

	finishedAt := d.FinishedAt
	d.Finish()
	assert.Equal(t, finishedAt, d.FinishedAt)
	assert.False(t, d.Aborted)
}

func TestDeploy_Abort_RunningDeploy(t *testing.T) {
	d := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Test deploy")
	d.Abort()
	assert.WithinDuration(t, time.Now(), d.FinishedAt, time.Second)
	assert.True(t, d.Aborted)

	finishedAt := d.FinishedAt
	d.Abort()
	assert.Equal(t, finishedAt, d.FinishedAt)
	assert.True(t, d.Aborted)
}

func TestDeploy_Abort_FinishedDeploy(t *testing.T) {
	d := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Test deploy")
	d.Finish()

	finishedAt := d.FinishedAt
	d.Abort()
	assert.Equal(t, finishedAt, d.FinishedAt)
	assert.False(t, d.Aborted)
}

func TestDeploy_Equal(t *testing.T) {
	d1 := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Test deploy")
	d1.StartedAt = time.Now().Add(-30 * time.Minute)
	d1.FinishedAt = time.Now().Add(-15 * time.Minute)

	var d2 deploy.Deploy
	assert.False(t, d1.Equal(d2))

	d2.User = d1.User
	assert.False(t, d1.Equal(d2))

	d2.Subject = d1.Subject
	assert.False(t, d1.Equal(d2))

	d2.StartedAt = d1.StartedAt
	assert.False(t, d1.Equal(d2))

	d2.FinishedAt = d1.FinishedAt
	assert.True(t, d1.Equal(d2))
}
