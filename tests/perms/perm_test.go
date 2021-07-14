package perms

import (
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/interfaces"
	"testing"
)

type PermTestSuite struct {
	uadmin.UadminTestSuite
}

func (suite *PermTestSuite) TestIntegration() {
	db := interfaces.GetDB()
	contentType := interfaces.ContentType{BlueprintName: "user", ModelName: "user"}
	db.Create(&contentType)
	permission := usermodels.Permission{ContentType: contentType, PermissionBits: interfaces.RevertPermBit}
	db.Create(&permission)
	g1 := usermodels.UserGroup{GroupName: "usergroup"}
	db.Create(&g1)
	db.Model(&g1).Association("Permissions").Append(&permission)
	db.Save(&g1)
	permission = usermodels.Permission{ContentType: contentType, PermissionBits: interfaces.AddPermBit}
	db.Create(&permission)
	u := usermodels.User{Username: "dsadas", Email: "ffsdfsd@example.com"}
	db.Create(&u)
	db.Model(&u).Association("Permissions").Append(&permission)
	db.Model(&u).Association("UserGroups").Append(&g1)
	db.Save(&u)
	var u1 usermodels.User
	db.Model(&usermodels.User{}).First(&u1)
	permRegistry := u1.BuildPermissionRegistry()
	userPerm := permRegistry.GetPermissionForBlueprint("user", "user")
	assert.True(suite.T(), userPerm.HasRevertPermission())
	assert.True(suite.T(), userPerm.HasAddPermission())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPermissionSystem(t *testing.T) {
	uadmin.Run(t, new(PermTestSuite))
}
