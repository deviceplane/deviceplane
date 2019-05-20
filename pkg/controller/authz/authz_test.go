package authz

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: test this
// TODO: yeah we totally don't handle this properly
//		ag, _ := Evaluate("applications", "CreateMembership", []Config{AdminAllRole})
//		require.True(t, ag)

func TestEvaluate(t *testing.T) {
	t.Run("read", func(t *testing.T) {
		ag, _ := Evaluate("applications", "GetApplication", []Config{ReadAllRole})
		require.True(t, ag)
		ag, _ = Evaluate("applications", "CreateApplication", []Config{ReadAllRole})
		require.False(t, ag)
		ag, _ = Evaluate("memberships", "CreateMembership", []Config{ReadAllRole})
		require.False(t, ag)
	})
	t.Run("write", func(t *testing.T) {
		ag, _ := Evaluate("applications", "GetApplication", []Config{WriteAllRole})
		require.True(t, ag)
		ag, _ = Evaluate("applications", "CreateApplication", []Config{WriteAllRole})
		require.True(t, ag)
		ag, _ = Evaluate("memberships", "CreateMembership", []Config{WriteAllRole})
		require.False(t, ag)
	})
	t.Run("admin", func(t *testing.T) {
		ag, _ := Evaluate("applications", "GetApplication", []Config{AdminAllRole})
		require.True(t, ag)
		ag, _ = Evaluate("applications", "CreateApplication", []Config{AdminAllRole})
		require.True(t, ag)
		ag, _ = Evaluate("memberships", "CreateMembership", []Config{AdminAllRole})
		require.True(t, ag)
	})
}
