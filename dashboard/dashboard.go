package dashboard

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/andrewslotin/michael/dashboard/formatters"
	"github.com/andrewslotin/michael/deploy"
)

type Dashboard struct {
	repo deploy.Repository
}

func New(repo deploy.Repository) *Dashboard {
	return &Dashboard{
		repo: repo,
	}
}

func (h *Dashboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	channelID := ChannelIDFromRequest(r)
	if channelID == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var history []deploy.Deploy
	if v := r.FormValue("since"); v != "" {
		timeSince, err := time.Parse(time.RFC3339, v)
		if err != nil {
			if err = Responder(r).RespondWithError(w, errors.New("Malformed time in `since` parameter"), http.StatusBadRequest); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		history = h.repo.Since(channelID, timeSince)
	} else {
		history = h.repo.All(channelID)
	}

	if err := Responder(r).RespondWithHistory(w, history); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ChannelIDFromRequest extracts and returns channelID from request URL.
func ChannelIDFromRequest(r *http.Request) string {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		return ""
	}

	if n := strings.IndexByte(path, '/'); n >= 0 {
		return path[:n]
	} else if n := strings.LastIndexByte(path, '.'); n >= 0 {
		return path[:n]
	}

	return path
}

// Responder returns a formatters.ResponseFormatter according to the extension in URL path.
func Responder(r *http.Request) formatters.ResponseFormatter {
	switch {
	case strings.HasSuffix(r.URL.Path, ".json"):
		return formatters.JSON
	case strings.HasSuffix(r.URL.Path, ".txt"):
		return formatters.PlainText
	default:
		return formatters.PlainText
	}
}
