package formatters

import (
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/andrewslotin/michael/deploy"
)

var (
	PlainText plainTextFormatter

	dashboardTemplate = template.Must(
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
{{ end }}`)))
)

type plainTextFormatter struct{}

func (plainTextFormatter) RespondWithHistory(w http.ResponseWriter, history []deploy.Deploy) error {
	w.Header().Set("Content-Type", "text/plain")
	return dashboardTemplate.Execute(w, history)
}

func (plainTextFormatter) RespondWithError(w http.ResponseWriter, err error, statusCode int) error {
	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, err.Error(), statusCode)
	return nil
}
