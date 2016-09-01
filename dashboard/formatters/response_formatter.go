package formatters

import (
	"net/http"

	"github.com/andrewslotin/michael/deploy"
)

type ResponseFormatter interface {
	RespondWithHistory(http.ResponseWriter, []deploy.Deploy) error
	RespondWithError(http.ResponseWriter, error, int) error
}
