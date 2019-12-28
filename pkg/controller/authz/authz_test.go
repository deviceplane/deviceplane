package authz

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvaluate(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		for _, scenario := range []struct {
			resource Resource
			action   Action
			configs  []Config
		}{
			{
				ResourceApplications,
				ActionGetApplication,
				[]Config{ReadAllRole},
			},
			{
				ResourceApplications,
				ActionGetApplication,
				[]Config{WriteAllRole},
			},
			{
				ResourceApplications,
				ActionCreateApplication,
				[]Config{WriteAllRole},
			},
			{
				ResourceApplications,
				ActionGetApplication,
				[]Config{AdminAllRole},
			},
			{
				ResourceApplications,
				ActionCreateApplication,
				[]Config{AdminAllRole},
			},
			{
				ResourceMemberships,
				ActionCreateMembership,
				[]Config{AdminAllRole},
			},
		} {
			require.True(t, Evaluate(scenario.resource, scenario.action, scenario.configs))
		}

		for _, scenario := range []struct {
			resource Resource
			action   Action
			configs  []Config
		}{
			{
				ResourceApplications,
				ActionCreateApplication,
				[]Config{ReadAllRole},
			},
			{
				ResourceMemberships,
				ActionCreateMembership,
				[]Config{ReadAllRole},
			},
			{
				ResourceMemberships,
				ActionCreateMembership,
				[]Config{WriteAllRole},
			},
		} {
			require.False(t, Evaluate(scenario.resource, scenario.action, scenario.configs))
		}
	})

	t.Run("custom", func(t *testing.T) {
		for _, scenario := range []struct {
			resource Resource
			action   Action
			configs  []Config
		}{
			{
				ResourceApplications,
				ActionGetApplication,
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []Resource{ResourceApplications},
								Actions:   []Action{ActionGetApplication},
							},
						},
					},
				},
			},
			{
				ResourceApplications,
				ActionGetApplication,
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []Resource{ResourceAny},
								Actions:   []Action{ActionGetApplication},
							},
						},
					},
				},
			},
			{
				ResourceApplications,
				ActionCreateApplication,
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []Resource{ResourceAny},
								Actions:   []Action{ActionWriteAll},
							},
							{
								Resources: []Resource{ResourceAny},
								Actions:   []Action{ActionReadAll},
								Effect:    EffectDeny,
							},
						},
					},
				},
			},
		} {
			require.True(t, Evaluate(scenario.resource, scenario.action, scenario.configs))
		}

		for _, scenario := range []struct {
			resource Resource
			action   Action
			configs  []Config
		}{
			{
				ResourceApplications,
				ActionCreateApplication,
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []Resource{ResourceApplications},
								Actions:   []Action{ActionGetApplication},
							},
						},
					},
				},
			},
			{
				ResourceApplications,
				ActionGetApplication,
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []Resource{ResourceReleases},
								Actions:   []Action{ActionGetApplication},
							},
						},
					},
				},
			},
			{
				ResourceApplications,
				ActionGetApplication,
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []Resource{ResourceAny},
								Actions:   []Action{ActionGetRelease},
							},
						},
					},
				},
			},
			{
				ResourceApplications,
				ActionGetApplication,
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []Resource{ResourceReleases},
								Actions:   []Action{ActionReadAll},
							},
						},
					},
				},
			},
			{
				ResourceApplications,
				ActionGetApplication,
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []Resource{ResourceApplications},
								Actions:   []Action{ActionGetApplication},
								Effect:    EffectDeny,
							},
						},
					},
				},
			},
			{
				ResourceApplications,
				ActionGetApplication,
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []Resource{ResourceAny},
								Actions:   []Action{ActionReadAll},
							},
							{
								Resources: []Resource{ResourceApplications},
								Actions:   []Action{ActionGetApplication},
								Effect:    EffectDeny,
							},
						},
					},
				},
			},
			{
				ResourceApplications,
				ActionGetApplication,
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []Resource{ResourceAny},
								Actions:   []Action{ActionWriteAll},
							},
							{
								Resources: []Resource{ResourceAny},
								Actions:   []Action{ActionReadAll},
								Effect:    EffectDeny,
							},
						},
					},
				},
			},
		} {
			require.False(t, Evaluate(scenario.resource, scenario.action, scenario.configs))
		}
	})
}
