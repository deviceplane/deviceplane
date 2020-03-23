package service

import (
	"net/http"

	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/gorilla/mux"
)

func (s *Service) metrics(w http.ResponseWriter, r *http.Request) {
	withPort(w, r, func(port int) {
		withPath(w, r, func(path string) {
			vars := mux.Vars(r)
			applicationID := vars["application"]
			service := vars["service"]

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
		})
	})
}
