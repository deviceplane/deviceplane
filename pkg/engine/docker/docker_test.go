package docker

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/require"
)

func TestGetRegistryAuthConfig(t *testing.T) {
	_, err := getProcessedRegistryAuth("")
	require.NotNil(t, err)

	processedRegistryAuth, err := getProcessedRegistryAuth("invalidbase64")
	require.NotNil(t, err)

	processedRegistryAuth, err = getProcessedRegistryAuth("dXNlcm5hbWU=")
	require.NotNil(t, err)

	processedRegistryAuth, err = getProcessedRegistryAuth("dXNlcm5hbWU6cGFzc3dvcmQ=")
	require.Nil(t, err)

	decodedProcessedRegistryAuth, err := base64.StdEncoding.DecodeString(processedRegistryAuth)
	require.Nil(t, err)

	var authConfig types.AuthConfig
	err = json.Unmarshal(decodedProcessedRegistryAuth, &authConfig)
	require.Nil(t, err)

	require.Equal(t, "username", authConfig.Username)
	require.Equal(t, "password", authConfig.Password)
}
