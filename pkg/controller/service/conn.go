package service

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Service) initiateDeviceConnection(w http.ResponseWriter, r *http.Request, projectID, deviceID string) {
	withHijackedConnection(w, func(conn net.Conn) {
		s.connKing.Set(projectID+deviceID, conn)
	})
}

func (s *Service) initiateSSH(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	withHijackedConnection(w, func(conn net.Conn) {
		vars := mux.Vars(r)
		deviceID := vars["device"]

		s.connKing.Join(projectID+deviceID, conn)
	})
}

func withHijackedConnection(w http.ResponseWriter, f func(net.Conn)) {
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
