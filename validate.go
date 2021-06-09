package uadmin

import (
	"github.com/asaskevich/govalidator"
	"github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/dialect"
)

func init() {
	// after
	govalidator.CustomTypeTagMap.Set("username-unique", func(i interface{}, o interface{}) bool {
		db := dialect.GetDB("dialect")
		var cUsers int64
		db.Where(&models.User{Username: i.(string)}).Count(&cUsers)
		return cUsers > 0
	})
	govalidator.CustomTypeTagMap.Set("email-unique", func(i interface{}, o interface{}) bool {
		db := dialect.GetDB("dialect")
		var cUsers int64
		db.Where(&models.User{Email: i.(string)}).Count(&cUsers)
		return cUsers > 0
	})
	govalidator.CustomTypeTagMap.Set("username-uadmin", func(i interface{}, o interface{}) bool {
		minLength := appInstance.Config.D.Auth.MinUsernameLength
		maxLength := appInstance.Config.D.Auth.MaxUsernameLength
		currentUsername := i.(string)
		if maxLength < len(currentUsername) || len(currentUsername) < minLength {
			return false
		}
		return true
	})
}
