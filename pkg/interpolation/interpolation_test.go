package interpolation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testInterpolate(t *testing.T, expected, in string, variables map[string]string) {
	out, err := Interpolate(in, func(variable string) string {
		return variables[variable]
	})
	require.Nil(t, err)
	require.Equal(t, expected, out)
}

func testInterpolateMissingVariable(t *testing.T, in, variable string, variables map[string]string) {
	_, err := Interpolate(in, func(variable string) string {
		return variables[variable]
	})
	errVariableNotDefined, ok := err.(*variableNotDefinedError)
	require.True(t, ok)
	require.Equal(t, variable, errVariableNotDefined.variable)
}

func testInvalidInterpolate(t *testing.T, in string) {
	_, err := Interpolate(in, func(string) string {
		return ""
	})
	require.Equal(t, errInvalidInterpolation, err)
}

func TestParseLine(t *testing.T) {
	variables := map[string]string{
		"A":           "ABC",
		"X":           "XYZ",
		"E":           "",
		"lower":       "WORKED",
		"MiXeD":       "WORKED",
		"split_VaLue": "WORKED",
	}

	testInterpolate(t, "WORKED", "$lower", variables)
	testInterpolate(t, "WORKED", "${MiXeD}", variables)
	testInterpolate(t, "WORKED", "${split_VaLue}", variables)

	testInterpolate(t, "ABC", "$A", variables)
	testInterpolate(t, "ABC", "${A}", variables)

	testInterpolate(t, "ABC DE", "$A DE", variables)
	testInterpolate(t, "ABCDE", "${A}DE", variables)

	testInterpolate(t, "$A", "$$A", variables)
	testInterpolate(t, "${A}", "$${A}", variables)

	testInterpolate(t, "$ABC", "$$${A}", variables)
	testInterpolate(t, "$ABC", "$$$A", variables)

	testInterpolate(t, "ABC XYZ", "$A $X", variables)
	testInterpolate(t, "ABCXYZ", "$A$X", variables)
	testInterpolate(t, "ABCXYZ", "${A}${X}", variables)

	testInterpolateMissingVariable(t, "$B", "B", variables)
	testInterpolateMissingVariable(t, "${B}", "B", variables)
	testInterpolateMissingVariable(t, "$ADE", "ADE", variables)

	testInterpolateMissingVariable(t, "$E", "E", variables)
	testInterpolateMissingVariable(t, "${E}", "E", variables)

	testInvalidInterpolate(t, "${")
	testInvalidInterpolate(t, "$}")
	testInvalidInterpolate(t, "${}")
	testInvalidInterpolate(t, "${ }")
	testInvalidInterpolate(t, "${A }")
	testInvalidInterpolate(t, "${ A}")
	testInvalidInterpolate(t, "${A!}")
	testInvalidInterpolate(t, "$!")
}
