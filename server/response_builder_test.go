package server_test

import (
	"errors"
	"testing"
	"time"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/server"
	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/stretchr/testify/assert"
)

func TestResponseBuilder_HelpMessage(t *testing.T) {
	b := server.NewResponseBuilder()
	response := b.HelpMessage()

	assert.Equal(t, slack.ResponseTypeEphemeral, response.ResponseType)
	for _, cmd := range [...]string{"&lt;subject&gt;", "done", "status", "help"} {
		assert.Contains(t, response.Text, "/deploy "+cmd+" ")
	}
}

func TestResponseBuilder_ErrorMessage(t *testing.T) {
	b := server.NewResponseBuilder()
	response := b.ErrorMessage("/command do", errors.New("error message"))

	assert.Equal(t, slack.ResponseTypeEphemeral, response.ResponseType)
	assert.Contains(t, response.Text, "/command do")
	assert.Contains(t, response.Text, "error message")
}

func TestResponseBuilder_NoRunningDeploysMessage(t *testing.T) {
	b := server.NewResponseBuilder()
	response := b.NoRunningDeploysMessage()

	assert.Equal(t, slack.ResponseTypeEphemeral, response.ResponseType)
	assert.NotEmpty(t, response.Text)
}

func TestResponseBuilder_DeployStatusMessage(t *testing.T) {
	d := deploy.Deploy{
		User:      slack.User{ID: "abc123", Name: "user1"},
		Subject:   "deploy subject",
		StartedAt: time.Now(),
	}

	b := server.NewResponseBuilder()
	response := b.DeployStatusMessage(d)

	assert.Equal(t, slack.ResponseTypeEphemeral, response.ResponseType)
	assert.Contains(t, response.Text, d.User.String())
	assert.Contains(t, response.Text, d.Subject)
	assert.Contains(t, response.Text, d.StartedAt.Format(time.RFC822))
}

func TestResponseBuilder_DeployInProgressMessage(t *testing.T) {
	d := deploy.Deploy{
		User:      slack.User{ID: "abc123", Name: "user1"},
		Subject:   "deploy subject",
		StartedAt: time.Now(),
	}

	b := server.NewResponseBuilder()
	response := b.DeployInProgressMessage(d)

	assert.Equal(t, slack.ResponseTypeEphemeral, response.ResponseType)
	assert.Contains(t, response.Text, d.User.String())
}

func TestResponseBuilder_DeployInterruptedAnnouncement(t *testing.T) {
	d := deploy.Deploy{
		User:      slack.User{ID: "abc123", Name: "user1"},
		Subject:   "deploy subject",
		StartedAt: time.Now(),
	}
	user := slack.User{ID: "xyz456", Name: "user2"}

	b := server.NewResponseBuilder()
	response := b.DeployInterruptedAnnouncement(d, user)

	assert.Equal(t, slack.ResponseTypeInChannel, response.ResponseType)
	assert.Contains(t, response.Text, user.String())
	assert.Contains(t, response.Text, d.User.String())
}

func TestResponseBuilder_DeployAnnouncement(t *testing.T) {
	user := slack.User{ID: "abc123", Name: "user1"}

	b := server.NewResponseBuilder()
	response := b.DeployAnnouncement(user, "deploy subject")

	assert.Equal(t, slack.ResponseTypeInChannel, response.ResponseType)
	assert.Contains(t, response.Text, user.String())
	assert.Contains(t, response.Text, "deploy subject")
}

func TestResponseBuilder_DeployDoneAnnouncement(t *testing.T) {
	user := slack.User{ID: "abc123", Name: "user1"}

	b := server.NewResponseBuilder()
	response := b.DeployDoneAnnouncement(user)

	assert.Equal(t, slack.ResponseTypeInChannel, response.ResponseType)
	assert.Contains(t, response.Text, user.String())
}
