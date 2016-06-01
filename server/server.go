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
		respondToUser(w, html.EscapeString(fmt.Sprintf("%s is not supported", cmd)))
		return
	}

	channelID := r.PostFormValue("channel_id")
	user := slack.User{
		ID:   r.PostFormValue("user_id"),
		Name: r.PostFormValue("user_name"),
	}

	switch subject := r.PostFormValue("text"); subject {
	case "", "help":
		respondToUser(w, html.EscapeString(HelpMessage))
		return
	case "status":
		d, ok := s.deploys.Get(channelID)
		if !ok {
			respondToUser(w, NoRunningDeploysMessage)
			return
		}

		respondToUser(w, fmt.Sprintf(DeployStatusMessage, d.User, d.Subject, d.StartedAt.Format(time.RFC822)))
	case "done":
		d, _ := s.deploys.Del(channelID)

		var response string
		if d.User.ID == user.ID {
			response = fmt.Sprintf(DeployDoneMessage, user)
		} else {
			response = fmt.Sprintf(DeployInterruptedMessage, user, d.User)
		}

		go sendDelayedResponse(w, r, response)
	default:
		if d, ok := s.deploys.Get(channelID); ok && d.User.ID != user.ID {
			respondToUser(w, fmt.Sprintf(DeployConflictMessage, d.User, d.StartedAt.Format(time.RFC822)))
			return
		}

		s.deploys.Set(channelID, user, subject)
		w.Write(nil)

		go sendDelayedResponse(w, r, fmt.Sprintf(DeployAnnouncementMessage, user, html.EscapeString(subject)))
	}
}

func respondToUser(w http.ResponseWriter, text string) {
	response, err := json.Marshal(slack.NewEphemeralResponse(text))
	if err != nil {
		log.Printf("failed to respond to user with %q (%s)", text, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func sendDelayedResponse(w http.ResponseWriter, r *http.Request, text string) {
	responseURL := r.PostFormValue("response_url")
	if responseURL == "" {
		log.Printf("cannot send delayed response to a without without response_url")
		return
	}

	response, err := json.Marshal(slack.NewInChannelResponse(text))
	if err != nil {
		log.Printf("failed to respond in channel with %s (%s)", text, err)
		return
	}

	slackResponse, err := http.Post(responseURL, "application/json", bytes.NewReader(response))
	if err != nil {
		log.Printf("failed to sent in_channel response (%s)", err)
		return
	}
	slackResponse.Body.Close()
}
