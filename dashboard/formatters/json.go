package formatters

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/andrewslotin/michael/deploy"
)

var (
	JSON jsonFormatter
)

type jsonPresenter struct {
	Author     string    `json:"author"`
	Subject    string    `json:"subject"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
	Aborted    bool      `json:"aborted,omitempty"`
	Reason     string    `json:"reason,omitempty"`
}

type jsonFormatter struct{}

func (jsonFormatter) RespondWithHistory(w http.ResponseWriter, history []deploy.Deploy) error {
	w.Header().Set("Content-Type", "application/json")

	v := make([]jsonPresenter, len(history))
	for i, d := range history {
		v[i].Author = d.User.Name
		v[i].Subject = d.Subject
		v[i].StartedAt = d.StartedAt
		v[i].FinishedAt = d.FinishedAt
		v[i].Aborted = d.Aborted
		v[i].Reason = d.AbortReason
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (jsonFormatter) RespondWithError(w http.ResponseWriter, err error, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	data, err := json.Marshal(struct {
		Error string `json:"error"`
	}{err.Error()})

	_, err = w.Write(data)
	return err
}
