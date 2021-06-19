package models

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUserModelFromJson(t *testing.T) {
	jsonstr, _ := json.Marshal(User{Username: "test"})
	user, _ := NewUserModelFromJson(jsonstr)
	assert.Equal(t, user.Username, "test")
}
