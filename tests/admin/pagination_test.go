package admin

import (
	"fmt"
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin/core"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

type AdminPaginationTestSuite struct {
	uadmin.TestSuite
}

func (suite *AdminPaginationTestSuite) SetupTestData() {
	for i := range core.GenerateNumberSequence(1, 100) {
		userModel := &core.User{
			Email:     fmt.Sprintf("admin_%d@example.com", i),
			Username:  "admin_" + strconv.Itoa(i),
			FirstName: "firstname_" + strconv.Itoa(i),
			LastName:  "lastname_" + strconv.Itoa(i),
		}
		suite.UadminDatabase.Db.Create(&userModel)
	}
}

func (suite *AdminPaginationTestSuite) TestPagination() {
	suite.SetupTestData()
	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	var users []core.User
	adminRequestParams := core.NewAdminRequestParams()
	adminUserPage.GetQueryset(adminUserPage, adminRequestParams).GetPaginatedQuerySet().Find(&users)
	assert.Equal(suite.T(), len(users), core.CurrentConfig.D.Uadmin.AdminPerPage)
	adminRequestParams.Paginator.Offset = 88
	adminUserPage.GetQueryset(adminUserPage, adminRequestParams).GetPaginatedQuerySet().Find(&users)
	assert.Greater(suite.T(), len(users), core.CurrentConfig.D.Uadmin.AdminPerPage)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAdminPagination(t *testing.T) {
	uadmin.RunTests(t, new(AdminPaginationTestSuite))
}
