package slack_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/andrewslotin/michael/slack"
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

func TestWebAPI_ListUsers(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	var requestNum int
	mux.HandleFunc("/users.list", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))

		requestNum++
		w.Write([]byte(`{"ok":true,"members":[{"id":"U1","name":"user1"},{"id":"U2","name":"user2"}]}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	users, err := api.ListUsers()
	require.NoError(t, err)
	require.Equal(t, 1, requestNum)

	if assert.Len(t, users, 2) {
		assert.Contains(t, users, slack.User{ID: "U1", Name: "user1"})
		assert.Contains(t, users, slack.User{ID: "U2", Name: "user2"})
	}
}

func TestWebAPI_ListUsers_ErrorHandling(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	var requestNum int
	mux.HandleFunc("/users.list", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))

		requestNum++
		w.Write([]byte(`{"ok":false,"error":"no users"`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	_, err := api.ListUsers()
	require.Equal(t, 1, requestNum)
	assert.Error(t, err)
}

func TestWebAPI_PostMessage_WithoutAttachments(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	message := slack.Message{Text: "Test message"}

	var requestNum int
	mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))
		assert.Equal(t, "channel1", r.FormValue("channel"))
		assert.Equal(t, "1", r.FormValue("link_names"))
		assert.Equal(t, message.Text, r.FormValue("text"))

		requestNum++
		w.Write([]byte(`{"ok":true}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	require.NoError(t, api.PostMessage("channel1", message))
	assert.Equal(t, 1, requestNum)
}

func TestWebAPI_PostMessage_WithAttachments(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	message := slack.Message{
		Text: "Test message",
		Attachments: []slack.Attachment{
			{Text: "attachment 1"},
			{Text: "attachment 2", Markdown: true},
		},
	}

	var requestNum int
	mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))
		assert.Equal(t, "channel1", r.FormValue("channel"))
		assert.Equal(t, "1", r.FormValue("link_names"))
		assert.Equal(t, message.Text, r.FormValue("text"))

		if encodedAttachments := r.FormValue("attachments"); assert.NotEmpty(t, encodedAttachments) {
			var attachments []slack.Attachment
			require.NoError(t, json.Unmarshal([]byte(encodedAttachments), &attachments))
			assert.Equal(t, message.Attachments, attachments)
		}

		requestNum++
		w.Write([]byte(`{"ok":true}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	require.NoError(t, api.PostMessage("channel1", message))
	assert.Equal(t, 1, requestNum)
}

func TestWebAPI_OpenIMChannel(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	user := slack.User{ID: "123", Name: "user1"}

	var requestNum int
	mux.HandleFunc("/im.open", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))
		assert.Equal(t, user.ID, r.FormValue("user"))

		requestNum++
		w.Write([]byte(`{"ok":true,"channel":{"id":"channel1"}}`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	channelID, err := api.OpenIMChannel(user)
	require.NoError(t, err)
	require.Equal(t, 1, requestNum)

	assert.Equal(t, "channel1", channelID)
}

func TestWebAPI_OpenIMChannel_ErrorHandling(t *testing.T) {
	mux, baseURL, teardown := setup()
	defer teardown()

	user := slack.User{ID: "123", Name: "user1"}

	var requestNum int
	mux.HandleFunc("/im.open", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "xxxx-token-12345", r.FormValue("token"))
		assert.Equal(t, user.ID, r.FormValue("user"))

		requestNum++
		w.Write([]byte(`{"ok":false,"error":"user_not_found"`))
	})

	api := slack.NewWebAPI("xxxx-token-12345", nil)
	api.BaseURL = baseURL

	_, err := api.OpenIMChannel(user)
	require.Equal(t, 1, requestNum)
	assert.Error(t, err)
}

func setup() (mux *http.ServeMux, baseURL string, teardownFn func()) {
	mux = http.NewServeMux()
	ts := httptest.NewServer(mux)

	return mux, ts.URL, ts.Close
}
