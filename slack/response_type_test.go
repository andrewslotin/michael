package slack_test

import (
	"encoding/json"
	"testing"

	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSON(t *testing.T) {
	data, err := json.Marshal(slack.ResponseTypeEphemeral)
	require.NoError(t, err)
	assert.Equal(t, []byte(`"ephemeral"`), data)

	data, err = json.Marshal(slack.ResponseTypeInChannel)
	require.NoError(t, err)
	assert.Equal(t, []byte(`"in_channel"`), data)
}
