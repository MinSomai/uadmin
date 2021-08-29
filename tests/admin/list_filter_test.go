package admin

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/core"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

type AdminListFilterTestSuite struct {
	uadmin.UadminTestSuite
}

func (apts *AdminListFilterTestSuite) SetupTestData() {
	uadminDatabase := core.NewUadminDatabase()
	for i := range core.GenerateNumberSequence(101, 200) {
		userModel := &core.User{
			Email: fmt.Sprintf("admin_%d@example.com", i),
			Username: "admin_" + strconv.Itoa(i),
			FirstName: "firstname_" + strconv.Itoa(i),
			LastName: "lastname_" + strconv.Itoa(i),
		}
		uadminDatabase.Db.Create(&userModel)
	}
	uadminDatabase.Close()
}

func (suite *AdminListFilterTestSuite) TestFiltering() {
	suite.SetupTestData()
	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	var users []core.User
	adminRequestParams := core.NewAdminRequestParams()
	adminRequestParams.RequestURL = "http://127.0.0.1/?Username__exact=admin_101"
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	statement := &gorm.Statement{DB: uadminDatabase.Db}
	statement.Parse(&core.User{})
	listFilter := &core.ListFilter{
		UrlFilteringParam: "Username__exact",
	}
	adminUserPage.ListFilter.Add(listFilter)
	adminUserPage.GetQueryset(adminUserPage, adminRequestParams).GetPaginatedQuerySet().Find(&users)
	assert.Equal(suite.T(), len(users), 1)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestListFilter(t *testing.T) {
	uadmin.Run(t, new(AdminListFilterTestSuite))
}
