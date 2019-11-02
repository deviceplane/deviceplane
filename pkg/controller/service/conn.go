package service

import (
	"bufio"
	"io"
	"net"
	"net/http"

	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/function61/holepunch-server/pkg/wsconnadapter"
)

func (s *Service) initiateDeviceConnection(w http.ResponseWriter, r *http.Request, projectID, deviceID string) {
	s.withHijackedWebSocketConnection(w, r, func(conn net.Conn) {
		s.connman.Set(projectID+deviceID, conn)
	})
}

func (s *Service) initiateSSH(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	s.withHijackedWebSocketConnection(w, r, func(conn net.Conn) {
		s.withDeviceConnection(w, r, projectID, deviceID, func(deviceConn net.Conn) {
			// TODO: build a proper client for this API
			req, _ := http.NewRequest("POST", "/ssh", nil)

			if err := req.Write(deviceConn); err != nil {
				http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
				return
			}

			go io.Copy(deviceConn, conn)
			io.Copy(conn, deviceConn)
		})
	})
}

func (s *Service) execute(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	s.withDeviceConnection(w, r, projectID, deviceID, func(deviceConn net.Conn) {
		// TODO: build a proper client for this API
		req, _ := http.NewRequest("POST", "/execute", r.Body)

		if _, ok := r.URL.Query()["background"]; ok {
			query := req.URL.Query()
			query.Add("background", "")
			req.URL.RawQuery = query.Encode()
		}

		if err := req.Write(deviceConn); err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		resp, err := http.ReadResponse(bufio.NewReader(deviceConn), req)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponse(w, resp)
	})
}

func (s *Service) metrics(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	s.withDeviceConnection(w, r, projectID, deviceID, func(deviceConn net.Conn) {
		// TODO: build a proper client for this API
		req, _ := http.NewRequest("GET", "/metrics", nil)

		if err := req.Write(deviceConn); err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		resp, err := http.ReadResponse(bufio.NewReader(deviceConn), req)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponse(w, resp)
	})
}

func (s *Service) withHijackedWebSocketConnection(w http.ResponseWriter, r *http.Request, f func(net.Conn)) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	f(wsconnadapter.New(conn))
}

func (s *Service) withDeviceConnection(w http.ResponseWriter, r *http.Request, projectID, deviceID string, f func(net.Conn)) {
	deviceConn, err := s.connman.Dial(r.Context(), projectID+deviceID)
	if err != nil {
		http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
		return
	}
	f(deviceConn)
}
