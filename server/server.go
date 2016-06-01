package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

type Server struct {
	Addr string

	listener net.Listener
}

func New(host string, port int) *Server {
	return &Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
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
	fmt.Fprint(w, "Hi there!")
}
