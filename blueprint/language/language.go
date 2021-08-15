package language

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/blueprint/language/migrations"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	languageAdminPage := admin.NewGormAdminPage(nil, func() (interface{}, interface{}) {return nil, nil}, "")
	languageAdminPage.PageName = "Languages"
	languageAdminPage.Slug = "language"
	languageAdminPage.BlueprintName = "language"
	languageAdminPage.Router = mainRouter
	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(languageAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing language blueprint: %s", err))
	}
	languagemodelAdminPage := admin.NewGormAdminPage(languageAdminPage, func() (interface{}, interface{}) {return &interfaces.Language{}, &[]*interfaces.Language{}}, "language")
	languagemodelAdminPage.PageName = "Languages"
	languagemodelAdminPage.Slug = "language"
	languagemodelAdminPage.BlueprintName = "language"
	languagemodelAdminPage.Router = mainRouter
	adminContext := &interfaces.AdminContext{}
	languageForm := interfaces.NewFormFromModelFromGinContext(adminContext, &interfaces.Language{}, make([]string, 0), []string{}, true, "")
	languagemodelAdminPage.Form = languageForm
	languagemodelAdminPage.ListDisplay.ClearAllFields()
	codeField, _ := languageForm.FieldRegistry.GetByName("Code")
	codeListDisplay := interfaces.NewListDisplay(codeField)
	codeListDisplay.Ordering = 1
	languagemodelAdminPage.ListDisplay.AddField(codeListDisplay)
	nameField, _ := languageForm.FieldRegistry.GetByName("Name")
	nameListDisplay := interfaces.NewListDisplay(nameField)
	nameListDisplay.Ordering = 2
	languagemodelAdminPage.ListDisplay.AddField(nameListDisplay)
	englishNameField, _ := languageForm.FieldRegistry.GetByName("EnglishName")
	englishNameListDisplay := interfaces.NewListDisplay(englishNameField)
	englishNameListDisplay.Ordering = 3
	languagemodelAdminPage.ListDisplay.AddField(englishNameListDisplay)
	activeField, _ := languageForm.FieldRegistry.GetByName("Active")
	activeListDisplay := interfaces.NewListDisplay(activeField)
	activeListDisplay.Ordering = 4
	languagemodelAdminPage.ListDisplay.AddField(activeListDisplay)
	availableInGuiField, _ := languageForm.FieldRegistry.GetByName("AvailableInGui")
	availableInGuiListDisplay := interfaces.NewListDisplay(availableInGuiField)
	availableInGuiListDisplay.Ordering = 5
	languagemodelAdminPage.ListDisplay.AddField(availableInGuiListDisplay)
	err = languageAdminPage.SubPages.AddAdminPage(languagemodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing language blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &interfaces.Language{}})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "language",
		Description:       "Language blueprint is responsible for managing languages used in the project",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
