package authz

type Action string

const (
	ActionReadAll  = Action("read")
	ActionWriteAll = Action("write")
	ActionAdminAll = Action("admin")

	ActionGetProject                   = Action("GetProject")
	ActionGetRole                      = Action("GetRole")
	ActionListRoles                    = Action("ListRoles")
	ActionGetMembership                = Action("GetMembership")
	ActionListMembershipsByProject     = Action("ListMembershipsByProject")
	ActionGetMembershipRoleBinding     = Action("GetMembershipRoleBinding")
	ActionListMembershipRoleBindings   = Action("ListMembershipRoleBindings")
	ActionGetServiceAccount            = Action("GetServiceAccount")
	ActionListServiceAccounts          = Action("ListServiceAccounts")
	ActionGetApplication               = Action("GetApplication")
	ActionListApplications             = Action("ListApplications")
	ActionGetLatestRelease             = Action("GetLatestRelease")
	ActionGetRelease                   = Action("GetRelease")
	ActionListReleases                 = Action("ListReleases")
	ActionGetDevice                    = Action("GetDevice")
	ActionListDevices                  = Action("ListDevices")
	ActionGetImagePullProgress         = Action("GetImagePullProgress")
	ActionGetMetrics                   = Action("GetMetrics")
	ActionGetServiceMetrics            = Action("GetServiceMetrics")
	ActionGetDeviceRegistrationToken   = Action("GetDeviceRegistrationToken")
	ActionListDeviceRegistrationTokens = Action("ListDeviceRegistrationTokens")
	ActionGetProjectConfig             = Action("GetProjectConfig")

	ActionCreateApplication                  = Action("CreateApplication")
	ActionUpdateApplication                  = Action("UpdateApplication")
	ActionDeleteApplication                  = Action("DeleteApplication")
	ActionCreateRelease                      = Action("CreateRelease")
	ActionUpdateDevice                       = Action("UpdateDevice")
	ActionDeleteDevice                       = Action("DeleteDevice")
	ActionSSH                                = Action("SSH")
	ActionReboot                             = Action("Reboot")
	ActionSetDeviceLabel                     = Action("SetDeviceLabel")
	ActionDeleteDeviceLabel                  = Action("DeleteDeviceLabel")
	ActionCreateDeviceRegistrationToken      = Action("CreateDeviceRegistrationToken")
	ActionUpdateDeviceRegistrationToken      = Action("UpdateDeviceRegistrationToken")
	ActionDeleteDeviceRegistrationToken      = Action("DeleteDeviceRegistrationToken")
	ActionSetDeviceRegistrationTokenLabel    = Action("SetDeviceRegistrationTokenLabel")
	ActionDeleteDeviceRegistrationTokenLabel = Action("DeleteDeviceRegistrationTokenLabel")

	ActionUpdateProject                   = Action("UpdateProject")
	ActionDeleteProject                   = Action("DeleteProject")
	ActionCreateRole                      = Action("CreateRole")
	ActionUpdateRole                      = Action("UpdateRole")
	ActionDeleteRole                      = Action("DeleteRole")
	ActionCreateMembership                = Action("CreateMembership")
	ActionDeleteMembership                = Action("DeleteMembership")
	ActionCreateMembershipRoleBinding     = Action("CreateMembershipRoleBinding")
	ActionDeleteMembershipRoleBinding     = Action("DeleteMembershipRoleBinding")
	ActionCreateServiceAccount            = Action("CreateServiceAccount")
	ActionUpdateServiceAccount            = Action("UpdateServiceAccount")
	ActionDeleteServiceAccount            = Action("DeleteServiceAccount")
	ActionCreateServiceAccountAccessKey   = Action("CreateServiceAccountAccessKey")
	ActionGetServiceAccountAccessKey      = Action("GetServiceAccountAccessKey")
	ActionListServiceAccountAccessKeys    = Action("ListServiceAccountAccessKeys")
	ActionDeleteServiceAccountAccessKey   = Action("DeleteServiceAccountAccessKey")
	ActionCreateServiceAccountRoleBinding = Action("CreateServiceAccountRoleBinding")
	ActionGetServiceAccountRoleBinding    = Action("GetServiceAccountRoleBinding")
	ActionListServiceAccountRoleBinding   = Action("ListServiceAccountRoleBinding")
	ActionDeleteServiceAccountRoleBinding = Action("DeleteServiceAccountRoleBinding")
	ActionSetProjectConfig                = Action("SetProjectConfig")
)

var (
	readActions = []Action{
		ActionGetProject,
		ActionGetRole,
		ActionListRoles,
		ActionGetMembership,
		ActionListMembershipsByProject,
		ActionGetMembershipRoleBinding,
		ActionListMembershipRoleBindings,
		ActionGetServiceAccount,
		ActionListServiceAccounts,
		ActionGetApplication,
		ActionListApplications,
		ActionGetLatestRelease,
		ActionGetRelease,
		ActionListReleases,
		ActionGetDevice,
		ActionListDevices,
		ActionGetImagePullProgress,
		ActionGetMetrics,
		ActionGetServiceMetrics,
		ActionGetDeviceRegistrationToken,
		ActionListDeviceRegistrationTokens,
		ActionGetProjectConfig,
	}
	writeActions = append(readActions, []Action{
		ActionCreateApplication,
		ActionUpdateApplication,
		ActionDeleteApplication,
		ActionCreateRelease,
		ActionUpdateDevice,
		ActionDeleteDevice,
		ActionSSH,
		ActionReboot,
		ActionSetDeviceLabel,
		ActionDeleteDeviceLabel,
		ActionCreateDeviceRegistrationToken,
		ActionUpdateDeviceRegistrationToken,
		ActionDeleteDeviceRegistrationToken,
		ActionSetDeviceRegistrationTokenLabel,
		ActionDeleteDeviceRegistrationTokenLabel,
	}...)
	adminActions = append(writeActions, []Action{
		ActionUpdateProject,
		ActionDeleteProject,
		ActionCreateRole,
		ActionUpdateRole,
		ActionDeleteRole,
		ActionCreateMembership,
		ActionDeleteMembership,
		ActionCreateMembershipRoleBinding,
		ActionDeleteMembershipRoleBinding,
		ActionCreateServiceAccount,
		ActionUpdateServiceAccount,
		ActionDeleteServiceAccount,
		ActionCreateServiceAccountAccessKey,
		ActionGetServiceAccountAccessKey,
		ActionListServiceAccountAccessKeys,
		ActionDeleteServiceAccountAccessKey,
		ActionCreateServiceAccountRoleBinding,
		ActionDeleteServiceAccountRoleBinding,
		ActionSetProjectConfig,
	}...)

	actionsByResource = map[Resource][]Action{}
)
