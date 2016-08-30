package slack_test

import (
	"encoding/json"
	"testing"

	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttachment_MarshalJSON(t *testing.T) {
	attachment := sampleAttachment()

	data, err := json.Marshal(attachment)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &m), string(data))

	assert.Equal(t, attachment.AuthorName, m["author_name"])
	assert.Equal(t, attachment.Title, m["title"])
	assert.Equal(t, attachment.TitleLink, m["title_link"])
	assert.Equal(t, attachment.Text, m["text"])
	assert.Nil(t, m["mrkdwn_in"])
}

func TestAttachment_MarshalJSON_Markdown(t *testing.T) {
	attachment := sampleAttachment()
	attachment.Markdown = true

	data, err := json.Marshal(attachment)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &m), string(data))

	assert.Equal(t, attachment.AuthorName, m["author_name"])
	assert.Equal(t, attachment.Title, m["title"])
	assert.Equal(t, attachment.TitleLink, m["title_link"])
	assert.Equal(t, attachment.Text, m["text"])
	assert.Equal(t, []interface{}{"text"}, m["mrkdwn_in"])
}

func TestAttachment_UnmarshalJSON(t *testing.T) {
	attachment := sampleAttachment()
	attachment.Markdown = true

	data, err := json.Marshal(attachment)
	require.NoError(t, err)

	var unmarshaledAttachment slack.Attachment
	require.NoError(t, json.Unmarshal(data, &unmarshaledAttachment), string(data))
	assert.Equal(t, attachment, unmarshaledAttachment)
}

func TestAttachment_UnmarshalJSON_Pointer(t *testing.T) {
	attachment := sampleAttachment()
	attachment.Markdown = true

	data, err := json.Marshal(attachment)
	require.NoError(t, err)

	var unmarshaledAttachment slack.Attachment
	require.NoError(t, json.Unmarshal(data, &unmarshaledAttachment), string(data))

	if assert.NotNil(t, unmarshaledAttachment) {
		assert.Equal(t, attachment, unmarshaledAttachment)
	}
}

func sampleAttachment() slack.Attachment {
	return slack.Attachment{
		AuthorName: "test user",
		Title:      "test title",
		TitleLink:  "http://test.com/",
		Text:       "test body",
	}
}
