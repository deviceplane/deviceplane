package interpolation

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	errInvalidInterpolation = errors.New("invalid interpolation")
)

func newVariableNotDefinedError(variable string) error {
	return &variableNotDefinedError{
		variable: variable,
	}
}

type variableNotDefinedError struct {
	variable string
}

func (e *variableNotDefinedError) Error() string {
	return fmt.Sprintf("variable %s is not defined", e.variable)
}

func Interpolate(s string, getVariable func(string) string) (string, error) {
	s, success, err := interpolate(s, func(variable string) (string, error) {
		value := getVariable(variable)
		if value == "" {
			return "", newVariableNotDefinedError(variable)
		}
		return value, nil
	})
	if err != nil {
		return "", err
	}
	if !success {
		return "", errInvalidInterpolation
	}
	return s, nil
}

func interpolate(s string, getVariable func(string) (string, error)) (string, bool, error) {
	var buffer bytes.Buffer

	for pos := 0; pos < len(s); pos++ {
		c := s[pos]

		switch {
		case c == '$':
			var replaced string
			var success bool
			var err error
			replaced, pos, success, err = parseInterpolationExpression(s, pos+1, getVariable)
			if err != nil {
				return "", false, err
			}
			if !success {
				return "", false, err
			}
			buffer.WriteString(replaced)

		default:
			buffer.WriteByte(c)
		}
	}

	return buffer.String(), true, nil
}

func parseInterpolationExpression(s string, pos int, getVariable func(string) (string, error)) (string, int, bool, error) {
	c := s[pos]

	switch {
	case c == '$':
		return "$", pos, true, nil
	case c == '{':
		return parseVariableWithBraces(s, pos+1, getVariable)
	case !isNum(c) && validVariableNameChar(c):
		return parseVariable(s, pos, getVariable)
	default:
		return "", 0, false, nil
	}
}

func parseVariable(s string, pos int, getVariable func(string) (string, error)) (string, int, bool, error) {
	var buffer bytes.Buffer

	for ; pos < len(s); pos++ {
		c := s[pos]

		switch {
		case validVariableNameChar(c):
			buffer.WriteByte(c)
		default:
			value, err := getVariable(buffer.String())
			return value, pos - 1, true, err
		}
	}

	value, err := getVariable(buffer.String())
	return value, pos, true, err
}

func parseVariableWithBraces(s string, pos int, getVariable func(string) (string, error)) (string, int, bool, error) {
	var buffer bytes.Buffer

	for ; pos < len(s); pos++ {
		c := s[pos]

		switch {
		case c == '}':
			bufferString := buffer.String()
			if bufferString == "" {
				return "", 0, false, nil
			}
			value, err := getVariable(buffer.String())
			return value, pos, true, err

		case validVariableNameChar(c):
			buffer.WriteByte(c)

		default:
			return "", 0, false, nil
		}
	}

	return "", 0, false, nil
}

func validVariableNameChar(c uint8) bool {
	return c == '_' ||
		c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		isNum(c)
}

func isNum(c uint8) bool {
	return c >= '0' && c <= '9'
}
