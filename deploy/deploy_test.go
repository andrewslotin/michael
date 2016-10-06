package deploy_test

import (
	"testing"
	"time"

	"github.com/andrewslotin/michael/deploy"
	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
)

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

	finishedAt := d.FinishedAt
	d.Finish()
	assert.Equal(t, finishedAt, d.FinishedAt)
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
