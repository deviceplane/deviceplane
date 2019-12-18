package service

import (
	"io"
	"net"
	"net/http"

	"github.com/deviceplane/deviceplane/pkg/agent/service/client"
	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/function61/holepunch-server/pkg/wsconnadapter"
	"github.com/gorilla/mux"
)

func (s *Service) initiateDeviceConnection(w http.ResponseWriter, r *http.Request, project models.Project, device models.Device) {
	s.withHijackedWebSocketConnection(w, r, func(conn net.Conn) {
		s.connman.Set(project.ID+device.ID, conn)
	})
}

func (s *Service) initiateSSH(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	s.withHijackedWebSocketConnection(w, r, func(conn net.Conn) {
		s.withDeviceConnection(w, r, projectID, deviceID, func(deviceConn net.Conn) {
			err := client.InitiateSSH(deviceConn)
			if err != nil {
				http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
				return
			}

			go io.Copy(deviceConn, conn)
			io.Copy(conn, deviceConn)
		})
	})
}

func (s *Service) initiateReboot(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	s.withDeviceConnection(w, r, projectID, deviceID, func(deviceConn net.Conn) {
		resp, err := client.InitiateReboot(deviceConn)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponseFromDevice(w, resp)
	})
}

func (s *Service) imagePullProgress(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	vars := mux.Vars(r)
	applicationID := vars["application"]
	service := vars["service"]

	s.withDeviceConnection(w, r, projectID, deviceID, func(deviceConn net.Conn) {
		resp, err := client.GetImagePullProgress(deviceConn, applicationID, service)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponseFromDevice(w, resp)
	})
}

func (s *Service) hostMetrics(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	s.withDeviceConnection(w, r, projectID, deviceID, func(deviceConn net.Conn) {
		resp, err := client.GetHostMetrics(deviceConn)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponseFromDevice(w, resp)
	})
}

func (s *Service) agentMetrics(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	s.withDeviceConnection(w, r, projectID, deviceID, func(deviceConn net.Conn) {
		resp, err := client.GetAgentMetrics(deviceConn)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponseFromDevice(w, resp)
	})
}

func (s *Service) serviceMetrics(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	applicationID, deviceID string,
) {
	vars := mux.Vars(r)
	service := vars["service"]

	app, err := s.applications.GetApplication(r.Context(), applicationID, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	serviceMetricsConfig, exists := app.ServiceMetricsConfig[service]
	if !exists {
		serviceMetricsConfig.Port = 2112
		serviceMetricsConfig.Path = "/metrics"
	}

	s.withDeviceConnection(w, r, projectID, deviceID, func(deviceConn net.Conn) {
		resp, err := client.GetServiceMetrics(deviceConn, applicationID, service, serviceMetricsConfig.Path, serviceMetricsConfig.Port)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponseFromDevice(w, resp)
	})
}

func (s *Service) withHijackedWebSocketConnection(w http.ResponseWriter, r *http.Request, f func(net.Conn)) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// We should set conn.CloseHandler() here

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
