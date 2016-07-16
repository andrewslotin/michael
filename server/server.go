package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/andrewslotin/slack-deploy-command/deploy/stores"
	"github.com/andrewslotin/slack-deploy-command/github"
	"github.com/andrewslotin/slack-deploy-command/slack"
)

type DeployEventHandler interface {
	DeployStarted(channelID string)
	DeployCompleted(channelID string)
}

type Server struct {
	Addr string

	listener   net.Listener
	slackToken string
	deploys    stores.Store
	responses  *ResponseBuilder

	deployEventHandlers []DeployEventHandler
}

func New(host string, port int, slackToken, githubToken string, deploys stores.Store) *Server {
	return &Server{
		Addr:       fmt.Sprintf("%s:%d", host, port),
		slackToken: slackToken,
		deploys:    deploys,
		responses:  NewResponseBuilder(github.NewClient(githubToken, nil)),
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

func (s *Server) AddDeployEventHandler(h DeployEventHandler) {
	s.deployEventHandlers = append(s.deployEventHandlers, h)
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
		sendImmediateResponse(w, s.responses.ErrorMessage(cmd, errors.New("not supported")))
		return
	}

	channelID := r.PostFormValue("channel_id")
	user := slack.User{
		ID:   r.PostFormValue("user_id"),
		Name: r.PostFormValue("user_name"),
	}

	switch subject := r.PostFormValue("text"); subject {
	case "", "help":
		sendImmediateResponse(w, s.responses.HelpMessage())
	case "status":
		d, ok := s.deploys.Get(channelID)
		if !ok {
			sendImmediateResponse(w, s.responses.NoRunningDeploysMessage())
			return
		}

		sendImmediateResponse(w, s.responses.DeployStatusMessage(d))
	case "done":
		d, ok := s.deploys.Del(channelID)
		if !ok {
			sendImmediateResponse(w, s.responses.NoRunningDeploysMessage())
			return
		}

		if d.User.ID == user.ID {
			go sendDelayedResponse(w, r, s.responses.DeployDoneAnnouncement(user))
		} else {
			go sendDelayedResponse(w, r, s.responses.DeployInterruptedAnnouncement(d, user))
		}

		for _, h := range s.deployEventHandlers {
			go h.DeployCompleted(channelID)
		}
	default:
		if d, ok := s.deploys.Get(channelID); ok && d.User.ID != user.ID {
			sendImmediateResponse(w, s.responses.DeployInProgressMessage(d))
			return
		}

		s.deploys.Set(channelID, user, subject)
		w.Write(nil)

		go sendDelayedResponse(w, r, s.responses.DeployAnnouncement(user, subject))
		for _, h := range s.deployEventHandlers {
			go h.DeployStarted(channelID)
		}
	}
}

func sendImmediateResponse(w http.ResponseWriter, response *slack.Response) {
	body, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed to respond to user with %q (%s)", response.Text, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func sendDelayedResponse(w http.ResponseWriter, req *http.Request, response *slack.Response) {
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
