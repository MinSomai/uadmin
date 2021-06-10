package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigValues(t *testing.T) {
	uadminConfig := NewConfig("configs/test.yaml")
	assert.Equal(t, uadminConfig.D.Uadmin.Theme, "default")
}
