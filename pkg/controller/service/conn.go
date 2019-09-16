package service

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

func (s *Service) initiateDeviceConnection(w http.ResponseWriter, r *http.Request, projectID, deviceID string) {
	withHijackedHTTPConnection(w, func(conn net.Conn) {
		s.connman.Set(projectID+deviceID, conn)
	})
}

func (s *Service) initiateSSH(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	withHijackedHTTPConnection(w, func(conn net.Conn) {
		s.connman.Join(projectID+deviceID, conn)
	})
}

func withHijackedHTTPConnection(w http.ResponseWriter, f func(net.Conn)) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n")

	f(conn)
}

type rwc struct {
	r io.Reader
	c *websocket.Conn
}

func (c *rwc) Write(p []byte) (int, error) {
	err := c.c.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *rwc) Read(p []byte) (int, error) {
	for {
		if c.r == nil {
			// Advance to next message.
			var err error
			_, c.r, err = c.c.NextReader()
			if err != nil {
				return 0, err
			}
		}
		n, err := c.r.Read(p)
		if err == io.EOF {
			// At end of message.
			c.r = nil
			if n > 0 {
				return n, nil
			} else {
				// No data read, continue to next message.
				continue
			}
		}
		return n, err
	}
}

func (c *rwc) Close() error {
	return c.c.Close()
}

func (s *Service) initiateWebSocketSSH(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	s.withHijackedWebSocketConnection(w, r, func(conn *websocket.Conn) {
		s.connman.Join(projectID+deviceID, &rwc{
			c: conn,
		})
	})
}

func (s *Service) withHijackedWebSocketConnection(w http.ResponseWriter, r *http.Request, f func(*websocket.Conn)) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f(conn)
}
