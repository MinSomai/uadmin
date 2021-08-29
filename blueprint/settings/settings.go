package settings

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/settings/migrations"
	settingmodel "github.com/uadmin/uadmin/blueprint/settings/models"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	settingsAdminPage := interfaces.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) {return nil, nil},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {return nil},
	)
	settingsAdminPage.PageName = "Settings"
	settingsAdminPage.Slug = "setting"
	settingsAdminPage.BlueprintName = "setting"
	settingsAdminPage.Router = mainRouter
	err := interfaces.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(settingsAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
	settingmodelAdminPage := interfaces.NewGormAdminPage(
		settingsAdminPage,
		func() (interface{}, interface{}) {return &settingmodel.Setting{}, &[]*settingmodel.Setting{}},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {return nil},
	)
	settingmodelAdminPage.PageName = "Settings"
	settingmodelAdminPage.Slug = "setting"
	settingmodelAdminPage.BlueprintName = "setting"
	settingmodelAdminPage.Router = mainRouter
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
	settingcategoriesmodelAdminPage := interfaces.NewGormAdminPage(
		settingsAdminPage,
		func() (interface{}, interface{}) {return &settingmodel.SettingCategory{}, &[]*settingmodel.SettingCategory{}},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {
			fields := []string{"Name", "Icon"}
			form := interfaces.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
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
	interfaces.ProjectModels.RegisterModel(func() interface{} {return &settingmodel.SettingCategory{}})
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &settingmodel.Setting{}})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "settings",
		Description:       "Settings blueprint responsible for wide-project settings",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
