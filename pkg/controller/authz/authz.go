package authz

type Config struct {
	Rules []Rule `yaml:"rules,omitempty"`
}

type Rule struct {
	Resources       []string `yaml:"resources,omitempty"`
	Actions         []string `yaml:"actions,omitempty"`
	ParentResources []string `yaml:"parent_resources,omitempty"`
	Effect          string   `yaml:"effect,omitempty"`
}

var (
	ReadAllRole = Config{
		Rules: []Rule{
			{
				Resources: []string{"*"},
				Actions:   []string{"read"},
			},
		},
	}
	WriteAllRole = Config{
		Rules: []Rule{
			{
				Resources: []string{"*"},
				Actions:   []string{"write"},
			},
		},
	}
	AdminAllRole = Config{
		Rules: []Rule{
			{
				Resources: []string{"*"},
				Actions:   []string{"admin"},
			},
		},
	}
)

var (
	readActions = []string{
		"GetProject",
		"GetRole",
		"ListRoles",
		"GetMembership",
		"ListMembershipsByProject",
		"GetMembershipRoleBinding",
		"ListMembershipRoleBindings",
		"GetServiceAccount",
		"ListServiceAccounts",
		"GetApplication",
		"ListApplications",
		"GetLatestRelease",
		"GetRelease",
		"ListReleases",
		"GetDevice",
		"ListDevices",
		"GetDeviceLabel",
		"ListDeviceLabels",
		"GetDeviceRegistrationToken",
		"ListDeviceRegistrationTokens",
	}
	writeActions = append(readActions, []string{
		"CreateApplication",
		"UpdateApplication",
		"DeleteApplication",
		"CreateRelease",
		"UpdateDevice",
		"DeleteDevice",
		"SSH",
		"SetDeviceLabel",
		"DeleteDeviceLabel",
		"CreateDeviceRegistrationToken",
		"UpdateDeviceRegistrationToken",
		"DeleteDeviceRegistrationToken",
		"SetDeviceRegistrationTokenLabel",
		"DeleteDeviceRegistrationTokenLabel",
	}...)
	adminActions = append(writeActions, []string{
		"UpdateProject",
		"DeleteProject",
		"CreateRole",
		"UpdateRole",
		"DeleteRole",
		"CreateMembership",
		"DeleteMembership",
		"CreateMembershipRoleBinding",
		"DeleteMembershipRoleBinding",
		"CreateServiceAccount",
		"UpdateServiceAccount",
		"DeleteServiceAccount",
		"CreateServiceAccountAccessKey",
		"GetServiceAccountAccessKey",
		"ListServiceAccountAccessKeys",
		"DeleteServiceAccountAccessKey",
		"CreateServiceAccountRoleBinding",
		"DeleteServiceAccountRoleBinding",
	}...)
)

func Evaluate(requestedResource, requestedAction string, configs []Config) bool {
	oneAllow := false
	for _, config := range configs {
		for _, rule := range config.Rules {
			rule = resolveRule(rule)
			for _, ruleResource := range rule.Resources {
				if ruleResource == requestedResource || ruleResource == "*" {
					for _, ruleAction := range rule.Actions {
						if ruleAction == requestedAction {
							oneAllow = true
						}
					}
				}
			}
		}
	}
	return oneAllow
}

func resolveRule(rule Rule) Rule {
	var finalActions []string
	for _, action := range rule.Actions {
		switch action {
		case "read":
			finalActions = append(finalActions, readActions...)
		case "write":
			finalActions = append(finalActions, writeActions...)
		case "admin":
			finalActions = append(finalActions, adminActions...)
		default:
			finalActions = append(finalActions, action)
		}
	}
	rule.Actions = finalActions
	return rule
}
