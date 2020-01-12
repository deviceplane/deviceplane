package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/controller/authz"
	"github.com/deviceplane/deviceplane/pkg/controller/connman"
	"github.com/deviceplane/deviceplane/pkg/controller/middleware"
	"github.com/deviceplane/deviceplane/pkg/controller/query"
	"github.com/deviceplane/deviceplane/pkg/controller/spaserver"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/email"
	"github.com/deviceplane/deviceplane/pkg/hash"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/namesgenerator"
	"github.com/deviceplane/deviceplane/pkg/revdial"
	"github.com/deviceplane/deviceplane/pkg/spec"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"gopkg.in/yaml.v2"
)

const (
	sessionCookie = "dp_sess"
)

var (
	errEmailDomainNotAllowed = errors.New("email domain not allowed")
	errEmailAlreadyTaken     = errors.New("email already taken")
	errTokenExpired          = errors.New("token expired")
)

type Service struct {
	users                      store.Users
	passwordRecoveryTokens     store.PasswordRecoveryTokens
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
	deviceRegistrationTokens   store.DeviceRegistrationTokens
	devicesRegisteredWithToken store.DevicesRegisteredWithToken
	deviceAccessKeys           store.DeviceAccessKeys
	applications               store.Applications
	applicationDeviceCounts    store.ApplicationDeviceCounts
	releases                   store.Releases
	releaseDeviceCounts        store.ReleaseDeviceCounts
	deviceApplicationStatuses  store.DeviceApplicationStatuses
	deviceServiceStatuses      store.DeviceServiceStatuses
	metricConfigs              store.MetricConfigs
	email                      email.Interface
	emailFromName              string
	emailFromAddress           string
	allowedEmailDomains        []string
	st                         *statsd.Client
	connman                    *connman.ConnectionManager

	router   *mux.Router
	upgrader websocket.Upgrader
}

