package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSchemaRegistry(t *testing.T) {
	schemaData := map[string]string{
		"":    schemaDataV1,
		"2":   schemaDataV2,
		"2.0": schemaDataV2,
		"2.1": schemaDataV2_1,
	}

	schemaRegistry, err := NewSchemaRegistry(schemaData)
	if err != nil {
		t.Error(err)
	}

	for _, schema := range schemaRegistry {
		assert.NotNil(t, schema.Data)
		assert.NotNil(t, schema.Loader)
		assert.NotNil(t, schema.ConstraintsLoader)
	}
}
