package service

import (
	"net/http"
	"net/http/pprof"
	"net/url"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/deviceplane/deviceplane/pkg/controller/connman"
	"github.com/deviceplane/deviceplane/pkg/controller/spaserver"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/email"
	"github.com/deviceplane/deviceplane/pkg/revdial"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Service struct {
	users                      store.Users
	internalUsers              store.InternalUsers
	externalUsers              store.ExternalUsers
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
	connections                store.Connections
	applications               store.Applications
	applicationDeviceCounts    store.ApplicationDeviceCounts
	releases                   store.Releases
	releaseDeviceCounts        store.ReleaseDeviceCounts
	deviceApplicationStatuses  store.DeviceApplicationStatuses
	deviceServiceStatuses      store.DeviceServiceStatuses
	deviceServiceStates        store.DeviceServiceStates
	metricConfigs              store.MetricConfigs
	email                      email.Interface
	emailFromName              string
	emailFromAddress           string
	allowedEmailDomains        []string
	auth0Domain                *url.URL
	auth0Audience              string
	st                         *statsd.Client
	connman                    *connman.ConnectionManager
	router                     *mux.Router
	upgrader                   websocket.Upgrader
}

func NewService(
	users store.Users,
	internalUsers store.InternalUsers,
	externalUsers store.ExternalUsers,
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
	connections store.Connections,
	applications store.Applications,
	applicationDeviceCounts store.ApplicationDeviceCounts,
	releases store.Releases,
	releasesDeviceCounts store.ReleaseDeviceCounts,
	deviceApplicationStatuses store.DeviceApplicationStatuses,
	deviceServiceStatuses store.DeviceServiceStatuses,
	deviceServiceStates store.DeviceServiceStates,
	metricConfigs store.MetricConfigs,
	email email.Interface,
	emailFromName string,
	emailFromAddress string,
	allowedEmailDomains []string,
	auth0Domain *url.URL,
	auth0Audience string,
	fileSystem http.FileSystem,
	st *statsd.Client,
	connman *connman.ConnectionManager,
	allowedOrigins []url.URL,
) *Service {
	s := &Service{
		users:                      users,
		internalUsers:              internalUsers,
		externalUsers:              externalUsers,
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
		connections:                connections,
		applications:               applications,
		applicationDeviceCounts:    applicationDeviceCounts,
		releases:                   releases,
		releaseDeviceCounts:        releasesDeviceCounts,
		deviceApplicationStatuses:  deviceApplicationStatuses,
		deviceServiceStatuses:      deviceServiceStatuses,
		deviceServiceStates:        deviceServiceStates,
		metricConfigs:              metricConfigs,
		email:                      email,
		emailFromName:              emailFromName,
		emailFromAddress:           emailFromAddress,
		allowedEmailDomains:        allowedEmailDomains,
		auth0Domain:                auth0Domain,
		auth0Audience:              auth0Audience,
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

	apiRouter.HandleFunc("/register", s.registerInternalUser).Methods("POST")
	apiRouter.HandleFunc("/registersso", s.registerExternalUser).Methods("POST")
	apiRouter.HandleFunc("/completeregistration", s.confirmRegistration).Methods("POST")

	apiRouter.HandleFunc("/changepassword", s.changeInternalUserPassword).Methods("POST")
	apiRouter.HandleFunc("/recoverpassword", s.recoverPassword).Methods("POST")
	apiRouter.HandleFunc("/passwordrecoverytokens/{passwordrecoverytokenvalue}", s.getPasswordRecoveryToken).Methods("GET")

	apiRouter.HandleFunc("/login", s.loginInternalUser).Methods("POST")
	apiRouter.HandleFunc("/loginsso", s.loginExternalUser).Methods("POST")
	apiRouter.HandleFunc("/logout", s.logout).Methods("POST")

	apiRouter.HandleFunc("/me", s.getMe).Methods("GET")
	apiRouter.HandleFunc("/me", s.updateMe).Methods("PATCH")

	apiRouter.HandleFunc("/memberships", s.listMembershipsByUser).Methods("GET")

	apiRouter.HandleFunc("/useraccesskeys", s.listUserAccessKeys).Methods("GET")
	apiRouter.HandleFunc("/useraccesskeys", s.createUserAccessKey).Methods("POST")
	apiRouter.HandleFunc("/useraccesskeys/{useraccesskey}", s.getUserAccessKey).Methods("GET")
	apiRouter.HandleFunc("/useraccesskeys/{useraccesskey}", s.deleteUserAccessKey).Methods("DELETE")

	apiRouter.HandleFunc("/projects", s.createProject).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}", s.getProject).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}", s.updateProject).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}", s.deleteProject).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/roles", s.createRole).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/roles", s.listRoles).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/roles/{role}", s.getRole).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/roles/{role}", s.updateRole).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/roles/{role}", s.deleteRole).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/memberships", s.listMembershipsByProject).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/memberships", s.createMembership).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}", s.getMembership).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}", s.deleteMembership).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/memberships/{user}/membershiprolebindings", s.listMembershipRoleBindings).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}/roles/{role}/membershiprolebindings", s.createMembershipRoleBinding).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}/roles/{role}/membershiprolebindings", s.getMembershipRoleBinding).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}/roles/{role}/membershiprolebindings", s.deleteMembershipRoleBinding).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/serviceaccounts", s.listServiceAccounts).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts", s.createServiceAccount).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}", s.getServiceAccount).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}", s.updateServiceAccount).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}", s.deleteServiceAccount).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys", s.listServiceAccountAccessKeys).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys", s.createServiceAccountAccessKey).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys/{serviceaccountaccesskey}", s.getServiceAccountAccessKey).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys/{serviceaccountaccesskey}", s.deleteServiceAccountAccessKey).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountrolebindings", s.listServiceAccountRoleBindings).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/roles/{role}/serviceaccountrolebindings", s.createServiceAccountRoleBinding).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/roles/{role}/serviceaccountrolebindings", s.getServiceAccountRoleBinding).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/roles/{role}/serviceaccountrolebindings", s.deleteServiceAccountRoleBinding).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/connections", s.listConnections).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/connections", s.createConnection).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/connections/{connection}", s.getConnection).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/connections/{connection}", s.updateConnection).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/connections/{connection}", s.deleteConnection).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/applications", s.listApplications).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/applications", s.createApplication).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}", s.getApplication).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}", s.updateApplication).Methods("PATCH")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}", s.deleteApplication).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/applications/{application}/releases", s.createRelease).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}/releases/latest", s.getLatestRelease).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}/releases/{release}", s.getRelease).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}/releases", s.listReleases).Methods("GET")

	apiRouter.HandleFunc("/projects/{project}/devices/{device}", s.getDevice).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices", s.listDevices).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/previewscheduling/{application}", s.previewScheduledDevices).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}", s.updateDevice).Methods("PATCH")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}", s.deleteDevice).Methods("DELETE")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/ssh", s.initiateSSH).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/reboot", s.initiateReboot).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/imagepullprogress", s.imagePullProgress).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/metrics/host", s.hostMetrics).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/metrics/agent", s.agentMetrics).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/metrics", s.serviceMetrics).Methods("GET")
	apiRouter.PathPrefix("/projects/{project}/devices/{device}/debug/").HandlerFunc(s.deviceDebug)

	apiRouter.HandleFunc("/projects/{project}/devices/{device}/environmentvariables", s.setDeviceEnvironmentVariable).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/environmentvariables/{key}", s.deleteDeviceEnvironmentVariable).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/devices/{device}/labels", s.setDeviceLabel).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/labels/{key}", s.deleteDeviceLabel).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/devicelabels", s.listAllDeviceLabelKeys).Methods("GET")

	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens", s.listDeviceRegistrationTokens).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens", s.createDeviceRegistrationToken).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}", s.getDeviceRegistrationToken).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}", s.updateDeviceRegistrationToken).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}", s.deleteDeviceRegistrationToken).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}/environmentvariables", s.setDeviceRegistrationTokenEnvironmentVariable).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}/environmentvariables/{key}", s.deleteDeviceRegistrationTokenEnvironmentVariable).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}/labels", s.setDeviceRegistrationTokenLabel).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}/labels/{key}", s.deleteDeviceRegistrationTokenLabel).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/configs/{key}", s.getProjectConfig).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/configs/{key}", s.setProjectConfig).Methods("PUT")

	apiRouter.HandleFunc("/projects/{project}/devices/register", s.registerDevice).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/bundle", s.getBundle).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/info", s.setDeviceInfo).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/deviceapplicationstatuses", s.setDeviceApplicationStatus).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/deviceapplicationstatuses", s.deleteDeviceApplicationStatus).Methods("DELETE")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/deviceservicestatuses", s.setDeviceServiceStatus).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/deviceservicestatuses", s.deleteDeviceServiceStatus).Methods("DELETE")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/deviceservicestates", s.setDeviceServiceState).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/deviceservicestates", s.deleteDeviceServiceState).Methods("DELETE")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/forwardmetrics/service", s.forwardServiceMetrics).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/forwardmetrics/device", s.forwardDeviceMetrics).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/connection", s.initiateDeviceConnection).Methods("GET")

	apiRouter.Handle("/revdial", revdial.ConnHandler(s.upgrader)).Methods("GET")

	debugRouter := apiRouter.PathPrefix("/debug/").Subrouter()

	debugRouter.Use(s.GorillaSuperUserAuth)
	debugRouter.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	debugRouter.HandleFunc("/pprof/profile", pprof.Profile)
	debugRouter.HandleFunc("/pprof/symbol", pprof.Symbol)
	debugRouter.HandleFunc("/pprof/trace", pprof.Trace)
	debugRouter.PathPrefix("/pprof/").Handler(http.StripPrefix("/api", http.HandlerFunc(pprof.Index)))

	apiRouter.HandleFunc("/health", s.health).Methods("GET")
	apiRouter.HandleFunc("/500", s.intentional500).Methods("GET")

	s.router.PathPrefix("/api").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	s.router.PathPrefix("/").Handler(spaserver.NewSPAFileServer(fileSystem))

	return s
}
