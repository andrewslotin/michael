package slack_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/andrewslotin/slack-deploy-command/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebAPI_Call_WithParams(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	var requestNum int
	mux.HandleFunc("/methodName", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "value1", r.FormValue("key1"))
		assert.Equal(t, "value2", r.FormValue("key2"))
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))

		requestNum++
		w.Write([]byte(`{"ok":true}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	params := url.Values{}
	params.Add("key1", "value1")
	params.Add("key2", "value2")

	response, _, err := api.Call("methodName", params)
	require.NoError(t, err)
	require.Equal(t, 1, requestNum)
	assert.Equal(t, `{"ok":true}`, string(response))
}

func TestWebAPI_Call_WithoutParams(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	var requestNum int
	mux.HandleFunc("/methodName", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))

		requestNum++
		w.Write([]byte(`{"ok":true}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	response, _, err := api.Call("methodName", nil)
	require.NoError(t, err)
	require.Equal(t, 1, requestNum)
	assert.Equal(t, `{"ok":true}`, string(response))
}

func TestWebAPI_Call_WebAPIError(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	var requestNum int
	mux.HandleFunc("/methodName", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))

		requestNum++
		w.Write([]byte(`{"ok":false,"error":"an error occurred"}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	_, _, err := api.Call("methodName", url.Values{})
	require.Equal(t, 1, requestNum)

	if assert.Error(t, err) {
		var slackErr *slack.WebAPIError
		if assert.IsType(t, slackErr, err) {
			slackErr = err.(*slack.WebAPIError)
			assert.Equal(t, "methodName", slackErr.Method)
			assert.NotEmpty(t, slackErr.URL)
		}
	}
}

func TestWebAPI_Call_HTTPError(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	var requestNum int
	mux.HandleFunc("/methodName", func(w http.ResponseWriter, r *http.Request) {
		requestNum++
		http.Error(w, "An error occurred", http.StatusInternalServerError)
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	_, _, err := api.Call("methodName", url.Values{})
	require.Equal(t, 1, requestNum)

	if assert.Error(t, err) {
		var slackErr *slack.WebAPIError
		if assert.IsType(t, slackErr, err) {
			slackErr = err.(*slack.WebAPIError)
			assert.Equal(t, "methodName", slackErr.Method)
			assert.NotEmpty(t, slackErr.URL)
		}
	}
}

func TestWebAPI_ChannelsSetTopic(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	var requestNum int
	mux.HandleFunc("/channels.setTopic", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "CHANNELID1", r.FormValue("channel"))
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))
		assert.Equal(t, "Example topic", r.FormValue("topic"))

		requestNum++
		w.Write([]byte(`{"ok":true,"channel":{"topic":{"value":"Example topic"}}}`))
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
		assert.Equal(t, "CHANNELID1", r.FormValue("channel"))
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))
		assert.Equal(t, "Example topic", r.FormValue("topic"))

		requestNum++
		w.Write([]byte(`{"ok":false,"error":"channel not found"}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	err := api.SetChannelTopic("CHANNELID1", "Example topic")
	require.Equal(t, 1, requestNum)
	assert.Error(t, err)
}

func TestWebAPI_ChannelsGetTopic(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	var requestNum int
	mux.HandleFunc("/channels.info", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "CHANNELID1", r.FormValue("channel"))
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))

		requestNum++
		w.Write([]byte(`{"ok":true,"channel":{"topic":{"value":"Example topic"}}}`))
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
	mux.HandleFunc("/channels.info", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "CHANNELID1", r.FormValue("channel"))
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))

		requestNum++
		w.Write([]byte(`{"ok":false,"error":"channel not found"}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	_, err := api.GetChannelTopic("CHANNELID1")
	require.Equal(t, 1, requestNum)
	assert.Error(t, err)
}

func setup() (mux *http.ServeMux, baseURL string, teardownFn func()) {
	mux = http.NewServeMux()
	ts := httptest.NewServer(mux)

	return mux, ts.URL, ts.Close
}
