package deploy_test

import (
	"testing"
	"time"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/slack"
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
