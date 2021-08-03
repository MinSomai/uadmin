package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/blueprint/logging/migrations"
	logmodel "github.com/uadmin/uadmin/blueprint/logging/models"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	logAdminPage := admin.NewGormAdminPage(func() interface{} {return nil}, "")
	logAdminPage.PageName = "Logs"
	logAdminPage.Slug = "log"
	logAdminPage.BlueprintName = "logging"
	logAdminPage.Router = mainRouter
	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(logAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing log blueprint: %s", err))
	}
	logmodelAdminPage := admin.NewGormAdminPage(func() interface{} {return &logmodel.Log{}}, "log")
	logmodelAdminPage.PageName = "Logs"
	logmodelAdminPage.Slug = "log"
	logmodelAdminPage.BlueprintName = "logging"
	logmodelAdminPage.Router = mainRouter
	err = logAdminPage.SubPages.AddAdminPage(logmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing log blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	interfaces.ProjectModels.RegisterModel(&logmodel.Log{})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "logging",
		Description:       "Logging blueprint is responsible to store all actions made through admin panel",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
