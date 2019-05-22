package service

import (
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/gorilla/mux"
)

func (s *Service) withApplication(handler func(http.ResponseWriter, *http.Request, string, string, string)) func(http.ResponseWriter, *http.Request, string, string) {
	return func(w http.ResponseWriter, r *http.Request, projectID, userID string) {
		vars := mux.Vars(r)
		application := vars["application"]
		if application == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var applicationID string
		if strings.Contains(application, "_") {
			applicationID = application
		} else {
			application, err := s.applications.LookupApplication(r.Context(), application, projectID)
			if err == store.ErrApplicationNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup application")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			applicationID = application.ID
		}

		handler(w, r, projectID, userID, applicationID)
	}
}

func (s *Service) withRole(handler func(http.ResponseWriter, *http.Request, string, string, string)) func(http.ResponseWriter, *http.Request, string, string) {
	return func(w http.ResponseWriter, r *http.Request, projectID, userID string) {
		vars := mux.Vars(r)
		role := vars["role"]
		if role == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var roleID string
		if strings.Contains(role, "_") {
			roleID = role
		} else {
			role, err := s.roles.LookupRole(r.Context(), role, projectID)
			if err == store.ErrRoleNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup role")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			roleID = role.ID
		}

		handler(w, r, projectID, userID, roleID)
	}
}

func (s *Service) withServiceAccount(handler func(http.ResponseWriter, *http.Request, string, string, string)) func(http.ResponseWriter, *http.Request, string, string) {
	return func(w http.ResponseWriter, r *http.Request, projectID, userID string) {
		vars := mux.Vars(r)
		serviceAccount := vars["serviceaccount"]
		if serviceAccount == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var serviceAccountID string
		if strings.Contains(serviceAccount, "_") {
			serviceAccountID = serviceAccount
		} else {
			serviceAccount, err := s.serviceAccounts.LookupServiceAccount(r.Context(), serviceAccount, projectID)
			if err == store.ErrServiceAccountNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup service account")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			serviceAccountID = serviceAccount.ID
		}

		handler(w, r, projectID, userID, serviceAccountID)
	}
}
