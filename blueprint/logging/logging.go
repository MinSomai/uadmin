package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/blueprint/logging/migrations"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	logAdminPage := admin.NewAdminPage()
	logAdminPage.PageName = "Logs"
	logAdminPage.Slug = "log"
	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(logAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing log blueprint: %s", err))
	}
	logmodelAdminPage := admin.NewAdminPage()
	logmodelAdminPage.PageName = "Logs"
	logmodelAdminPage.Slug = "log"
	err = logAdminPage.SubPages.AddAdminPage(logmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing log blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "logging",
		Description:       "Logging blueprint is responsible to store all actions made through admin panel",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
