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

type AdminSearchFieldTestSuite struct {
	uadmin.TestSuite
}

func (suite *AdminSearchFieldTestSuite) SetupTestData() {
	uadminDatabase := core.NewUadminDatabase()
	for i := range core.GenerateNumberSequence(201, 300) {
		userModel := &core.User{
			Email:     fmt.Sprintf("admin_%d@example.com", i),
			Username:  "admin_" + strconv.Itoa(i),
			FirstName: "firstname_" + strconv.Itoa(i),
			LastName:  "lastname_" + strconv.Itoa(i),
		}
		uadminDatabase.Db.Create(&userModel)
	}
	uadminDatabase.Close()
}

func (suite *AdminSearchFieldTestSuite) TestFiltering() {
	suite.SetupTestData()
	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	var users []core.User
	adminRequestParams := core.NewAdminRequestParams()
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	statement := &gorm.Statement{DB: uadminDatabase.Db}
	statement.Parse(&core.User{})
	adminRequestParams.Search = "admin_202@example.com"
	adminUserPage.GetQueryset(adminUserPage, adminRequestParams).GetFullQuerySet().Find(&users)
	assert.Equal(suite.T(), len(users), 1)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSearchField(t *testing.T) {
	uadmin.Run(t, new(AdminSearchFieldTestSuite))
}
