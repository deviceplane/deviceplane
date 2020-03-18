package service

import (
	"net"
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/controller/authz"
	serviceutils "github.com/deviceplane/deviceplane/pkg/controller/service/utils"
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

func (s *Service) withHijackedWebSocketConnection(w http.ResponseWriter, r *http.Request, f func(clientConn net.Conn)) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: set conn.CloseHandler() here

	f(wsconnadapter.New(conn))
}

func (s *Service) withDeviceConnection(w http.ResponseWriter, r *http.Request, project *models.Project, device *models.Device, f func(deviceConn net.Conn)) {
	deviceConn, err := s.connman.Dial(r.Context(), project.ID+device.ID)
	if err != nil {
		http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
		return
	}
	defer deviceConn.Close()

	f(deviceConn)
}

func (s *Service) validateAuthorization(
	requestedResource authz.Resource,
	requestedAction authz.Action,
	w http.ResponseWriter,
	r *http.Request,
	user *models.User,
	serviceAccount *models.ServiceAccount,
	f func(project *models.Project),
) {
	if user == nil && serviceAccount == nil {
		log.WithError(ErrDependencyNotSupplied).Error("validating authorization")
		w.WriteHeader(http.StatusInternalServerError)
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

	var roles []string
	superAdmin := false
	if user != nil {
		if user.SuperAdmin {
			superAdmin = true
		} else {
			if _, err := s.memberships.GetMembership(r.Context(),
				user.ID, project.ID,
			); err == store.ErrMembershipNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				// TODO: better logging all around
				log.WithField("user_id", user.ID).
					WithField("project_id", project.ID).
					WithError(err).
					Error("get membership")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			roleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(
				r.Context(),
				user.ID,
				project.ID,
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
	} else if serviceAccount != nil {
		// Sanity check that this service account belongs to this project
		if _, err := s.serviceAccounts.GetServiceAccount(r.Context(),
			serviceAccount.ID, project.ID,
		); err == store.ErrServiceAccountNotFound {
			w.WriteHeader(http.StatusForbidden)
			return
		} else if err != nil {
			log.WithError(err).Error("get service account")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		roleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(),
			serviceAccount.ID, project.ID)
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
			role, err := s.roles.GetRole(r.Context(), roleID, project.ID)
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

	f(project)
}

func (s *Service) withDeviceAuth(w http.ResponseWriter, r *http.Request, f func(project *models.Project, device *models.Device)) {
	vars := mux.Vars(r)
	projectID := vars["project"]

	deviceAccessKeyValue, _, _ := r.BasicAuth()
	if deviceAccessKeyValue == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	deviceAccessKey, err := s.deviceAccessKeys.ValidateDeviceAccessKey(r.Context(), projectID, hash.Hash(deviceAccessKeyValue))
	if err == store.ErrDeviceAccessKeyNotFound {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.WithError(err).Error("validate device access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	device, err := s.devices.GetDevice(r.Context(), deviceAccessKey.DeviceID, projectID)
	if err != nil {
		log.WithError(err).Error("get device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	project, err := s.projects.GetProject(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("get project")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f(project, device)
}

func (s *Service) withRole(w http.ResponseWriter, r *http.Request, project *models.Project, f func(role *models.Role)) {
	if project == nil {
		log.WithError(ErrDependencyNotSupplied).Error("getting role")
		w.WriteHeader(http.StatusInternalServerError)
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
		role, err = s.roles.GetRole(r.Context(), roleIdentifier, project.ID)
	} else {
		role, err = s.roles.LookupRole(r.Context(), roleIdentifier, project.ID)
	}
	if err == store.ErrRoleNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get/lookup role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f(role)
}

func (s *Service) withUserOrServiceAccountAuth(w http.ResponseWriter, r *http.Request, f func(user *models.User, serviceAccount *models.ServiceAccount)) {
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

	if userID != "" {
		user, err := s.users.GetUser(r.Context(), userID)
		if err != nil {
			log.WithError(err).Error("get user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		f(user, nil)
		return
	}
	if serviceAccountAccessKey != nil {
		serviceAccount, err := s.serviceAccounts.GetServiceAccount(r.Context(),
			serviceAccountAccessKey.ServiceAccountID, serviceAccountAccessKey.ProjectID)
		if err != nil {
			log.WithError(err).Error("get service account")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		f(nil, serviceAccount)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	return
}

func (s *Service) withValidatedSsoJWT(w http.ResponseWriter, r *http.Request, f func(ssoJWT models.SsoJWT)) {
	var ssoRequest models.Auth0SsoRequest
	if err := read(r, &ssoRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if s.auth0Audience == "" || s.auth0Domain == nil {
		http.Error(w, "SSO is not enabled", http.StatusNotImplemented)
		return
	}
	_, claims, err := serviceutils.ParseAndValidateSignedJWT(s.auth0Domain, s.auth0Audience, ssoRequest.IdToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	get := func(key string) (string, error) {
		v, ok := claims[key]
		if !ok {
			return "", errors.New("expected JWT claim not found")
		}
		value, ok := v.(string)
		if !ok {
			return "", errors.New("expected JWT claim to be string")
		}
		return value, nil
	}

	email, err := get("email")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name, err := get("name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: validate nonce
	_, err = get("nonce")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := get("sub")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subParts := strings.Split(sub, "|")
	if len(subParts) != 2 {
		http.Error(w, "invalid number of subject parts", http.StatusBadRequest)
		return
	}
	subProvider := subParts[0]
	subID := subParts[1]

	f(models.SsoJWT{
		Email:    email,
		Name:     name,
		Provider: subProvider,
		Subject:  subID,
		Claims:   claims,
	})
}

func (s *Service) withUserAuth(w http.ResponseWriter, r *http.Request, f func(user *models.User)) {
	s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
		if user == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		f(user)
	})
}

func (s *Service) withSuperUserAuth(w http.ResponseWriter, r *http.Request, user *models.User, f func()) {
	if user == nil {
		log.WithError(ErrDependencyNotSupplied).Error("verifying super-user auth")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !user.SuperAdmin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	f()
}

func (s *Service) withServiceAccount(w http.ResponseWriter, r *http.Request, project *models.Project, f func(serviceAccount *models.ServiceAccount)) {
	if project == nil {
		log.WithError(ErrDependencyNotSupplied).Error("getting service account")
		w.WriteHeader(http.StatusInternalServerError)
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
		serviceAccount, err = s.serviceAccounts.GetServiceAccount(r.Context(), serviceAccountIdentifier, project.ID)
	} else {
		serviceAccount, err = s.serviceAccounts.LookupServiceAccount(r.Context(), serviceAccountIdentifier, project.ID)
	}
	if err == store.ErrServiceAccountNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("lookup service account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f(serviceAccount)
}

func (s *Service) withApplication(w http.ResponseWriter, r *http.Request, project *models.Project, f func(application *models.Application)) {
	if project == nil {
		log.WithError(ErrDependencyNotSupplied).Error("getting application")
		w.WriteHeader(http.StatusInternalServerError)
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
		application, err = s.applications.GetApplication(r.Context(), applicationIdentifier, project.ID)
	} else {
		application, err = s.applications.LookupApplication(r.Context(), applicationIdentifier, project.ID)
	}
	if err == store.ErrApplicationNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("lookup application")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f(application)
}

func (s *Service) withRelease(w http.ResponseWriter, r *http.Request, project *models.Project, application *models.Application, f func(release *models.Release)) {
	if application == nil || project == nil {
		log.WithError(ErrDependencyNotSupplied).Error("getting release")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	releaseIdentifier := vars["release"]
	if releaseIdentifier == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var release *models.Release
	var err error
	release, err = utils.GetReleaseByIdentifier(s.releases, r.Context(), project.ID, application.ID, releaseIdentifier)
	if err == store.ErrReleaseNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get/lookup release")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f(release)
}

func (s *Service) withDevice(w http.ResponseWriter, r *http.Request, project *models.Project, f func(device *models.Device)) {
	vars := mux.Vars(r)
	deviceIdentifier := vars["device"]
	if deviceIdentifier == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var device *models.Device
	var err error
	if strings.Contains(deviceIdentifier, "_") {
		device, err = s.devices.GetDevice(r.Context(), deviceIdentifier, project.ID)
	} else {
		device, err = s.devices.LookupDevice(r.Context(), deviceIdentifier, project.ID)
	}
	if err == store.ErrDeviceNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("lookup device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f(device)
}

func (s *Service) withDeviceRegistrationToken(w http.ResponseWriter, r *http.Request, project *models.Project, f func(deviceRegistrationToken *models.DeviceRegistrationToken)) {
	vars := mux.Vars(r)
	tokenIdentifier := vars["deviceregistrationtoken"]
	if tokenIdentifier == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var token *models.DeviceRegistrationToken
	var err error
	if strings.Contains(tokenIdentifier, "_") {
		token, err = s.deviceRegistrationTokens.GetDeviceRegistrationToken(r.Context(), tokenIdentifier, project.ID)
	} else {
		token, err = s.deviceRegistrationTokens.LookupDeviceRegistrationToken(r.Context(), tokenIdentifier, project.ID)
	}
	if err == store.ErrDeviceRegistrationTokenNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("lookup device registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f(token)
}

func (s *Service) GorillaSuperUserAuth(handler http.Handler) http.Handler {
	return http.Handler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s.withUserOrServiceAccountAuth(w, r, func(user *models.User, serviceAccount *models.ServiceAccount) {
				s.withSuperUserAuth(w, r, user, func() {
					handler.ServeHTTP(w, r)
				})
			})
		},
	))
}
