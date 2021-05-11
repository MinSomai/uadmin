package uadmin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetSchema is a unit testing function for getSchema() function
func TestSchema(t *testing.T) {
	schema, _ := getSchema(User{})
	assert.Equal(t, "users", schema.TableName)
}
