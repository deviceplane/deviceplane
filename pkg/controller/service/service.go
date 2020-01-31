package service

import (
	"net/http"
	"net/http/pprof"
	"net/url"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/deviceplane/deviceplane/pkg/controller/authz"
	"github.com/deviceplane/deviceplane/pkg/controller/connman"
	"github.com/deviceplane/deviceplane/pkg/controller/spaserver"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/email"
	"github.com/deviceplane/deviceplane/pkg/revdial"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
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

	apiRouter.HandleFunc("/register", s.initWith()(s.register)).Methods("POST")
	apiRouter.HandleFunc("/completeregistration", s.initWith()(s.confirmRegistration)).Methods("POST")

	apiRouter.HandleFunc("/recoverpassword", s.initWith()(s.recoverPassword)).Methods("POST")
	apiRouter.HandleFunc("/passwordrecoverytokens/{passwordrecoverytokenvalue}", s.initWith()(s.getPasswordRecoveryToken)).Methods("GET")
	apiRouter.HandleFunc("/changepassword", s.initWith()(s.changePassword)).Methods("POST")

	apiRouter.HandleFunc("/login", s.initWith()(s.login)).Methods("POST")
	apiRouter.HandleFunc("/logout", s.initWith()(s.logout)).Methods("POST")

	apiRouter.HandleFunc("/me", s.initWith(s.withUserOrServiceAccountAuth)(s.getMe)).Methods("GET")
	apiRouter.HandleFunc("/me", s.initWith(s.withUserOrServiceAccountAuth)(s.updateMe)).Methods("PATCH")

	apiRouter.HandleFunc("/memberships", s.initWith(s.withUserOrServiceAccountAuth)(s.listMembershipsByUser)).Methods("GET")

	apiRouter.HandleFunc("/useraccesskeys", s.initWith(s.withUserOrServiceAccountAuth)(s.createUserAccessKey)).Methods("POST")
	apiRouter.HandleFunc("/useraccesskeys/{useraccesskey}", s.initWith(s.withUserOrServiceAccountAuth)(s.getUserAccessKey)).Methods("GET")
	apiRouter.HandleFunc("/useraccesskeys", s.initWith(s.withUserOrServiceAccountAuth)(s.listUserAccessKeys)).Methods("GET")
	apiRouter.HandleFunc("/useraccesskeys/{useraccesskey}", s.initWith(s.withUserOrServiceAccountAuth)(s.deleteUserAccessKey)).Methods("DELETE")

	apiRouter.HandleFunc("/projects", s.initWith(s.withUserOrServiceAccountAuth)(s.createProject)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}", s.initWith(s.validateAuthorization(authz.ResourceProjects, authz.ActionGetProject))(s.getProject)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}", s.initWith(s.validateAuthorization(authz.ResourceProjects, authz.ActionUpdateProject))(s.updateProject)).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}", s.initWith(s.validateAuthorization(authz.ResourceProjects, authz.ActionDeleteProject))(s.deleteProject)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/roles", s.initWith(s.validateAuthorization(authz.ResourceRoles, authz.ActionCreateRole))(s.createRole)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/roles/{role}", s.initWith(s.validateAuthorization(authz.ResourceRoles, authz.ActionGetRole), s.withRole)(s.getRole)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/roles", s.initWith(s.validateAuthorization(authz.ResourceRoles, authz.ActionListRoles))(s.listRoles)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/roles/{role}", s.initWith(s.validateAuthorization(authz.ResourceRoles, authz.ActionUpdateRole))(s.updateRole)).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/roles/{role}", s.initWith(s.validateAuthorization(authz.ResourceRoles, authz.ActionDeleteRole))(s.deleteRole)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/memberships", s.initWith(s.validateAuthorization(authz.ResourceMemberships, authz.ActionCreateMembership))(s.createMembership)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}", s.initWith(s.validateAuthorization(authz.ResourceMemberships, authz.ActionGetMembership))(s.getMembership)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/memberships", s.initWith(s.validateAuthorization(authz.ResourceMemberships, authz.ActionListMembershipsByProject))(s.listMembershipsByProject)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}", s.initWith(s.validateAuthorization(authz.ResourceMemberships, authz.ActionDeleteMembership))(s.deleteMembership)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/memberships/{user}/roles/{role}/membershiprolebindings", s.initWith(s.validateAuthorization(authz.ResourceMembershipRoleBindings, authz.ActionCreateMembershipRoleBinding), s.withRole)(s.createMembershipRoleBinding)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}/roles/{role}/membershiprolebindings", s.initWith(s.validateAuthorization(authz.ResourceMembershipRoleBindings, authz.ActionGetMembershipRoleBinding), s.withRole)(s.getMembershipRoleBinding)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}/membershiprolebindings", s.initWith(s.validateAuthorization(authz.ResourceMembershipRoleBindings, authz.ActionListMembershipRoleBindings))(s.listMembershipRoleBindings)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/memberships/{user}/roles/{role}/membershiprolebindings", s.initWith(s.validateAuthorization(authz.ResourceMembershipRoleBindings, authz.ActionDeleteMembershipRoleBinding), s.withRole)(s.deleteMembershipRoleBinding)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/serviceaccounts", s.initWith(s.validateAuthorization(authz.ResourceServiceAccounts, authz.ActionCreateServiceAccount))(s.createServiceAccount)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}", s.initWith(s.validateAuthorization(authz.ResourceServiceAccounts, authz.ActionGetServiceAccount), s.withServiceAccount)(s.getServiceAccount)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts", s.initWith(s.validateAuthorization(authz.ResourceServiceAccounts, authz.ActionListServiceAccounts))(s.listServiceAccounts)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}", s.initWith(s.validateAuthorization(authz.ResourceServiceAccounts, authz.ActionUpdateServiceAccount), s.withServiceAccount)(s.updateServiceAccount)).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}", s.initWith(s.validateAuthorization(authz.ResourceServiceAccounts, authz.ActionDeleteServiceAccount), s.withServiceAccount)(s.deleteServiceAccount)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys", s.initWith(s.validateAuthorization(authz.ResourceServiceAccountAccessKeys, authz.ActionCreateServiceAccountAccessKey))(s.createServiceAccountAccessKey)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys/{serviceaccountaccesskey}", s.initWith(s.validateAuthorization(authz.ResourceServiceAccountAccessKeys, authz.ActionGetServiceAccountAccessKey))(s.getServiceAccountAccessKey)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys", s.initWith(s.validateAuthorization(authz.ResourceServiceAccountAccessKeys, authz.ActionListServiceAccountAccessKeys))(s.listServiceAccountAccessKeys)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountaccesskeys/{serviceaccountaccesskey}", s.initWith(s.validateAuthorization(authz.ResourceServiceAccountAccessKeys, authz.ActionDeleteServiceAccountAccessKey))(s.deleteServiceAccountAccessKey)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/roles/{role}/serviceaccountrolebindings", s.initWith(s.validateAuthorization(authz.ResourceServiceAccountRoleBindings, authz.ActionCreateServiceAccountRoleBinding), s.withRole)(s.createServiceAccountRoleBinding)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/roles/{role}/serviceaccountrolebindings", s.initWith(s.validateAuthorization(authz.ResourceServiceAccountRoleBindings, authz.ActionGetServiceAccountRoleBinding), s.withRole)(s.getServiceAccountRoleBinding)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/serviceaccountrolebindings", s.initWith(s.validateAuthorization(authz.ResourceServiceAccountRoleBindings, authz.ActionListServiceAccountRoleBinding))(s.listServiceAccountRoleBindings)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/serviceaccounts/{serviceaccount}/roles/{role}/serviceaccountrolebindings", s.initWith(s.validateAuthorization(authz.ResourceServiceAccountRoleBindings, authz.ActionDeleteServiceAccountRoleBinding), s.withRole)(s.deleteServiceAccountRoleBinding)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/applications", s.initWith(s.validateAuthorization(authz.ResourceApplications, authz.ActionCreateApplication))(s.createApplication)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}", s.initWith(s.validateAuthorization(authz.ResourceApplications, authz.ActionGetApplication), s.withApplication)(s.getApplication)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/applications", s.initWith(s.validateAuthorization(authz.ResourceApplications, authz.ActionListApplications))(s.listApplications)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}", s.initWith(s.validateAuthorization(authz.ResourceApplications, authz.ActionUpdateApplication), s.withApplication)(s.updateApplication)).Methods("PATCH")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}", s.initWith(s.validateAuthorization(authz.ResourceApplications, authz.ActionDeleteApplication), s.withApplication)(s.deleteApplication)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/applications/{application}/releases", s.initWith(s.validateAuthorization(authz.ResourceReleases, authz.ActionCreateRelease), s.withApplication)(s.createRelease)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}/releases/{release}", s.initWith(s.validateAuthorization(authz.ResourceReleases, authz.ActionGetRelease), s.withApplication, s.withRelease)(s.getRelease)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/applications/{application}/releases", s.initWith(s.validateAuthorization(authz.ResourceReleases, authz.ActionListReleases), s.withApplication)(s.listReleases)).Methods("GET")

	apiRouter.HandleFunc("/projects/{project}/devices/{device}", s.initWith(s.validateAuthorization(authz.ResourceDevices, authz.ActionGetDevice), s.withDevice)(s.getDevice)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices", s.initWith(s.validateAuthorization(authz.ResourceDevices, authz.ActionListDevices))(s.listDevices)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/previewscheduling/{application}", s.initWith(s.validateAuthorization(authz.ResourceDevices, authz.ActionPreviewApplicationScheduling))(s.previewScheduledDevices)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}", s.initWith(s.validateAuthorization(authz.ResourceDevices, authz.ActionUpdateDevice), s.withDevice)(s.updateDevice)).Methods("PATCH")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}", s.initWith(s.validateAuthorization(authz.ResourceDevices, authz.ActionDeleteDevice), s.withDevice)(s.deleteDevice)).Methods("DELETE")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/ssh", s.initWith(s.validateAuthorization(authz.ResourceDevices, authz.ActionSSH), s.withDevice, s.withHijackedWebSocketConnection, s.withDeviceConnection)(s.initiateSSH)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/reboot", s.initWith(s.validateAuthorization(authz.ResourceDevices, authz.ActionReboot), s.withDevice, s.withDeviceConnection)(s.initiateReboot)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/imagepullprogress", s.initWith(s.validateAuthorization(authz.ResourceDevices, authz.ActionGetImagePullProgress), s.withDevice, s.withDeviceConnection)(s.imagePullProgress)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/metrics/host", s.initWith(s.validateAuthorization(authz.ResourceDevices, authz.ActionGetMetrics), s.withDevice, s.withDeviceConnection)(s.hostMetrics)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/metrics/agent", s.initWith(s.validateAuthorization(authz.ResourceDevices, authz.ActionGetMetrics), s.withDevice, s.withDeviceConnection)(s.agentMetrics)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/metrics", s.initWith(s.validateAuthorization(authz.ResourceDevices, authz.ActionGetServiceMetrics), s.withApplication, s.withDevice, s.withDeviceConnection)(s.serviceMetrics)).Methods("GET")
	apiRouter.PathPrefix("/projects/{project}/devices/{device}/debug/").HandlerFunc(s.initWith(s.validateAuthorization(authz.ResourceDevices, authz.ActionGetMetrics), s.withDevice, s.withDeviceConnection)(s.deviceDebug))

	apiRouter.HandleFunc("/projects/{project}/devices/{device}/labels", s.initWith(s.validateAuthorization(authz.ResourceDeviceLabels, authz.ActionSetDeviceLabel), s.withDevice)(s.setDeviceLabel)).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/labels/{key}", s.initWith(s.validateAuthorization(authz.ResourceDeviceLabels, authz.ActionDeleteDeviceLabel), s.withDevice)(s.deleteDeviceLabel)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/devicelabels", s.initWith(s.validateAuthorization(authz.ResourceDeviceLabels, authz.ActionListAllDeviceLabels))(s.listAllDeviceLabelKeys)).Methods("GET")

	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens", s.initWith(s.validateAuthorization(authz.ResourceDeviceRegistrationTokens, authz.ActionListDeviceRegistrationTokens))(s.listDeviceRegistrationTokens)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens", s.initWith(s.validateAuthorization(authz.ResourceDeviceRegistrationTokens, authz.ActionCreateDeviceRegistrationToken))(s.createDeviceRegistrationToken)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}", s.initWith(s.validateAuthorization(authz.ResourceDeviceRegistrationTokens, authz.ActionGetDeviceRegistrationToken), s.withDeviceRegistrationToken)(s.getDeviceRegistrationToken)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}", s.initWith(s.validateAuthorization(authz.ResourceDeviceRegistrationTokens, authz.ActionUpdateDeviceRegistrationToken), s.withDeviceRegistrationToken)(s.updateDeviceRegistrationToken)).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}", s.initWith(s.validateAuthorization(authz.ResourceDeviceRegistrationTokens, authz.ActionDeleteDeviceRegistrationToken), s.withDeviceRegistrationToken)(s.deleteDeviceRegistrationToken)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}/labels", s.initWith(s.validateAuthorization(authz.ResourceDeviceRegistrationTokenLabels, authz.ActionSetDeviceRegistrationTokenLabel), s.withDeviceRegistrationToken)(s.setDeviceRegistrationTokenLabel)).Methods("PUT")
	apiRouter.HandleFunc("/projects/{project}/deviceregistrationtokens/{deviceregistrationtoken}/labels/{key}", s.initWith(s.validateAuthorization(authz.ResourceDeviceRegistrationTokenLabels, authz.ActionDeleteDeviceRegistrationTokenLabel), s.withDeviceRegistrationToken)(s.deleteDeviceRegistrationTokenLabel)).Methods("DELETE")

	apiRouter.HandleFunc("/projects/{project}/configs/{key}", s.initWith(s.validateAuthorization(authz.ResourceProjectConfigs, authz.ActionGetProjectConfig))(s.getProjectConfig)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/configs/{key}", s.initWith(s.validateAuthorization(authz.ResourceProjectConfigs, authz.ActionSetProjectConfig))(s.setProjectConfig)).Methods("PUT")

	apiRouter.HandleFunc("/projects/{project}/devices/register", s.initWith()(s.registerDevice)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/bundle", s.initWith(s.withDeviceAuth)(s.getBundle)).Methods("GET")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/info", s.initWith(s.withDeviceAuth)(s.setDeviceInfo)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/deviceapplicationstatuses", s.initWith(s.withDeviceAuth)(s.setDeviceApplicationStatus)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/deviceapplicationstatuses", s.initWith(s.withDeviceAuth)(s.deleteDeviceApplicationStatus)).Methods("DELETE")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/deviceservicestatuses", s.initWith(s.withDeviceAuth)(s.setDeviceServiceStatus)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/applications/{application}/services/{service}/deviceservicestatuses", s.initWith(s.withDeviceAuth)(s.deleteDeviceServiceStatus)).Methods("DELETE")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/forwardmetrics/service", s.InitWith(s.withDeviceAuth)(s.forwardServiceMetrics)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/forwardmetrics/device", s.InitWith(s.withDeviceAuth)(s.forwardDeviceMetrics)).Methods("POST")
	apiRouter.HandleFunc("/projects/{project}/devices/{device}/connection", s.InitWith(s.withDeviceAuth)(s.initiateDeviceConnection)).Methods("GET")

	apiRouter.Handle("/revdial", revdial.ConnHandler(s.upgrader)).Methods("GET")

	debugRouter := apiRouter.PathPrefix("/debug/").Subrouter()
	debugRouter.HandleFunc("/pprof/cmdline", s.initWith(s.withSuperUserAuth)(HandlerFunc(pprof.Cmdline)))
	debugRouter.HandleFunc("/pprof/profile", s.initWith(s.withSuperUserAuth)(HandlerFunc(pprof.Profile)))
	debugRouter.HandleFunc("/pprof/symbol", s.initWith(s.withSuperUserAuth)(HandlerFunc(pprof.Symbol)))
	debugRouter.HandleFunc("/pprof/trace", s.initWith(s.withSuperUserAuth)(HandlerFunc(pprof.Trace)))
	debugRouter.PathPrefix("/pprof/").Handler(http.StripPrefix("/api", http.HandlerFunc(s.initWith(s.withSuperUserAuth)(HandlerFunc(pprof.Index)))))

	apiRouter.HandleFunc("/health", s.initWith()(s.health)).Methods("GET")
	apiRouter.HandleFunc("/500", s.initWith(s.withSuperUserAuth)(s.intentional500)).Methods("GET")

	s.router.PathPrefix("/api").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	s.router.PathPrefix("/").Handler(spaserver.NewSPAFileServer(fileSystem))

	return s
}
