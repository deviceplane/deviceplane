package updater

import (
	"testing"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/yamltypes"
	"github.com/stretchr/testify/require"
)

func TestWithCommandInterpolation(t *testing.T) {
	for _, scenario := range []struct {
		expected  []string
		command   []string
		projectID string
	}{
		{
			expected: []string{"a"},
			command:  []string{"a"},
		},
		{
			expected:  []string{"a", "bprj"},
			command:   []string{"a", "b$PROJECT"},
			projectID: "prj",
		},
		{
			expected:  []string{"prj"},
			command:   []string{"$PROJECT"},
			projectID: "prj",
		},
	} {
		s := models.Service{
			Entrypoint: []string{"a"},
			Command:    scenario.command,
		}
		s = withCommandInterpolation(s, scenario.projectID)
		require.Equal(t, yamltypes.Command(scenario.expected), s.Command)
		require.Equal(t, yamltypes.Command([]string{"a"}), s.Entrypoint)
	}

}
