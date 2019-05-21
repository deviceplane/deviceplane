package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
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
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/gorilla/mux"
)

const (
	sessionCookie = "dp_sess"
)

type Service struct {
	users                     store.Users
	registrationTokens        store.RegistrationTokens
	accessKeys                store.AccessKeys
	sessions                  store.Sessions
	projects                  store.Projects
	projectDeviceCounts       store.ProjectDeviceCounts
	projectApplicationCounts  store.ProjectApplicationCounts
	roles                     store.Roles
	memberships               store.Memberships
	membershipRoleBindings    store.MembershipRoleBindings
	devices                   store.Devices
	deviceStatuses            store.DeviceStatuses
	deviceLabels              store.DeviceLabels
	deviceRegistrationTokens  store.DeviceRegistrationTokens
	deviceAccessKeys          store.DeviceAccessKeys
	applications              store.Applications
	releases                  store.Releases
	releaseDeviceCounts       store.ReleaseDeviceCounts
	deviceApplicationStatuses store.DeviceApplicationStatuses
	deviceServiceStatuses     store.DeviceServiceStatuses
	email                     email.Interface
	router                    *mux.Router
	cookieDomain              string
	cookieSecure              bool
}

func NewService(
	users store.Users,
	registrationTokens store.RegistrationTokens,
	sessions store.Sessions,
	accessKeys store.AccessKeys,
	projects store.Projects,
	projectDeviceCounts store.ProjectDeviceCounts,
	projectApplicationCounts store.ProjectApplicationCounts,
	roles store.Roles,
	memberships store.Memberships,
	membershipRoleBindings store.MembershipRoleBindings,
	devices store.Devices,
	deviceStatuses store.DeviceStatuses,
	deviceLabels store.DeviceLabels,
	deviceRegistrationTokens store.DeviceRegistrationTokens,
	deviceAccessKeys store.DeviceAccessKeys,
	applications store.Applications,
	releases store.Releases,
	releasesDeviceCounts store.ReleaseDeviceCounts,
	deviceApplicationStatuses store.DeviceApplicationStatuses,
	deviceServiceStatuses store.DeviceServiceStatuses,
	email email.Interface,
	cookieDomain string,
	cookieSecure bool,
) *Service {
	s := &Service{
		users:                     users,
		registrationTokens:        registrationTokens,
		sessions:                  sessions,
		accessKeys:                accessKeys,
		projects:                  projects,
		projectDeviceCounts:       projectDeviceCounts,
		projectApplicationCounts:  projectApplicationCounts,
		roles:                     roles,
		memberships:               memberships,
		membershipRoleBindings:    membershipRoleBindings,
		devices:                   devices,
		deviceStatuses:            deviceStatuses,
		deviceLabels:              deviceLabels,
		deviceRegistrationTokens:  deviceRegistrationTokens,
		deviceAccessKeys:          deviceAccessKeys,
		applications:              applications,
		releases:                  releases,
		releaseDeviceCounts:       releasesDeviceCounts,
		deviceApplicationStatuses: deviceApplicationStatuses,
		deviceServiceStatuses:     deviceServiceStatuses,
		email:                     email,
		cookieDomain:              cookieDomain,
		cookieSecure:              cookieSecure,
		router:                    mux.NewRouter(),
	}

	s.router.HandleFunc("/health", s.health).Methods("GET")

	s.router.HandleFunc("/register", s.register).Methods("POST")
	s.router.HandleFunc("/completeregistration", s.confirmRegistration).Methods("POST")
	s.router.HandleFunc("/login", s.login).Methods("POST")
	s.router.HandleFunc("/logout", s.logout).Methods("POST")
	s.router.HandleFunc("/me", s.withUserAuth(s.me)).Methods("GET")

	s.router.HandleFunc("/users/{user}/memberships", s.withUserAuth(s.listMembershipsByUser)).Methods("GET")
	s.router.HandleFunc("/users/{user}/memberships/full", s.withUserAuth(s.listMembershipsByUserFull)).Methods("GET")

	s.router.HandleFunc("/projects", s.withUserAuth(s.createProject)).Methods("POST")
	s.router.HandleFunc("/projects/{project}", s.validateAuthorization("projects", "GetProject", s.getProject)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/roles", s.validateAuthorization("roles", "CreateRole", s.createRole)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/roles/{role}", s.validateAuthorization("roles", "GetRole", s.withRole(s.getRole))).Methods("GET")
	s.router.HandleFunc("/projects/{project}/roles", s.validateAuthorization("roles", "ListRoles", s.listRoles)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/memberships", s.validateAuthorization("memberships", "CreateMembership", s.createMembership)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/memberships", s.validateAuthorization("memberships", "ListMembershipsByProject", s.listMembershipsByProject)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/memberships/{membership}/roles/{role}/membershiprolebindings", s.validateAuthorization("membershiprolebindings", "CreateMembershipRoleBinding", s.withRole(s.createMembershipRoleBinding))).Methods("POST")
	s.router.HandleFunc("/projects/{project}/memberships/{membership}/roles/{role}/membershiprolebindings", s.validateAuthorization("membershiprolebindings", "GetMembershipRoleBinding", s.withRole(s.getMembershipRoleBinding))).Methods("GET")
	s.router.HandleFunc("/projects/{project}/memberships/{membership}/membershiprolebindings", s.validateAuthorization("membershiprolebindings", "ListMembershipRoleBindings", s.listMembershipRoleBindings)).Methods("GET")

	s.router.HandleFunc("/projects/{project}/applications", s.validateAuthorization("applications", "CreateApplication", s.createApplication)).Methods("POST")
	s.router.HandleFunc("/projects/{project}/applications/{application}", s.validateAuthorization("applications", "GetApplication", s.withApplication(s.getApplication))).Methods("GET")
	s.router.HandleFunc("/projects/{project}/applications", s.validateAuthorization("applications", "ListApplications", s.listApplications)).Methods("GET")
	s.router.HandleFunc("/projects/{project}/applications/{application}", s.validateAuthorization("applications", "UpdateApplication", s.withApplication(s.updateApplication))).Methods("PUT")

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

func (s *Service) withUserAuth(handler func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID string

		sessionValue, err := r.Cookie(sessionCookie)

		switch err {
		case nil:
			session, err := s.sessions.ValidateSession(r.Context(), hash(sessionValue.Value))
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

			accessKey, err := s.accessKeys.ValidateAccessKey(r.Context(), hash(accessKeyValue))
			if err == store.ErrAccessKeyNotFound {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if err != nil {
				log.WithError(err).Error("validate access key")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			userID = accessKey.UserID
		default:
			log.WithError(err).Error("get session cookie")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if userID == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		user, err := s.users.GetUser(r.Context(), userID)
		if err == store.ErrUserNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("get user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !user.RegistrationCompleted {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		handler(w, r, userID)
	}
}

func (s *Service) validateAuthorization(requestedResource, requestedAction string, handler func(http.ResponseWriter, *http.Request, string, string)) func(http.ResponseWriter, *http.Request) {
	return s.withUserAuth(func(w http.ResponseWriter, r *http.Request, userID string) {
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

		membership, err := s.memberships.GetMembership(r.Context(), userID, projectID)
		if err != nil {
			log.WithError(err).Error("get membership")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if membership == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(), membership.ID, projectID)
		if err != nil {
			log.WithError(err).Error("get membership role bindings")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var configs []authz.Config
		for _, membershipRoleBinding := range membershipRoleBindings {
			role, err := s.roles.GetRole(r.Context(), membershipRoleBinding.RoleID, projectID)
			if err == store.ErrRoleNotFound {
				continue
			} else if err != nil {
				log.WithError(err).Error("get role")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			configs = append(configs, role.Config)
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

		handler(w, r, projectID, userID)
	})
}

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

	user, err := s.users.CreateUser(r.Context(), registerRequest.Email, hash(registerRequest.Password),
		registerRequest.FirstName, registerRequest.LastName)
	if err != nil {
		log.WithError(err).Error("create user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	registrationTokenValue := ksuid.New().String()

	if _, err := s.registrationTokens.CreateRegistrationToken(r.Context(), user.ID, hash(registrationTokenValue)); err != nil {
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
		hash(confirmRegistrationRequest.RegistrationTokenValue))
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

	user, err := s.users.ValidateUser(r.Context(), loginRequest.Email, hash(loginRequest.Password))
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

	if _, err := s.sessions.CreateSession(r.Context(), userID, hash(sessionValue)); err != nil {
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
		session, err := s.sessions.ValidateSession(r.Context(), hash(sessionValue.Value))
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

func (s *Service) me(w http.ResponseWriter, r *http.Request, authenticatedUserID string) {
	user, err := s.users.GetUser(r.Context(), authenticatedUserID)
	if err != nil {
		log.WithError(err).Error("get user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *Service) listMembershipsByUser(w http.ResponseWriter, r *http.Request, authenticatedUserID string) {
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

func (s *Service) listMembershipsByUserFull(w http.ResponseWriter, r *http.Request, authenticatedUserID string) {
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
			User: *user,
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

func (s *Service) createProject(w http.ResponseWriter, r *http.Request, userID string) {
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

	membership, err := s.memberships.CreateMembership(r.Context(), userID, project.ID)
	if err != nil {
		log.WithError(err).Error("create membership")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	adminRole, err := s.roles.CreateRole(r.Context(), project.ID, "default", "", authz.AdminAllRole)
	if err != nil {
		log.WithError(err).Error("create role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = s.membershipRoleBindings.CreateMembershipRoleBinding(r.Context(),
		membership.ID, adminRole.ID, project.ID,
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

	var config authz.Config
	if err := yaml.Unmarshal([]byte(createRoleRequest.Config), &config); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := s.roles.CreateRole(r.Context(), projectID, createRoleRequest.Name,
		createRoleRequest.Description, config)
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

			membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(), membership.ID, projectID)
			if err != nil {
				log.WithError(err).Error("list membership role bindings")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var roles []models.Role
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
				User:  *user,
				Roles: roles,
			})
		}

		ret = membershipsFull
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (s *Service) createMembershipRoleBinding(w http.ResponseWriter, r *http.Request, projectID, userID, roleID string) {
	vars := mux.Vars(r)
	membershipID := vars["membership"]

	membershipRoleBinding, err := s.membershipRoleBindings.CreateMembershipRoleBinding(r.Context(), membershipID, roleID, projectID)
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
	membershipID := vars["membership"]

	membershipRoleBinding, err := s.membershipRoleBindings.GetMembershipRoleBinding(r.Context(), membershipID, roleID, projectID)
	if err != nil {
		log.WithError(err).Error("get membership role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(membershipRoleBinding)
}

func (s *Service) listMembershipRoleBindings(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	vars := mux.Vars(r)
	membershipID := vars["membership"]

	membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(), projectID, membershipID)
	if err != nil {
		log.WithError(err).Error("list membership role bindings")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(membershipRoleBindings)
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
	if err != nil {
		log.WithError(err).Error("get application")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(application)
}

func (s *Service) listApplications(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	applications, err := s.applications.ListApplications(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list applications")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(applications)
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

		deviceAccessKey, err := s.deviceAccessKeys.ValidateDeviceAccessKey(r.Context(), projectID, hash(deviceAccessKeyValue))
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

	deviceAccessKey, err := s.deviceAccessKeys.CreateDeviceAccessKey(r.Context(), projectID, device.ID, hash(deviceAccessKeyValue))
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

func hash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}
