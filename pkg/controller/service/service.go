package service

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/segmentio/ksuid"
	"gopkg.in/yaml.v2"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/controller/authz"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/email"
	"github.com/deviceplane/deviceplane/pkg/hash"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/gorilla/mux"
)

const (
	sessionCookie = "dp_sess"
)

type Service struct {
	users                      store.Users
	registrationTokens         store.RegistrationTokens
	userAccessKeys             store.UserAccessKeys
	sessions                   store.Sessions
	projects                   store.Projects
	projectDeviceCounts        store.ProjectDeviceCounts
	projectApplicationCounts   store.ProjectApplicationCounts
	roles                      store.Roles
	memberships                store.Memberships
	membershipRoleBindings     store.MembershipRoleBindings
	serviceAccounts            store.ServiceAccounts
	serviceAccountAccessKeys   store.ServiceAccountAccessKeys
	serviceAccountRoleBindings store.ServiceAccountRoleBindings
	devices                    store.Devices
	deviceStatuses             store.DeviceStatuses
	deviceLabels               store.DeviceLabels
	deviceRegistrationTokens   store.DeviceRegistrationTokens
	deviceAccessKeys           store.DeviceAccessKeys
	applications               store.Applications
	applicationDeviceCounts    store.ApplicationDeviceCounts
	releases                   store.Releases
	releaseDeviceCounts        store.ReleaseDeviceCounts
	deviceApplicationStatuses  store.DeviceApplicationStatuses
	deviceServiceStatuses      store.DeviceServiceStatuses
	email                      email.Interface
	router                     *mux.Router
	cookieDomain               string
	cookieSecure               bool
}

