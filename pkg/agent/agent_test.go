package agent

import (
	"encoding/json"
	"testing"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestMergeBundleClean(t *testing.T) {
	old := models.Bundle{
		EnvironmentVariables: map[string]string{
			"AAAA": "AAAA",
		},
		DesiredAgentVersion: "1",
	}

	new := models.Bundle{
		EnvironmentVariables: map[string]string{
			"ASDF": "WASDF",
		},
		DesiredAgentVersion: "1.2",
	}

	newB, err := json.Marshal(new)
	assert.NoError(t, err)

	merged := mergeBundle(&old, newB)
	assert.Equal(t, new, *merged)
}

func TestMergeBundleIncompatible(t *testing.T) {
	old := models.Bundle{
		EnvironmentVariables: map[string]string{
			"AAAA": "AAAA",
		},
		DesiredAgentVersion: "1",
	}

	new := map[string]interface{}{
		"environmentVariables": map[string]string{
			"ASDF": "WASDF",
		},
		"applications": map[string]string{
			"ASDF":  "WASDF",
			"WASDF": "ASDF",
		},
		"desiredAgentVersion": "2",
	}

	newB, err := json.Marshal(new)
	assert.NoError(t, err)

	merged := mergeBundle(&old, newB)
	assert.NotEqual(t, new, *merged)
	assert.Equal(t, new["desiredAgentVersion"], merged.DesiredAgentVersion)

	// Change this so we can easily compare everything else
	merged.DesiredAgentVersion = old.DesiredAgentVersion

	// Marshal old and new so we can compare bytes and not types
	oldB, err := json.Marshal(old)
	assert.NoError(t, err)
	mergedB, err := json.Marshal(merged)
	assert.NoError(t, err)

	assert.Equal(t, string(oldB), string(mergedB))
}

func TestMergeBundleIncompatibleWithEmptyOld(t *testing.T) {
	var old *models.Bundle
	new := map[string]interface{}{
		"environmentVariables": map[string]string{
			"ASDF": "WASDF",
		},
		"applications": map[string]string{
			"ASDF":  "WASDF",
			"WASDF": "ASDF",
		},
		"desiredAgentVersion": "2",
	}

	newB, err := json.Marshal(new)
	assert.NoError(t, err)

	merged := mergeBundle(old, newB)
	assert.NotEqual(t, new, *merged)
	assert.Equal(t, new["desiredAgentVersion"], merged.DesiredAgentVersion)
}