func NewService(
	users store.Users,
	registrationTokens store.RegistrationTokens,
	passwordRecoveryTokens store.PasswordRecoveryTokens,
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
	deviceRegistrationTokens store.DeviceRegistrationTokens,
	devicesRegisteredWithToken store.DevicesRegisteredWithToken,
	deviceAccessKeys store.DeviceAccessKeys,
	applications store.Applications,
	applicationDeviceCounts store.ApplicationDeviceCounts,
	releases store.Releases,
	releasesDeviceCounts store.ReleaseDeviceCounts,
	deviceApplicationStatuses store.DeviceApplicationStatuses,
	deviceServiceStatuses store.DeviceServiceStatuses,
	metricConfigs store.MetricConfigs,
	email email.Interface,
	emailFromName string,
	emailFromAddress string,
	allowedEmailDomains []string,
	fileSystem http.FileSystem,
	st *statsd.Client,
	connman *connman.ConnectionManager,
	allowedOrigins []url.URL,
) *Service {
	s := &Service{
		users:                      users,
		registrationTokens:         registrationTokens,
		passwordRecoveryTokens:     passwordRecoveryTokens,
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
		deviceRegistrationTokens:   deviceRegistrationTokens,
		devicesRegisteredWithToken: devicesRegisteredWithToken,
		deviceAccessKeys:           deviceAccessKeys,
		applications:               applications,
		applicationDeviceCounts:    applicationDeviceCounts,
		releases:                   releases,
		releaseDeviceCounts:        releasesDeviceCounts,
		deviceApplicationStatuses:  deviceApplicationStatuses,
		deviceServiceStatuses:      deviceServiceStatuses,
		metricConfigs:              metricConfigs,
		email:                      email,
		emailFromName:              emailFromName,
		emailFromAddress:           emailFromAddress,
		allowedEmailDomains:        allowedEmailDomains,
		st:                         st,
		connman:                    connman,

		router: mux.NewRouter(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			Subprotocols:    []string{"binary"},
			CheckOrigin: func(r *http.Request) bool {
				return utils.CheckSameOrAllowedOrigin(r, allowedOrigins)
			},
		},
	}

	apiRouter := s.router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/register", s.register).Methods("POST")
	apiRouter.HandleFunc("/completeregistration", s.confirmRegistration).Methods("POST")

	apiRouter.HandleFunc("/recoverpassword", s.recoverPassword).Methods("POST")
	apiRouter.HandleFunc("/passwordrecoverytokens/{passwordrecoverytokenvalue}", s.getPasswordRecoveryToken).Methods("GET")
	apiRouter.HandleFunc("/changepassword", s.changePassword).Methods("POST")

	apiRouter.HandleFunc("/login", s.login).Methods("POST")
	apiRouter.HandleFunc("/logout", s.logout).Methods("POST")

	apiRouter.HandleFunc("/me", s.withUserOrServiceAccountAuth(s.getMe)).Methods("GET")
	apiRouter.HandleFunc("/me", s.withUserOrServiceAccountAuth(s.updateMe)).Methods("PATCH")

	apiRouter.HandleFunc("/memberships", s.withUserOrServiceAccountAuth(s.listMembershipsByUser)).Methods("GET")

	apiRouter.HandleFunc("/useraccesskeys", s.withUserOrServiceAccountAuth(s.createUserAccessKey)).Methods("POST")
	apiRouter.HandleFunc("/useraccesskeys/{useraccesskey}", s.withUserOrServiceAccountAuth(s.getUserAccessKey)).Methods("GET")
	apiRouter.HandleFunc("/useraccesskeys", s.withUserOrServiceAccountAuth(s.listUserAccessKeys)).Methods("GET")
	apiRouter.HandleFunc("/useraccesskeys/{useraccesskey}", s.withUserOrServiceAccountAuth(s.deleteUserAccessKey)).Methods("DELETE")

	apiRouter.HandleFunc("/projects", s.withUserOrServiceAccountAuth(s.createProject)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}", s.validateAuthorization(authz.ResourceProjects, authz.ActionGetProject, s.getProject)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}", s.validateAuthorization(authz.ResourceProjects, authz.ActionUpdateProject, s.updateProject)).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}", s.validateAuthorization(authz.ResourceProjects, authz.ActionDeleteProject, s.deleteProject)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/roles", s.validateAuthorization(authz.ResourceRoles, authz.ActionCreateRole, s.createRole)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/roles/{role}", s.validateAuthorization(authz.ResourceRoles, authz.ActionGetRole, s.withRole(s.getRole))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/roles", s.validateAuthorization(authz.ResourceRoles, authz.ActionListRoles, s.listRoles)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/roles/{role}", s.validateAuthorization(authz.ResourceRoles, authz.ActionUpdateRole, s.updateRole)).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/roles/{role}", s.validateAuthorization(authz.ResourceRoles, authz.ActionDeleteRole, s.deleteRole)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/memberships", s.validateAuthorization(authz.ResourceMemberships, authz.ActionCreateMembership, s.createMembership)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}", s.validateAuthorization(authz.ResourceMemberships, authz.ActionGetMembership, s.getMembership)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/memberships", s.validateAuthorization(authz.ResourceMemberships, authz.ActionListMembershipsByProject, s.listMembershipsByProject)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}", s.validateAuthorization(authz.ResourceMemberships, authz.ActionDeleteMembership, s.deleteMembership)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/memberships/{user}/roles/{role}/membershiprolebindings", s.validateAuthorization(authz.ResourceMembershipRoleBindings, authz.ActionCreateMembershipRoleBinding, s.withRole(s.createMembershipRoleBinding))).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}/roles/{role}/membershiprolebindings", s.validateAuthorization(authz.ResourceMembershipRoleBindings, authz.ActionGetMembershipRoleBinding, s.withRole(s.getMembershipRoleBinding))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}/membershiprolebindings", s.validateAuthorization(authz.ResourceMembershipRoleBindings, authz.ActionListMembershipRoleBindings, s.listMembershipRoleBindings)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}/roles/{role}/membershiprolebindings", s.validateAuthorization(authz.ResourceMembershipRoleBindings, authz.ActionDeleteMembershipRoleBinding, s.withRole(s.deleteMembershipRoleBinding))).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/serviceaccounts", s.validateAuthorization(authz.ResourceServiceAccounts, authz.ActionCreateServiceAccount, s.createServiceAccount)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}", s.validateAuthorization(authz.ResourceServiceAccounts, authz.ActionGetServiceAccount, s.withServiceAccount(s.getServiceAccount))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts", s.validateAuthorization(authz.ResourceServiceAccounts, authz.ActionListServiceAccounts, s.listServiceAccounts)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}", s.validateAuthorization(authz.ResourceServiceAccounts, authz.ActionUpdateServiceAccount, s.withServiceAccount(s.updateServiceAccount))).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}", s.validateAuthorization(authz.ResourceServiceAccounts, authz.ActionDeleteServiceAccount, s.withServiceAccount(s.deleteServiceAccount))).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys", s.validateAuthorization(authz.ResourceServiceAccountAccessKeys, authz.ActionCreateServiceAccountAccessKey, s.createServiceAccountAccessKey)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys/{serviceaccountaccesskey}", s.validateAuthorization(authz.ResourceServiceAccountAccessKeys, authz.ActionGetServiceAccountAccessKey, s.getServiceAccountAccessKey)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys", s.validateAuthorization(authz.ResourceServiceAccountAccessKeys, authz.ActionListServiceAccountAccessKeys, s.listServiceAccountAccessKeys)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys/{serviceaccountaccesskey}", s.validateAuthorization(authz.ResourceServiceAccountAccessKeys, authz.ActionDeleteServiceAccountAccessKey, s.deleteServiceAccountAccessKey)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/roles/{role}/serviceaccountrolebindings", s.validateAuthorization(authz.ResourceServiceAccountRoleBindings, authz.ActionCreateServiceAccountRoleBinding, s.withRole(s.createServiceAccountRoleBinding))).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/roles/{role}/serviceaccountrolebindings", s.validateAuthorization(authz.ResourceServiceAccountRoleBindings, authz.ActionGetServiceAccountRoleBinding, s.withRole(s.getServiceAccountRoleBinding))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountrolebindings", s.validateAuthorization(authz.ResourceServiceAccountRoleBindings, authz.ActionListServiceAccountRoleBinding, s.listServiceAccountRoleBindings)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/roles/{role}/serviceaccountrolebindings", s.validateAuthorization(authz.ResourceServiceAccountRoleBindings, authz.ActionDeleteServiceAccountRoleBinding, s.withRole(s.deleteServiceAccountRoleBinding))).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/applications", s.validateAuthorization(authz.ResourceApplications, authz.ActionCreateApplication, s.createApplication)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}", s.validateAuthorization(authz.ResourceApplications, authz.ActionGetApplication, s.withApplication(s.getApplication))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/applications", s.validateAuthorization(authz.ResourceApplications, authz.ActionListApplications, s.listApplications)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}", s.validateAuthorization(authz.ResourceApplications, authz.ActionUpdateApplication, s.withApplication(s.updateApplication))).Methods("PATCH")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}", s.validateAuthorization(authz.ResourceApplications, authz.ActionDeleteApplication, s.withApplication(s.deleteApplication))).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/applications/{application}/releases", s.validateAuthorization(authz.ResourceReleases, authz.ActionCreateRelease, s.withApplication(s.createRelease))).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}/releases/latest", s.validateAuthorization(authz.ResourceReleases, authz.ActionGetLatestRelease, s.withApplication(s.getLatestRelease))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}/releases/{release}", s.validateAuthorization(authz.ResourceReleases, authz.ActionGetRelease, s.withApplication(s.getRelease))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}/releases", s.validateAuthorization(authz.ResourceReleases, authz.ActionListReleases, s.withApplication(s.listReleases))).Methods("GET")

	apiRouter.HandleFunc("/projects/{project}/devices/{device}", s.validateAuthorization(authz.ResourceDevices, authz.ActionGetDevice, s.withDevice(s.getDevice))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices", s.validateAuthorization(authz.ResourceDevices, authz.ActionListDevices, s.listDevices)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}", s.validateAuthorization(authz.ResourceDevices, authz.ActionUpdateDevice, s.withDevice(s.updateDevice))).Methods("PATCH")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}", s.validateAuthorization(authz.ResourceDevices, authz.ActionDeleteDevice, s.withDevice(s.deleteDevice))).Methods("DELETE")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/ssh", s.validateAuthorization(authz.ResourceDevices, authz.ActionSSH, s.withDevice(s.initiateSSH))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/reboot", s.validateAuthorization(authz.ResourceDevices, authz.ActionReboot, s.withDevice(s.initiateReboot))).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/imagepullprogress", s.validateAuthorization(authz.ResourceDevices, authz.ActionGetImagePullProgress, s.withDevice(s.imagePullProgress))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/metrics/host", s.validateAuthorization(authz.ResourceDevices, authz.ActionGetMetrics, s.withDevice(s.hostMetrics))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/metrics/agent", s.validateAuthorization(authz.ResourceDevices, authz.ActionGetMetrics, s.withDevice(s.agentMetrics))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/metrics", s.validateAuthorization(authz.ResourceDevices, authz.ActionGetServiceMetrics, s.withApplicationAndDevice(s.serviceMetrics))).Methods("GET")

	apiRouter.HandleFunc("/projects/{project}/devices/{device}/labels", s.validateAuthorization(authz.ResourceDeviceLabels, authz.ActionSetDeviceLabel, s.withDevice(s.setDeviceLabel))).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/labels/{key}", s.validateAuthorization(authz.ResourceDeviceLabels, authz.ActionDeleteDeviceLabel, s.withDevice(s.deleteDeviceLabel))).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/devicelabels", s.validateAuthorization(authz.ResourceDeviceLabels, authz.ActionListAllDeviceLabels, s.listAllDeviceLabelKeys)).Methods("GET")

	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens", s.validateAuthorization(authz.ResourceDeviceRegistrationTokens, authz.ActionListDeviceRegistrationTokens, s.listDeviceRegistrationTokens)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens", s.validateAuthorization(authz.ResourceDeviceRegistrationTokens, authz.ActionCreateDeviceRegistrationToken, s.createDeviceRegistrationToken)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}", s.validateAuthorization(authz.ResourceDeviceRegistrationTokens, authz.ActionGetDeviceRegistrationToken, s.withDeviceRegistrationToken(s.getDeviceRegistrationToken))).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}", s.validateAuthorization(authz.ResourceDeviceRegistrationTokens, authz.ActionUpdateDeviceRegistrationToken, s.withDeviceRegistrationToken(s.updateDeviceRegistrationToken))).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}", s.validateAuthorization(authz.ResourceDeviceRegistrationTokens, authz.ActionDeleteDeviceRegistrationToken, s.withDeviceRegistrationToken(s.deleteDeviceRegistrationToken))).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}/labels", s.validateAuthorization(authz.ResourceDeviceRegistrationTokenLabels, authz.ActionSetDeviceRegistrationTokenLabel, s.withDeviceRegistrationToken(s.setDeviceRegistrationTokenLabel))).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}/labels/{key}", s.validateAuthorization(authz.ResourceDeviceRegistrationTokenLabels, authz.ActionDeleteDeviceRegistrationTokenLabel, s.withDeviceRegistrationToken(s.deleteDeviceRegistrationTokenLabel))).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/configs/{key}", s.validateAuthorization(authz.ResourceProjectConfigs, authz.ActionGetProjectConfig, s.getProjectConfig)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/configs/{key}", s.validateAuthorization(authz.ResourceProjectConfigs, authz.ActionSetProjectConfig, s.setProjectConfig)).Methods("PUT")

	apiRouter.HandleFunc("/projects/{project}/devices/register", s.registerDevice).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/bundle", s.withDeviceAuth(s.getBundle)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/info", s.withDeviceAuth(s.setDeviceInfo)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/deviceapplicationstatuses", s.withDeviceAuth(s.setDeviceApplicationStatus)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/deviceapplicationstatuses", s.withDeviceAuth(s.deleteDeviceApplicationStatus)).Methods("DELETE")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/deviceservicestatuses", s.withDeviceAuth(s.setDeviceServiceStatus)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/deviceservicestatuses", s.withDeviceAuth(s.deleteDeviceServiceStatus)).Methods("DELETE")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/connection", s.withDeviceAuth(s.initiateDeviceConnection)).Methods("GET")

	apiRouter.Handle("/revdial", revdial.ConnHandler(s.upgrader)).Methods("GET")

	apiRouter.HandleFunc("/health", s.health).Methods("GET")
	apiRouter.HandleFunc("/500", s.withUserOrServiceAccountAuth(s.intentional500)).Methods("GET")

	s.router.PathPrefix("/api").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	s.router.PathPrefix("/").Handler(spaserver.NewSPAFileServer(fileSystem))

	return s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) health(w http.ResponseWriter, r *http.Request) {
	s.st.Incr("health", nil, 1)
}

