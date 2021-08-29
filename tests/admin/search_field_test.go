package admin

import (
	"fmt"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/interfaces"
	"strconv"
	"testing"
)

type AdminSearchFieldTestSuite struct {
	uadmin.UadminTestSuite
}

func (apts *AdminSearchFieldTestSuite) SetupTestData() {
	uadminDatabase := interfaces.NewUadminDatabase()
	for i := range interfaces.GenerateNumberSequence(201, 300) {
		userModel := &interfaces.User{
			Email: fmt.Sprintf("admin_%d@example.com", i),
			Username: "admin_" + strconv.Itoa(i),
			FirstName: "firstname_" + strconv.Itoa(i),
			LastName: "lastname_" + strconv.Itoa(i),
		}
		uadminDatabase.Db.Create(&userModel)
	}
	uadminDatabase.Close()
}

func (suite *AdminSearchFieldTestSuite) TestFiltering() {
	// @todo, uncomment when fix search by different fields in the admin
	//suite.SetupTestData()
	//adminUserBlueprintPage, _ := interfaces.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	//adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	//var users []interfaces.User
	//adminRequestParams := interfaces.NewAdminRequestParams()
	//uadminDatabase := interfaces.NewUadminDatabase()
	//defer uadminDatabase.Close()
	//statement := &gorm.Statement{DB: uadminDatabase.Db}
	//statement.Parse(&interfaces.User{})
	//adminRequestParams.Search = "admin_202@example.com"
	//adminUserPage.GetQueryset(adminUserPage, adminRequestParams).GormQuerySet.Find(&users)
	//assert.Equal(suite.T(), len(users), 1)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSearchField(t *testing.T) {
	uadmin.Run(t, new(AdminSearchFieldTestSuite))
}
