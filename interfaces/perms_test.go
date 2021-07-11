package interfaces

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPerm(t *testing.T) {
	perm := NewPerm(ReadPermBit | AddPermBit | EditPermBit | DeletePermBit | PublishPermBit | RevertPermBit)
	assert.True(t, perm.HasAddPermission())
	assert.True(t, perm.HasReadPermission())
	assert.True(t, perm.HasEditPermission())
	assert.True(t, perm.HasDeletePermission())
	assert.True(t, perm.HasPublishPermission())
	assert.True(t, perm.HasRevertPermission())
	perm = NewPerm(ReadPermBit, "test")
	assert.False(t, perm.HasAddPermission())
	assert.True(t, perm.DoesUserHaveRightFor(CustomPermission("test")))
}

