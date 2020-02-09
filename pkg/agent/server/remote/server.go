package remote

import (
	"context"
	"net/http"
	"time"

	"github.com/deviceplane/deviceplane/pkg/agent/client"
	"github.com/deviceplane/deviceplane/pkg/agent/server/conncontext"
	dpcontext "github.com/deviceplane/deviceplane/pkg/context"
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
	ctx, cancel := dpcontext.New(context.Background(), time.Minute)
	defer cancel()

	conn, err := s.client.InitiateDeviceConnection(ctx)
	if err != nil {
		return errors.Wrap(err, "initiate connection")
	}

	listener := revdial.NewListener(conn, s.revdial)
	defer listener.Close()

	return s.httpServer.Serve(listener)
}

func (s *Server) revdial(ctx context.Context, path string) (*websocket.Conn, *http.Response, error) {
	dpctx, cancel := dpcontext.New(ctx, time.Minute)
	defer cancel()

	conn, resp, err := s.client.Revdial(dpctx, path)
	if err != nil {
		return nil, nil, err
	}

	return conn.Conn, resp.Response, nil
}
