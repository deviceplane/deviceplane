package interpolation

import (
	"bytes"
	"fmt"
)

func Interpolate(config map[string]interface{}, getenv func(string) string) error {
	for k, v := range config {
		if err := parseConfig(k, &v, func(s string) (string, error) {
			value := getenv(s)
			if value == "" {
				return "", fmt.Errorf("variable %s is not defined", s)
			}
			return value, nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func parseConfig(service string, data *interface{}, mapping func(string) (string, error)) error {
	switch typedData := (*data).(type) {
	case string:
		var success bool
		var err error
		if *data, success, err = parseLine(typedData, mapping); err != nil {
			return err
		} else if !success {
			return fmt.Errorf("invalid interpolation format in line '%s'", typedData)
		}

	case []interface{}:
		for k, v := range typedData {
			if err := parseConfig(service, &v, mapping); err != nil {
				return err
			}
			typedData[k] = v
		}

	case map[interface{}]interface{}:
		for k, v := range typedData {
			if err := parseConfig(service, &v, mapping); err != nil {
				return err
			}
			typedData[k] = v
		}
	}

	return nil
}

func parseLine(line string, mapping func(string) (string, error)) (string, bool, error) {
	var buffer bytes.Buffer

	for pos := 0; pos < len(line); pos++ {
		c := line[pos]

		switch {
		case c == '$':
			var replaced string
			var success bool
			var err error
			replaced, pos, success, err = parseInterpolationExpression(line, pos+1, mapping)
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

func parseInterpolationExpression(line string, pos int, mapping func(string) (string, error)) (string, int, bool, error) {
	c := line[pos]

	switch {
	case c == '$':
		return "$", pos, true, nil
	case c == '{':
		return parseVariableWithBraces(line, pos+1, mapping)
	case !isNum(c) && validVariableNameChar(c):
		return parseVariable(line, pos, mapping)
	default:
		return "", 0, false, nil
	}
}

func parseVariable(line string, pos int, mapping func(string) (string, error)) (string, int, bool, error) {
	var buffer bytes.Buffer

	for ; pos < len(line); pos++ {
		c := line[pos]

		switch {
		case validVariableNameChar(c):
			buffer.WriteByte(c)
		default:
			value, err := mapping(buffer.String())
			return value, pos - 1, true, err
		}
	}

	value, err := mapping(buffer.String())
	return value, pos, true, err
}

func parseVariableWithBraces(line string, pos int, mapping func(string) (string, error)) (string, int, bool, error) {
	var buffer bytes.Buffer

	for ; pos < len(line); pos++ {
		c := line[pos]

		switch {
		case c == '}':
			bufferString := buffer.String()
			if bufferString == "" {
				return "", 0, false, nil
			}
			value, err := mapping(buffer.String())
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
