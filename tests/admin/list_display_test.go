package admin

import (
	"github.com/stretchr/testify/assert"
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin/core"
	"testing"
)

type AdminListDisplayTestSuite struct {
	uadmin.TestSuite
}

func (suite *AdminListDisplayTestSuite) TestListDisplay() {
	userModel := &core.User{Username: "admin", FirstName: "firstname", LastName: "lastname"}
	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	listDisplayUsername, _ := adminUserPage.ListDisplay.GetFieldByDisplayName("Username")
	assert.Equal(suite.T(), listDisplayUsername.GetValue(userModel), "admin")
	compositeField := core.NewListDisplay(nil)
	compositeField.MethodName = "FullName"
	compositeField = core.NewListDisplay(nil)
	compositeField.Populate = func(m interface{}) string {
		return m.(*core.User).FullName()
	}
	assert.Equal(suite.T(), compositeField.GetValue(userModel), "firstname lastname")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAdminListDisplay(t *testing.T) {
	uadmin.Run(t, new(AdminListDisplayTestSuite))
}
