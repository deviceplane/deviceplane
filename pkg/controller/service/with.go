package service

import (
	"net"
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/controller/authz"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/hash"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/function61/holepunch-server/pkg/wsconnadapter"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var ErrDependencyNotSupplied = errors.New("internal route dependency not supplied")

type FetchObject struct {
	Project                 *models.Project
	Role                    *models.Role
	User                    *models.User
	ServiceAccount          *models.ServiceAccount
	Application             *models.Application
	Release                 *models.Release
	Device                  *models.Device
	DeviceRegistrationToken *models.DeviceRegistrationToken
	DeviceConn              net.Conn
	ClientConn              net.Conn
}

func HandlerFunc(hf http.HandlerFunc) Handler {
	return func(with *FetchObject) http.HandlerFunc {
		return hf
	}
}

type Handler func(with *FetchObject) http.HandlerFunc
type Middleware func(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc

func (s *Service) initWith(middlewares ...Middleware) func(handler Handler) http.HandlerFunc {
	return func(handler Handler) http.HandlerFunc {
		with := FetchObject{}

		if len(middlewares) == 0 {
			return handler(&with)
		}

		var chain http.HandlerFunc = handler(&with)
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]
			chain = middleware(&with, chain)
		}

		return chain
	}
}

func (s *Service) withHijackedWebSocketConnection(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: set conn.CloseHandler() here

		with.ClientConn = wsconnadapter.New(conn)
		hf.ServeHTTP(w, r)
	})
}

func (s *Service) withDeviceConnection(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deviceConn, err := s.connman.Dial(r.Context(), with.Project.ID+with.Device.ID)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}
		defer deviceConn.Close()

		with.DeviceConn = deviceConn
		hf.ServeHTTP(w, r)
	})
}

