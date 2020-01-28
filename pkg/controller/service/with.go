package service

import (
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/gorilla/mux"
)

func (s *Service) withRole(handler func(http.ResponseWriter, *http.Request, string, string, string, string)) func(http.ResponseWriter, *http.Request, string, string, string) {
	return func(w http.ResponseWriter, r *http.Request, projectID, authenticatedUserID, authenticatedServiceAccountID string) {
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
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup role")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			roleID = role.ID
		}

		handler(w, r, projectID, authenticatedUserID, authenticatedServiceAccountID, roleID)
	}
}

func (s *Service) withServiceAccount(handler func(http.ResponseWriter, *http.Request, string, string, string, string)) func(http.ResponseWriter, *http.Request, string, string, string) {
	return func(w http.ResponseWriter, r *http.Request, projectID, authenticatedUserID, authenticatedServiceAccountID string) {
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
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup service account")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			serviceAccountID = serviceAccount.ID
		}

		handler(w, r, projectID, authenticatedUserID, authenticatedServiceAccountID, serviceAccountID)
	}
}

func (s *Service) withApplication(handler func(http.ResponseWriter, *http.Request, string, string, string, string)) func(http.ResponseWriter, *http.Request, string, string, string) {
	return func(w http.ResponseWriter, r *http.Request, projectID, authenticatedUserID, authenticatedServiceAccountID string) {
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
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup application")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			applicationID = application.ID
		}

		handler(w, r, projectID, authenticatedUserID, authenticatedServiceAccountID, applicationID)
	}
}

func (s *Service) withApplicationAndRelease(handler func(http.ResponseWriter, *http.Request, string, string, string, *models.Application, *models.Release)) func(http.ResponseWriter, *http.Request, string, string, string) {
	return func(w http.ResponseWriter, r *http.Request, projectID, authenticatedUserID, authenticatedServiceAccountID string) {
		vars := mux.Vars(r)
		applicationMuxVar := vars["application"]
		releaseMuxVar := vars["release"]

		if applicationMuxVar == "" || releaseMuxVar == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var application *models.Application
		var err error
		if strings.Contains(applicationMuxVar, "_") {
			application, err = s.applications.GetApplication(r.Context(), applicationMuxVar, projectID)
		} else {
			application, err = s.applications.LookupApplication(r.Context(), applicationMuxVar, projectID)
		}
		if err == store.ErrApplicationNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("get application")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var release *models.Release
		release, err = utils.GetReleaseByIdentifier(s.releases, r.Context(), projectID, application.ID, releaseMuxVar)
		if err == store.ErrReleaseNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("get release")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler(w, r, projectID, authenticatedUserID, authenticatedServiceAccountID, application, release)
	}
}

func (s *Service) withDevice(handler func(http.ResponseWriter, *http.Request, string, string, string, string)) func(http.ResponseWriter, *http.Request, string, string, string) {
	return func(w http.ResponseWriter, r *http.Request, projectID, authenticatedUserID, authenticatedServiceAccountID string) {
		vars := mux.Vars(r)
		device := vars["device"]
		if device == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var deviceID string
		if strings.Contains(device, "_") {
			deviceID = device
		} else {
			device, err := s.devices.LookupDevice(r.Context(), device, projectID)
			if err == store.ErrDeviceNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup device")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			deviceID = device.ID
		}

		handler(w, r, projectID, authenticatedUserID, authenticatedServiceAccountID, deviceID)
	}
}

func (s *Service) withApplicationAndDevice(handler func(http.ResponseWriter, *http.Request, string, string, string, string, string)) func(http.ResponseWriter, *http.Request, string, string, string) {
	return func(w http.ResponseWriter, r *http.Request, projectID, authenticatedUserID, authenticatedServiceAccountID string) {
		vars := mux.Vars(r)
		application := vars["application"]
		if application == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		device := vars["device"]
		if device == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var applicationID string
		if strings.Contains(application, "_") {
			applicationID = application
		} else {
			application, err := s.applications.LookupApplication(r.Context(), application, projectID)
			if err == store.ErrApplicationNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup application")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			applicationID = application.ID
		}

		var deviceID string
		if strings.Contains(device, "_") {
			deviceID = device
		} else {
			device, err := s.devices.LookupDevice(r.Context(), device, projectID)
			if err == store.ErrDeviceNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup device")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			deviceID = device.ID
		}

		handler(w, r, projectID, authenticatedUserID, authenticatedServiceAccountID, applicationID, deviceID)
	}
}

func (s *Service) withDeviceRegistrationToken(handler func(http.ResponseWriter, *http.Request, string, string, string, string)) func(http.ResponseWriter, *http.Request, string, string, string) {
	return func(w http.ResponseWriter, r *http.Request, projectID, authenticatedUserID, authenticatedServiceAccountID string) {
		vars := mux.Vars(r)
		token := vars["deviceregistrationtoken"]
		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var tokenID string
		if strings.Contains(token, "_") {
			tokenID = token
		} else {
			token, err := s.deviceRegistrationTokens.LookupDeviceRegistrationToken(r.Context(), token, projectID)
			if err == store.ErrDeviceRegistrationTokenNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup device registration token")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			tokenID = token.ID
		}

		handler(w, r, projectID, authenticatedUserID, authenticatedServiceAccountID, tokenID)
	}
}
