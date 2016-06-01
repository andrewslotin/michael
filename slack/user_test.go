package slack_test

import (
	"testing"

	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/stretchr/testify/assert"
)

func TestUserString(t *testing.T) {
	u := slack.User{
		ID:   "user1",
		Name: "Test User",
	}

	assert.Equal(t, "<@user1|Test User>", u.String())
}
