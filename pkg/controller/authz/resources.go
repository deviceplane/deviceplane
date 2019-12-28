package authz

type Resource string

const (
	ResourceAny = Resource("*")

	ResourceProjects                      = Resource("projects")
	ResourceRoles                         = Resource("roles")
	ResourceMemberships                   = Resource("memberships")
	ResourceMembershipRoleBindings        = Resource("membershiprolebindings")
	ResourceServiceAccounts               = Resource("serviceaccounts")
	ResourceServiceAccountAccessKeys      = Resource("serviceaccountaccesskeys")
	ResourceServiceAccountRoleBindings    = Resource("serviceaccountrolebindings")
	ResourceApplications                  = Resource("applications")
	ResourceReleases                      = Resource("releases")
	ResourceDevices                       = Resource("devices")
	ResourceDeviceLabels                  = Resource("devicelabels")
	ResourceDeviceRegistrationTokens      = Resource("deviceregistrationtokens")
	ResourceDeviceRegistrationTokenLabels = Resource("deviceregistrationtokenlabels")
	ResourceProjectConfigs                = Resource("projectconfigs")
)
