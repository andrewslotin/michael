package dashboard

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/andrewslotin/slack-deploy-command/deploy"
)

var dashboardTemplate = template.Must(
	template.New("dashboard").
		Funcs(template.FuncMap{
			"ftime": func(t time.Time) string { return t.Format(time.RFC822) },
		}).
		Parse(strings.TrimSpace(`
Deploy history
--------------

{{ range . -}}
{{ if not .FinishedAt.IsZero -}}
  * {{ .User.Name }} was deploying {{ .Subject }} since {{ .StartedAt | ftime }} until {{ .FinishedAt | ftime }}
{{ else -}}
  * {{ .User.Name }} is currently deploying {{ .Subject }} since {{ .StartedAt | ftime }}
{{ end -}}
{{ else -}}
  No deploys in channel so far
{{ end }}
`)))

type Dashboard struct {
	repo deploy.Repository
}

func New(repo deploy.Repository) *Dashboard {
	return &Dashboard{
		repo: repo,
	}
}

func (h *Dashboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	channelID := r.FormValue("channel_id")
	if channelID == "" {
		http.Error(w, "Missing channel_id", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	dashboardTemplate.Execute(w, h.repo.All(channelID))
}
