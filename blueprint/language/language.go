package language

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/blueprint/language/migrations"
	"github.com/uadmin/uadmin/blueprint/language/models"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	languageAdminPage := admin.NewGormAdminPage(func() interface{} {return nil}, "")
	languageAdminPage.PageName = "Languages"
	languageAdminPage.Slug = "language"
	languageAdminPage.BlueprintName = "language"
	languageAdminPage.Router = mainRouter
	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(languageAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing language blueprint: %s", err))
	}
	languagemodelAdminPage := admin.NewGormAdminPage(func() interface{} {return &models.Language{}}, "language")
	languagemodelAdminPage.PageName = "Languages"
	languagemodelAdminPage.Slug = "language"
	languagemodelAdminPage.BlueprintName = "language"
	languagemodelAdminPage.Router = mainRouter
	err = languageAdminPage.SubPages.AddAdminPage(languagemodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing language blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	interfaces.ProjectModels.RegisterModel(&models.Language{})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "language",
		Description:       "Language blueprint is responsible for managing languages used in the project",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
