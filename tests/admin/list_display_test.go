package admin

import (
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/admin"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"testing"
)

type AdminListDisplayTestSuite struct {
	uadmin.UadminTestSuite
}

func (suite *AdminListDisplayTestSuite) TestListDisplay() {
	userModel := &usermodels.User{Username: "admin", FirstName: "firstname", LastName: "lastname"}
	adminUserBlueprintPage, _ := admin.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	listDisplayUsername, _ := adminUserPage.ListDisplay.GetFieldByDisplayName("Username")
	assert.Equal(suite.T(), listDisplayUsername.GetValue(userModel), "admin")
	compositeField := admin.NewListDisplay(nil)
	compositeField.MethodName = "FullName"
	compositeField = admin.NewListDisplay(nil)
	compositeField.Populate = func(m interface{}) string {
		return m.(*usermodels.User).FullName()
	}
	assert.Equal(suite.T(), compositeField.GetValue(userModel), "firstname lastname")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAdminListDisplay(t *testing.T) {
	uadmin.Run(t, new(AdminListDisplayTestSuite))
}
