package admin

import (
	"fmt"
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin/core"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

type AdminListFilterTestSuite struct {
	uadmin.TestSuite
}

func (suite *AdminListFilterTestSuite) SetupTestData() {
	for i := range core.GenerateNumberSequence(101, 200) {
		userModel := &core.User{
			Email:     fmt.Sprintf("admin_%d@example.com", i),
			Username:  "admin_" + strconv.Itoa(i),
			FirstName: "firstname_" + strconv.Itoa(i),
			LastName:  "lastname_" + strconv.Itoa(i),
		}
		suite.UadminDatabase.Db.Create(&userModel)
	}
}

func (suite *AdminListFilterTestSuite) TestFiltering() {
	suite.SetupTestData()
	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	var users []core.User
	adminRequestParams := core.NewAdminRequestParams()
	adminRequestParams.RequestURL = "http://127.0.0.1/?Username__exact=admin_101"
	statement := &gorm.Statement{DB: suite.UadminDatabase.Db}
	statement.Parse(&core.User{})
	listFilter := &core.ListFilter{
		URLFilteringParam: "Username__exact",
	}
	adminUserPage.ListFilter.Add(listFilter)
	adminUserPage.GetQueryset(adminUserPage, adminRequestParams).GetPaginatedQuerySet().Find(&users)
	assert.Equal(suite.T(), len(users), 1)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestListFilter(t *testing.T) {
	uadmin.RunTests(t, new(AdminListFilterTestSuite))
}
