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
	settingsAdminPage := admin.NewGormAdminPage(nil, func() (interface{}, interface{}) {return nil, nil}, "")
	settingsAdminPage.PageName = "Settings"
	settingsAdminPage.Slug = "setting"
	settingsAdminPage.BlueprintName = "setting"
	settingsAdminPage.Router = mainRouter
	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(settingsAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
	settingmodelAdminPage := admin.NewGormAdminPage(settingsAdminPage, func() (interface{}, interface{}) {return &settingmodel.Setting{}, &[]*settingmodel.Setting{}}, "setting")
	settingmodelAdminPage.PageName = "Settings"
	settingmodelAdminPage.Slug = "setting"
	settingmodelAdminPage.BlueprintName = "setting"
	settingmodelAdminPage.Router = mainRouter
	adminContext := &interfaces.AdminContext{}
	settingForm := interfaces.NewFormFromModelFromGinContext(adminContext, &settingmodel.Setting{}, make([]string, 0), []string{}, true, "")
	settingmodelAdminPage.Form = settingForm
	settingNameField, _ := settingForm.FieldRegistry.GetByName("Name")
	settingNameListDisplay := interfaces.NewListDisplay(settingNameField)
	settingmodelAdminPage.ListDisplay.AddField(settingNameListDisplay)
	settingNameListDisplay.Ordering = 1
	valueField, _ := settingForm.FieldRegistry.GetByName("Value")
	valueListDisplay := interfaces.NewListDisplay(valueField)
	settingmodelAdminPage.ListDisplay.AddField(valueListDisplay)
	valueListDisplay.Ordering = 2
	defaultValueField, _ := settingForm.FieldRegistry.GetByName("DefaultValue")
	defaultValueListDisplay := interfaces.NewListDisplay(defaultValueField)
	defaultValueListDisplay.Ordering = 3
	settingmodelAdminPage.ListDisplay.AddField(defaultValueListDisplay)
	dataTypeField, _ := settingForm.FieldRegistry.GetByName("DataType")
	dataTypeListDisplay := interfaces.NewListDisplay(dataTypeField)
	dataTypeListDisplay.Ordering = 4
	dataTypeListDisplay.Populate = func(m interface{}) string {
		return settingmodel.HumanizeDataType(m.(*settingmodel.Setting).DataType)
	}
	settingmodelAdminPage.ListDisplay.AddField(dataTypeListDisplay)
	helpField, _ := settingForm.FieldRegistry.GetByName("Help")
	helpListDisplay := interfaces.NewListDisplay(helpField)
	helpListDisplay.Ordering = 5
	settingmodelAdminPage.ListDisplay.AddField(helpListDisplay)
	categoryField, _ := settingForm.FieldRegistry.GetByName("Category")
	categoryListDisplay := interfaces.NewListDisplay(categoryField)
	categoryListDisplay.Populate = func(m interface{}) string {
		return m.(*settingmodel.Setting).Category.Name
	}
	categoryListDisplay.Ordering = 6
	settingmodelAdminPage.ListDisplay.AddField(categoryListDisplay)
	err = settingsAdminPage.SubPages.AddAdminPage(settingmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
	settingcategoriesmodelAdminPage := admin.NewGormAdminPage(settingsAdminPage, func() (interface{}, interface{}) {return &settingmodel.SettingCategory{}, &[]*settingmodel.SettingCategory{}}, "settingcategory")
	settingcategoriesmodelAdminPage.PageName = "Setting categories"
	settingcategoriesmodelAdminPage.Slug = "settingcategory"
	settingcategoriesmodelAdminPage.BlueprintName = "setting"
	settingcategoriesmodelAdminPage.Router = mainRouter
	adminContext = &interfaces.AdminContext{}
	settingCategoryForm := interfaces.NewFormFromModelFromGinContext(adminContext, &settingmodel.SettingCategory{}, make([]string, 0), []string{}, true, "")
	settingcategoriesmodelAdminPage.Form = settingCategoryForm
	settingCategoryNameField, _ := settingCategoryForm.FieldRegistry.GetByName("Name")
	settingCategoryNameListDisplay := interfaces.NewListDisplay(settingCategoryNameField)
	settingcategoriesmodelAdminPage.ListDisplay.AddField(settingCategoryNameListDisplay)
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
