package local

import (
	"net"
	"net/http"

	"github.com/deviceplane/deviceplane/pkg/agent/server/conncontext"
)

type Server struct {
	httpServer *http.Server
	listener   net.Listener
}

func NewServer(service http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Handler:     service,
			ConnContext: conncontext.SaveConn,
		},
	}
}

func (s *Server) SetListener(listener net.Listener) {
	s.listener = listener
}

func (s *Server) Serve() error {
	return s.httpServer.Serve(s.listener)
}
