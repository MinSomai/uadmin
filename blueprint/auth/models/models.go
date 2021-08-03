package models

import (
	"github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/interfaces"
	"time"
)

// Session !
type UserAuthToken struct {
	interfaces.Model
	User       models.User
	UserID     uint
	Token      string
	SessionDuration  time.Duration
}

