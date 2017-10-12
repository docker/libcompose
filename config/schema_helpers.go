package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-connections/nat"
	"github.com/xeipuuv/gojsonschema"
)

// Schema represents the various gojsonschema utilites we'll use with our parsed
// JSON schema data.
type Schema struct {
	Loader            gojsonschema.JSONLoader
	ConstraintsLoader gojsonschema.JSONLoader
	Data              map[string]interface{}
}

// SchemaRegistry represents our map of Compose versions to Schema objects.
type SchemaRegistry map[string]*Schema

// GetVersion
func (reg SchemaRegistry) GetVersion(version string) (*Schema, error) {
	if schema, ok := reg[version]; ok {
		return schema, nil
	}
	return nil, fmt.Errorf("couldn't load JSON schema '%s'")
}

var schemaRegistry SchemaRegistry

// Initialize our app's schemaRegistry which maps Compose version numbers to
// various JSON schemas.
func init() {
	var err error
	schemaRegistry, err = NewSchemaRegistry(map[string]string{
		"":    schemaDataV1,
		"2":   schemaDataV2,
		"2.0": schemaDataV2,
		"2.1": schemaDataV2_1,
	})
	if err != nil {
		logrus.Fatalf("can't init JSON schema registry:", err)
	}
}

// NewSchemaRegistry creates a new SchemaRegistry which converts a map of
// versions to JSON schemas into versions and proper Schema objects.
func NewSchemaRegistry(schemaData map[string]string) (map[string]*Schema, error) {
	schemaRegistry := make(map[string]*Schema)
	for version, data := range schemaData {
		schema, err := NewSchema(data)
		if err != nil {
			return nil, fmt.Errorf("can't load JSON schema:", err)
		}
		schemaRegistry[version] = schema
	}
	return schemaRegistry, nil
}

type (
	environmentFormatChecker struct{}
	portsFormatChecker       struct{}
)

func (checker environmentFormatChecker) IsFormat(input string) bool {
	// If the value is a boolean, a warning should be given
	//
	// However, we can't determine type since gojsonschema converts the value to
	// a string
	//
	// Adding a function with an interface{} parameter to gojsonschema is
	// probably the best way to handle this
	return true
}

func (checker portsFormatChecker) IsFormat(input string) bool {
	_, _, err := nat.ParsePortSpecs([]string{input})
	return err == nil
}

// NewSchema constructs a new Schema including parsing and initializing JSON
// schema loaders used to validate sections of a Compose file.
func NewSchema(schemaJSON string) (*Schema, error) {
	var schemaParsed interface{}
	err := json.Unmarshal([]byte(schemaJSON), &schemaParsed)
	if err != nil {
		return nil, err
	}

	data := schemaParsed.(map[string]interface{})

	gojsonschema.FormatCheckers.Add("environment", environmentFormatChecker{})
	gojsonschema.FormatCheckers.Add("ports", portsFormatChecker{})
	gojsonschema.FormatCheckers.Add("expose", portsFormatChecker{})
	loader := gojsonschema.NewGoLoader(schemaParsed)

	definitions := data["definitions"].(map[string]interface{})
	constraints := definitions["constraints"].(map[string]interface{})
	service := constraints["service"].(map[string]interface{})
	constraintsLoader := gojsonschema.NewGoLoader(service)

	return &Schema{
		Loader:            loader,
		ConstraintsLoader: constraintsLoader,
		Data:              data,
	}, nil
}

// gojsonschema doesn't provide a list of valid types for a property
// This parses the schema manually to find all valid types
func parseValidTypesFromSchema(schema map[string]interface{}, context string) []string {
	contextSplit := strings.Split(context, ".")
	key := contextSplit[len(contextSplit)-1]

	definitions := schema["definitions"].(map[string]interface{})
	service := definitions["service"].(map[string]interface{})
	properties := service["properties"].(map[string]interface{})
	property := properties[key].(map[string]interface{})

	var validTypes []string

	if val, ok := property["oneOf"]; ok {
		validConditions := val.([]interface{})

		for _, validCondition := range validConditions {
			condition := validCondition.(map[string]interface{})
			validTypes = append(validTypes, condition["type"].(string))
		}
	} else if val, ok := property["$ref"]; ok {
		reference := val.(string)
		if reference == "#/definitions/string_or_list" {
			return []string{"string", "array"}
		} else if reference == "#/definitions/list_of_strings" {
			return []string{"array"}
		} else if reference == "#/definitions/list_or_dict" {
			return []string{"array", "object"}
		}
	}

	return validTypes
}
