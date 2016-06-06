package github_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrewslotin/slack-deploy-command/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientGetPullRequest_WithToken(t *testing.T) {
	baseURL, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/user1/repo1/pulls/123", func(w http.ResponseWriter, req *http.Request) {
		if assert.Equal(t, "token abc123", req.Header.Get("Authorization")) {
			w.Write([]byte(`{"title":"Test PR","body":"PR description","number":123,"user":{"login":"author1"}}`))
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	})

	c := github.NewClient("abc123", nil)
	require.NotNil(t, c)

	c.BaseURL = baseURL
	pr, err := c.GetPullRequest("user1/repo1", "123")
	require.NoError(t, err)

	assert.Equal(t, "Test PR", pr.Title)
	assert.Equal(t, "PR description", pr.Body)
	assert.Equal(t, 123, pr.Number)
	assert.Equal(t, "author1", pr.Author.Name)
}

func TestClientGetPullRequest_NoToken(t *testing.T) {
	baseURL, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/user1/repo1/pulls/123", func(w http.ResponseWriter, req *http.Request) {
		if assert.Equal(t, "", req.Header.Get("Authorization")) {
			w.Write([]byte(`{"title":"Test PR","body":"PR description","number":123,"user":{"login":"author1"}}`))
		} else {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
	})

	c := github.NewClient("", nil)
	require.NotNil(t, c)

	c.BaseURL = baseURL
	pr, err := c.GetPullRequest("user1/repo1", "123")
	require.NoError(t, err)

	assert.Equal(t, "Test PR", pr.Title)
	assert.Equal(t, "PR description", pr.Body)
	assert.Equal(t, 123, pr.Number)
	assert.Equal(t, "author1", pr.Author.Name)
}

func TestClientGetPullRequest_NoSuchRepository(t *testing.T) {
	baseURL, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/user1/repo1/pulls/123", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	c := github.NewClient("", nil)
	require.NotNil(t, c)

	c.BaseURL = baseURL
	_, err := c.GetPullRequest("user1/repo1", "123")
	require.Error(t, err)
}

func setup() (baseURL string, mux *http.ServeMux, teardownFn func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	return server.URL, mux, server.Close
}
