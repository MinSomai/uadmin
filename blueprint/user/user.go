package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/user/migrations"
	"github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/dialect"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

type PasswordValidationStruct struct {
	Password string `valid:"password-uadmin"`
	ConfirmedPassword string
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
}

func (b Blueprint) Init(config *config.UadminConfig) {
	govalidator.CustomTypeTagMap.Set("username-unique", func(i interface{}, o interface{}) bool {
		db := dialect.GetDB()
		var cUsers int64
		db.Model(&models.User{}).Where(&models.User{Username: i.(string)}).Count(&cUsers)
		return cUsers == 0
	})
	govalidator.CustomTypeTagMap.Set("email-unique", func(i interface{}, o interface{}) bool {
		db := dialect.GetDB()
		var cUsers int64
		db.Model(&models.User{}).Where(&models.User{Email: i.(string)}).Count(&cUsers)
		return cUsers == 0
	})
	govalidator.CustomTypeTagMap.Set("username-uadmin", func(i interface{}, o interface{}) bool {
		minLength := config.D.Auth.MinUsernameLength
		maxLength := config.D.Auth.MaxUsernameLength
		currentUsername := i.(string)
		if maxLength < len(currentUsername) || len(currentUsername) < minLength {
			return false
		}
		return true
	})
	govalidator.CustomTypeTagMap.Set("password-uadmin", func(i interface{}, o interface{}) bool {
		passwordStruct := o.(PasswordValidationStruct)
		if passwordStruct.Password != passwordStruct.ConfirmedPassword {
			return false
		}
		if len(passwordStruct.Password) < config.D.Auth.MinPasswordLength {
			return false
		}
		return true
	})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "user",
		Description:       "this blueprint is about users",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
