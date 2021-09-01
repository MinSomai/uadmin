package settings

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/settings/migrations"
	settingmodel "github.com/uadmin/uadmin/blueprint/settings/models"
	"github.com/uadmin/uadmin/core"
)

type Blueprint struct {
	core.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	settingsAdminPage := core.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) { return nil, nil },
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	settingsAdminPage.PageName = "Settings"
	settingsAdminPage.Slug = "setting"
	settingsAdminPage.BlueprintName = "setting"
	settingsAdminPage.Router = mainRouter
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(settingsAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
	settingmodelAdminPage := core.NewGormAdminPage(
		settingsAdminPage,
		func() (interface{}, interface{}) { return &settingmodel.Setting{}, &[]*settingmodel.Setting{} },
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			return nil
		},
	)
	settingmodelAdminPage.PageName = "Settings"
	settingmodelAdminPage.Slug = "setting"
	settingmodelAdminPage.BlueprintName = "setting"
	settingmodelAdminPage.Router = mainRouter
	settingmodelAdminPage.NoPermissionToEdit = true
	settingmodelAdminPage.NoPermissionToAddNew = true
	dataTypeListDisplay, _ := settingmodelAdminPage.ListDisplay.GetFieldByDisplayName("DataType")
	dataTypeListDisplay.Populate = func(m interface{}) string {
		return settingmodel.HumanizeDataType(m.(*settingmodel.Setting).DataType)
	}
	categoryListDisplay, _ := settingmodelAdminPage.ListDisplay.GetFieldByDisplayName("Category")
	categoryListDisplay.Field.FieldConfig.Widget.SetReadonly(true)
	//categoryListDisplay.Populate = func(m interface{}) string {
	//	return m.(*settingmodel.Setting).Category.Name
	//}
	settingmodelAdminPage.NoPermissionToAddNew = true
	err = settingsAdminPage.SubPages.AddAdminPage(settingmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
	settingcategoriesmodelAdminPage := core.NewGormAdminPage(
		settingsAdminPage,
		func() (interface{}, interface{}) {
			return &settingmodel.SettingCategory{}, &[]*settingmodel.SettingCategory{}
		},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"Name", "Icon"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			return form
		},
	)
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
	core.ProjectModels.RegisterModel(func() interface{} { return &settingmodel.SettingCategory{} })
	core.ProjectModels.RegisterModel(func() interface{} { return &settingmodel.Setting{} })
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "settings",
		Description:       "Settings blueprint responsible for wide-project settings",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
