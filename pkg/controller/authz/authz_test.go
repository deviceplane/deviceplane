package authz

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvaluate(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		for _, scenario := range []struct {
			resource string
			action   string
			configs  []Config
		}{
			{
				"applications",
				"GetApplication",
				[]Config{ReadAllRole},
			},
			{
				"applications",
				"GetApplication",
				[]Config{WriteAllRole},
			},
			{
				"applications",
				"CreateApplication",
				[]Config{WriteAllRole},
			},
			{
				"applications",
				"GetApplication",
				[]Config{AdminAllRole},
			},
			{
				"applications",
				"CreateApplication",
				[]Config{AdminAllRole},
			},
			{
				"memberships",
				"CreateMembership",
				[]Config{AdminAllRole},
			},
		} {
			require.True(t, Evaluate(scenario.resource, scenario.action, scenario.configs))
		}

		for _, scenario := range []struct {
			resource string
			action   string
			configs  []Config
		}{
			{
				"applications",
				"CreateApplication",
				[]Config{ReadAllRole},
			},
			{
				"memberships",
				"CreateMembership",
				[]Config{ReadAllRole},
			},
			{
				"memberships",
				"CreateMembership",
				[]Config{WriteAllRole},
			},
		} {
			require.False(t, Evaluate(scenario.resource, scenario.action, scenario.configs))
		}
	})

	t.Run("custom", func(t *testing.T) {
		for _, scenario := range []struct {
			resource string
			action   string
			configs  []Config
		}{
			{
				"applications",
				"GetApplication",
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []string{"applications"},
								Actions:   []string{"GetApplication"},
							},
						},
					},
				},
			},
		} {
			require.True(t, Evaluate(scenario.resource, scenario.action, scenario.configs))
		}

		for _, scenario := range []struct {
			resource string
			action   string
			configs  []Config
		}{
			{
				"applications",
				"CreateApplication",
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []string{"applications"},
								Actions:   []string{"GetApplication"},
							},
						},
					},
				},
			},
			{
				"applications",
				"GetApplication",
				[]Config{
					{
						Rules: []Rule{
							{
								Resources: []string{"releases"},
								Actions:   []string{"GetApplication"},
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
