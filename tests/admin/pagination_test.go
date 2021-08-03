package admin

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/admin"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/interfaces"
	"strconv"
	"testing"
)

type AdminPaginationTestSuite struct {
	uadmin.UadminTestSuite
}

func (apts *AdminPaginationTestSuite) SetupTestData() {
	uadminDatabase := interfaces.NewUadminDatabase()
	for i := range interfaces.GenerateNumberSequence(1, 100) {
		userModel := &usermodels.User{
			Email: fmt.Sprintf("admin_%d@example.com", i),
			Username: "admin_" + strconv.Itoa(i),
			FirstName: "firstname_" + strconv.Itoa(i),
			LastName: "lastname_" + strconv.Itoa(i),
		}
		uadminDatabase.Db.Create(&userModel)
	}
	uadminDatabase.Close()
}

func (suite *AdminPaginationTestSuite) TestPagination() {
	suite.SetupTestData()
	adminUserBlueprintPage, _ := admin.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	var users []usermodels.User
	adminRequestParams := interfaces.NewAdminRequestParams()
	adminUserPage.GetQueryset(adminUserPage, adminRequestParams).PaginatedGormQuerySet.Find(&users)
	assert.Equal(suite.T(), len(users), interfaces.CurrentConfig.D.Uadmin.AdminPerPage)
	adminRequestParams.Paginator.Offset = 88
	adminUserPage.GetQueryset(adminUserPage, adminRequestParams).PaginatedGormQuerySet.Find(&users)
	assert.Greater(suite.T(), len(users), interfaces.CurrentConfig.D.Uadmin.AdminPerPage)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAdminPagination(t *testing.T) {
	uadmin.Run(t, new(AdminPaginationTestSuite))
}
