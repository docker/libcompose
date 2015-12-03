package project

import (
	"fmt"
	"strings"

	"github.com/flynn/go-shlex"
)

// stringer converts ints, strings and bools to a string
func stringer(v interface{}) (string, error) {
	switch v.(type) {
	case string, int64, int32, int, bool:
		return fmt.Sprint(v), nil
	default:
		return "", fmt.Errorf("Value of type %T can't be converted to a string", v)
	}
}

func sliceStringer(value []interface{}) ([]string, error) {
	slice := make([]string, len(value))
	for k, v := range value {
		if vstr, err := stringer(v); err != nil {
			return nil, err
		} else {
			slice[k] = vstr
		}
	}
	return slice, nil
}

func mapStringer(value map[interface{}]interface{}) (map[string]string, error) {
	parts := map[string]string{}
	for k, v := range value {
		kstr, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("Map must have string keys only, had %T", k)
		}
		if vstr, err := stringer(v); err != nil {
			return nil, err
		} else {
			parts[kstr] = vstr
		}
	}
	return parts, nil
}

func mapToSlice(m map[string]string, joinStr string) []string {
	slice := []string{}
	for k, v := range m {
		slice = append(slice, strings.Join([]string{k, v}, joinStr))
	}
	return slice
}

// Stringorslice represents a string or an array of strings.
// TODO use docker/docker/pkg/stringutils.StrSlice once 1.9.x is released.
type Stringorslice struct {
	parts []string
}

// MarshalYAML implements the Marshaller interface.
func (s Stringorslice) MarshalYAML() (tag string, value interface{}, err error) {
	return "", s.parts, nil
}

// UnmarshalYAML implements the Unmarshaller interface.
func (s *Stringorslice) UnmarshalYAML(tag string, value interface{}) error {
	var err error
	switch value := value.(type) {
	case []interface{}:
		s.parts, err = sliceStringer(value)
	case string:
		s.parts = []string{value}
	default:
		return fmt.Errorf("Failed to unmarshal Stringorslice: %#v", value)
	}
	return err
}

// Len returns the number of parts of the Stringorslice.
func (s *Stringorslice) Len() int {
	if s == nil {
		return 0
	}
	return len(s.parts)
}

// Slice gets the parts of the StrSlice as a Slice of string.
func (s *Stringorslice) Slice() []string {
	if s == nil {
		return nil
	}
	return s.parts
}

// NewStringorslice creates an Stringorslice based on the specified parts (as strings).
func NewStringorslice(parts ...string) Stringorslice {
	return Stringorslice{parts}
}

// Command represents a docker command, can be a string or an array of strings.
// FIXME why not use Stringorslice (type Command struct { Stringorslice }
type Command struct {
	parts []string
}

// MarshalYAML implements the Marshaller interface.
func (s Command) MarshalYAML() (tag string, value interface{}, err error) {
	return "", s.parts, nil
}

// UnmarshalYAML implements the Unmarshaller interface.
func (s *Command) UnmarshalYAML(tag string, value interface{}) error {
	var err error
	switch value := value.(type) {
	case []interface{}:
		s.parts, err = sliceStringer(value)
	case string:
		s.parts, err = shlex.Split(value)
	default:
		return fmt.Errorf("Failed to unmarshal Command: %#v", value)
	}
	return err
}

// ToString returns the parts of the command as a string (joined by spaces).
func (s *Command) ToString() string {
	return strings.Join(s.parts, " ")
}

// Slice gets the parts of the Command as a Slice of string.
func (s *Command) Slice() []string {
	return s.parts
}

// NewCommand create a Command based on the specified parts (as strings).
func NewCommand(parts ...string) Command {
	return Command{parts}
}

// SliceorMap represents a slice or a map of strings.
type SliceorMap struct {
	parts map[string]string
}

// MarshalYAML implements the Marshaller interface.
func (s SliceorMap) MarshalYAML() (tag string, value interface{}, err error) {
	return "", s.parts, nil
}

// UnmarshalYAML implements the Unmarshaller interface.
func (s *SliceorMap) UnmarshalYAML(tag string, value interface{}) error {
	var err error
	switch value := value.(type) {
	case map[interface{}]interface{}:
		s.parts, err = mapStringer(value)
	case []interface{}:
		parts := map[string]string{}
		values, err := sliceStringer(value)
		if err != nil {
			return err
		}
		for _, v := range values {
			keyValueSlice := strings.SplitN(strings.TrimSpace(v), "=", 2)
			key := keyValueSlice[0]
			val := ""
			if len(keyValueSlice) == 2 {
				val = keyValueSlice[1]
			}
			parts[key] = val
		}
		s.parts = parts
	default:
		return fmt.Errorf("Failed to unmarshal SliceorMap: %#v", value)
	}
	return err
}

