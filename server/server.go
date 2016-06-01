package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

type Server struct {
	Addr string

	listener   net.Listener
	slackToken string
}

type Response struct {
	ResponseType string `json:"response_type,omitempty"`
	Text         string `json:"text"`
}

func New(host string, port int, slackToken string) *Server {
	return &Server{
		Addr:       fmt.Sprintf("%s:%d", host, port),
		slackToken: slackToken,
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
		respondToUser(w, fmt.Sprintf("%s is not supported", cmd))
		return
	}

	subject := r.PostFormValue("text")
	if subject == "" {
		respondToUser(w, "Please specify what you are about to deploy.\nExample:\n```\n/deploy repository-name#1 repository-name#2\n```")
		return
	}

	respondInChannel(w, fmt.Sprintf("%s is about to deploy %s", userLink(r.PostFormValue("user_id"), r.PostFormValue("user_name")), strings.Replace(subject, " ", ", ", -1)))
}

func userLink(userID, userName string) string {
	return "<@" + userID + "|" + userName + ">"
}

func respondToUser(w http.ResponseWriter, text string) {
	response, err := json.Marshal(Response{
		ResponseType: "ephemeral",
		Text:         text,
	})
	if err != nil {
		log.Printf("failed to respond to user with %q (%s)", text, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func respondInChannel(w http.ResponseWriter, text string) {
	response, err := json.Marshal(Response{
		ResponseType: "in_channel",
		Text:         text,
	})
	if err != nil {
		log.Printf("failed to respond in channel with %q (%s)", text, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
