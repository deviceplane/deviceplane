package service

import (
	"bufio"
	"net"
	"net/http"

	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *Service) metrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	applicationID := vars["application"]
	service := vars["service"]

	containerID, ok := s.supervisorLookup.GetContainerID(applicationID, service)
	if !ok {
		w.WriteHeader(codes.StatusMetricsNotAvailable)
		return
	}

	if err := s.netnsManager.RunInContainerNamespace(r.Context(), containerID, func() {
		conn, err := net.Dial("tcp", "127.0.0.1:2112")
		if err != nil {
			http.Error(w, err.Error(), codes.StatusMetricsNotAvailable)
			return
		}

		req, _ := http.NewRequest("GET", "/metrics", nil)

		if err := req.Write(conn); err != nil {
			http.Error(w, err.Error(), codes.StatusMetricsNotAvailable)
			return
		}

		resp, err := http.ReadResponse(bufio.NewReader(conn), req)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusMetricsNotAvailable)
			return
		}

		utils.ProxyResponse(w, resp)
	}); err != nil {
		http.Error(w, errors.Wrap(err, "run in container namespace").Error(), codes.StatusInternalDeviceError)
		return
	}
}
