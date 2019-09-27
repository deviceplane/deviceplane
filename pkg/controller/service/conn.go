package service

import (
	"net/http"

	"github.com/function61/holepunch-server/pkg/wsconnadapter"
	"github.com/gorilla/websocket"
)

func (s *Service) initiateDeviceConnection(w http.ResponseWriter, r *http.Request, projectID, deviceID string) {
	s.withHijackedWebSocketConnection(w, r, func(conn *websocket.Conn) {
		s.connman.Set(projectID+deviceID, wsconnadapter.New(conn))
	})
}

func (s *Service) initiateSSH(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	s.withHijackedWebSocketConnection(w, r, func(conn *websocket.Conn) {
		s.connman.Join(projectID+deviceID, wsconnadapter.New(conn))
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
