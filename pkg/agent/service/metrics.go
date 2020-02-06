package service

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/gorilla/mux"
)

func (s *Service) metrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	applicationID := vars["application"]
	service := vars["service"]

	query := r.URL.Query()

	portRaw := query.Get("port")
	port, err := strconv.Atoi(portRaw)
	if err != nil {
		http.Error(w, "invalid port", 400)
		return
	}

	pathRaw := query.Get("path")
	path, err := base64.RawURLEncoding.DecodeString(pathRaw)
	if err != nil {
		http.Error(w, "invalid base64 encoded path", 400)
		return
	}

	resp, err := s.serviceMetricsFetcher.ContainerServiceMetrics(
		r.Context(),
		applicationID,
		service,
		port,
		string(path),
	)
	if err != nil {
		http.Error(w, err.Error(), codes.StatusMetricsNotAvailable)
		return
	}
	defer resp.Body.Close()

	utils.ProxyResponse(w, resp)
}
