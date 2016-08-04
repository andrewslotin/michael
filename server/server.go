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

func (s *Server) Start(h http.Handler) error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s (%s)", s.Addr, err)
	}
	s.listener = listener

	go func(srv *http.Server, ln net.Listener) {
		err := srv.Serve(ln)
		if err != nil {
			log.Fatal(err)
		}
	}(&http.Server{Handler: h}, s.listener)

	return nil
}

func (s *Server) Shutdown() {
	s.listener.Close()
}
