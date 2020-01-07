package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInternalTitledTitledRegex(t *testing.T) {
	for _, valid := range []string{
		"0_9",
		"a_aaaadsf13242ASDF9123",
		"ref_123",
	} {
		require.True(t, internalTitleRegex.Match([]byte(valid)))
	}

	for _, invalid := range []string{
		"",
		"a",
		"a b",
		"a(b)",
		"a_b_c",
	} {
		require.False(t, internalTitleRegex.Match([]byte(invalid)))
	}

}

func TestUserTitledRegex(t *testing.T) {
	for _, valid := range []string{
		"a",
		"abc",
		"a-b-c",
	} {
		require.True(t, userTitleRegex.Match([]byte(valid)))
	}

	for _, invalid := range []string{
		"",
		"a b",
		"a(b)",
		"a_b_c",
	} {
		require.False(t, userTitleRegex.Match([]byte(invalid)))
	}

}
