package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
)

var defaultValues = make(map[string]string)

func isNum(c uint8) bool {
	return c >= '0' && c <= '9'
}

func validVariableDefault(c uint8, line string, pos int) bool {
	return (c == ':' && line[pos+1] == '-') || (c == '-')
}

func validVariableNameChar(c uint8) bool {
	return c == '_' ||
		c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		isNum(c)
}

func parseVariable(line string, pos int, mapping func(string) string) (string, int, bool) {
	var buffer bytes.Buffer

	for ; pos < len(line); pos++ {
		c := line[pos]

		switch {
		case validVariableNameChar(c):
			buffer.WriteByte(c)
		default:
			return mapping(buffer.String()), pos - 1, true
		}
	}

	return mapping(buffer.String()), pos, true
}

func parseDefaultValue(line string, pos int, mapping func(string) string) (string, int, bool) {
	var buffer bytes.Buffer
	success := false
	replaced := ""

	fmt.Printf("INCOMING: %v\n", string(line[pos]))
	// only skip :, :- and - at the beginning
	for ; pos < len(line); pos++ {
		c := line[pos]
		fmt.Printf("SKIP: %v\n", string(c))
		if c == ':' || c == '-' {
			continue
		}
		break
	}
	for ; pos < len(line); pos++ {
		c := line[pos]
		switch {
		case c == '$':
			if line[pos-1] != '\\' {
				fmt.Printf("NESTED: %v pos: %v\n", string(c), pos)
				replaced, pos, success = parseInterpolationExpression(line, pos + 1, mapping)
				_, err := buffer.WriteString(replaced)
				fmt.Printf("REPLACED: %v pos: %v\n", replaced, pos)
				if success == false || err != nil {
					return "", 0, false
				}
			}
		case c == '}':
			fmt.Printf("ENDIT: %v\n", string(c))
			return buffer.String(), pos - 1, true
		default:
			fmt.Printf("APPEND: %v\n", string(c))
			// err := buffer.WriteByte(c)
			if err := buffer.WriteByte(c); err != nil {
				return "", pos, false
			}
			fmt.Printf("APPEND: %v\n", string(c))
		}
		fmt.Printf("BUFFER: %v\n", string(buffer.String()))
	}
	return "", 0, false
}

func parseVariableWithBraces(line string, pos int, mapping func(string) string) (string, int, bool) {
	var buffer bytes.Buffer

	for ; pos < len(line); pos++ {
		c := line[pos]
		fmt.Printf("val: %v\n", string(c))

		switch {
		case c == '}':
			bufferString := buffer.String()
			if bufferString == "" {
				return "", 0, false
			}
			return mapping(buffer.String()), pos, true
		case validVariableNameChar(c):
			buffer.WriteByte(c)
		case validVariableDefault(c, line, pos):
			defaultValue := ""
			defaultValue, pos, _ = parseDefaultValue(line, pos, mapping)
			defaultValues[buffer.String()] = defaultValue
		default:
			return "", 0, false
		}
	}

	return "", 0, false
}

func parseInterpolationExpression(line string, pos int, mapping func(string) string) (string, int, bool) {
	c := line[pos]

	switch {
	case c == '$':
		return "$", pos, true
	case c == '{':
		return parseVariableWithBraces(line, pos+1, mapping)
	case !isNum(c) && validVariableNameChar(c):
		// Variables can't start with a number
		return parseVariable(line, pos, mapping)
	default:
		return "", 0, false
	}
}

func parseLine(line string, mapping func(string) string) (string, bool) {
	var buffer bytes.Buffer

	for pos := 0; pos < len(line); pos++ {
		c := line[pos]
		switch {
		case c == '$':
			var replaced string
			var success bool

			replaced, pos, success = parseInterpolationExpression(line, pos+1, mapping)

			if !success {
				return "", false
			}

			buffer.WriteString(replaced)
		default:
			buffer.WriteByte(c)
		}
	}

	return buffer.String(), true
}

func parseConfig(key string, data *interface{}, mapping func(string) string) error {
	switch typedData := (*data).(type) {
	case string:
		var success bool

		*data, success = parseLine(typedData, mapping)

		if !success {
			return fmt.Errorf("Invalid interpolation format for key \"%s\": \"%s\"", key, typedData)
		}
	case []interface{}:
		for k, v := range typedData {
			err := parseConfig(key, &v, mapping)

			if err != nil {
				return err
			}

			typedData[k] = v
		}
	case map[interface{}]interface{}:
		for k, v := range typedData {
			err := parseConfig(key, &v, mapping)

			if err != nil {
				return err
			}

			typedData[k] = v
		}
	}

	return nil
}

// Interpolate replaces variables in a map entry
func Interpolate(key string, data *interface{}, environmentLookup EnvironmentLookup) error {
	return parseConfig(key, data, func(s string) string {
		values := environmentLookup.Lookup(s, nil)

		if len(values) == 0 {
			if val, ok := defaultValues[s]; ok {
				return val
			}
			logrus.Warnf("The %s variable is not set. Substituting a blank string.", s)
			return ""
		}

		if strings.SplitN(values[0], "=", 2)[1] == "" {
			if val, ok := defaultValues[s]; ok {
				return val
			}
		}

		// Use first result if many are given
		value := values[0]

		// Environment variables come in key=value format
		// Return everything past first '='
		return strings.SplitN(value, "=", 2)[1]
	})
}
