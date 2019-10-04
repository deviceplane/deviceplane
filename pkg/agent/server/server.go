package server

import (
	"net"
	"net/http"

	"github.com/deviceplane/deviceplane/pkg/agent/service"
)

type Server struct {
	httpServer *http.Server
	listener   net.Listener
}

func NewServer() *Server {
	return &Server{
		httpServer: &http.Server{
			Handler: service.NewService(),
		},
	}
}

func (s *Server) SetListener(listener net.Listener) {
	s.listener = listener
}

func (s *Server) Serve() error {
	return s.httpServer.Serve(s.listener)
}
