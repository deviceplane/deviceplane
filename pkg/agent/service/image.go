package service

import (
	"net/http"

	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/gorilla/mux"
)

func (s *Service) imagePullProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	applicationID := vars["application"]
	service := vars["service"]

	progress, ok := s.supervisorLookup.GetImagePullProgress(applicationID, service)
	if !ok {
		w.WriteHeader(codes.StatusImagePullProgressNotAvailable)
		return
	}

	utils.Respond(w, progress)
}