func (s *Service) intentional500(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	// TODO
	if authenticatedUserID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := s.users.GetUser(r.Context(), authenticatedUserID)
	if err != nil {
		log.WithError(err).Error("get user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.SuperAdmin {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
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

func (s *Service) validateAuthorization(requestedResource authz.Resource, requestedAction authz.Action, handler func(http.ResponseWriter, *http.Request, string, string, string)) func(http.ResponseWriter, *http.Request) {
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
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if err != nil {
				log.WithError(err).Error("lookup project")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			projectID = project.ID
		}

		var roles []string
		superAdmin := false
		if authenticatedUserID != "" {
			authenticatedUser, err := s.users.GetUser(r.Context(), authenticatedUserID)
			if err != nil {
				log.WithError(err).Error("get user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if authenticatedUser.SuperAdmin {
				superAdmin = true
			} else {
				if _, err := s.memberships.GetMembership(r.Context(),
					authenticatedUserID, projectID,
				); err == store.ErrMembershipNotFound {
					http.Error(w, err.Error(), http.StatusNotFound)
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
		if superAdmin {
			configs = []authz.Config{
				authz.AdminAllRole,
			}
		} else {
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
		}

		if !authz.Evaluate(requestedResource, requestedAction, configs) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		handler(w, r, projectID, authenticatedUserID, authenticatedServiceAccountID)
	})
}

func (s *Service) register(w http.ResponseWriter, r *http.Request) {
	utils.WithReferrer(w, r, func(referrer *url.URL) {
		var registerRequest struct {
			Email     string `json:"email" validate:"email"`
			Password  string `json:"password" validate:"password"`
			FirstName string `json:"firstName" validate:"required,min=1,max=100"`
			LastName  string `json:"lastName" validate:"required,min=1,max=100"`
			Company   string `json:"company" validate:"max=100"`
		}
		if err := read(r, &registerRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		emailDomain, err := utils.GetDomainFromEmail(registerRequest.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		emailDomainAllowed := false
		if len(s.allowedEmailDomains) == 0 {
			emailDomainAllowed = true
		} else {
			for _, allowedEmailDomain := range s.allowedEmailDomains {
				if allowedEmailDomain == emailDomain {
					emailDomainAllowed = true
					break
				}
			}
		}

		if !emailDomainAllowed {
			http.Error(w, errEmailDomainNotAllowed.Error(), http.StatusBadRequest)
			return
		}

		if _, err := s.users.LookupUser(r.Context(), registerRequest.Email); err == nil {
			http.Error(w, errEmailAlreadyTaken.Error(), http.StatusBadRequest)
			return
		} else if err != store.ErrUserNotFound {
			log.WithError(err).Error("lookup user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := s.users.CreateUser(r.Context(), registerRequest.Email, hash.Hash(registerRequest.Password),
			registerRequest.FirstName, registerRequest.LastName, registerRequest.Company)
		if err != nil {
			log.WithError(err).Error("create user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if s.email == nil {
			// If email provider is nil then skip the registration workflow
			if _, err := s.users.MarkRegistrationCompleted(r.Context(), user.ID); err != nil {
				log.WithError(err).Error("mark registration completed")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			registrationTokenValue := ksuid.New().String()

			if _, err := s.registrationTokens.CreateRegistrationToken(r.Context(), user.ID, hash.Hash(registrationTokenValue)); err != nil {
				log.WithError(err).Error("create registration token")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			content := fmt.Sprintf(
				"Please go to the following URL to complete registration. %s://%s/confirm/%s",
				referrer.Scheme,
				referrer.Host,
				registrationTokenValue,
			)

			if err := s.email.Send(email.Request{
				FromName:    s.emailFromName,
				FromAddress: s.emailFromAddress,
				ToName:      user.FirstName + " " + user.LastName,
				ToAddress:   user.Email,
				Subject:     "Deviceplane Registration Confirmation",
				Body:        content,
			}); err != nil {
				log.WithError(err).Error("send registration email")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		utils.Respond(w, user)
	})
}

func (s *Service) confirmRegistration(w http.ResponseWriter, r *http.Request) {
	var confirmRegistrationRequest struct {
		RegistrationTokenValue string `json:"registrationTokenValue" validate:"required,min=1,max=1000"`
	}
	if err := read(r, &confirmRegistrationRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func (s *Service) recoverPassword(w http.ResponseWriter, r *http.Request) {
	utils.WithReferrer(w, r, func(referrer *url.URL) {
		var recoverPasswordRequest struct {
			Email string `json:"email" validate:"email"`
		}
		if err := read(r, &recoverPasswordRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := s.users.LookupUser(r.Context(), recoverPasswordRequest.Email)
		if err == store.ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			log.WithError(err).Error("lookup user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		passwordRecoveryTokenValue := ksuid.New().String()

		if _, err := s.passwordRecoveryTokens.CreatePasswordRecoveryToken(r.Context(),
			user.ID, hash.Hash(passwordRecoveryTokenValue)); err != nil {
			log.WithError(err).Error("create password recovery token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		content := fmt.Sprintf(
			"Please go to the following URL to recover your password. %s://%s/recover/%s",
			referrer.Scheme,
			referrer.Host,
			passwordRecoveryTokenValue,
		)

		if err := s.email.Send(email.Request{
			FromName:    s.emailFromName,
			FromAddress: s.emailFromAddress,
			ToName:      user.FirstName + " " + user.LastName,
			ToAddress:   user.Email,
			Subject:     "Deviceplane Password Recovery",
			Body:        content,
		}); err != nil {
			log.WithError(err).Error("send recovery email")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Service) getPasswordRecoveryToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	passwordRecoveryTokenValue := vars["passwordrecoverytokenvalue"]

	passwordRecoveryToken, err := s.passwordRecoveryTokens.ValidatePasswordRecoveryToken(r.Context(), hash.Hash(passwordRecoveryTokenValue))
	if err == store.ErrPasswordRecoveryTokenNotFound {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.WithError(err).Error("validate password recovery token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, passwordRecoveryToken)
}

func (s *Service) changePassword(w http.ResponseWriter, r *http.Request) {
	var changePasswordRequest struct {
		PasswordRecoveryTokenValue string `json:"passwordRecoveryTokenValue" validate:"required,min=1,max=1000"`
		Password                   string `json:"password" validate:"password"`
	}
	if err := read(r, &changePasswordRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	passwordRecoveryToken, err := s.passwordRecoveryTokens.ValidatePasswordRecoveryToken(r.Context(),
		hash.Hash(changePasswordRequest.PasswordRecoveryTokenValue))
	if err == store.ErrPasswordRecoveryTokenNotFound {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.WithError(err).Error("validate password recovery token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if time.Now().After(passwordRecoveryToken.ExpiresAt) {
		http.Error(w, errTokenExpired.Error(), http.StatusForbidden)
		return
	}

	if _, err := s.users.UpdatePasswordHash(r.Context(),
		passwordRecoveryToken.UserID, hash.Hash(changePasswordRequest.Password)); err != nil {
		log.WithError(err).Error("update password hash")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := s.users.GetUser(r.Context(), passwordRecoveryToken.UserID)
	if err != nil {
		log.WithError(err).Error("get user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, user)
}

func (s *Service) login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email" validate:"email"`
		Password string `json:"password" validate:"password"`
	}
	if err := read(r, &loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.users.ValidateUserWithEmail(r.Context(), loginRequest.Email, hash.Hash(loginRequest.Password))
	if err == store.ErrUserNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("validate user with email")
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
	utils.WithReferrer(w, r, func(referrer *url.URL) {
		sessionValue := ksuid.New().String()

		if _, err := s.sessions.CreateSession(r.Context(), userID, hash.Hash(sessionValue)); err != nil {
			log.WithError(err).Error("create session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var secure bool
		switch referrer.Scheme {
		case "http":
			secure = false
		case "https":
			secure = true
		}

		cookie := &http.Cookie{
			Name:  sessionCookie,
			Value: sessionValue,

			Expires: time.Now().AddDate(0, 1, 0),

			Domain:   r.Host,
			Secure:   secure,
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)
	})
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
		return
	default:
		log.WithError(err).Error("get session cookie")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) getMe(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	// TODO
	if authenticatedUserID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := s.users.GetUser(r.Context(), authenticatedUserID)
	if err != nil {
		log.WithError(err).Error("get user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, user)
}

func (s *Service) updateMe(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	// TODO
	if authenticatedUserID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: validation
	var updateUserRequest struct {
		Password        *string `json:"password"`
		CurrentPassword *string `json:"currentPassword"`
		FirstName       *string `json:"firstName"`
		LastName        *string `json:"lastName"`
		Company         *string `json:"company"`
	}
	if err := read(r, &updateUserRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updateUserRequest.Password != nil {
		if updateUserRequest.CurrentPassword == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if _, err := s.users.ValidateUser(r.Context(), authenticatedUserID, hash.Hash(*updateUserRequest.CurrentPassword)); err == store.ErrUserNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			log.WithError(err).Error("validate user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := s.users.UpdatePasswordHash(r.Context(), authenticatedUserID, hash.Hash(*updateUserRequest.Password)); err != nil {
			log.WithError(err).Error("update password hash")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if updateUserRequest.FirstName != nil {
		if _, err := s.users.UpdateFirstName(r.Context(), authenticatedUserID, *updateUserRequest.FirstName); err != nil {
			log.WithError(err).Error("update first name")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if updateUserRequest.LastName != nil {
		if _, err := s.users.UpdateLastName(r.Context(), authenticatedUserID, *updateUserRequest.LastName); err != nil {
			log.WithError(err).Error("update last name")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if updateUserRequest.Company != nil {
		if _, err := s.users.UpdateCompany(r.Context(), authenticatedUserID, *updateUserRequest.Company); err != nil {
			log.WithError(err).Error("update company")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	user, err := s.users.GetUser(r.Context(), authenticatedUserID)
	if err != nil {
		log.WithError(err).Error("get user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, user)
}

func (s *Service) listMembershipsByUser(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	// TODO
	if authenticatedUserID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	memberships, err := s.memberships.ListMembershipsByUser(r.Context(), authenticatedUserID)
	if err != nil {
		log.WithError(err).Error("list memberships by user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = memberships
	if _, ok := r.URL.Query()["full"]; ok {
		membershipsFull := make([]models.MembershipFull1, 0)

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

		ret = membershipsFull
	}

	utils.Respond(w, ret)
}

func (s *Service) createUserAccessKey(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	// TODO
	if authenticatedUserID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var createUserAccessKeyRequest struct {
		Description string `json:"description" validate:"description"`
	}
	if err := read(r, &createUserAccessKeyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userAccessKeyValue := "u" + ksuid.New().String()

	user, err := s.userAccessKeys.CreateUserAccessKey(r.Context(),
		authenticatedUserID, hash.Hash(userAccessKeyValue), createUserAccessKeyRequest.Description)
	if err != nil {
		log.WithError(err).Error("create user access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, models.UserAccessKeyWithValue{
		UserAccessKey: *user,
		Value:         userAccessKeyValue,
	})
}

// TODO: verify that the user owns this access key
func (s *Service) getUserAccessKey(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	// TODO
	if authenticatedUserID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	userAccessKeyID := vars["useraccesskey"]

	userAccessKey, err := s.userAccessKeys.GetUserAccessKey(r.Context(), userAccessKeyID)
	if err == store.ErrUserAccessKeyNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get user access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, userAccessKey)
}

func (s *Service) listUserAccessKeys(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	// TODO
	if authenticatedUserID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userAccessKeys, err := s.userAccessKeys.ListUserAccessKeys(r.Context(), authenticatedUserID)
	if err != nil {
		log.WithError(err).Error("list users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, userAccessKeys)
}

// TODO: verify that the user owns this access key
func (s *Service) deleteUserAccessKey(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	vars := mux.Vars(r)
	userAccessKeyID := vars["useraccesskey"]

	if err := s.userAccessKeys.DeleteUserAccessKey(r.Context(), userAccessKeyID); err != nil {
		log.WithError(err).Error("delete user access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) createProject(w http.ResponseWriter, r *http.Request, authenticatedUserID, authenticatedServiceAccountID string) {
	var createProjectRequest struct {
		Name string `json:"name" validate:"name"`
	}
	if err := read(r, &createProjectRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := s.projects.LookupProject(r.Context(), createProjectRequest.Name); err == nil {
		http.Error(w, store.ErrProjectNameAlreadyInUse.Error(), http.StatusBadRequest)
		return
	} else if err != nil && err != store.ErrProjectNotFound {
		log.WithError(err).Error("lookup project")
		w.WriteHeader(http.StatusInternalServerError)
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

	// Create default device registration token.
	// It is named "default" and has an unlimited device registration cap.
	_, err = s.deviceRegistrationTokens.CreateDeviceRegistrationToken(
		r.Context(),
		project.ID,
		"default",
		// This copy can be improved:
		"This default registration token is used for provisioning new devices from the UI.",
		nil,
	)
	if err != nil {
		log.WithError(err).Error("create default registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, project)
}

func (s *Service) getProject(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	project, err := s.projects.GetProject(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("get project")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, project)
}

func (s *Service) updateProject(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	var updateProjectRequest struct {
		Name          string `json:"name" validate:"name"`
		DatadogApiKey string `json:"datadogApiKey"`
	}
	if err := read(r, &updateProjectRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if project, err := s.projects.LookupProject(r.Context(),
		updateProjectRequest.Name); err == nil && project.ID != projectID {
		http.Error(w, store.ErrProjectNameAlreadyInUse.Error(), http.StatusBadRequest)
		return
	} else if err != nil && err != store.ErrProjectNotFound {
		log.WithError(err).Error("lookup project")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	project, err := s.projects.UpdateProject(r.Context(), projectID, updateProjectRequest.Name, updateProjectRequest.DatadogApiKey)
	if err != nil {
		log.WithError(err).Error("update project")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, project)
}

func (s *Service) deleteProject(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	if err := s.projects.DeleteProject(r.Context(), projectID); err != nil {
		log.WithError(err).Error("delete project")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) createRole(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	var createRoleRequest struct {
		Name        string `json:"name" validate:"name"`
		Description string `json:"description" validate:"description"`
		Config      string `json:"config" validate:"config"`
	}
	if err := read(r, &createRoleRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := s.roles.LookupRole(r.Context(), createRoleRequest.Name, projectID); err == nil {
		http.Error(w, store.ErrRoleNameAlreadyInUse.Error(), http.StatusBadRequest)
		return
	} else if err != nil && err != store.ErrRoleNotFound {
		log.WithError(err).Error("lookup role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var roleConfig authz.Config
	if err := yaml.UnmarshalStrict([]byte(createRoleRequest.Config), &roleConfig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	role, err := s.roles.CreateRole(r.Context(), projectID, createRoleRequest.Name,
		createRoleRequest.Description, createRoleRequest.Config)
	if err != nil {
		log.WithError(err).Error("create role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, role)
}

func (s *Service) getRole(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	roleID string,
) {
	role, err := s.roles.GetRole(r.Context(), roleID, projectID)
	if err == store.ErrRoleNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, role)
}

func (s *Service) listRoles(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	roles, err := s.roles.ListRoles(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list roles")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, roles)
}

func (s *Service) updateRole(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	vars := mux.Vars(r)
	roleID := vars["role"]

	var updateRoleRequest struct {
		Name        string `json:"name" validate:"name"`
		Description string `json:"description" validate:"description"`
		Config      string `json:"config" validate:"config"`
	}
	if err := read(r, &updateRoleRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var roleConfig authz.Config
	if err := yaml.UnmarshalStrict([]byte(updateRoleRequest.Config), &roleConfig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if role, err := s.roles.LookupRole(r.Context(),
		updateRoleRequest.Name, projectID); err == nil && role.ID != roleID {
		http.Error(w, store.ErrRoleNameAlreadyInUse.Error(), http.StatusBadRequest)
		return
	} else if err != nil && err != store.ErrRoleNotFound {
		log.WithError(err).Error("lookup role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	role, err := s.roles.UpdateRole(r.Context(), roleID, projectID, updateRoleRequest.Name,
		updateRoleRequest.Description, updateRoleRequest.Config)
	if err != nil {
		log.WithError(err).Error("update role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, role)
}

func (s *Service) deleteRole(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	vars := mux.Vars(r)
	roleID := vars["role"]

	if err := s.roles.DeleteRole(r.Context(), roleID, projectID); err != nil {
		log.WithError(err).Error("delete role")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) createServiceAccount(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	var createServiceAccountRequest struct {
		Name        string `json:"name" validate:"name"`
		Description string `json:"description" validate:"description"`
	}
	if err := read(r, &createServiceAccountRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := s.serviceAccounts.LookupServiceAccount(r.Context(), createServiceAccountRequest.Name, projectID); err == nil {
		http.Error(w, store.ErrServiceAccountNameAlreadyInUse.Error(), http.StatusBadRequest)
		return
	} else if err != nil && err != store.ErrServiceAccountNotFound {
		log.WithError(err).Error("lookup service account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	serviceAccount, err := s.serviceAccounts.CreateServiceAccount(r.Context(), projectID, createServiceAccountRequest.Name,
		createServiceAccountRequest.Description)
	if err != nil {
		log.WithError(err).Error("create service account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, serviceAccount)
}

func (s *Service) getServiceAccount(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	serviceAccountID string,
) {
	serviceAccount, err := s.serviceAccounts.GetServiceAccount(r.Context(), serviceAccountID, projectID)
	if err == store.ErrServiceAccountNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
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

	utils.Respond(w, ret)
}

func (s *Service) listServiceAccounts(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
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

	utils.Respond(w, ret)
}

func (s *Service) updateServiceAccount(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	serviceAccountID string,
) {
	var updateServiceAccountRequest struct {
		Name        string `json:"name" validate:"name"`
		Description string `json:"description" validate:"description"`
	}
	if err := read(r, &updateServiceAccountRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if serviceAccount, err := s.serviceAccounts.LookupServiceAccount(r.Context(),
		updateServiceAccountRequest.Name, projectID); err == nil && serviceAccount.ID != serviceAccountID {
		http.Error(w, store.ErrServiceAccountNameAlreadyInUse.Error(), http.StatusBadRequest)
		return
	} else if err != nil && err != store.ErrServiceAccountNotFound {
		log.WithError(err).Error("lookup service account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	serviceAccount, err := s.serviceAccounts.UpdateServiceAccount(r.Context(), serviceAccountID, projectID,
		updateServiceAccountRequest.Name, updateServiceAccountRequest.Description)
	if err != nil {
		log.WithError(err).Error("update service account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, serviceAccount)
}

func (s *Service) deleteServiceAccount(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	serviceAccountID string,
) {
	if err := s.serviceAccounts.DeleteServiceAccount(r.Context(), serviceAccountID, projectID); err != nil {
		log.WithError(err).Error("delete service account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) createServiceAccountAccessKey(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	vars := mux.Vars(r)
	serviceAccountID := vars["serviceaccount"]

	var createServiceAccountAccessKeyRequest struct {
		Description string `json:"description" validate:"description"`
	}
	if err := read(r, &createServiceAccountAccessKeyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	serviceAccountAccessKeyValue := "s" + ksuid.New().String()

	serviceAccount, err := s.serviceAccountAccessKeys.CreateServiceAccountAccessKey(r.Context(),
		projectID, serviceAccountID, hash.Hash(serviceAccountAccessKeyValue), createServiceAccountAccessKeyRequest.Description)
	if err != nil {
		log.WithError(err).Error("create service account access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, models.ServiceAccountAccessKeyWithValue{
		ServiceAccountAccessKey: *serviceAccount,
		Value:                   serviceAccountAccessKeyValue,
	})
}

func (s *Service) getServiceAccountAccessKey(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	vars := mux.Vars(r)
	serviceAccountAccessKeyID := vars["serviceaccountaccesskey"]

	serviceAccountAccessKey, err := s.serviceAccountAccessKeys.GetServiceAccountAccessKey(r.Context(), serviceAccountAccessKeyID, projectID)
	if err == store.ErrServiceAccountAccessKeyNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get service account access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, serviceAccountAccessKey)
}

func (s *Service) listServiceAccountAccessKeys(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	vars := mux.Vars(r)
	serviceAccountID := vars["serviceaccount"]

	serviceAccountAccessKeys, err := s.serviceAccountAccessKeys.ListServiceAccountAccessKeys(r.Context(), projectID, serviceAccountID)
	if err != nil {
		log.WithError(err).Error("list service accounts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, serviceAccountAccessKeys)
}

func (s *Service) deleteServiceAccountAccessKey(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	vars := mux.Vars(r)
	serviceAccountAccessKeyID := vars["serviceaccountaccesskey"]

	if err := s.serviceAccountAccessKeys.DeleteServiceAccountAccessKey(r.Context(), serviceAccountAccessKeyID, projectID); err != nil {
		log.WithError(err).Error("delete service account access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) createServiceAccountRoleBinding(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	roleID string,
) {
	vars := mux.Vars(r)
	serviceAccountID := vars["serviceaccount"]

	serviceAccountRoleBinding, err := s.serviceAccountRoleBindings.CreateServiceAccountRoleBinding(r.Context(), serviceAccountID, roleID, projectID)
	if err != nil {
		log.WithError(err).Error("create service account role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, serviceAccountRoleBinding)
}

func (s *Service) getServiceAccountRoleBinding(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	roleID string,
) {
	vars := mux.Vars(r)
	serviceAccountID := vars["serviceaccount"]

	serviceAccountRoleBinding, err := s.serviceAccountRoleBindings.GetServiceAccountRoleBinding(r.Context(), serviceAccountID, roleID, projectID)
	if err != nil {
		log.WithError(err).Error("get service account role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, serviceAccountRoleBinding)
}

func (s *Service) listServiceAccountRoleBindings(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	vars := mux.Vars(r)
	serviceAccountID := vars["serviceaccount"]

	serviceAccountRoleBindings, err := s.serviceAccountRoleBindings.ListServiceAccountRoleBindings(r.Context(), projectID, serviceAccountID)
	if err != nil {
		log.WithError(err).Error("list service account role bindings")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, serviceAccountRoleBindings)
}

func (s *Service) deleteServiceAccountRoleBinding(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	roleID string,
) {
	vars := mux.Vars(r)
	serviceAccountID := vars["serviceaccount"]

	if err := s.serviceAccountRoleBindings.DeleteServiceAccountRoleBinding(r.Context(), serviceAccountID, roleID, projectID); err != nil {
		log.WithError(err).Error("delete service account role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) createMembership(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	var createMembershipRequest struct {
		Email string `json:"email" validate:"email"`
	}
	if err := read(r, &createMembershipRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.users.LookupUser(r.Context(), createMembershipRequest.Email)
	if err == store.ErrUserNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("lookup user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	membership, err := s.memberships.CreateMembership(r.Context(), user.ID, projectID)
	if err != nil {
		log.WithError(err).Error("create membership")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, membership)
}

func (s *Service) getMembership(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	vars := mux.Vars(r)
	userID := vars["user"]

	membership, err := s.memberships.GetMembership(r.Context(), userID, projectID)
	if err == store.ErrMembershipNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
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

		membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(), userID, projectID)
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

	utils.Respond(w, ret)
}

func (s *Service) listMembershipsByProject(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	memberships, err := s.memberships.ListMembershipsByProject(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list memberships by project")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = memberships
	if _, ok := r.URL.Query()["full"]; ok {
		membershipsFull := make([]models.MembershipFull2, 0)

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

	utils.Respond(w, ret)
}

func (s *Service) deleteMembership(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	vars := mux.Vars(r)
	userID := vars["user"]

	if err := s.memberships.DeleteMembership(r.Context(), userID, projectID); err != nil {
		log.WithError(err).Error("delete membership")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) createMembershipRoleBinding(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	roleID string,
) {
	vars := mux.Vars(r)
	userID := vars["user"]

	membershipRoleBinding, err := s.membershipRoleBindings.CreateMembershipRoleBinding(r.Context(), userID, roleID, projectID)
	if err != nil {
		log.WithError(err).Error("create membership role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, membershipRoleBinding)
}

func (s *Service) getMembershipRoleBinding(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	roleID string,
) {
	vars := mux.Vars(r)
	userID := vars["user"]

	membershipRoleBinding, err := s.membershipRoleBindings.GetMembershipRoleBinding(r.Context(), userID, roleID, projectID)
	if err == store.ErrMembershipRoleBindingNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get membership role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, membershipRoleBinding)
}

func (s *Service) listMembershipRoleBindings(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	vars := mux.Vars(r)
	userID := vars["user"]

	membershipRoleBindings, err := s.membershipRoleBindings.ListMembershipRoleBindings(r.Context(), userID, projectID)
	if err != nil {
		log.WithError(err).Error("list membership role bindings")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, membershipRoleBindings)
}

func (s *Service) deleteMembershipRoleBinding(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	roleID string,
) {
	vars := mux.Vars(r)
	userID := vars["user"]

	if err := s.membershipRoleBindings.DeleteMembershipRoleBinding(r.Context(), userID, roleID, projectID); err != nil {
		log.WithError(err).Error("delete membership role binding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) createApplication(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	var createApplicationRequest struct {
		Name        string `json:"name" validate:"name"`
		Description string `json:"description" validate:"description"`
	}
	if err := read(r, &createApplicationRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := s.applications.LookupApplication(r.Context(), createApplicationRequest.Name, projectID); err == nil {
		http.Error(w, store.ErrApplicationNameAlreadyInUse.Error(), http.StatusBadRequest)
		return
	} else if err != nil && err != store.ErrApplicationNotFound {
		log.WithError(err).Error("lookup application")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	application, err := s.applications.CreateApplication(
		r.Context(),
		projectID,
		createApplicationRequest.Name,
		createApplicationRequest.Description)
	if err != nil {
		log.WithError(err).Error("create application")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, application)
}

func (s *Service) getApplication(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	applicationID string,
) {
	application, err := s.applications.GetApplication(r.Context(), applicationID, projectID)
	if err == store.ErrApplicationNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
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

		ret = models.ApplicationFull1{
			Application:   *application,
			LatestRelease: latestRelease,
			DeviceCounts:  *applicationDeviceCounts,
		}
	}

	utils.Respond(w, ret)
}

func (s *Service) listApplications(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	applications, err := s.applications.ListApplications(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list applications")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = applications
	if _, ok := r.URL.Query()["full"]; ok {
		applicationsFull := make([]models.ApplicationFull1, 0)

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

			applicationsFull = append(applicationsFull, models.ApplicationFull1{
				Application:   application,
				LatestRelease: latestRelease,
				DeviceCounts:  *applicationDeviceCounts,
			})
		}

		ret = applicationsFull
	}

	utils.Respond(w, ret)
}

func (s *Service) updateApplication(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	applicationID string,
) {
	var updateApplicationRequest struct {
		Name                  *string                                 `json:"name" validate:"name,omitempty"`
		Description           *string                                 `json:"description" validate:"description,omitempty"`
		SchedulingRule        *models.Query                           `json:"schedulingRule"`
		MetricEndpointConfigs *map[string]models.MetricEndpointConfig `json:"metricEndpointConfigs"`
	}
	if err := read(r, &updateApplicationRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var application *models.Application
	var err error
	if updateApplicationRequest.Name != nil {
		if application, err = s.applications.LookupApplication(r.Context(),
			*updateApplicationRequest.Name, projectID); err == nil && application.ID != applicationID {
			http.Error(w, store.ErrApplicationNameAlreadyInUse.Error(), http.StatusBadRequest)
			return
		} else if err != nil && err != store.ErrApplicationNotFound {
			log.WithError(err).Error("lookup application")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if application, err = s.applications.UpdateApplicationName(r.Context(), applicationID, projectID, *updateApplicationRequest.Name); err != nil {
			log.WithError(err).Error("update application name")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if updateApplicationRequest.Description != nil {
		if application, err = s.applications.UpdateApplicationDescription(r.Context(), applicationID, projectID, *updateApplicationRequest.Description); err != nil {
			log.WithError(err).Error("update application description")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if updateApplicationRequest.SchedulingRule != nil {
		err = query.ValidateQuery(*updateApplicationRequest.SchedulingRule)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if application, err = s.applications.UpdateApplicationSchedulingRule(r.Context(), applicationID, projectID, *updateApplicationRequest.SchedulingRule); err != nil {
			log.WithError(err).Error("update application scheduling rule")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if updateApplicationRequest.MetricEndpointConfigs != nil {
		if application, err = s.applications.UpdateApplicationMetricEndpointConfigs(r.Context(), applicationID, projectID, *updateApplicationRequest.MetricEndpointConfigs); err != nil {
			log.WithError(err).Error("update application service metrics config")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	utils.Respond(w, application)
}

func (s *Service) deleteApplication(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	applicationID string,
) {
	if err := s.applications.DeleteApplication(r.Context(), applicationID, projectID); err != nil {
		log.WithError(err).Error("delete application")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// TODO: this has a vulnerability!
func (s *Service) createRelease(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	applicationID string,
) {
	var createReleaseRequest models.CreateReleaseRequest
	if err := read(r, &createReleaseRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := spec.Validate([]byte(createReleaseRequest.RawConfig)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var applicationConfig map[string]models.Service
	if err := yaml.UnmarshalStrict([]byte(createReleaseRequest.RawConfig), &applicationConfig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonApplicationConfig, err := json.Marshal(applicationConfig)
	if err != nil {
		log.WithError(err).Error("marshal json application config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	release, err := s.releases.CreateRelease(
		r.Context(),
		projectID,
		applicationID,
		createReleaseRequest.RawConfig,
		string(jsonApplicationConfig),
		authenticatedUserID,
		authenticatedServiceAccountID,
	)
	if err != nil {
		log.WithError(err).Error("create release")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, release)
}

func (s *Service) getRelease(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	applicationID string,
) {
	vars := mux.Vars(r)
	releaseID := vars["release"]

	release, err := s.releases.GetRelease(r.Context(), releaseID, projectID, applicationID)
	if err == store.ErrReleaseNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get release")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = release
	if _, ok := r.URL.Query()["full"]; ok {
		ret, err = s.getReleaseFull(r.Context(), *release)
		if err != nil {
			log.WithError(err).Error("get release full")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	utils.Respond(w, ret)
}

func (s *Service) getLatestRelease(w http.ResponseWriter, r *http.Request,
	projectID string, authenticatedUserID, authenticatedServiceAccountID,
	applicationID string,
) {
	release, err := s.releases.GetLatestRelease(r.Context(), projectID, applicationID)
	if err == store.ErrReleaseNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get latest release")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = release
	if _, ok := r.URL.Query()["full"]; ok {
		ret, err = s.getReleaseFull(r.Context(), *release)
		if err != nil {
			log.WithError(err).Error("get release full")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	utils.Respond(w, ret)
}

func (s *Service) listReleases(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	applicationID string,
) {
	releases, err := s.releases.ListReleases(r.Context(), projectID, applicationID)
	if err != nil {
		log.WithError(err).Error("list releases")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = releases
	if _, ok := r.URL.Query()["full"]; ok {
		releasesFull := make([]models.ReleaseFull, 0)
		for _, release := range releases {
			releaseFull, err := s.getReleaseFull(r.Context(), release)
			if err != nil {
				log.WithError(err).Error("get release full")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			releasesFull = append(releasesFull, *releaseFull)
		}
		ret = releasesFull
	}

	utils.Respond(w, ret)
}

func (s *Service) listDevices(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	searchQuery := r.URL.Query().Get("search")

	devices, err := s.devices.ListDevices(r.Context(), projectID, searchQuery)
	if err != nil {
		log.WithError(err).Error("list devices")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filters, err := query.FiltersFromQuery(r.URL.Query())
	if err != nil {
		http.Error(w, errors.Wrap(err, "get filters from query").Error(), http.StatusBadRequest)
		return
	}

	if len(filters) != 0 {
		devices, err = query.FilterDevices(devices, filters)
		if err != nil {
			http.Error(w, errors.Wrap(err, "filter devices").Error(), http.StatusBadRequest)
			return
		}
	}

	ds := make([]interface{}, len(devices))
	for i := range devices {
		ds[i] = devices[i]
	}
	middleware.SortAndPaginateAndRespond(*r, w, ds)
}

func (s *Service) getDevice(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	device, err := s.devices.GetDevice(r.Context(), deviceID, projectID)
	if err == store.ErrDeviceNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.WithError(err).Error("get device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = device
	if _, ok := r.URL.Query()["full"]; ok {
		applications, err := s.applications.ListApplications(r.Context(), projectID)
		if err != nil {
			log.WithError(err).Error("list applications")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		allApplicationStatusInfo := make([]models.DeviceApplicationStatusInfo, 0)
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
			} else {
				log.WithError(err).Error("get device service statuses")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			allApplicationStatusInfo = append(allApplicationStatusInfo, applicationStatusInfo)
		}

		ret = models.DeviceFull{
			Device:                *device,
			ApplicationStatusInfo: allApplicationStatusInfo,
		}
	}

	utils.Respond(w, ret)
}

func (s *Service) updateDevice(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	var updateDeviceRequest struct {
		Name string `json:"name" validate:"name"`
	}
	if err := read(r, &updateDeviceRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if device, err := s.devices.LookupDevice(r.Context(),
		updateDeviceRequest.Name, projectID); err == nil && device.ID != deviceID {
		http.Error(w, store.ErrDeviceNameAlreadyInUse.Error(), http.StatusBadRequest)
		return
	} else if err != nil && err != store.ErrDeviceNotFound {
		log.WithError(err).Error("lookup device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	device, err := s.devices.UpdateDeviceName(r.Context(), deviceID, projectID, updateDeviceRequest.Name)
	if err != nil {
		log.WithError(err).Error("update device name")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, device)
}

func (s *Service) deleteDevice(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	if err := s.devices.DeleteDevice(r.Context(), deviceID, projectID); err != nil {
		log.WithError(err).Error("delete device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) listAllDeviceLabelKeys(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	deviceLabels, err := s.devices.ListAllDeviceLabelKeys(
		r.Context(),
		projectID,
	)
	if err != nil {
		log.WithError(err).Error("list device labels")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, deviceLabels)
}

func (s *Service) setDeviceLabel(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	var setDeviceLabelRequest struct {
		Key   string `json:"key" validate:"labelkey"`
		Value string `json:"value" validate:"labelvalue"`
	}
	if err := read(r, &setDeviceLabelRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deviceLabel, err := s.devices.SetDeviceLabel(
		r.Context(),
		deviceID,
		projectID,
		setDeviceLabelRequest.Key,
		setDeviceLabelRequest.Value,
	)
	if err != nil {
		log.WithError(err).Error("set device label")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, deviceLabel)
}

func (s *Service) deleteDeviceLabel(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceID string,
) {
	vars := mux.Vars(r)
	key := vars["key"]

	if err := s.devices.DeleteDeviceLabel(r.Context(), deviceID, projectID, key); err != nil {
		log.WithError(err).Error("delete device label")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) createDeviceRegistrationToken(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	var createDeviceRegistrationTokenRequest struct {
		Name             string `json:"name" validate:"name"`
		Description      string `json:"description" validate:"description"`
		MaxRegistrations *int   `json:"maxRegistrations"`
	}
	if err := read(r, &createDeviceRegistrationTokenRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := s.deviceRegistrationTokens.LookupDeviceRegistrationToken(
		r.Context(),
		createDeviceRegistrationTokenRequest.Name,
		projectID,
	); err == nil {
		http.Error(w, store.ErrDeviceRegistrationTokenNameAlreadyInUse.Error(), http.StatusBadRequest)
		return
	} else if err != nil && err != store.ErrDeviceRegistrationTokenNotFound {
		log.WithError(err).Error("lookup device registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deviceRegistrationToken, err := s.deviceRegistrationTokens.CreateDeviceRegistrationToken(
		r.Context(),
		projectID,
		createDeviceRegistrationTokenRequest.Name,
		createDeviceRegistrationTokenRequest.Description,
		createDeviceRegistrationTokenRequest.MaxRegistrations)
	if err != nil {
		log.WithError(err).Error("create device registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, deviceRegistrationToken)
}

func (s *Service) getDeviceRegistrationToken(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID, tokenID string,
) {
	deviceRegistrationToken, err := s.deviceRegistrationTokens.GetDeviceRegistrationToken(r.Context(), tokenID, projectID)
	if err != nil {
		log.WithError(err).Error("get device registration token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var ret interface{} = deviceRegistrationToken
	if _, ok := r.URL.Query()["full"]; ok {
		devicesRegistered, err := s.devicesRegisteredWithToken.GetDevicesRegisteredWithTokenCount(r.Context(), deviceRegistrationToken.ID, projectID)
		if err != nil {
			log.WithError(err).Error("get registered device count")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ret = models.DeviceRegistrationTokenFull{
			DeviceRegistrationToken: *deviceRegistrationToken,
			DeviceCounts:            *devicesRegistered,
		}
	}

	utils.Respond(w, ret)
}

func (s *Service) updateDeviceRegistrationToken(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID, tokenID string,
) {
	var updateDeviceRegistrationTokenRequest struct {
		Name             string `json:"name" validate:"name"`
		Description      string `json:"description" validate:"description"`
		MaxRegistrations *int   `json:"maxRegistrations"`
	}
	if err := read(r, &updateDeviceRegistrationTokenRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if deviceRegistrationToken, err := s.deviceRegistrationTokens.LookupDeviceRegistrationToken(r.Context(),
		updateDeviceRegistrationTokenRequest.Name, projectID); err == nil && deviceRegistrationToken.ID != tokenID {
		http.Error(w, store.ErrDeviceRegistrationTokenNameAlreadyInUse.Error(), http.StatusBadRequest)
		return
	} else if err != nil && err != store.ErrDeviceRegistrationTokenNotFound {
		log.WithError(err).Error("lookup device registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deviceRegistrationToken, err := s.deviceRegistrationTokens.UpdateDeviceRegistrationToken(
		r.Context(),
		tokenID,
		projectID,
		updateDeviceRegistrationTokenRequest.Name,
		updateDeviceRegistrationTokenRequest.Description,
		updateDeviceRegistrationTokenRequest.MaxRegistrations,
	)
	if err != nil {
		log.WithError(err).Error("update device registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, deviceRegistrationToken)
}

func (s *Service) deleteDeviceRegistrationToken(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	tokenID string,
) {
	if err := s.deviceRegistrationTokens.DeleteDeviceRegistrationToken(r.Context(), tokenID, projectID); err != nil {
		log.WithError(err).Error("delete device registration token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) listDeviceRegistrationTokens(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	deviceRegistrationTokens, err := s.deviceRegistrationTokens.ListDeviceRegistrationTokens(r.Context(), projectID)
	if err != nil {
		log.WithError(err).Error("list device registration tokens")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ret interface{} = deviceRegistrationTokens
	if _, ok := r.URL.Query()["full"]; ok {
		deviceRegistrationTokensFull := make([]models.DeviceRegistrationTokenFull, 0)

		for _, deviceRegistrationToken := range deviceRegistrationTokens {
			deviceRegistrationTokenCounts, err := s.devicesRegisteredWithToken.GetDevicesRegisteredWithTokenCount(r.Context(), deviceRegistrationToken.ID, projectID)
			if err != nil {
				log.WithError(err).Error("get count of devices registered with device registration token")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			deviceRegistrationTokensFull = append(deviceRegistrationTokensFull, models.DeviceRegistrationTokenFull{
				DeviceRegistrationToken: deviceRegistrationToken,
				DeviceCounts:            *deviceRegistrationTokenCounts,
			})
		}

		ret = deviceRegistrationTokensFull
	}

	utils.Respond(w, ret)
}

func (s *Service) setDeviceRegistrationTokenLabel(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceRegistrationTokenID string,
) {
	var setLabelRequest struct {
		Key   string `json:"key" validate:"labelkey"`
		Value string `json:"value" validate:"labelvalue"`
	}
	if err := read(r, &setLabelRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	label, err := s.deviceRegistrationTokens.SetDeviceRegistrationTokenLabel(
		r.Context(),
		deviceRegistrationTokenID,
		projectID,
		setLabelRequest.Key,
		setLabelRequest.Value,
	)
	if err != nil {
		log.WithError(err).Error("set device registration token label")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, label)
}

func (s *Service) deleteDeviceRegistrationTokenLabel(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID,
	deviceRegistrationTokenID string,
) {
	vars := mux.Vars(r)
	key := vars["key"]

	if err := s.deviceRegistrationTokens.DeleteDeviceRegistrationTokenLabel(r.Context(), deviceRegistrationTokenID, projectID, key); err != nil {
		log.WithError(err).Error("delete device registration token label")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) getProjectConfig(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	vars := mux.Vars(r)
	key := vars["key"]

	var value interface{}
	var err error
	switch key {
	case string(models.ProjectMetricsConfigKey):
		value, err = s.metricConfigs.GetProjectMetricsConfig(r.Context(), projectID)
	case string(models.DeviceMetricsConfigKey):
		value, err = s.metricConfigs.GetDeviceMetricsConfig(r.Context(), projectID)
	case string(models.ServiceMetricsConfigKey):
		value, err = s.metricConfigs.GetServiceMetricsConfigs(r.Context(), projectID)
	default:
		http.Error(w, store.ErrProjectConfigNotFound.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		log.WithError(err).Error("get project config with key " + key)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, value)
}

func (s *Service) setProjectConfig(w http.ResponseWriter, r *http.Request,
	projectID, authenticatedUserID, authenticatedServiceAccountID string,
) {
	vars := mux.Vars(r)
	key := vars["key"]

	var err error
	switch key {
	case string(models.ProjectMetricsConfigKey):
		var value models.ProjectMetricsConfig
		if err := read(r, &value); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.metricConfigs.SetProjectMetricsConfig(r.Context(), projectID, value)
	case string(models.DeviceMetricsConfigKey):
		var value models.DeviceMetricsConfig
		if err := read(r, &value); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.metricConfigs.SetDeviceMetricsConfig(r.Context(), projectID, value)
	case string(models.ServiceMetricsConfigKey):
		var values []models.ServiceMetricsConfig
		// TODO: use read() here
		if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.metricConfigs.SetServiceMetricsConfigs(r.Context(), projectID, values)
	default:
		http.Error(w, store.ErrProjectConfigNotFound.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		log.WithError(err).Error("set project config with key " + key)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func (s *Service) withDeviceAuth(handler func(http.ResponseWriter, *http.Request, models.Project, models.Device)) func(http.ResponseWriter, *http.Request) {
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

		project, err := s.projects.GetProject(r.Context(), projectID)
		if err != nil {
			log.WithError(err).Error("get project")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		device, err := s.devices.GetDevice(r.Context(), deviceAccessKey.DeviceID, projectID)
		if err != nil {
			log.WithError(err).Error("get device")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler(w, r, *project, *device)
	}
}

// TODO: verify project ID
func (s *Service) registerDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project"]

	var registerDeviceRequest models.RegisterDeviceRequest
	if err := read(r, &registerDeviceRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deviceRegistrationToken, err := s.deviceRegistrationTokens.GetDeviceRegistrationToken(r.Context(), registerDeviceRequest.DeviceRegistrationTokenID, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if deviceRegistrationToken.MaxRegistrations != nil {
		devicesRegisteredCount, err := s.devicesRegisteredWithToken.GetDevicesRegisteredWithTokenCount(r.Context(), registerDeviceRequest.DeviceRegistrationTokenID, projectID)
		if err != nil {
			log.WithError(err).Error("get devices registered with token count")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if devicesRegisteredCount.AllCount >= *deviceRegistrationToken.MaxRegistrations {
			log.WithError(err).Error("device allocation limit reached")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	device, err := s.devices.CreateDevice(r.Context(), projectID, namesgenerator.GetRandomName(), deviceRegistrationToken.ID, deviceRegistrationToken.Labels)
	if err != nil {
		log.WithError(err).Error("create device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deviceAccessKeyValue := ksuid.New().String()

	_, err = s.deviceAccessKeys.CreateDeviceAccessKey(r.Context(), projectID, device.ID, hash.Hash(deviceAccessKeyValue))
	if err != nil {
		log.WithError(err).Error("create device access key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, models.RegisterDeviceResponse{
		DeviceID:             device.ID,
		DeviceAccessKeyValue: deviceAccessKeyValue,
	})
}

func (s *Service) getBundle(w http.ResponseWriter, r *http.Request, project models.Project, device models.Device) {
	s.st.Incr("get_bundle", []string{
		fmt.Sprintf("project_id:%s", project.ID),
		fmt.Sprintf("project_name:%s", project.Name),
	}, 1)

	if err := s.devices.UpdateDeviceLastSeenAt(r.Context(), device.ID, project.ID); err != nil {
		log.WithError(err).Error("update device last seen at")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	applications, err := s.applications.ListApplications(r.Context(), project.ID)
	if err != nil {
		log.WithError(err).Error("list applications")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bundle := models.Bundle{
		DesiredAgentSpec:    device.DesiredAgentSpec,
		DesiredAgentVersion: device.DesiredAgentVersion,
	}

	for i, application := range applications {
		match, err := query.DeviceMatchesQuery(device, application.SchedulingRule)
		if err != nil {
			log.WithError(err).Error("evaluate application scheduling rule")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !match {
			continue
		}

		release, err := s.releases.GetLatestRelease(r.Context(), project.ID, application.ID)
		if err == store.ErrReleaseNotFound {
			continue
		} else if err != nil {
			log.WithError(err).Error("get latest release")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bundle.Applications = append(bundle.Applications, models.ApplicationFull2{
			Application:   applications[i],
			LatestRelease: *release,
		})
	}

	deviceApplicationStatuses, err := s.deviceApplicationStatuses.ListDeviceApplicationStatuses(
		r.Context(), project.ID, device.ID)
	if err == nil {
		bundle.ApplicationStatuses = deviceApplicationStatuses
	} else {
		log.WithError(err).Error("list device application statuses")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deviceServiceStatuses, err := s.deviceServiceStatuses.ListDeviceServiceStatuses(
		r.Context(), project.ID, device.ID)
	if err == nil {
		bundle.ServiceStatuses = deviceServiceStatuses
	} else {
		log.WithError(err).Error("list device service statuses")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Respond(w, bundle)
}

func (s *Service) setDeviceInfo(w http.ResponseWriter, r *http.Request, project models.Project, device models.Device) {
	var setDeviceInfoRequest models.SetDeviceInfoRequest
	if err := read(r, &setDeviceInfoRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := s.devices.SetDeviceInfo(r.Context(), device.ID, project.ID, setDeviceInfoRequest.DeviceInfo); err != nil {
		log.WithError(err).Error("set device info")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) setDeviceApplicationStatus(w http.ResponseWriter, r *http.Request, project models.Project, device models.Device) {
	vars := mux.Vars(r)
	applicationID := vars["application"]

	var setDeviceApplicationStatusRequest models.SetDeviceApplicationStatusRequest
	if err := read(r, &setDeviceApplicationStatusRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.deviceApplicationStatuses.SetDeviceApplicationStatus(r.Context(), project.ID, device.ID,
		applicationID, setDeviceApplicationStatusRequest.CurrentReleaseID,
	); err != nil {
		log.WithError(err).Error("set device application status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) deleteDeviceApplicationStatus(w http.ResponseWriter, r *http.Request, project models.Project, device models.Device) {
	vars := mux.Vars(r)
	applicationID := vars["application"]

	if err := s.deviceApplicationStatuses.DeleteDeviceApplicationStatus(r.Context(),
		project.ID, device.ID, applicationID,
	); err != nil {
		log.WithError(err).Error("delete device application status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) setDeviceServiceStatus(w http.ResponseWriter, r *http.Request, project models.Project, device models.Device) {
	vars := mux.Vars(r)
	applicationID := vars["application"]
	service := vars["service"]

	var setDeviceServiceStatusRequest models.SetDeviceServiceStatusRequest
	if err := read(r, &setDeviceServiceStatusRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.deviceServiceStatuses.SetDeviceServiceStatus(r.Context(), project.ID, device.ID,
		applicationID, service, setDeviceServiceStatusRequest.CurrentReleaseID,
	); err != nil {
		log.WithError(err).Error("set device service status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) deleteDeviceServiceStatus(w http.ResponseWriter, r *http.Request, project models.Project, device models.Device) {
	vars := mux.Vars(r)
	applicationID := vars["application"]
	service := vars["service"]

	if err := s.deviceServiceStatuses.DeleteDeviceServiceStatus(r.Context(),
		project.ID, device.ID, applicationID, service,
	); err != nil {
		log.WithError(err).Error("delete device service status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
