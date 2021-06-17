package models

import (
	"github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/model"
	"time"
)

// Session !
type UserAuthToken struct {
	model.Model
	User       models.User
	UserID     uint
	Token      string
	SessionDuration  *time.Duration
}

