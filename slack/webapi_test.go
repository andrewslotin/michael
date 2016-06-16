package slack_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebAPI_ChannelsSetTopic(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	var requestNum int
	mux.HandleFunc("/channels.setTopic", func(w http.ResponseWriter, r *http.Request) {
		requestNum++

		assert.Equal(t, "CHANNELID1", r.FormValue("channel_id"))
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))
		assert.Equal(t, "Example topic", r.FormValue("topic"))

		w.Write([]byte(`{"ok":true,"topic":"Example topic"}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	require.NoError(t, api.SetChannelTopic("CHANNELID1", "Example topic"))
	require.Equal(t, 1, requestNum)
}

func TestWebAPI_ChannelsSetTopic_ErrorHandling(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	var requestNum int
	mux.HandleFunc("/channels.setTopic", func(w http.ResponseWriter, r *http.Request) {
		requestNum++

		assert.Equal(t, "CHANNELID1", r.FormValue("channel_id"))
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))
		assert.Equal(t, "Example topic", r.FormValue("topic"))

		w.Write([]byte(`{"ok":false,"error":"channel not found"}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	err := api.SetChannelTopic("CHANNELID1", "Example topic")
	require.Equal(t, 1, requestNum)

	if assert.EqualError(t, err, "WebAPI returned error (channel not found)") {
		var slackErr *slack.WebAPIError
		if assert.IsType(t, slackErr, err) {
			slackErr = err.(*slack.WebAPIError)
			assert.Equal(t, "channels.setTopic", slackErr.Method)
			assert.NotEmpty(t, slackErr.URL)
		}
	}
}

func TestWebAPI_ChannelsSetTopic_HTTPFailure(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/channels.setTopic", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "An error occurred", http.StatusInternalServerError)
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	err := api.SetChannelTopic("CHANNELID1", "Example topic")
	if assert.EqualError(t, err, `WebAPI responded with HTTP 500 "An error occurred\n"`) {
		var slackErr *slack.WebAPIError
		if assert.IsType(t, slackErr, err) {
			slackErr = err.(*slack.WebAPIError)
			assert.Equal(t, "channels.setTopic", slackErr.Method)
			assert.NotEmpty(t, slackErr.URL)
		}
	}
}

func TestWebAPI_ChannelsGetTopic(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	var requestNum int
	mux.HandleFunc("/channels.getTopic", func(w http.ResponseWriter, r *http.Request) {
		requestNum++

		assert.Equal(t, "CHANNELID1", r.FormValue("channel_id"))
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))

		w.Write([]byte(`{"ok":true,"topic":"Example topic"}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	topic, err := api.GetChannelTopic("CHANNELID1")
	require.NoError(t, err)
	require.Equal(t, 1, requestNum)

	assert.Equal(t, "Example topic", topic)
}

func TestWebAPI_ChannelsGetTopic_ErrorHandling(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	var requestNum int
	mux.HandleFunc("/channels.getTopic", func(w http.ResponseWriter, r *http.Request) {
		requestNum++

		assert.Equal(t, "CHANNELID1", r.FormValue("channel_id"))
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))

		w.Write([]byte(`{"ok":false,"error":"channel not found"}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	_, err := api.GetChannelTopic("CHANNELID1")
	require.Equal(t, 1, requestNum)

	if assert.EqualError(t, err, "WebAPI returned error (channel not found)") {
		var slackErr *slack.WebAPIError
		if assert.IsType(t, slackErr, err) {
			slackErr = err.(*slack.WebAPIError)
			assert.Equal(t, "channels.getTopic", slackErr.Method)
			assert.NotEmpty(t, slackErr.URL)
		}
	}
}

func TestWebAPI_ChannelsGetTopic_HTTPFailure(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/channels.getTopic", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "An error occurred", http.StatusInternalServerError)
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	_, err := api.GetChannelTopic("CHANNELID1")
	if assert.EqualError(t, err, `WebAPI responded with HTTP 500 "An error occurred\n"`) {
		var slackErr *slack.WebAPIError
		if assert.IsType(t, slackErr, err) {
			slackErr = err.(*slack.WebAPIError)
			assert.Equal(t, "channels.getTopic", slackErr.Method)
			assert.NotEmpty(t, slackErr.URL)
		}
	}
}

func setup() (mux *http.ServeMux, baseURL string, teardownFn func()) {
	mux = http.NewServeMux()
	ts := httptest.NewServer(mux)

	return mux, ts.URL, ts.Close
}
