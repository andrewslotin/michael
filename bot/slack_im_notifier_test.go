package bot_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrewslotin/michael/bot"
	"github.com/andrewslotin/michael/deploy"
	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
)

func TestSlackIMNotifier_DeployCompleted(t *testing.T) {
	const webAPIToken = "xxxxx-token1"

	d := deploy.Deploy{
		User:    slack.User{ID: "U1", Name: "author"},
		Subject: "Deploy subject",
		InterestedUsers: []deploy.UserReference{
			{"recipient1"},
			{"nonExistingRecipient"},
			{"recipient2"},
		},
	}

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	var (
		requestNum struct{ UsersList, IMOpen int }
		receivers  []string
	)

	mux.HandleFunc("/users.list", func(w http.ResponseWriter, r *http.Request) {
		requestNum.UsersList++
		assert.Equal(t, webAPIToken, r.FormValue("token"))

		fmt.Fprint(w, `{"ok":true,"members":[{"id":"R1","name":"recipient1"},{"id":"R2","name":"recipient2"},{"id":"R3","name":"recipient3"}]}`)
	})
	mux.HandleFunc("/im.open", func(w http.ResponseWriter, r *http.Request) {
		requestNum.IMOpen++
		assert.Equal(t, webAPIToken, r.FormValue("token"))

		if userID := r.FormValue("user"); assert.NotEmpty(t, userID) {
			fmt.Fprintf(w, `{"ok":true,"channel":{"id":"DM%s"}}`, userID)
		} else {
			fmt.Fprint(w, `{"ok":false,"error":"user_not_found"}`)
		}
	})
	mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, webAPIToken, r.FormValue("token"))

		if channelID := r.FormValue("channel"); assert.NotEmpty(t, channelID) {
			receivers = append(receivers, channelID)
		}

		if msg := r.FormValue("text"); assert.NotEmpty(t, msg) {
			assert.Contains(t, msg, d.User.String())
			assert.Contains(t, msg, d.Subject)
		}

		fmt.Fprint(w, `{"ok":true}`)
	})

	api := slack.NewWebAPI(webAPIToken, nil)
	api.BaseURL = server.URL

	notifier := bot.NewSlackIMNotifier(api)
	notifier.DeployCompleted("", d)

	assert.Equal(t, 2, requestNum.UsersList) // nonExistingRecipient will not hit the cache
	assert.Equal(t, 2, requestNum.IMOpen)

	if assert.Len(t, receivers, 2) {
		assert.Contains(t, receivers, "DMR1")
		assert.Contains(t, receivers, "DMR2")
	}

	// Retry to check user list caching
	receivers = receivers[:0]
	notifier.DeployCompleted("", d)
	assert.Equal(t, 3, requestNum.UsersList) // +1 request because of nonExistingRecipient
	assert.Equal(t, 2, requestNum.IMOpen)    // no new channels are expected to be open

	if assert.Len(t, receivers, 2) {
		assert.Contains(t, receivers, "DMR1")
		assert.Contains(t, receivers, "DMR2")
	}
}
