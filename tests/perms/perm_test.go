package perms

import (
	"github.com/stretchr/testify/assert"
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin/core"
	"testing"
)

type PermTestSuite struct {
	uadmin.TestSuite
}

func (suite *PermTestSuite) TestIntegration() {
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	db := uadminDatabase.Db
	contentType := core.ContentType{BlueprintName: "user", ModelName: "user"}
	db.Create(&contentType)
	permission := core.Permission{ContentType: contentType, PermissionBits: core.RevertPermBit}
	db.Create(&permission)
	g1 := core.UserGroup{GroupName: "usergroup"}
	db.Create(&g1)
	db.Model(&g1).Association("Permissions").Append(&permission)
	db.Save(&g1)
	permission = core.Permission{ContentType: contentType, PermissionBits: core.AddPermBit}
	db.Create(&permission)
	u := core.User{Username: "dsadas", Email: "ffsdfsd@example.com"}
	db.Create(&u)
	db.Model(&u).Association("Permissions").Append(&permission)
	db.Model(&u).Association("UserGroups").Append(&g1)
	db.Save(&u)
	var u1 core.User
	db.Model(&core.User{}).First(&u1)
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
