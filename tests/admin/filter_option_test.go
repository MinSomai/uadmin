package admin

import (
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/core"
	"testing"
	"time"
)

type FilterOptionTestSuite struct {
	uadmin.UadminTestSuite
}

func (suite *FilterOptionTestSuite) TestFilterOptionByYear() {
	userModel := &core.User{Username: "admin", Email: "admin@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Now().Add((-10 * 12 * 86400 * 30) * time.Second)
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	uadminDatabase.Db.Create(userModel)
	userModel = &core.User{Username: "admin1", Email: "admin1@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Now().Add((-5 * 12 * 86400 * 30) * time.Second)
	uadminDatabase.Db.Create(userModel)
	userModel = &core.User{Username: "admin2", Email: "admin2@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Now().Add((-3 * 12 * 86400 * 30) * time.Second)
	uadminDatabase.Db.Create(userModel)
	userModel = &core.User{Username: "admin3", Email: "admin3@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Now().Add((-1 * 12 * 86400 * 30) * time.Second)
	uadminDatabase.Db.Create(userModel)
	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	newFilterOption := core.NewFilterOption()
	newFilterOption.FetchOptions = func (afo core.IAdminFilterObjects) []*core.DisplayFilterOption {
		return core.FetchOptionsFromGormModelFromDateTimeField(afo, "created_at")
	}
	adminUserPage.FilterOptions.AddFilterOption(newFilterOption)
	assert.True(suite.T(), len(adminUserPage.FetchFilterOptions()) > 0)
}

func (suite *FilterOptionTestSuite) TestFilterOptionByMonth() {
	userModel := &core.User{Username: "admin", Email: "admin@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	// ((-10 * 12 * 86400 * 30) * time.Second)
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	uadminDatabase.Db.Create(userModel)
	userModel = &core.User{Username: "admin1", Email: "admin1@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Date(2020, time.February, 1, 0, 0, 0, 0, time.UTC)
	uadminDatabase.Db.Create(userModel)
	userModel = &core.User{Username: "admin2", Email: "admin2@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Date(2020, time.March, 1, 0, 0, 0, 0, time.UTC)
	uadminDatabase.Db.Create(userModel)
	userModel = &core.User{Username: "admin3", Email: "admin3@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC)
	uadminDatabase.Db.Create(userModel)
	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	newFilterOption := core.NewFilterOption()
	newFilterOption.FetchOptions = func (afo core.IAdminFilterObjects) []*core.DisplayFilterOption {
		return core.FetchOptionsFromGormModelFromDateTimeField(afo, "created_at")
	}
	adminUserPage.FilterOptions = core.NewFilterOptionsRegistry()
	adminUserPage.FilterOptions.AddFilterOption(newFilterOption)
	assert.True(suite.T(), len(adminUserPage.FetchFilterOptions()) > 0)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFilterOption(t *testing.T) {
	uadmin.Run(t, new(FilterOptionTestSuite))
}
