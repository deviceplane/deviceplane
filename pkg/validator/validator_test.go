package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStandardRegex(t *testing.T) {
	for _, valid := range []string{
		"a",
		"abc",
		"a-b-c",
		"a_b_c",
	} {
		require.True(t, standardRegex.Match([]byte(valid)))
	}

	for _, invalid := range []string{
		"",
		"a b",
		"a(b)",
	} {
		require.False(t, standardRegex.Match([]byte(invalid)))
	}

}
