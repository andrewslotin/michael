package bot_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andrewslotin/michael/bot"
	"github.com/andrewslotin/michael/deploy"
	"github.com/andrewslotin/michael/github"
	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
)

func TestResponseBuilder_HelpMessage(t *testing.T) {
	b := bot.NewResponseBuilder(github.NewClient("", nil))
	response := b.HelpMessage()

	assert.Equal(t, slack.ResponseTypeEphemeral, response.ResponseType)
	for _, cmd := range [...]string{"&lt;subject&gt;", "done", "status", "help"} {
		assert.Contains(t, response.Text, "/deploy "+cmd+" ")
	}
}

func TestResponseBuilder_ErrorMessage(t *testing.T) {
	b := bot.NewResponseBuilder(github.NewClient("", nil))
	response := b.ErrorMessage("/command do", errors.New("error message"))

	assert.Equal(t, slack.ResponseTypeEphemeral, response.ResponseType)
	assert.Contains(t, response.Text, "/command do")
	assert.Contains(t, response.Text, "error message")
}

func TestResponseBuilder_NoRunningDeploysMessage(t *testing.T) {
	b := bot.NewResponseBuilder(github.NewClient("", nil))
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

	b := bot.NewResponseBuilder(github.NewClient("", nil))
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

	b := bot.NewResponseBuilder(github.NewClient("", nil))
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

	b := bot.NewResponseBuilder(github.NewClient("", nil))
	response := b.DeployInterruptedAnnouncement(d, user)

	assert.Equal(t, slack.ResponseTypeInChannel, response.ResponseType)
	assert.Contains(t, response.Text, user.String())
	assert.Contains(t, response.Text, d.User.String())
}

func TestResponseBuilder_DeployAnnouncement(t *testing.T) {
	baseURL, mux, teardown := setupGitHubTestServer()
	defer teardown()

	mux.HandleFunc("/repos/user1/repo1/pulls/123", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`{"number":123,"title":"Hello","body":"PR description","html_url":"http://xyz.abc","user":{"login":"andrewslotin"}}`))
	})
	mux.HandleFunc("/repos/user2/repo2/pulls/234", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	githubClient := github.NewClient("", nil)
	githubClient.BaseURL = baseURL

	d := deploy.Deploy{
		User:    slack.User{ID: "abc123", Name: "user1"},
		Subject: "new feature",
		PullRequests: []deploy.PullRequestReference{
			{ID: "123", Repository: "user1/repo1"},
			{ID: "234", Repository: "user2/repo2"},
		},
		StartedAt: time.Now(),
	}

	b := bot.NewResponseBuilder(githubClient)
	response := b.DeployAnnouncement(d)

	assert.Equal(t, slack.ResponseTypeInChannel, response.ResponseType)
	assert.Contains(t, response.Text, d.User.String())
	assert.Contains(t, response.Text, d.Subject)

	if assert.Len(t, response.Attachments, 2) {
		assert.Equal(t, "PR #123: Hello", response.Attachments[0].Title)
		assert.Equal(t, "http://xyz.abc", response.Attachments[0].TitleLink)
		assert.Equal(t, "PR description", response.Attachments[0].Text)
		assert.Equal(t, "andrewslotin", response.Attachments[0].AuthorName)

		assert.Equal(t, "user2/repo2#234", response.Attachments[1].Title)
		assert.Equal(t, "https://github.com/user2/repo2/pulls/234", response.Attachments[1].TitleLink)
		assert.Empty(t, response.Attachments[1].Text)
		assert.Empty(t, response.Attachments[1].AuthorName)
	}
}

func TestResponseBuilder_DeployDoneAnnouncement(t *testing.T) {
	user := slack.User{ID: "abc123", Name: "user1"}

	b := bot.NewResponseBuilder(github.NewClient("", nil))
	response := b.DeployDoneAnnouncement(user)

	assert.Equal(t, slack.ResponseTypeInChannel, response.ResponseType)
	assert.Contains(t, response.Text, user.String())
}

func TestResponseBuilder_DeployHistoryLink_WithAuthToken(t *testing.T) {
	b := bot.NewResponseBuilder(github.NewClient("", nil))
	response := b.DeployHistoryLink("www.example.com:8080", "abc 123", "secret token")

	assert.Equal(t, slack.ResponseTypeEphemeral, response.ResponseType)
	assert.Contains(t, response.Text, "https://www.example.com:8080/abc%20123?token=secret+token")
}

func TestResponseBuilder_DeployHistoryLink_EmptyAuthToken(t *testing.T) {
	b := bot.NewResponseBuilder(github.NewClient("", nil))
	response := b.DeployHistoryLink("www.example.com:8080", "abc 123", "")

	assert.Equal(t, slack.ResponseTypeEphemeral, response.ResponseType)
	assert.Contains(t, response.Text, "https://www.example.com:8080/abc%20123")
}

func TestResponseBuilder_DeployHistoryLink_StandardPorts(t *testing.T) {
	standardPorts := [...]string{"80", "443"}

	b := bot.NewResponseBuilder(github.NewClient("", nil))

	for _, port := range standardPorts {
		response := b.DeployHistoryLink("www.example.com:"+port, "abc 123", "")
		assert.Contains(t, response.Text, "https://www.example.com/abc%20123", "port: %s", port)
	}
}

func setupGitHubTestServer() (baseURL string, mux *http.ServeMux, teardownFn func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	return server.URL, mux, server.Close
}