// MapParts get the parts of the SliceorMap as a Map of string.
func (s *SliceorMap) MapParts() map[string]string {
	if s == nil {
		return nil
	}
	return s.parts
}

// NewSliceorMap creates a new SliceorMap based on the specified parts (as map of string).
func NewSliceorMap(parts map[string]string) SliceorMap {
	return SliceorMap{parts}
}

// MaporEqualSlice represents a slice of strings that gets unmarshal from a
// YAML map into 'key=value' string.
type MaporEqualSlice struct {
	parts []string
}

// MarshalYAML implements the Marshaller interface.
func (s MaporEqualSlice) MarshalYAML() (tag string, value interface{}, err error) {
	return "", s.parts, nil
}

// UnmarshalYAML implements the Unmarshaller interface.
func (s *MaporEqualSlice) UnmarshalYAML(tag string, value interface{}) error {
	var err error
	switch value := value.(type) {
	case []interface{}:
		s.parts, err = sliceStringer(value)
	case map[interface{}]interface{}:
		parts, err := mapStringer(value)
		if err != nil {
			return err
		}
		s.parts = mapToSlice(parts, "=")
	default:
		return fmt.Errorf("Failed to unmarshal MaporEqualSlice: %#v", value)
	}
	return err
}

// Slice gets the parts of the MaporEqualSlice as a Slice of string.
func (s *MaporEqualSlice) Slice() []string {
	return s.parts
}

// NewMaporEqualSlice creates a new MaporEqualSlice based on the specified parts.
func NewMaporEqualSlice(parts []string) MaporEqualSlice {
	return MaporEqualSlice{parts}
}

// MaporColonSlice represents a slice of strings that gets unmarshal from a
// YAML map into 'key:value' string.
type MaporColonSlice struct {
	parts []string
}

// MarshalYAML implements the Marshaller interface.
func (s MaporColonSlice) MarshalYAML() (tag string, value interface{}, err error) {
	return "", s.parts, nil
}

// UnmarshalYAML implements the Unmarshaller interface.
func (s *MaporColonSlice) UnmarshalYAML(tag string, value interface{}) error {
	var err error
	switch value := value.(type) {
	case []interface{}:
		s.parts, err = sliceStringer(value)
	case map[interface{}]interface{}:
		parts, err := mapStringer(value)
		if err != nil {
			return err
		}
		s.parts = mapToSlice(parts, ":")
	default:
		return fmt.Errorf("Failed to unmarshal MaporColonSlice: %#v", value)
	}
	return err
}

// Slice gets the parts of the MaporColonSlice as a Slice of string.
func (s *MaporColonSlice) Slice() []string {
	return s.parts
}

// NewMaporColonSlice creates a new MaporColonSlice based on the specified parts.
func NewMaporColonSlice(parts []string) MaporColonSlice {
	return MaporColonSlice{parts}
}

// MaporSpaceSlice represents a slice of strings that gets unmarshal from a
// YAML map into 'key value' string.
type MaporSpaceSlice struct {
	parts []string
}

// MarshalYAML implements the Marshaller interface.
func (s MaporSpaceSlice) MarshalYAML() (tag string, value interface{}, err error) {
	return "", s.parts, nil
}

// UnmarshalYAML implements the Unmarshaller interface.
func (s *MaporSpaceSlice) UnmarshalYAML(tag string, value interface{}) error {
	var err error
	switch value := value.(type) {
	case []interface{}:
		s.parts, err = sliceStringer(value)
	case map[interface{}]interface{}:
		parts, err := mapStringer(value)
		if err != nil {
			return err
		}
		s.parts = mapToSlice(parts, " ")
	default:
		return fmt.Errorf("Failed to unmarshal MaporSpaceSlice: %#v", value)
	}
	return err
}

// Slice gets the parts of the MaporSpaceSlice as a Slice of string.
func (s *MaporSpaceSlice) Slice() []string {
	return s.parts
}

// NewMaporSpaceSlice creates a new MaporSpaceSlice based on the specified parts.
func NewMaporSpaceSlice(parts []string) MaporSpaceSlice {
	return MaporSpaceSlice{parts}
}
