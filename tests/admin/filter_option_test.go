package admin

import (
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin/core"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type FilterOptionTestSuite struct {
	uadmin.TestSuite
}

func (suite *FilterOptionTestSuite) BeforeTest(suiteName string, testMethod string) {
	if testMethod == "TestFilterOptionByYear" {
	} else {
	}
}

func (suite *FilterOptionTestSuite) Test() {
	userModel := &core.User{Username: "adminfilteroption", Email: "adminfilteroption@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Now().Add((-10 * 12 * 86400 * 30) * time.Second)
	suite.UadminDatabase.Db.Create(userModel)
	userModel = &core.User{Username: "adminfilteroption1", Email: "adminfilteroption1@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Now().Add((-5 * 12 * 86400 * 30) * time.Second)
	suite.UadminDatabase.Db.Create(userModel)
	userModel = &core.User{Username: "adminfilteroption2", Email: "adminfilteroption2@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Now().Add((-3 * 12 * 86400 * 30) * time.Second)
	suite.UadminDatabase.Db.Create(userModel)
	userModel = &core.User{Username: "adminfilteroption3", Email: "adminfilteroption3@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Now().Add((-1 * 12 * 86400 * 30) * time.Second)
	suite.UadminDatabase.Db.Create(userModel)
	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	newFilterOption := core.NewFilterOption()
	newFilterOption.FetchOptions = func(afo core.IAdminFilterObjects) []*core.DisplayFilterOption {
		return core.FetchOptionsFromGormModelFromDateTimeField(afo, "created_at")
	}
	adminUserPage.FilterOptions.AddFilterOption(newFilterOption)
	assert.True(suite.T(), len(adminUserPage.FetchFilterOptions()) > 0)
	suite.UadminDatabase.Db.Unscoped().Where("1 = 1").Delete(&core.User{})
	userModel = &core.User{Username: "adminfilteroptionA", Email: "adminfilteroption@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	// ((-10 * 12 * 86400 * 30) * time.Second)
	suite.UadminDatabase.Db.Create(userModel)
	userModel = &core.User{Username: "adminfilteroption1A", Email: "adminfilteroption1@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Date(2020, time.February, 1, 0, 0, 0, 0, time.UTC)
	suite.UadminDatabase.Db.Create(userModel)
	userModel = &core.User{Username: "adminfilteroption2A", Email: "adminfilteroption2@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Date(2020, time.March, 1, 0, 0, 0, 0, time.UTC)
	suite.UadminDatabase.Db.Create(userModel)
	userModel = &core.User{Username: "adminfilteroption3A", Email: "adminfilteroption3@example.com", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	userModel.CreatedAt = time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC)
	suite.UadminDatabase.Db.Create(userModel)
	adminUserBlueprintPage, _ = core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ = adminUserBlueprintPage.SubPages.GetBySlug("user")
	newFilterOption = core.NewFilterOption()
	newFilterOption.FetchOptions = func(afo core.IAdminFilterObjects) []*core.DisplayFilterOption {
		return core.FetchOptionsFromGormModelFromDateTimeField(afo, "created_at")
	}
	adminUserPage.FilterOptions = core.NewFilterOptionsRegistry()
	adminUserPage.FilterOptions.AddFilterOption(newFilterOption)
	assert.True(suite.T(), len(adminUserPage.FetchFilterOptions()) > 0)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFilterOption(t *testing.T) {
	uadmin.RunTests(t, new(FilterOptionTestSuite))
}