// Prefix with withUserOrServiceAccountAuth
func (s *Service) validateAuthorization(requestedResource authz.Resource, requestedAction authz.Action) func(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
	return func(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if with.User == nil || with.ServiceAccount == nil {
				http.Error(w, ErrDependencyNotSupplied.Error(), http.StatusInternalServerError)
				return
			}

			vars := mux.Vars(r)
			projectIdentifier := vars["project"]
			if projectIdentifier == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			var project *models.Project
			var err error
			if strings.Contains(projectIdentifier, "_") {
				project, err = s.projects.GetProject(r.Context(), projectIdentifier)
			} else {
				project, err = s.projects.LookupProject(r.Context(), projectIdentifier)
			}
			if err == store.ErrProjectNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup project")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			with.Project = project

			var roles []string
			superAdmin := false
			if with.User != nil {
				if with.User.SuperAdmin {
					superAdmin = true
				} else {
					if _, err := s.memberships.GetMembership(r.Context(),
						with.User.ID, with.Project.ID,
					); err == store.ErrMembershipNotFound {
						http.Error(w, err.Error(), http.StatusNotFound)
						return
					} else if err != nil {
						// TODO: better logging all around
						log.WithField("user_id", with.User.ID).
							WithField("project_id", with.Project.ID).
							WithError(err).
							Error("get membership")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					roleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(
						r.Context(),
						with.User.ID,
						with.Project.ID,
					)
					if err != nil {
						log.WithError(err).Error("list membership role bindings")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					for _, roleBinding := range roleBindings {
						roles = append(roles, roleBinding.RoleID)
					}
				}
			} else if with.ServiceAccount != nil {
				// Sanity check that this service account belongs to this project
				if _, err := s.serviceAccounts.GetServiceAccount(r.Context(),
					with.ServiceAccount.ID, with.Project.ID,
				); err == store.ErrServiceAccountNotFound {
					w.WriteHeader(http.StatusForbidden)
					return
				} else if err != nil {
					log.WithError(err).Error("get service account")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				roleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(),
					with.ServiceAccount.ID, with.Project.ID)
				if err != nil {
					log.WithError(err).Error("list service account role bindings")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				for _, roleBinding := range roleBindings {
					roles = append(roles, roleBinding.RoleID)
				}
			}

			var configs []authz.Config
			if superAdmin {
				configs = []authz.Config{
					authz.AdminAllRole,
				}
			} else {
				for _, roleID := range roles {
					role, err := s.roles.GetRole(r.Context(), roleID, with.Project.ID)
					if err == store.ErrRoleNotFound {
						continue
					} else if err != nil {
						log.WithError(err).Error("get role")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					var config authz.Config
					if err := yaml.Unmarshal([]byte(role.Config), &config); err != nil {
						log.WithError(err).Error("unmarshal role config")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					configs = append(configs, config)
				}
			}

			if !authz.Evaluate(requestedResource, requestedAction, configs) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			hf.ServeHTTP(w, r)
		})
	}
}

func (s *Service) withDeviceAuth(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deviceAccessKeyValue, _, _ := r.BasicAuth()
		if deviceAccessKeyValue == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		deviceAccessKey, err := s.deviceAccessKeys.ValidateDeviceAccessKey(r.Context(), with.Project.ID, hash.Hash(deviceAccessKeyValue))
		if err == store.ErrDeviceAccessKeyNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			log.WithError(err).Error("validate device access key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		device, err := s.devices.GetDevice(r.Context(), deviceAccessKey.DeviceID, with.Project.ID)
		if err != nil {
			log.WithError(err).Error("get device")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		with.Device = device
		hf.ServeHTTP(w, r)
	})
}

func (s *Service) withRole(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if with.Project == nil {
			http.Error(w, ErrDependencyNotSupplied.Error(), http.StatusInternalServerError)
			return
		}

		roleIdentifier := mux.Vars(r)["role"]
		if roleIdentifier == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var role *models.Role
		var err error
		if strings.Contains(roleIdentifier, "_") {
			role, err = s.roles.GetRole(r.Context(), roleIdentifier, with.Project.ID)
		} else {
			role, err = s.roles.LookupRole(r.Context(), roleIdentifier, with.Project.ID)
		}
		if err == store.ErrRoleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("get/lookup role")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		with.Role = role
		hf.ServeHTTP(w, r)
	})
}

func (s *Service) withUserOrServiceAccountAuth(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string
		var serviceAccountAccessKey *models.ServiceAccountAccessKey

		sessionValue, err := r.Cookie(sessionCookie)

		switch err {
		case nil:
			session, err := s.sessions.ValidateSession(r.Context(), hash.Hash(sessionValue.Value))
			if err == store.ErrSessionNotFound {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if err != nil {
				log.WithError(err).Error("validate session")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			userID = session.UserID
		case http.ErrNoCookie:
			accessKeyValue, _, _ := r.BasicAuth()
			if accessKeyValue == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if strings.HasPrefix(accessKeyValue, "u") {
				userAccessKey, err := s.userAccessKeys.ValidateUserAccessKey(r.Context(), hash.Hash(accessKeyValue))
				if err == store.ErrUserAccessKeyNotFound {
					w.WriteHeader(http.StatusUnauthorized)
					return
				} else if err != nil {
					log.WithError(err).Error("validate user access key")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				userID = userAccessKey.UserID
			} else if strings.HasPrefix(accessKeyValue, "s") {
				serviceAccountAccessKey, err = s.serviceAccountAccessKeys.ValidateServiceAccountAccessKey(r.Context(), hash.Hash(accessKeyValue))
				if err == store.ErrServiceAccountAccessKeyNotFound {
					w.WriteHeader(http.StatusUnauthorized)
					return
				} else if err != nil {
					log.WithError(err).Error("validate service account access key")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		default:
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if userID == "" && serviceAccountAccessKey == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var user *models.User
		var serviceAccount *models.ServiceAccount
		if userID != "" {
			user, err = s.users.GetUser(r.Context(), userID)
			if err != nil {
				log.WithError(err).Error("get user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !user.RegistrationCompleted {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		} else if serviceAccountAccessKey != nil {
			serviceAccount, err = s.serviceAccounts.GetServiceAccount(r.Context(),
				serviceAccountAccessKey.ServiceAccountID, serviceAccountAccessKey.ProjectID)
			if err != nil {
				log.WithError(err).Error("get service account")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		with.User = user
		with.ServiceAccount = serviceAccount
		hf.ServeHTTP(w, r)
	})
}

func (s *Service) withSuperUserAuth(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if with.User == nil {
			http.Error(w, ErrDependencyNotSupplied.Error(), http.StatusInternalServerError)
			return
		}

		if !with.User.SuperAdmin {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		hf.ServeHTTP(w, r)
	})
}

func (s *Service) withServiceAccount(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if with.Project == nil {
			http.Error(w, ErrDependencyNotSupplied.Error(), http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)
		serviceAccountIdentifier := vars["serviceaccount"]
		if serviceAccountIdentifier == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var serviceAccount *models.ServiceAccount
		var err error
		if strings.Contains(serviceAccountIdentifier, "_") {
			serviceAccount, err = s.serviceAccounts.GetServiceAccount(r.Context(), serviceAccountIdentifier, with.Project.ID)
		} else {
			serviceAccount, err = s.serviceAccounts.LookupServiceAccount(r.Context(), serviceAccountIdentifier, with.Project.ID)
		}
		if err == store.ErrServiceAccountNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("lookup service account")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		with.ServiceAccount = serviceAccount
		hf.ServeHTTP(w, r)
	})
}

func (s *Service) withApplication(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if with.Project == nil {
			http.Error(w, ErrDependencyNotSupplied.Error(), http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)
		applicationIdentifier := vars["application"]
		if applicationIdentifier == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var application *models.Application
		var err error
		if strings.Contains(applicationIdentifier, "_") {
			application, err = s.applications.GetApplication(r.Context(), applicationIdentifier, with.Project.ID)
		} else {
			application, err = s.applications.LookupApplication(r.Context(), applicationIdentifier, with.Project.ID)
		}
		if err == store.ErrApplicationNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("lookup application")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		with.Application = application
		hf.ServeHTTP(w, r)
	})
}

func (s *Service) withRelease(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if with.Application == nil {
			http.Error(w, ErrDependencyNotSupplied.Error(), http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)
		releaseIdentifier := vars["application"]
		if releaseIdentifier == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var release *models.Release
		var err error
		release, err = utils.GetReleaseByIdentifier(s.releases, r.Context(), with.Project.ID, with.Application.ID, releaseIdentifier)
		if err == store.ErrReleaseNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("get/lookup release")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		with.Release = release
		hf.ServeHTTP(w, r)
	})
}

func (s *Service) withDevice(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		deviceIdentifier := vars["device"]
		if deviceIdentifier == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var device *models.Device
		var err error
		if strings.Contains(deviceIdentifier, "_") {
			device, err = s.devices.GetDevice(r.Context(), deviceIdentifier, with.Project.ID)
		} else {
			device, err = s.devices.LookupDevice(r.Context(), deviceIdentifier, with.Project.ID)
		}
		if err == store.ErrDeviceNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("lookup device")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		with.Device = device
		hf.ServeHTTP(w, r)
	})
}

func (s *Service) withDeviceRegistrationToken(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenIdentifier := vars["deviceregistrationtoken"]
		if tokenIdentifier == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var token *models.DeviceRegistrationToken
		var err error
		if strings.Contains(tokenIdentifier, "_") {
			token, err = s.deviceRegistrationTokens.GetDeviceRegistrationToken(r.Context(), tokenIdentifier, with.Project.ID)
		} else {
			token, err = s.deviceRegistrationTokens.LookupDeviceRegistrationToken(r.Context(), tokenIdentifier, with.Project.ID)
		}
		if err == store.ErrDeviceRegistrationTokenNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("lookup device registration token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		with.DeviceRegistrationToken = token
		hf.ServeHTTP(w, r)
	})
}
