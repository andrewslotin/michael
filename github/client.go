package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	BaseURL string // to use in tests

	authHeader string
	client     *http.Client
}

func NewClient(token string, client *http.Client) *Client {
	c := &Client{
		BaseURL: "https://api.github.com",
	}

	if token != "" {
		c.authHeader = "token " + token
	}

	if client != nil {
		c.client = client
	} else {
		c.client = http.DefaultClient
	}

	return c
}

func (c *Client) GetPullRequest(repo string, number string) (pr PullRequest, err error) {
	url := c.BaseURL + "/repos/" + repo + "/pulls/" + number
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return pr, fmt.Errorf("failed to build a request to %s (%s)", url, err)
	}

	if c.authHeader != "" {
		req.Header.Set("Authorization", c.authHeader)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return pr, fmt.Errorf("request to %s failed (%s)", url, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pr, fmt.Errorf("failed to read %s response body (%s)", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return pr, fmt.Errorf("got HTTP %d response from %s: %q", resp.StatusCode, url, body)
	}

	if err = json.Unmarshal(body, &pr); err != nil {
		return pr, err
	}

	return pr, nil
}
