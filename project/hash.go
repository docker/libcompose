package project

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"reflect"
	"sort"
	"strings"
)

func GetServiceHash(name string, config ServiceConfig, ignore map[string]bool) string {
	hash := sha1.New()

	io.WriteString(hash, fmt.Sprintln(name))

	//Get values of Service through reflection
	val := reflect.ValueOf(config).Elem()

	//Create slice to sort the keys in Service Config, which allow constant hash ordering
	serviceKeys := []string{}

	//Create a data structure of map of values keyed by a string
	unsortedKeyValue := make(map[string]interface{})

	//Get all keys and values in Service Configuration
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		keyField := val.Type().Field(i)

		serviceKeys = append(serviceKeys, keyField.Name)
		unsortedKeyValue[keyField.Name] = valueField.Interface()
	}

	//Sort serviceKeys alphabetically
	sort.Strings(serviceKeys)

	//Go through keys and write hash
	for _, serviceKey := range serviceKeys {
		if ignore[strings.ToLower(serviceKey)] {
			continue
		}

		serviceValue := unsortedKeyValue[serviceKey]

		io.WriteString(hash, fmt.Sprintf("\n  %v: ", serviceKey))

		switch s := serviceValue.(type) {
		case SliceorMap:
			writeMap(hash, s.MapParts())
		case MaporEqualSlice:
			writeSlice(hash, s.Slice())
		case MaporColonSlice:
			writeSlice(hash, s.Slice())
		case MaporSpaceSlice:
			writeSlice(hash, s.Slice())
		case Command:
			writeSlice(hash, s.Slice())
		case Stringorslice:
			writeSlice(hash, s.Slice())
		case []string:
			writeSlice(hash, s)
		default:
			writeString(hash, fmt.Sprintf("%v", serviceValue))
		}
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func writeSlice(h hash.Hash, data []string) {
	for _, part := range data {
		writeString(h, fmt.Sprintf("%s", part))
		h.Write([]byte{0})
	}
}

func writeMap(h hash.Hash, data map[string]string) {
	keys := []string{}
	for key, _ := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		writeString(h, fmt.Sprintf("%s=%v", key, data[key]))
		h.Write([]byte{0})
	}
}

func writeString(h hash.Hash, val string) {
	io.WriteString(h, val)
}
