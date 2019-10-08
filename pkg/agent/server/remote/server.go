package remote

import (
	"context"
	"net/http"

	"github.com/deviceplane/deviceplane/pkg/agent/client"
	"github.com/deviceplane/deviceplane/pkg/agent/server/conncontext"
	"github.com/deviceplane/deviceplane/pkg/revdial"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Server struct {
	client     *client.Client
	httpServer *http.Server
}

func NewServer(client *client.Client, service http.Handler) *Server {
	return &Server{
		client: client,
		httpServer: &http.Server{
			Handler:     service,
			ConnContext: conncontext.SaveConn,
		},
	}
}

func (s *Server) Serve() error {
	conn, err := s.client.InitiateDeviceConnection(context.TODO())
	if err != nil {
		return errors.Wrap(err, "initiate connection")
	}

	listener := revdial.NewListener(conn, func(ctx context.Context, path string) (*websocket.Conn, *http.Response, error) {
		return s.client.Revdial(ctx, path)
	})
	defer listener.Close()

	return s.httpServer.Serve(listener)
}
