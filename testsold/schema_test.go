package testsold

import (
	"github.com/uadmin/uadmin/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetSchema is a unit testing function for getSchema() function
func TestSchema(t *testing.T) {
	schema, _ := model.GetSchema(User{})
	assert.Equal(t, "users", schema.TableName)
}
