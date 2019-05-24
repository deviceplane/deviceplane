package authz

type Config struct {
	Rules []Rule
}

type Rule struct {
	Resources       []string
	Actions         []string
	ParentResources []string
	Effect          string
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
	}
	writeActions = append(readActions, []string{
		"CreateApplication",
		"UpdateApplication",
		"CreateRelease",
		"SetDeviceLabel",
		"DeleteDeviceLabel",
		"CreateDeviceRegistrationToken",
	}...)
	adminActions = append(writeActions, []string{
		"CreateRole",
		"UpdateRole",
		"CreateMembership",
		"CreateMembershipRoleBinding",
		"CreateServiceAccount",
	}...)
)

func Evaluate(requestedResource, requestedAction string, configs []Config) (bool, error) {
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
	return oneAllow, nil
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
