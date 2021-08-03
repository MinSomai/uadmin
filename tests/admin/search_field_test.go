package admin

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/admin"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/interfaces"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

type AdminSearchFieldTestSuite struct {
	uadmin.UadminTestSuite
}

func (apts *AdminSearchFieldTestSuite) SetupTestData() {
	uadminDatabase := interfaces.NewUadminDatabase()
	for i := range interfaces.GenerateNumberSequence(201, 300) {
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

func (suite *AdminSearchFieldTestSuite) TestFiltering() {
	suite.SetupTestData()
	adminUserBlueprintPage, _ := admin.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	var users []usermodels.User
	adminRequestParams := interfaces.NewAdminRequestParams()
	uadminDatabase := interfaces.NewUadminDatabase()
	statement := &gorm.Statement{DB: uadminDatabase.Db}
	statement.Parse(&usermodels.User{})
	adminRequestParams.Search = "admin_202@example.com"
	adminUserPage.GetQueryset(adminUserPage, adminRequestParams).GormQuerySet.Find(&users)
	assert.Equal(suite.T(), len(users), 1)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSearchField(t *testing.T) {
	uadmin.Run(t, new(AdminSearchFieldTestSuite))
}
