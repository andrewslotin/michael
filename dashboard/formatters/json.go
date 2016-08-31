package formatters

import (
	"encoding/json"
	"net/http"

	"github.com/andrewslotin/michael/deploy"
)

var (
	JSON jsonFormatter
)

type jsonFormatter struct{}

func (jsonFormatter) RespondWithHistory(w http.ResponseWriter, history []deploy.Deploy) error {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(history)
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
