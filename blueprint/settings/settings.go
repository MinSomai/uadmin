package settings

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/blueprint/settings/migrations"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	settingsAdminPage := admin.NewAdminPage()
	settingsAdminPage.PageName = "Settings"
	settingsAdminPage.Slug = "setting"
	settingsAdminPage.BlueprintName = "setting"
	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(settingsAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
	settingmodelAdminPage := admin.NewAdminPage()
	settingmodelAdminPage.PageName = "Settings"
	settingmodelAdminPage.Slug = "setting"
	settingmodelAdminPage.BlueprintName = "setting"
	err = settingsAdminPage.SubPages.AddAdminPage(settingmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
	settingcategoriesmodelAdminPage := admin.NewAdminPage()
	settingcategoriesmodelAdminPage.PageName = "Setting categories"
	settingcategoriesmodelAdminPage.Slug = "settingcategory"
	settingcategoriesmodelAdminPage.BlueprintName = "setting"
	err = settingsAdminPage.SubPages.AddAdminPage(settingcategoriesmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "settings",
		Description:       "Settings blueprint responsible for wide-project settings",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
