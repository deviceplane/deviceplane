package spec

import (
	"testing"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestValidate(t *testing.T) {
	t.Run("full", func(t *testing.T) {
		full, _ := yaml.Marshal(map[string]models.Service{
			"s": fullService(),
		})
		require.NoError(t, Validate(full))
	})
}
