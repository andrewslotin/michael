package bot_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andrewslotin/michael/bot"
	"github.com/andrewslotin/michael/deploy"
	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
)

const webAPIToken = "xxxx-token-abc12"

type SlackChannel struct {
	ID, Topic string
}

func TestSlackTopicManager_DeployStarted_NoRunningDeploy(t *testing.T) {
	baseURL, channel, teardown := setupSlackWebAPITestServer(t)
	defer teardown()

	channel.ID = "CHANNELID1"
	channel.Topic = "-=:poop:" + strings.Repeat(bot.DeployDoneEmotion, 3) + ":poop:=-"

	webAPI := slack.NewWebAPI(webAPIToken, nil)
	webAPI.BaseURL = baseURL

	mgr := bot.NewSlackTopicManager(webAPI)
	mgr.DeployStarted(channel.ID, deploy.Deploy{})

	assert.Equal(t, "-=:poop:"+strings.Repeat(bot.DeployInProgressEmotion, 3)+":poop:=-", channel.Topic)
}

func TestSlackTopicManager_DeployStarted_NoTopicNotification(t *testing.T) {
	baseURL, channel, teardown := setupSlackWebAPITestServer(t)
	defer teardown()

	channel.ID = "CHANNELID1"
	channel.Topic = "-=:poop:=-"

	webAPI := slack.NewWebAPI(webAPIToken, nil)
	webAPI.BaseURL = baseURL

	mgr := bot.NewSlackTopicManager(webAPI)
	mgr.DeployStarted(channel.ID, deploy.Deploy{})

	assert.Equal(t, "-=:poop:=-", channel.Topic)
}

func TestSlackTopicManager_DeployStarted_InProgress(t *testing.T) {
	baseURL, channel, teardown := setupSlackWebAPITestServer(t)
	defer teardown()

	channel.ID = "CHANNELID1"
	channel.Topic = "-=:poop:" + strings.Repeat(bot.DeployInProgressEmotion, 3) + ":poop:=-"

	webAPI := slack.NewWebAPI(webAPIToken, nil)
	webAPI.BaseURL = baseURL

	mgr := bot.NewSlackTopicManager(webAPI)
	mgr.DeployStarted(channel.ID, deploy.Deploy{})

	assert.Equal(t, "-=:poop:"+strings.Repeat(bot.DeployInProgressEmotion, 3)+":poop:=-", channel.Topic)
}

func TestSlackTopicManager_DeployCompleted_InProgress(t *testing.T) {
	baseURL, channel, teardown := setupSlackWebAPITestServer(t)
	defer teardown()

	channel.ID = "CHANNELID1"
	channel.Topic = "-=:poop:" + strings.Repeat(bot.DeployInProgressEmotion, 3) + ":poop:=-"

	webAPI := slack.NewWebAPI(webAPIToken, nil)
	webAPI.BaseURL = baseURL

	mgr := bot.NewSlackTopicManager(webAPI)
	mgr.DeployCompleted(channel.ID, deploy.Deploy{})

	assert.Equal(t, "-=:poop:"+strings.Repeat(bot.DeployDoneEmotion, 3)+":poop:=-", channel.Topic)
}

func TestSlackTopicManager_DeployCompleted_NoTopicNotification(t *testing.T) {
	baseURL, channel, teardown := setupSlackWebAPITestServer(t)
	defer teardown()

	channel.ID = "CHANNELID1"
	channel.Topic = "-=:poop:=-"

	webAPI := slack.NewWebAPI(webAPIToken, nil)
	webAPI.BaseURL = baseURL

	mgr := bot.NewSlackTopicManager(webAPI)
	mgr.DeployCompleted(channel.ID, deploy.Deploy{})

	assert.Equal(t, "-=:poop:=-", channel.Topic)
}

func TestSlackTopicManager_DeployCompleted_NoRunningDeploy(t *testing.T) {
	baseURL, channel, teardown := setupSlackWebAPITestServer(t)
	defer teardown()

	channel.ID = "CHANNELID1"
	channel.Topic = "-=:poop:" + strings.Repeat(bot.DeployDoneEmotion, 3) + ":poop:=-"

	webAPI := slack.NewWebAPI(webAPIToken, nil)
	webAPI.BaseURL = baseURL

	mgr := bot.NewSlackTopicManager(webAPI)
	mgr.DeployCompleted(channel.ID, deploy.Deploy{})

	assert.Equal(t, "-=:poop:"+strings.Repeat(bot.DeployDoneEmotion, 3)+":poop:=-", channel.Topic)
}

func setupSlackWebAPITestServer(t *testing.T) (baseURL string, channel *SlackChannel, teardownFn func()) {
	channel = &SlackChannel{}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/channels.info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"ok":true,"channel":{"topic":{"value":"%s"}}}`, channel.Topic)
	})

	mux.HandleFunc("/channels.setTopic", func(w http.ResponseWriter, r *http.Request) {
		if token := r.FormValue("token"); !assert.Equal(t, webAPIToken, token) {
			fmt.Fprintf(w, `{"ok":false,"error":"wrong token %q"}`, token)
			return
		}

		if ch := r.FormValue("channel"); !assert.Equal(t, channel.ID, ch) {
			fmt.Fprintf(w, `{"ok":false,"error":"wrong channel %q"}`, ch)
			return
		}

		channel.Topic = r.FormValue("topic")
		fmt.Fprint(w, `{"ok":true}`)
	})

	return server.URL, channel, server.Close
}
