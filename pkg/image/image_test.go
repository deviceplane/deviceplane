package image

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToCanonical(t *testing.T) {
	require.Equal(t, "docker.io/library/ubuntu", ToCanonical("ubuntu"))
	require.Equal(t, "docker.io/deviceplane/deviceplane", ToCanonical("deviceplane/deviceplane"))
	require.Equal(t, "docker.io/deviceplane/deviceplane", ToCanonical("docker.io/deviceplane/deviceplane"))
}
