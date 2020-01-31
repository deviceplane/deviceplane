package service

import (
	"bufio"
	"io"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/deviceplane/deviceplane/pkg/agent/service/client"
	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/gorilla/mux"
)

func (s *Service) initiateDeviceConnection(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.connman.Set(with.Project.ID+with.Device.ID, with.ClientConn)
	})
}

var currentSSHCount int64

const currentSSHCountName = "internal.current_ssh_connection_count"

func (s *Service) initiateSSH(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := client.InitiateSSH(r.Context(), with.ClientConn)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		sshCount := atomic.AddInt64(&currentSSHCount, 1)
		s.st.Gauge(currentSSHCountName, float64(sshCount), utils.InternalTags(with.Project.ID), 1)
		defer func() {
			sshCount := atomic.AddInt64(&currentSSHCount, -1)
			s.st.Gauge(currentSSHCountName, float64(sshCount), utils.InternalTags(with.Project.ID), 1)
		}()

		go io.Copy(with.DeviceConn, with.ClientConn)
		io.Copy(with.ClientConn, with.DeviceConn)
	})
}

func (s *Service) initiateReboot(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.InitiateReboot(r.Context(), with.DeviceConn)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponseFromDevice(w, resp)
	})
}

func (s *Service) deviceDebug(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()
		dIndex := strings.Index(path, "/debug/")
		if dIndex == -1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		debugPath := path[dIndex:]

		req, err := http.NewRequestWithContext(
			r.Context(),
			r.Method,
			debugPath,
			r.Body,
		)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		if err := req.Write(with.DeviceConn); err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		resp, err := http.ReadResponse(bufio.NewReader(with.DeviceConn), req)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponseFromDevice(w, resp)
	})
}

func (s *Service) imagePullProgress(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		applicationID := vars["application"]
		service := vars["service"]

		resp, err := client.GetImagePullProgress(r.Context(), with.DeviceConn, applicationID, service)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponseFromDevice(w, resp)
	})
}

func (s *Service) hostMetrics(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.GetDeviceMetrics(r.Context(), with.DeviceConn)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponseFromDevice(w, resp)
	})
}

func (s *Service) agentMetrics(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.GetAgentMetrics(r.Context(), with.DeviceConn)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponseFromDevice(w, resp)
	})
}

func (s *Service) serviceMetrics(with *FetchObject) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		service := vars["service"]

		serviceMetricEndpointConfig, exists := with.Application.MetricEndpointConfigs[service]
		if !exists {
			serviceMetricEndpointConfig.Port = models.DefaultMetricPort
			serviceMetricEndpointConfig.Path = models.DefaultMetricPath
		}

		resp, err := client.GetServiceMetrics(
			r.Context(), with.DeviceConn, with.Application.ID, service,
			serviceMetricEndpointConfig.Path, serviceMetricEndpointConfig.Port,
		)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponseFromDevice(w, resp)
	})
}
