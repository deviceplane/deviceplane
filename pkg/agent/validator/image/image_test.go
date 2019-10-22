package image

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidation(t *testing.T) {
	require.True(t,
		isValid("redis", []string{}),
		"Should pass on empty files",
	)

	require.True(t,
		isValid("redis", []string{"redis"}),
		"Should pass on matching images",
	)

	require.True(t,
		isValid("deviceplane/agent:latest",
			[]string{
				"redis",
				"deviceplane/",
			},
		),
		"Should pass on matching org prefix",
	)

	require.False(t,
		isValid("deviceplaneexploit/agent:latest",
			[]string{
				"redis",
				"deviceplane/",
			},
		),
		"Should fail on non-matching org prefix",
	)

	require.False(t,
		isValid("postgres",
			[]string{
				"redis",
				"deviceplane/",
			},
		),
		"Should fail on non-matching image",
	)
}
