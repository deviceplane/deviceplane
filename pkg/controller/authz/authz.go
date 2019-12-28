package authz

type Config struct {
	Rules []Rule `yaml:"rules,omitempty"`
}

type Rule struct {
	Resources []Resource `yaml:"resources,omitempty"`
	Actions   []Action   `yaml:"actions,omitempty"`
	Effect    Effect     `yaml:"effect,omitempty"`
}

var (
	ReadAllRole = Config{
		Rules: []Rule{
			{
				Resources: []Resource{ResourceAny},
				Actions:   []Action{ActionReadAll},
			},
		},
	}
	WriteAllRole = Config{
		Rules: []Rule{
			{
				Resources: []Resource{ResourceAny},
				Actions:   []Action{ActionWriteAll},
			},
		},
	}
	AdminAllRole = Config{
		Rules: []Rule{
			{
				Resources: []Resource{ResourceAny},
				Actions:   []Action{ActionAdminAll},
			},
		},
	}
)

func Evaluate(requestedResource Resource, requestedAction Action, configs []Config) bool {
	oneAllow := false
	oneDeny := false
	for _, config := range configs {
		for _, rule := range config.Rules {
			rule = resolveRule(rule)
			for _, ruleResource := range rule.Resources {
				if ruleResource == requestedResource || ruleResource == ResourceAny {
					for _, ruleAction := range rule.Actions {
						if ruleAction == requestedAction {
							if rule.Effect == EffectDeny {
								oneDeny = true
							} else {
								oneAllow = true
							}
						}
					}
				}
			}
		}
	}
	return oneAllow && !oneDeny
}

func resolveRule(rule Rule) Rule {
	var finalActions []Action
	for _, action := range rule.Actions {
		switch action {
		case ActionReadAll:
			finalActions = append(finalActions, readActions...)
		case ActionWriteAll:
			finalActions = append(finalActions, writeActions...)
		case ActionAdminAll:
			finalActions = append(finalActions, adminActions...)
		default:
			finalActions = append(finalActions, action)
		}
	}
	rule.Actions = finalActions
	return rule
}
