package dashboard_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andrewslotin/michael/dashboard"
	"github.com/andrewslotin/michael/deploy"
	"github.com/andrewslotin/michael/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

/*      Test objects      */
type repoMock struct {
	mock.Mock
}

func (m repoMock) All(key string) []deploy.Deploy {
	return m.Called(key).Get(0).([]deploy.Deploy)
}

func (m repoMock) Since(key string, t time.Time) []deploy.Deploy {
	return m.Called(key, t).Get(0).([]deploy.Deploy)
}

/*          Tests         */
func TestDashboard_OneDeploy(t *testing.T) {
	url, mux, teardown := setup()
	defer teardown()

	d := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Test deploy")
	d.StartedAt, _ = time.Parse(time.RFC822, "04 Aug 16 09:28 CEST")
	d.FinishedAt, _ = time.Parse(time.RFC822, "04 Aug 16 09:38 CEST")

	var repo repoMock
	repo.On("All", "key1").Return([]deploy.Deploy{d})

	mux.Handle("/", dashboard.New(repo))

	response, err := http.Get(url + "/key1")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "text/plain", response.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	expected := "" +
		"Deploy history\n" +
		"--------------\n" +
		"\n" +
		"* Test User was deploying Test deploy since 04 Aug 16 09:28 CEST until 04 Aug 16 09:38 CEST"

	assert.Equal(t, expected, string(bytes.TrimSpace(body)))
}

func TestDashboard_MultipleDeploys(t *testing.T) {
	url, mux, teardown := setup()
	defer teardown()

	d1 := deploy.New(slack.User{ID: "1", Name: "Test User"}, "First deploy")
	d1.StartedAt, _ = time.Parse(time.RFC822, "04 Aug 16 09:28 CEST")
	d1.FinishedAt, _ = time.Parse(time.RFC822, "04 Aug 16 09:38 CEST")

	d2 := deploy.New(slack.User{ID: "1", Name: "Test User"}, "Second deploy")
	d2.StartedAt, _ = time.Parse(time.RFC822, "04 Aug 16 09:39 CEST")
	d2.FinishedAt, _ = time.Parse(time.RFC822, "04 Aug 16 09:40 CEST")

	d3 := deploy.New(slack.User{ID: "2", Name: "Another User"}, "Third deploy")
	d3.StartedAt, _ = time.Parse(time.RFC822, "04 Aug 16 09:50 CEST")

	var repo repoMock
	repo.On("All", "key1").Return([]deploy.Deploy{d1, d2, d3})

	mux.Handle("/", dashboard.New(repo))

	response, err := http.Get(url + "/key1")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "text/plain", response.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	expected := "" +
		"Deploy history\n" +
		"--------------\n" +
		"\n" +
		"* Test User was deploying First deploy since 04 Aug 16 09:28 CEST until 04 Aug 16 09:38 CEST\n" +
		"* Test User was deploying Second deploy since 04 Aug 16 09:39 CEST until 04 Aug 16 09:40 CEST\n" +
		"* Another User is currently deploying Third deploy since 04 Aug 16 09:50 CEST"

	assert.Equal(t, expected, string(bytes.TrimSpace(body)))
}

func TestDashboard_NoDeploys(t *testing.T) {
	url, mux, teardown := setup()
	defer teardown()

	var repo repoMock
	repo.On("All", "key1").Return([]deploy.Deploy(nil))

	mux.Handle("/", dashboard.New(repo))

	response, err := http.Get(url + "/key1")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "text/plain", response.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	expected := "" +
		"Deploy history\n" +
		"--------------\n" +
		"\n" +
		"No deploys in channel so far"

	assert.Equal(t, expected, string(bytes.TrimSpace(body)))
}

func TestDashboard_MissingChannelID(t *testing.T) {
	url, mux, teardown := setup()
	defer teardown()

	var repo repoMock
	repo.On("All", "key1").Return([]deploy.Deploy(nil))

	mux.Handle("/", dashboard.New(repo))

	response, err := http.Get(url + "/")
	require.NoError(t, err)
	response.Body.Close()

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestChannelIDFromRequest(t *testing.T) {
	examples := map[string]string{
		"/channel1":                        "channel1",
		"/channel2?key=val":                "channel2",
		"/channel3/hello":                  "channel3",
		"/channel4/hello/world/":           "channel4",
		"/channel5/hello/world/?key=val":   "channel5",
		"/channel6/notchannel.txt":         "channel6",
		"/channel7/notchannel.txt?key=val": "channel7",
		"/channel8.txt":                    "channel8",
		"/":                                "",
		"/?key=val":                        "",
	}

	for path, expectedID := range examples {
		req, err := http.NewRequest("GET", path, nil)
		if !assert.NoError(t, err) {
			continue
		}

		assert.Equal(t, expectedID, dashboard.ChannelIDFromRequest(req))
	}
}

func setup() (url string, mux *http.ServeMux, teardownFn func()) {
	mux = http.NewServeMux()
	srv := httptest.NewServer(mux)

	return srv.URL, mux, srv.Close
}
