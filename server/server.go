package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/andrewslotin/slack-deploy-command/deploy"
	"github.com/andrewslotin/slack-deploy-command/slack"
)

const (
	HelpMessage = `Available commands:

/deploy help — print help (this message)
/deploy <subject> — announce deploy of <subject> in channel
/deploy status — show deploy status in channel
/deploy done — finish deploy`
	NoRunningDeploysMessage   = "No one is deploying at the moment"
	DeployStatusMessage       = "%s is deploying %s since %s"
	DeployConflictMessage     = "%s is deploying since %s. You can type `/deploy done` if you think this deploy is finished."
	DeployAnnouncementMessage = "%s is about to deploy %s"
	DeployDoneMessage         = "%s done deploying"
	DeployInterruptedMessage  = "%s has finished the deploy started by %s"
)

type Server struct {
	Addr string

	listener   net.Listener
	slackToken string
	deploys    *deploy.Store
}

func New(host string, port int, slackToken string, deploys *deploy.Store) *Server {
	return &Server{
		Addr:       fmt.Sprintf("%s:%d", host, port),
		slackToken: slackToken,
		deploys:    deploys,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s (%s)", s.Addr, err)
	}
	s.listener = listener

	srv := http.Server{
		Handler: s,
	}

	go func() {
		err := srv.Serve(s.listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func (s *Server) Shutdown() {
	s.listener.Close()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are supported", http.StatusBadRequest)
		return
	}

	if r.PostFormValue("token") != s.slackToken {
		http.Error(w, "Invalid token", http.StatusForbidden)
		return
	}

	if cmd := r.PostFormValue("command"); cmd != "/deploy" {
		sendImmediateResponse(w, slack.NewEphemeralResponse(fmt.Sprintf("%s is not supported", html.EscapeString(cmd))))
		return
	}

	channelID := r.PostFormValue("channel_id")
	user := slack.User{
		ID:   r.PostFormValue("user_id"),
		Name: r.PostFormValue("user_name"),
	}

	switch subject := r.PostFormValue("text"); subject {
	case "", "help":
		sendImmediateResponse(w, slack.NewEphemeralResponse(html.EscapeString(HelpMessage)))
		return
	case "status":
		d, ok := s.deploys.Get(channelID)
		if !ok {
			sendImmediateResponse(w, slack.NewEphemeralResponse(NoRunningDeploysMessage))
			return
		}

		responseText := fmt.Sprintf(DeployStatusMessage, d.User, d.Subject, d.StartedAt.Format(time.RFC822))
		sendImmediateResponse(w, slack.NewEphemeralResponse(responseText))
	case "done":
		d, _ := s.deploys.Del(channelID)

		var responseText string
		if d.User.ID == user.ID {
			responseText = fmt.Sprintf(DeployDoneMessage, user)
		} else {
			responseText = fmt.Sprintf(DeployInterruptedMessage, user, d.User)
		}

		go sendDelayedResponse(w, r, slack.NewInChannelResponse(responseText))
	default:
		var responseText string
		if d, ok := s.deploys.Get(channelID); ok && d.User.ID != user.ID {
			responseText = fmt.Sprintf(DeployConflictMessage, d.User, d.StartedAt.Format(time.RFC822))
			sendImmediateResponse(w, slack.NewEphemeralResponse(responseText))
			return
		}

		s.deploys.Set(channelID, user, subject)
		w.Write(nil)

		responseText = fmt.Sprintf(DeployAnnouncementMessage, user, html.EscapeString(subject))

		response := slack.NewInChannelResponse(responseText)
		for _, ref := range deploy.FindReferences(subject) {
			response.Attachments = append(response.Attachments, slack.Attachment{
				Title:     ref.Repository + "#" + ref.ID,
				TitleLink: "https://github.com/" + ref.Repository + "/pulls/" + ref.ID,
			})
		}

		go sendDelayedResponse(w, r, response)
	}
}

func sendImmediateResponse(w http.ResponseWriter, response slack.Response) {
	body, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed to respond to user with %q (%s)", response.Text, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func sendDelayedResponse(w http.ResponseWriter, req *http.Request, response slack.Response) {
	responseURL := req.PostFormValue("response_url")
	if responseURL == "" {
		log.Printf("cannot send delayed response to a without without response_url")
		return
	}

	body, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed to respond in channel with %s (%s)", response.Text, err)
		return
	}

	slackResponse, err := http.Post(responseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("failed to sent in_channel response (%s)", err)
		return
	}
	slackResponse.Body.Close()
}
