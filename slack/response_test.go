package slack_test

import (
	"encoding/json"
	"testing"

	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEphemeralResponse(t *testing.T) {
	r := slack.NewEphemeralResponse("test message")
	data, err := json.Marshal(r)
	require.NoError(t, err)

	var m map[string]string
	require.NoError(t, json.Unmarshal(data, &m))
	assert.Empty(t, m["response_type"])
	assert.Equal(t, "test message", m["text"])
	assert.Len(t, m, 1)
}

func TestNewInChannelResponse(t *testing.T) {
	r := slack.NewInChannelResponse("message")
	data, err := json.Marshal(r)
	require.NoError(t, err)

	var m map[string]string
	require.NoError(t, json.Unmarshal(data, &m))
	assert.Equal(t, "in_channel", m["response_type"])
	assert.Equal(t, "message", m["text"])
	assert.Len(t, m, 2)
}
