package settings

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/blueprint/settings/migrations"
	settingmodel "github.com/uadmin/uadmin/blueprint/settings/models"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	settingsAdminPage := admin.NewGormAdminPage(func() interface{} {return nil}, "")
	settingsAdminPage.PageName = "Settings"
	settingsAdminPage.Slug = "setting"
	settingsAdminPage.BlueprintName = "setting"
	settingsAdminPage.Router = mainRouter
	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(settingsAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
	settingmodelAdminPage := admin.NewGormAdminPage(func() interface{} {return &settingmodel.Setting{}}, "setting")
	settingmodelAdminPage.PageName = "Settings"
	settingmodelAdminPage.Slug = "setting"
	settingmodelAdminPage.BlueprintName = "setting"
	settingmodelAdminPage.Router = mainRouter
	err = settingsAdminPage.SubPages.AddAdminPage(settingmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
	settingcategoriesmodelAdminPage := admin.NewGormAdminPage(func() interface{} {return &settingmodel.SettingCategory{}}, "settingcategory")
	settingcategoriesmodelAdminPage.PageName = "Setting categories"
	settingcategoriesmodelAdminPage.Slug = "settingcategory"
	settingcategoriesmodelAdminPage.BlueprintName = "setting"
	settingcategoriesmodelAdminPage.Router = mainRouter
	err = settingsAdminPage.SubPages.AddAdminPage(settingcategoriesmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	interfaces.ProjectModels.RegisterModel(&settingmodel.SettingCategory{})
	interfaces.ProjectModels.RegisterModel(&settingmodel.Setting{})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "settings",
		Description:       "Settings blueprint responsible for wide-project settings",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
