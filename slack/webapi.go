package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const SlackWebAPIEndpoint = "https://slack.com/api"

type WebAPIError struct {
	Method, URL, Response string
}

func (e *WebAPIError) Error() string {
	return e.Response
}

type WebAPI struct {
	c     *http.Client
	token string

	BaseURL string
}

func NewWebAPI(token string, httpClient *http.Client) *WebAPI {
	api := &WebAPI{
		token:   token,
		c:       httpClient,
		BaseURL: SlackWebAPIEndpoint,
	}

	if api.c == nil {
		api.c = http.DefaultClient
	}

	return api
}

func (api *WebAPI) SetChannelTopic(channelID, topic string) error {
	const method = "channels.setTopic"

	req, err := api.newRequestWithToken(method)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("channel_id", channelID)
	q.Add("topic", topic)
	req.URL.RawQuery = q.Encode()

	resp, err := api.c.Do(req)
	if err != nil {
		return wrapError(fmt.Errorf("failed to call method (%s)", err), method, req.URL)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return wrapError(fmt.Errorf("failed to read response body (%s)", err), method, req.URL)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return wrapError(fmt.Errorf("WebAPI responded with HTTP %d %q", resp.StatusCode, body), method, req.URL)
	}

	var v struct {
		Ok    bool   `json:"ok"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &v); err != nil {
		return wrapError(fmt.Errorf("failed to decode response body %q (%s)", body, err), method, req.URL)
	}

	if !v.Ok {
		return wrapError(fmt.Errorf("WebAPI returned error (%s)", v.Error), method, req.URL)
	}

	return nil
}

func (api *WebAPI) GetChannelTopic(channelID string) (string, error) {
	const method = "channels.getTopic"

	req, err := api.newRequestWithToken(method)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("channel_id", channelID)
	req.URL.RawQuery = q.Encode()

	resp, err := api.c.Do(req)
	if err != nil {
		return "", wrapError(fmt.Errorf("failed to call method (%s)", err), method, req.URL)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", wrapError(fmt.Errorf("failed to read response body (%s)", err), method, req.URL)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return "", wrapError(fmt.Errorf("WebAPI responded with HTTP %d %q", resp.StatusCode, body), method, req.URL)
	}

	var v struct {
		Ok    bool   `json:"ok"`
		Topic string `json:"topic"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &v); err != nil {
		return "", wrapError(fmt.Errorf("failed to decode response body %q (%s)", body, err), method, req.URL)
	}

	if !v.Ok {
		return "", wrapError(fmt.Errorf("WebAPI returned error (%s)", v.Error), method, req.URL)
	}

	return v.Topic, nil
}

func (api *WebAPI) newRequestWithToken(method string) (*http.Request, error) {
	req, err := http.NewRequest("GET", api.BaseURL+"/"+method, nil)
	if err != nil {
		return nil, wrapError(fmt.Errorf("failed to build WebAPI request (%s)", err), method, nil)
	}

	q := req.URL.Query()
	q.Add("token", api.token)
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func wrapError(err error, method string, url *url.URL) *WebAPIError {
	e := &WebAPIError{
		Method:   method,
		URL:      url.String(),
		Response: err.Error(),
	}

	if url != nil {
		url.Query().Set("token", "[hidden]")
		e.URL = url.String()
	}

	return e
}