func NewService(
	users store.Users,
	registrationTokens store.RegistrationTokens,
	sessions store.Sessions,
	userAccessKeys store.UserAccessKeys,
	projects store.Projects,
	projectDeviceCounts store.ProjectDeviceCounts,
	projectApplicationCounts store.ProjectApplicationCounts,
	roles store.Roles,
	memberships store.Memberships,
	membershipRoleBindings store.MembershipRoleBindings,
	serviceAccounts store.ServiceAccounts,
	serviceAccountAccessKeys store.ServiceAccountAccessKeys,
	serviceAccountRoleBindings store.ServiceAccountRoleBindings,
	devices store.Devices,
	deviceStatuses store.DeviceStatuses,
	deviceLabels store.DeviceLabels,
	deviceRegistrationTokens store.DeviceRegistrationTokens,
	deviceAccessKeys store.DeviceAccessKeys,
	applications store.Applications,
	applicationDeviceCounts store.ApplicationDeviceCounts,
	releases store.Releases,
	releasesDeviceCounts store.ReleaseDeviceCounts,
	deviceApplicationStatuses store.DeviceApplicationStatuses,
	deviceServiceStatuses store.DeviceServiceStatuses,
	email email.Interface,
	cookieDomain string,
	cookieSecure bool,
) *Service {
	s := &Service{
		users:                      users,
		registrationTokens:         registrationTokens,
		sessions:                   sessions,
		userAccessKeys:             userAccessKeys,
		projects:                   projects,
		projectDeviceCounts:        projectDeviceCounts,
		projectApplicationCounts:   projectApplicationCounts,
		roles:                      roles,
		memberships:                memberships,
		membershipRoleBindings:     membershipRoleBindings,
		serviceAccounts:            serviceAccounts,
		serviceAccountAccessKeys:   serviceAccountAccessKeys,
		serviceAccountRoleBindings: serviceAccountRoleBindings,
		devices:                    devices,
		deviceStatuses:             deviceStatuses,
		deviceLabels:               deviceLabels,
		deviceRegistrationTokens:   deviceRegistrationTokens,
		deviceAccessKeys:           deviceAccessKeys,
		applications:               applications,
		applicationDeviceCounts:    applicationDeviceCounts,
		releases:                   releases,
		releaseDeviceCounts:        releasesDeviceCounts,
		deviceApplicationStatuses:  deviceApplicationStatuses,
		deviceServiceStatuses:      deviceServiceStatuses,
		email:                      email,
		cookieDomain:               cookieDomain,
		cookieSecure:               cookieSecure,
		router:                     mux.NewRouter(),
	}

	s.router.HandleFunc("/health", s.health).Methods("GET")

	s.router.HandleFunc("/register", s.register).Methods("POST")
	s.router.HandleFunc("/completeregistration", s.confirmRegistration).Methods("POST")
	s.router.HandleFunc("/login", s.login).Methods("POST")
	s.router.HandleFunc("/logout", s.logout).Methods("POST")
	s.router.HandleFunc("/me", s.withUserOrServiceAccountAuth(s.me)).Methods("GET")

	s.router.HandleFunc("/users/{user}/memberships", s.withUserOrServiceAccountAuth(s.listMembershipsByUser)).Methods("GET")
	s.router.HandleFunc("/users/{user}/memberships/full", s.withUserOrServiceAccountAuth(s.listMembershipsByUserFull)).Methods("GET")

	s.router.HandleFunc("/projects", s.withUserOrServiceAccountAuth(s.createProject)).Methods("POST")
	s.router.HandleFunc("/projects/{project}", s.validateAuthorization("projects", "GetProject", s.getProject)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/roles", s.validateAuthorization("roles", "CreateRole", s.createRole)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/roles/{role}", s.validateAuthorization("roles", "GetRole", s.withRole(s.getRole))).Methods("GET")
	s.router.HandleFunc("/projects/{project}/roles", s.validateAuthorization("roles", "ListRoles", s.listRoles)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/roles/{role}", s.validateAuthorization("roles", "UpdateRole", s.updateRole)).Methods("PUT")
	s.router.HandleFunc("/projects/{project}/roles/{role}", s.validateAuthorization("roles", "DeleteRole", s.deleteRole)).Methods("DELETE")

	s.router.HandleFunc("/projects/{project}/memberships", s.validateAuthorization("memberships", "CreateMembership", s.createMembership)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/memberships/{user}", s.validateAuthorization("memberships", "GetMembership", s.getMembership)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/memberships", s.validateAuthorization("memberships", "ListMembershipsByProject", s.listMembershipsByProject)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/memberships/{user}/roles/{role}/membershiprolebindings", s.validateAuthorization("membershiprolebindings", "CreateMembershipRoleBinding", s.withRole(s.createMembershipRoleBinding))).Methods("POST")
	s.router.HandleFunc("/projects/{project}/memberships/{user}/roles/{role}/membershiprolebindings", s.validateAuthorization("membershiprolebindings", "GetMembershipRoleBinding", s.withRole(s.getMembershipRoleBinding))).Methods("GET")
	s.router.HandleFunc("/projects/{project}/memberships/{user}/membershiprolebindings", s.validateAuthorization("membershiprolebindings", "ListMembershipRoleBindings", s.listMembershipRoleBindings)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/memberships/{user}/roles/{role}/membershiprolebindings", s.validateAuthorization("membershiprolebindings", "DeleteMembershipRoleBinding", s.withRole(s.deleteMembershipRoleBinding))).Methods("DELETE")

	s.router.HandleFunc("/projects/{project}/serviceaccounts", s.validateAuthorization("serviceaccounts", "CreateServiceAccount", s.createServiceAccount)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}", s.validateAuthorization("serviceaccounts", "GetServiceAccount", s.withServiceAccount(s.getServiceAccount))).Methods("GET")
	s.router.HandleFunc("/projects/{project}/serviceaccounts", s.validateAuthorization("serviceaccounts", "ListServiceAccounts", s.listServiceAccounts)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}", s.validateAuthorization("serviceaccounts", "UpdateServiceAccount", s.withServiceAccount(s.updateServiceAccount))).Methods("PUT")
	s.router.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}", s.validateAuthorization("serviceaccounts", "DeleteServiceAccount", s.withServiceAccount(s.deleteServiceAccount))).Methods("DELETE")

	s.router.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys", s.validateAuthorization("serviceaccountaccesskeys", "CreateServiceAccountAccessKey", s.createServiceAccountAccessKey)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys/{serviceaccountaccesskey}", s.validateAuthorization("serviceaccountaccesskeys", "GetServiceAccountAccessKey", s.getServiceAccountAccessKey)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys", s.validateAuthorization("serviceaccountsaccesskeys", "ListServiceAccountAccessKeys", s.listServiceAccountAccessKeys)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys/{serviceaccountaccesskey}", s.validateAuthorization("serviceaccountaccesskeys", "DeleteServiceAccountAccessKey", s.deleteServiceAccountAccessKey)).Methods("DELETE")

	s.router.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/roles/{role}/serviceaccountrolebindings", s.validateAuthorization("serviceaccountrolebindings", "CreateServiceAccountRoleBinding", s.withRole(s.createServiceAccountRoleBinding))).Methods("POST")
	s.router.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/roles/{role}/serviceaccountrolebindings", s.validateAuthorization("serviceaccountrolebindings", "GetServiceAccountRoleBinding", s.withRole(s.getServiceAccountRoleBinding))).Methods("GET")
	s.router.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountrolebindings", s.validateAuthorization("serviceaccountrolebindings", "ListMembershipRoleBindings", s.listServiceAccountRoleBindings)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/roles/{role}/serviceaccountrolebindings", s.validateAuthorization("serviceaccountrolebindings", "DeleteServiceAccountRoleBinding", s.withRole(s.deleteServiceAccountRoleBinding))).Methods("DELETE")

	s.router.HandleFunc("/projects/{project}/applications", s.validateAuthorization("applications", "CreateApplication", s.createApplication)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/applications/{application}", s.validateAuthorization("applications", "GetApplication", s.withApplication(s.getApplication))).Methods("GET")
	s.router.HandleFunc("/projects/{project}/applications", s.validateAuthorization("applications", "ListApplications", s.listApplications)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/applications/{application}", s.validateAuthorization("applications", "UpdateApplication", s.withApplication(s.updateApplication))).Methods("PUT")
	s.router.HandleFunc("/projects/{project}/applications/{application}", s.validateAuthorization("applications", "DeleteApplication", s.withApplication(s.deleteApplication))).Methods("DELETE")

	s.router.HandleFunc("/projects/{project}/applications/{application}/releases", s.validateAuthorization("releases", "CreateRelease", s.withApplication(s.createRelease))).Methods("POST")
	s.router.HandleFunc("/projects/{project}/applications/{application}/releases/latest", s.validateAuthorization("releases", "GetLatestRelease", s.withApplication(s.getLatestRelease))).Methods("GET")
	s.router.HandleFunc("/projects/{project}/applications/{application}/releases/{release}", s.validateAuthorization("releases", "GetRelease", s.withApplication(s.getRelease))).Methods("GET")
	s.router.HandleFunc("/projects/{project}/applications/{application}/releases", s.validateAuthorization("releases", "ListReleases", s.withApplication(s.listReleases))).Methods("GET")

	s.router.HandleFunc("/projects/{project}/devices/{device}", s.validateAuthorization("devices", "GetDevice", s.getDevice)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/devices", s.validateAuthorization("devices", "ListDevices", s.listDevices)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/devices/{device}/labels", s.validateAuthorization("devicelabels", "SetDeviceLabel", s.setDeviceLabel)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/devices/{device}/labels/{key}", s.validateAuthorization("devicelabels", "GetDeviceLabel", s.getDeviceLabel)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/devices/{device}/labels", s.validateAuthorization("devicelabels", "ListDeviceLabels", s.listDeviceLabels)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/devices/{device}/labels/{key}", s.validateAuthorization("devicelabels", "DeleteDeviceLabel", s.deleteDeviceLabel)).Methods("DELETE")

	s.router.HandleFunc("/projects/{project}/deviceregistrationtokens", s.validateAuthorization("deviceregistrationtokens", "CreateDeviceRegistrationToken", s.createDeviceRegistrationToken)).Methods("POST")

	s.router.HandleFunc("/projects/{project}/devices/register", s.registerDevice).Methods("POST")
	s.router.HandleFunc("/projects/{project}/devices/{device}/bundle", s.withDeviceAuth(s.getBundle)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/devices/{device}/info", s.withDeviceAuth(s.setDeviceInfo)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/deviceapplicationstatuses", s.withDeviceAuth(s.setDeviceApplicationStatus)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/deviceservicestatuses", s.withDeviceAuth(s.setDeviceServiceStatus)).Methods("POST")

	return s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Service) withUserOrServiceAccountAuth(handler func(http.ResponseWriter, *http.Request, string, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

			if !user.RegistrationCompleted {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			handler(w, r, userID, "")
		} else if serviceAccountAccessKey != nil {
			serviceAccount, err := s.serviceAccounts.GetServiceAccount(r.Context(),
				serviceAccountAccessKey.ServiceAccountID, serviceAccountAccessKey.ProjectID)
			if err != nil {
				log.WithError(err).Error("get service account")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			handler(w, r, "", serviceAccount.ID)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (s *Service) validateAuthorization(requestedResource, requestedAction string, handler func(http.ResponseWriter, *http.Request, string, string)) func(http.ResponseWriter, *http.Request) {
	return s.withUserOrServiceAccountAuth(func(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
		vars := mux.Vars(r)
		project := vars["project"]
		if project == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: move this logic into a middleware function
		var projectID string
		if strings.Contains(project, "_") {
			projectID = project
		} else {
			project, err := s.projects.LookupProject(r.Context(), project)
			if err == store.ErrProjectNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup project")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			projectID = project.ID
		}

		var roles []string
		if authenticatedUserID != "" {
			if _, err := s.memberships.GetMembership(r.Context(),
				authenticatedUserID, projectID,
			); err == store.ErrMembershipNotFound {
				w.WriteHeader(http.StatusForbidden)
				return
			} else if err != nil {
				// TODO: better logging all around
				log.WithField("user_id", authenticatedUserID).
					WithField("project_id", projectID).
					WithError(err).
					Error("get membership")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			roleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(),
				authenticatedUserID, projectID)
			if err != nil {
				log.WithError(err).Error("list membership role bindings")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			for _, roleBinding := range roleBindings {
				roles = append(roles, roleBinding.RoleID)
			}
		} else if authenticatedServiceAccountID != "" {
			// Sanity check that this service account belongs to this project
			if _, err := s.serviceAccounts.GetServiceAccount(r.Context(),
				authenticatedServiceAccountID, projectID,
			); err == store.ErrServiceAccountNotFound {
				w.WriteHeader(http.StatusForbidden)
				return
			} else if err != nil {
				log.WithError(err).Error("get service account")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			roleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(),
				authenticatedServiceAccountID, projectID)
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
		for _, roleID := range roles {
			role, err := s.roles.GetRole(r.Context(), roleID, projectID)
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

		accessGranted, err := authz.Evaluate(requestedResource, requestedAction, configs)
		if err != nil {
			log.WithError(err).Error("evaluate authz")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !accessGranted {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		handler(w, r, projectID, authenticatedUserID)
	})
}

func (s *Service) register(w http.ResponseWriter, r *http.Request) {
	var registerRequest struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := s.users.CreateUser(r.Context(), registerRequest.Email, hash.Hash(registerRequest.Password),
		registerRequest.FirstName, registerRequest.LastName)
	if err != nil {
		log.WithError(err).Error("create user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	registrationTokenValue := ksuid.New().String()

	if _, err := s.registrationTokens.CreateRegistrationToken(r.Context(), user.ID, hash.Hash(registrationTokenValue)); err != nil {
		log.WithError(err).Error("create registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	name := user.FirstName + " " + user.LastName

	if err := s.email.Send(email.Request{
		FromName:         "Device Plane",
		FromAddress:      "noreply@deviceplane.io",
		ToName:           name,
		ToAddress:        user.Email,
		Subject:          "Device Plane Registration Confirmation",
		PlainTextContent: "Please go to the following URL to complete registration. https://app.deviceplane.io/confirm/" + registrationTokenValue,
		HTMLContent:      "Please go to the following URL to complete registration. https://app.deviceplane.io/confirm/" + registrationTokenValue,
	}); err != nil {
		log.WithError(err).Error("send registration email")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) confirmRegistration(w http.ResponseWriter, r *http.Request) {
	var confirmRegistrationRequest struct {
		RegistrationTokenValue string `json:"registrationTokenValue"`
	}
	if err := json.NewDecoder(r.Body).Decode(&confirmRegistrationRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	registrationToken, err := s.registrationTokens.ValidateRegistrationToken(r.Context(),
		hash.Hash(confirmRegistrationRequest.RegistrationTokenValue))
	if err != nil {
		log.WithError(err).Error("validate registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := s.users.MarkRegistrationCompleted(r.Context(), registrationToken.UserID); err != nil {
		log.WithError(err).Error("mark registration completed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.newSession(w, r, registrationToken.UserID)
}

func (s *Service) login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := s.users.ValidateUser(r.Context(), loginRequest.Email, hash.Hash(loginRequest.Password))
	if err == store.ErrUserNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("validate user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !user.RegistrationCompleted {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	s.newSession(w, r, user.ID)
}

func (s *Service) newSession(w http.ResponseWriter, r *http.Request, userID string) {
	sessionValue := ksuid.New().String()

	if _, err := s.sessions.CreateSession(r.Context(), userID, hash.Hash(sessionValue)); err != nil {
		log.WithError(err).Error("create session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  sessionCookie,
		Value: sessionValue,

		Domain:  s.cookieDomain,
		Expires: time.Now().AddDate(0, 1, 0),

		Secure:   s.cookieSecure,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}

func (s *Service) logout(w http.ResponseWriter, r *http.Request) {
	sessionValue, err := r.Cookie(sessionCookie)

	switch err {
	case nil:
		session, err := s.sessions.ValidateSession(r.Context(), hash.Hash(sessionValue.Value))
		if err == store.ErrSessionNotFound {
			w.WriteHeader(http.StatusForbidden)
			return
		} else if err != nil {
			log.WithError(err).Error("validate session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.sessions.DeleteSession(r.Context(), session.ID); err != nil {
			log.WithError(err).Error("delete session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.ErrNoCookie:
		w.WriteHeader(http.StatusOK)
		return
	default:
		log.WithError(err).Error("get session cookie")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) me(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	user, err := s.users.GetUser(r.Context(), authenticatedUserID)
	if err != nil {
		log.WithError(err).Error("get user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *Service) listMembershipsByUser(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	vars := mux.Vars(r)
	userID := vars["user"]

	if userID != authenticatedUserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	memberships, err := s.memberships.ListMembershipsByUser(r.Context(), userID)
	if err != nil {
		log.WithError(err).Error("list memberships by user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(memberships)
}

func (s *Service) listMembershipsByUserFull(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	vars := mux.Vars(r)
	userID := vars["user"]

	if userID != authenticatedUserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	memberships, err := s.memberships.ListMembershipsByUser(r.Context(), userID)
	if err != nil {
		log.WithError(err).Error("list memberships by user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var membershipsFull []models.MembershipFull1

	for _, membership := range memberships {
		user, err := s.users.GetUser(r.Context(), membership.UserID)
		if err != nil {
			log.WithError(err).Error("get user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		project, err := s.projects.GetProject(r.Context(), membership.ProjectID)
		if err != nil {
			log.WithError(err).Error("get project")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		projectDeviceCounts, err := s.projectDeviceCounts.GetProjectDeviceCounts(r.Context(), membership.ProjectID)
		if err != nil {
			log.WithError(err).Error("get project device counts")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		projectApplicationCounts, err := s.projectApplicationCounts.GetProjectApplicationCounts(r.Context(), membership.ProjectID)
		if err != nil {
			log.WithError(err).Error("get project application counts")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		membershipsFull = append(membershipsFull, models.MembershipFull1{
			Membership: membership,
			User:       *user,
			Project: models.ProjectFull{
				Project:           *project,
				DeviceCounts:      *projectDeviceCounts,
				ApplicationCounts: *projectApplicationCounts,
			},
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(membershipsFull)
}

func (s *Service) createProject(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	var createProjectRequest struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&createProjectRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	project, err := s.projects.CreateProject(r.Context(), createProjectRequest.Name)
	if err != nil {
		log.WithError(err).Error("create project")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = s.memberships.CreateMembership(r.Context(), authenticatedUserID, project.ID); err != nil {
		log.WithError(err).Error("create membership")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	adminAllRoleBytes, err := yaml.Marshal(authz.AdminAllRole)
	if err != nil {
		log.WithError(err).Error("marshal admin role config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	adminRole, err := s.roles.CreateRole(r.Context(), project.ID, "admin-all", "", string(adminAllRoleBytes))
	if err != nil {
		log.WithError(err).Error("create admin role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeAllRoleBytes, err := yaml.Marshal(authz.WriteAllRole)
	if err != nil {
		log.WithError(err).Error("marshal write role config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = s.roles.CreateRole(r.Context(), project.ID, "write-all", "", string(writeAllRoleBytes)); err != nil {
		log.WithError(err).Error("create write role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	readAllRoleBytes, err := yaml.Marshal(authz.ReadAllRole)
	if err != nil {
		log.WithError(err).Error("marshal read role config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = s.roles.CreateRole(r.Context(), project.ID, "read-all", "", string(readAllRoleBytes)); err != nil {
		log.WithError(err).Error("create read role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = s.membershipRoleBindings.CreateMembershipRoleBinding(r.Context(),
		authenticatedUserID, adminRole.ID, project.ID,
	); err != nil {
		log.WithError(err).Error("create membership role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

func (s *Service) getProject(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	project, err := s.projects.GetProject(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("get project")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

func (s *Service) createRole(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	var createRoleRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Config      string `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&createRoleRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := s.roles.CreateRole(r.Context(), projectID, createRoleRequest.Name,
		createRoleRequest.Description, createRoleRequest.Config)
	if err != nil {
		log.WithError(err).Error("create role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(role)
}

func (s *Service) getRole(w http.ResponseWriter, r *http.Request, projectID, userID, roleID string) {
	role, err := s.roles.GetRole(r.Context(), roleID, projectID)
	if err == store.ErrRoleNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(role)
}

func (s *Service) listRoles(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	roles, err := s.roles.ListRoles(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list roles")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roles)
}

func (s *Service) updateRole(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	roleID := vars["role"]

	var updateRoleRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Config      string `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateRoleRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := s.roles.UpdateRole(r.Context(), roleID, projectID, updateRoleRequest.Name,
		updateRoleRequest.Description, updateRoleRequest.Config)
	if err != nil {
		log.WithError(err).Error("update role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(role)
}

func (s *Service) deleteRole(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	roleID := vars["role"]

	if err := s.roles.DeleteRole(r.Context(), roleID, projectID); err != nil {
		log.WithError(err).Error("delete role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) createServiceAccount(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	var createServiceAccountRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&createServiceAccountRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	serviceAccount, err := s.serviceAccounts.CreateServiceAccount(r.Context(), projectID, createServiceAccountRequest.Name,
		createServiceAccountRequest.Description)
	if err != nil {
		log.WithError(err).Error("create service account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceAccount)
}

func (s *Service) getServiceAccount(w http.ResponseWriter, r *http.Request, projectID, userID, serviceAccountID string) {
	serviceAccount, err := s.serviceAccounts.GetServiceAccount(r.Context(), serviceAccountID, projectID)
	if err == store.ErrServiceAccountNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get service account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = serviceAccount
	if _, ok := r.URL.Query()["full"]; ok {
		serviceAccountRoleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(), serviceAccountID, projectID)
		if err != nil {
			log.WithError(err).Error("list service account role bindings")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		roles := make([]models.Role, 0)
		for _, serviceAccountRoleBinding := range serviceAccountRoleBindings {
			role, err := s.roles.GetRole(r.Context(), serviceAccountRoleBinding.RoleID, projectID)
			if err != nil {
				log.WithError(err).Error("get role")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			roles = append(roles, *role)
		}

		ret = models.ServiceAccountFull{
			ServiceAccount: *serviceAccount,
			Roles:          roles,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (s *Service) listServiceAccounts(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	serviceAccounts, err := s.serviceAccounts.ListServiceAccounts(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list service accounts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = serviceAccounts
	if _, ok := r.URL.Query()["full"]; ok {
		serviceAccountsFull := make([]models.ServiceAccountFull, 0)

		for _, serviceAccount := range serviceAccounts {
			serviceAccountRoleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(), serviceAccount.ID, projectID)
			if err != nil {
				log.WithError(err).Error("list service account role bindings")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			roles := make([]models.Role, 0)
			for _, serviceAccountRoleBinding := range serviceAccountRoleBindings {
				role, err := s.roles.GetRole(r.Context(), serviceAccountRoleBinding.RoleID, projectID)
				if err != nil {
					log.WithError(err).Error("get role")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				roles = append(roles, *role)
			}

			serviceAccountsFull = append(serviceAccountsFull, models.ServiceAccountFull{
				ServiceAccount: serviceAccount,
				Roles:          roles,
			})
		}

		ret = serviceAccountsFull
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (s *Service) updateServiceAccount(w http.ResponseWriter, r *http.Request, projectID, userID, serviceAccountID string) {
	var updateServiceAccountRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateServiceAccountRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	serviceAccount, err := s.serviceAccounts.UpdateServiceAccount(r.Context(), serviceAccountID, projectID,
		updateServiceAccountRequest.Name, updateServiceAccountRequest.Description)
	if err != nil {
		log.WithError(err).Error("update service account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceAccount)
}

func (s *Service) deleteServiceAccount(w http.ResponseWriter, r *http.Request, projectID, userID, serviceAccountID string) {
	if err := s.serviceAccounts.DeleteServiceAccount(r.Context(), serviceAccountID, projectID); err != nil {
		log.WithError(err).Error("delete service account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) createServiceAccountAccessKey(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	serviceAccountID := vars["serviceaccount"]

	serviceAccountAccessKeyValue := "s" + ksuid.New().String()

	serviceAccount, err := s.serviceAccountAccessKeys.CreateServiceAccountAccessKey(r.Context(),
		projectID, serviceAccountID, serviceAccountAccessKeyValue)
	if err != nil {
		log.WithError(err).Error("create service account access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceAccount)
}

func (s *Service) getServiceAccountAccessKey(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	serviceAccountAccessKeyID := vars["serviceaccountaccesskey"]

	serviceAccountAccessKey, err := s.serviceAccountAccessKeys.GetServiceAccountAccessKey(r.Context(), serviceAccountAccessKeyID, projectID)
	if err == store.ErrServiceAccountAccessKeyNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get service account access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceAccountAccessKey)
}

func (s *Service) listServiceAccountAccessKeys(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	serviceAccountID := vars["serviceaccount"]

	serviceAccountAccessKeys, err := s.serviceAccountAccessKeys.ListServiceAccountAccessKeys(r.Context(), projectID, serviceAccountID)
	if err != nil {
		log.WithError(err).Error("list service accounts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceAccountAccessKeys)
}

func (s *Service) deleteServiceAccountAccessKey(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	serviceAccountAccessKeyID := vars["serviceaccountaccesskey"]

	if err := s.serviceAccountAccessKeys.DeleteServiceAccountAccessKey(r.Context(), serviceAccountAccessKeyID, projectID); err != nil {
		log.WithError(err).Error("delete service account access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) createServiceAccountRoleBinding(w http.ResponseWriter, r *http.Request, projectID, userID, roleID string) {
	vars := mux.Vars(r)
	serviceAccountID := vars["serviceaccount"]

	serviceAccountRoleBinding, err := s.serviceAccountRoleBindings.CreateServiceAccountRoleBinding(r.Context(), serviceAccountID, roleID, projectID)
	if err != nil {
		log.WithError(err).Error("create service account role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceAccountRoleBinding)
}

func (s *Service) getServiceAccountRoleBinding(w http.ResponseWriter, r *http.Request, projectID, userID, roleID string) {
	vars := mux.Vars(r)
	serviceAccountID := vars["serviceaccount"]

	serviceAccountRoleBinding, err := s.serviceAccountRoleBindings.GetServiceAccountRoleBinding(r.Context(), serviceAccountID, roleID, projectID)
	if err != nil {
		log.WithError(err).Error("get service account role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceAccountRoleBinding)
}

func (s *Service) listServiceAccountRoleBindings(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	serviceAccountID := vars["serviceaccount"]

	serviceAccountRoleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(), projectID, serviceAccountID)
	if err != nil {
		log.WithError(err).Error("list service account role bindings")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceAccountRoleBindings)
}

func (s *Service) deleteServiceAccountRoleBinding(w http.ResponseWriter, r *http.Request, projectID, userID, roleID string) {
	vars := mux.Vars(r)
	serviceAccountID := vars["serviceaccount"]

	if err := s.serviceAccountRoleBindings.DeleteServiceAccountRoleBinding(r.Context(), serviceAccountID, roleID, projectID); err != nil {
		log.WithError(err).Error("delete service account role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) createMembership(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	var createMembershipRequest struct {
		UserID string `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&createMembershipRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	membership, err := s.memberships.CreateMembership(r.Context(),
		createMembershipRequest.UserID, projectID)
	if err != nil {
		log.WithError(err).Error("create membership")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(membership)
}

func (s *Service) getMembership(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	// TODO: rename this to userID and change all instances of the other to authenticatedUserUD
	membershipUserID := vars["user"]

	membership, err := s.memberships.GetMembership(r.Context(), membershipUserID, projectID)
	if err == store.ErrMembershipNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get membership")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = membership
	if _, ok := r.URL.Query()["full"]; ok {
		user, err := s.users.GetUser(r.Context(), membership.UserID)
		if err != nil {
			log.WithError(err).Error("get user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(), membershipUserID, projectID)
		if err != nil {
			log.WithError(err).Error("list membership role bindings")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		roles := make([]models.Role, 0)
		for _, membershipRoleBinding := range membershipRoleBindings {
			role, err := s.roles.GetRole(r.Context(), membershipRoleBinding.RoleID, projectID)
			if err != nil {
				log.WithError(err).Error("get role")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			roles = append(roles, *role)
		}

		ret = models.MembershipFull2{
			Membership: *membership,
			User:       *user,
			Roles:      roles,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (s *Service) listMembershipsByProject(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	memberships, err := s.memberships.ListMembershipsByProject(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list memberships by project")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = memberships
	if _, ok := r.URL.Query()["full"]; ok {
		var membershipsFull []models.MembershipFull2

		for _, membership := range memberships {
			user, err := s.users.GetUser(r.Context(), membership.UserID)
			if err != nil {
				log.WithError(err).Error("get user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(), membership.UserID, projectID)
			if err != nil {
				log.WithError(err).Error("list membership role bindings")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			roles := make([]models.Role, 0)
			for _, membershipRoleBinding := range membershipRoleBindings {
				role, err := s.roles.GetRole(r.Context(), membershipRoleBinding.RoleID, projectID)
				if err != nil {
					log.WithError(err).Error("get role")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				roles = append(roles, *role)
			}

			membershipsFull = append(membershipsFull, models.MembershipFull2{
				Membership: membership,
				User:       *user,
				Roles:      roles,
			})
		}

		ret = membershipsFull
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (s *Service) createMembershipRoleBinding(w http.ResponseWriter, r *http.Request, projectID, userID, roleID string) {
	vars := mux.Vars(r)
	// TODO: rename this to userID and change all instances of the other to authenticatedUserUD
	membershipUserID := vars["user"]

	membershipRoleBinding, err := s.membershipRoleBindings.CreateMembershipRoleBinding(r.Context(), membershipUserID, roleID, projectID)
	if err != nil {
		log.WithError(err).Error("create membership role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(membershipRoleBinding)
}

func (s *Service) getMembershipRoleBinding(w http.ResponseWriter, r *http.Request, projectID, userID, roleID string) {
	vars := mux.Vars(r)
	// TODO: rename this to userID and change all instances of the other to authenticatedUserUD
	membershipUserID := vars["user"]

	membershipRoleBinding, err := s.membershipRoleBindings.GetMembershipRoleBinding(r.Context(), membershipUserID, roleID, projectID)
	if err == store.ErrMembershipRoleBindingNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get membership role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(membershipRoleBinding)
}

func (s *Service) listMembershipRoleBindings(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	// TODO: rename this to userID and change all instances of the other to authenticatedUserUD
	membershipUserID := vars["user"]

	membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(), membershipUserID, projectID)
	if err != nil {
		log.WithError(err).Error("list membership role bindings")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(membershipRoleBindings)
}

func (s *Service) deleteMembershipRoleBinding(w http.ResponseWriter, r *http.Request, projectID, userID, roleID string) {
	vars := mux.Vars(r)
	// TODO: rename this to userID and change all instances of the other to authenticatedUserUD
	membershipUserID := vars["user"]

	if err := s.membershipRoleBindings.DeleteMembershipRoleBinding(r.Context(), membershipUserID, roleID, projectID); err != nil {
		log.WithError(err).Error("delete membership role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) createApplication(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	var createApplicationRequest struct {
		Name        string                     `json:"name"`
		Description string                     `json:"description"`
		Settings    models.ApplicationSettings `json:"settings"`
	}
	if err := json.NewDecoder(r.Body).Decode(&createApplicationRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	application, err := s.applications.CreateApplication(r.Context(), projectID, createApplicationRequest.Name,
		createApplicationRequest.Description, createApplicationRequest.Settings)
	if err != nil {
		log.WithError(err).Error("create application")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(application)
}

func (s *Service) getApplication(w http.ResponseWriter, r *http.Request, projectID, userID, applicationID string) {
	application, err := s.applications.GetApplication(r.Context(), applicationID, projectID)
	if err == store.ErrApplicationNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get application")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = application
	if _, ok := r.URL.Query()["full"]; ok {
		latestRelease, err := s.releases.GetLatestRelease(r.Context(), projectID, applicationID)
		if err != nil && err != store.ErrReleaseNotFound {
			log.WithError(err).Error("get latest release")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		applicationDeviceCounts, err := s.applicationDeviceCounts.GetApplicationDeviceCounts(r.Context(), projectID, applicationID)
		if err != nil {
			log.WithError(err).Error("get application device counts")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ret = models.ApplicationFull{
			Application:   *application,
			LatestRelease: latestRelease,
			DeviceCounts:  *applicationDeviceCounts,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (s *Service) listApplications(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	applications, err := s.applications.ListApplications(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list applications")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = applications
	if _, ok := r.URL.Query()["full"]; ok {
		var applicationsFull []models.ApplicationFull

		for _, application := range applications {
			latestRelease, err := s.releases.GetLatestRelease(r.Context(), projectID, application.ID)
			if err != nil && err != store.ErrReleaseNotFound {
				log.WithError(err).Error("get latest release")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			applicationDeviceCounts, err := s.applicationDeviceCounts.GetApplicationDeviceCounts(r.Context(), projectID, application.ID)
			if err != nil {
				log.WithError(err).Error("get application device counts")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			applicationsFull = append(applicationsFull, models.ApplicationFull{
				Application:   application,
				LatestRelease: latestRelease,
				DeviceCounts:  *applicationDeviceCounts,
			})
		}

		ret = applicationsFull
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (s *Service) updateApplication(w http.ResponseWriter, r *http.Request, projectID, userID, applicationID string) {
	var updateApplicationRequest struct {
		Name        string                     `json:"name"`
		Description string                     `json:"description"`
		Settings    models.ApplicationSettings `json:"settings"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateApplicationRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	application, err := s.applications.UpdateApplication(r.Context(), applicationID, projectID, updateApplicationRequest.Name,
		updateApplicationRequest.Description, updateApplicationRequest.Settings)
	if err != nil {
		log.WithError(err).Error("update application")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(application)
}

func (s *Service) deleteApplication(w http.ResponseWriter, r *http.Request, projectID, userID, applicationID string) {
	if err := s.applications.DeleteApplication(r.Context(), applicationID, projectID); err != nil {
		log.WithError(err).Error("delete application")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// TOOD: this has a vulnerability!
func (s *Service) createRelease(w http.ResponseWriter, r *http.Request, projectID, userID, applicationID string) {
	var createReleaseRequest struct {
		Config string `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&createReleaseRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	release, err := s.releases.CreateRelease(r.Context(), projectID, applicationID, createReleaseRequest.Config)
	if err != nil {
		log.WithError(err).Error("create release")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(release)
}

func (s *Service) getRelease(w http.ResponseWriter, r *http.Request, projectID, userID, applicationID string) {
	vars := mux.Vars(r)
	releaseID := vars["release"]

	release, err := s.releases.GetRelease(r.Context(), releaseID, projectID, applicationID)
	if err == store.ErrReleaseNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get release")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = release
	if _, ok := r.URL.Query()["full"]; ok {
		releaseDeviceCounts, err := s.releaseDeviceCounts.GetReleaseDeviceCounts(r.Context(), projectID, applicationID, releaseID)
		if err != nil {
			log.WithError(err).Error("get release device counts")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ret = models.ReleaseFull{
			Release:      *release,
			DeviceCounts: *releaseDeviceCounts,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (s *Service) getLatestRelease(w http.ResponseWriter, r *http.Request, projectID string, userID, applicationID string) {
	release, err := s.releases.GetLatestRelease(r.Context(), projectID, applicationID)
	if err == store.ErrReleaseNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get latest release")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = release
	if _, ok := r.URL.Query()["full"]; ok {
		releaseDeviceCounts, err := s.releaseDeviceCounts.GetReleaseDeviceCounts(r.Context(), projectID, applicationID, release.ID)
		if err != nil {
			log.WithError(err).Error("get release device counts")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ret = models.ReleaseFull{
			Release:      *release,
			DeviceCounts: *releaseDeviceCounts,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (s *Service) listReleases(w http.ResponseWriter, r *http.Request, projectID, userID, applicationID string) {
	releases, err := s.releases.ListReleases(r.Context(), projectID, applicationID)
	if err != nil {
		log.WithError(err).Error("list releases")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = releases
	if _, ok := r.URL.Query()["full"]; ok {
		var releasesFull []models.ReleaseFull

		for _, release := range releases {
			releaseDeviceCounts, err := s.releaseDeviceCounts.GetReleaseDeviceCounts(r.Context(), projectID, applicationID, release.ID)
			if err != nil {
				log.WithError(err).Error("get release device counts")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			releasesFull = append(releasesFull, models.ReleaseFull{
				Release:      release,
				DeviceCounts: *releaseDeviceCounts,
			})
		}

		ret = releasesFull
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (s *Service) listDevices(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	devices, err := s.devices.ListDevices(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list devices")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = devices
	if _, ok := r.URL.Query()["full"]; ok {
		var deviceIDs []string
		for _, device := range devices {
			deviceIDs = append(deviceIDs, device.ID)
		}

		deviceStatuses, err := s.deviceStatuses.GetDeviceStatuses(r.Context(), deviceIDs)
		if err != nil {
			log.WithError(err).Error("get device statuses")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var devicesFull []models.DeviceFull1
		for i, device := range devices {
			devicesFull = append(devicesFull, models.DeviceFull1{
				Device: device,
				Status: deviceStatuses[i],
			})
		}

		ret = devicesFull
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (s *Service) getDevice(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	deviceID := vars["device"]

	device, err := s.devices.GetDevice(r.Context(), deviceID, projectID)
	if err != nil {
		log.WithError(err).Error("get device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = device
	if _, ok := r.URL.Query()["full"]; ok {
		deviceStatus, err := s.deviceStatuses.GetDeviceStatus(r.Context(), device.ID)
		if err != nil {
			log.WithError(err).Error("get device status")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		applications, err := s.applications.ListApplications(r.Context(), projectID)
		if err != nil {
			log.WithError(err).Error("list applications")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var allApplicationStatusInfo []models.DeviceApplicationStatusInfo
		for _, application := range applications {
			applicationStatusInfo := models.DeviceApplicationStatusInfo{
				Application: application,
			}

			deviceApplicationStatus, err := s.deviceApplicationStatuses.GetDeviceApplicationStatus(
				r.Context(), projectID, device.ID, application.ID)
			if err == nil {
				applicationStatusInfo.ApplicationStatus = deviceApplicationStatus
			} else if err != store.ErrDeviceApplicationStatusNotFound {
				log.WithError(err).Error("get device application status")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			deviceServiceStatuses, err := s.deviceServiceStatuses.GetDeviceServiceStatuses(
				r.Context(), projectID, device.ID, application.ID)
			if err == nil {
				applicationStatusInfo.ServiceStatuses = deviceServiceStatuses
			} else if err != store.ErrDeviceServiceStatusNotFound {
				log.WithError(err).Error("get device service statuses")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			allApplicationStatusInfo = append(allApplicationStatusInfo, applicationStatusInfo)
		}

		ret = models.DeviceFull2{
			Device:                *device,
			Status:                deviceStatus,
			ApplicationStatusInfo: allApplicationStatusInfo,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (s *Service) setDeviceLabel(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	deviceID := vars["device"]

	var setDeviceLabelRequest struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&setDeviceLabelRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deviceLabel, err := s.deviceLabels.SetDeviceLabel(r.Context(), setDeviceLabelRequest.Key,
		deviceID, projectID, setDeviceLabelRequest.Value)
	if err != nil {
		log.WithError(err).Error("set device label")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceLabel)
}

func (s *Service) getDeviceLabel(w http.ResponseWriter, r *http.Request, projectID string, userID string) {
	vars := mux.Vars(r)
	deviceID := vars["device"]
	key := vars["key"]

	deviceLabel, err := s.deviceLabels.GetDeviceLabel(r.Context(), key, deviceID, projectID)
	if err != nil {
		log.WithError(err).Error("get device label")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceLabel)
}

func (s *Service) listDeviceLabels(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	deviceID := vars["device"]

	deviceLabels, err := s.deviceLabels.ListDeviceLabels(r.Context(), deviceID, projectID)
	if err != nil {
		log.WithError(err).Error("list device labels")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceLabels)
}

func (s *Service) deleteDeviceLabel(w http.ResponseWriter, r *http.Request, projectID string, userID string) {
	vars := mux.Vars(r)
	deviceID := vars["device"]
	key := vars["key"]

	if err := s.deviceLabels.DeleteDeviceLabel(r.Context(), key, deviceID, projectID); err != nil {
		log.WithError(err).Error("delete device label")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) createDeviceRegistrationToken(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	deviceRegistrationToken, err := s.deviceRegistrationTokens.CreateDeviceRegistrationToken(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("create device registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceRegistrationToken)
}

func (s *Service) withDeviceAuth(handler func(http.ResponseWriter, *http.Request, string, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		handler(w, r, projectID, deviceAccessKey.DeviceID)
	}
}

// TODO: verify project ID
func (s *Service) registerDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project"]

	var registerDeviceRequest models.RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&registerDeviceRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deviceRegistrationToken, err := s.deviceRegistrationTokens.GetDeviceRegistrationToken(r.Context(), registerDeviceRequest.DeviceRegistrationTokenID, projectID)
	if err != nil {
		log.WithError(err).Error("get device registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if deviceRegistrationToken.DeviceAccessKeyID != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	device, err := s.devices.CreateDevice(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("create device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deviceAccessKeyValue := ksuid.New().String()

	deviceAccessKey, err := s.deviceAccessKeys.CreateDeviceAccessKey(r.Context(), projectID, device.ID, hash.Hash(deviceAccessKeyValue))
	if err != nil {
		log.WithError(err).Error("create device access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = s.deviceRegistrationTokens.BindDeviceRegistrationToken(r.Context(), deviceRegistrationToken.ID, projectID, deviceAccessKey.ID); err != nil {
		log.WithError(err).Error("bind device registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.RegisterDeviceResponse{
		DeviceID:             device.ID,
		DeviceAccessKeyValue: deviceAccessKeyValue,
	})
}

func (s *Service) getBundle(w http.ResponseWriter, r *http.Request, projectID, deviceID string) {
	if err := s.deviceStatuses.ResetDeviceStatus(r.Context(), deviceID, time.Minute); err != nil {
		log.WithError(err).Error("reset device status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	applications, err := s.applications.ListApplications(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list applications")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deviceLabels, err := s.deviceLabels.ListDeviceLabels(r.Context(), deviceID, projectID)
	if err != nil {
		log.WithError(err).Error("list device labels")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var bundle models.Bundle
	for i, application := range applications {
		release, err := s.releases.GetLatestRelease(r.Context(), projectID, application.ID)
		if err == store.ErrReleaseNotFound {
			continue
		} else if err != nil {
			log.WithError(err).Error("get latest release")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if application.Settings.SchedulingRule != "" {
			expression, err := govaluate.NewEvaluableExpression(application.Settings.SchedulingRule)
			if err != nil {
				log.WithError(err).Error("parse application scheduling rule")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			result, err := expression.Eval(deviceLabelParameters(deviceLabels))
			if err != nil {
				log.WithError(err).Error("evaluate application scheduling rule")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			booleanResult, ok := result.(bool)
			if !ok {
				log.Error("invalid scheduling rule evaluation result")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !booleanResult {
				continue
			}
		}

		bundle.Applications = append(bundle.Applications, models.ApplicationAndLatestRelease{
			Application:   applications[i],
			LatestRelease: release,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bundle)
}

func (s *Service) setDeviceInfo(w http.ResponseWriter, r *http.Request, projectID, deviceID string) {
	var setDeviceInfoRequest models.SetDeviceInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&setDeviceInfoRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := s.devices.SetDeviceInfo(r.Context(), deviceID, projectID, setDeviceInfoRequest.DeviceInfo); err != nil {
		log.WithError(err).Error("set device info")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) setDeviceApplicationStatus(w http.ResponseWriter, r *http.Request, projectID, deviceID string) {
	vars := mux.Vars(r)
	applicationID := vars["application"]

	var setDeviceApplicationStatusRequest models.SetDeviceApplicationStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&setDeviceApplicationStatusRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.deviceApplicationStatuses.SetDeviceApplicationStatus(r.Context(), projectID, deviceID,
		applicationID, setDeviceApplicationStatusRequest.CurrentReleaseID,
	); err != nil {
		log.WithError(err).Error("set device application status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) setDeviceServiceStatus(w http.ResponseWriter, r *http.Request, projectID, deviceID string) {
	vars := mux.Vars(r)
	applicationID := vars["application"]
	service := vars["service"]

	var setDeviceServiceStatusRequest models.SetDeviceServiceStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&setDeviceServiceStatusRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.deviceServiceStatuses.SetDeviceServiceStatus(r.Context(), projectID, deviceID,
		applicationID, service, setDeviceServiceStatusRequest.CurrentReleaseID,
	); err != nil {
		log.WithError(err).Error("set device service status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
