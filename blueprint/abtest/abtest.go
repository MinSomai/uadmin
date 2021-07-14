package abtest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/blueprint/abtest/migrations"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	abTestAdminPage := admin.NewAdminPage("")
	abTestAdminPage.PageName = "AB Tests"
	abTestAdminPage.Slug = "abtest"
	abTestAdminPage.BlueprintName = "abtest"
	abTestAdminPage.Router = group
	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(abTestAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
	abtestmodelAdminPage := admin.NewAdminPage("abtest")
	abtestmodelAdminPage.PageName = "AB Tests"
	abtestmodelAdminPage.Slug = "abtest"
	abtestmodelAdminPage.BlueprintName = "abtest"
	abtestmodelAdminPage.Router = group
	err = abTestAdminPage.SubPages.AddAdminPage(abtestmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
	abtestvaluemodelAdminPage := admin.NewAdminPage("abtestvalue")
	abtestvaluemodelAdminPage.PageName = "AB Test Values"
	abtestvaluemodelAdminPage.Slug = "abtestvalue"
	abtestvaluemodelAdminPage.BlueprintName = "abtest"
	abtestvaluemodelAdminPage.Router = group
	err = abTestAdminPage.SubPages.AddAdminPage(abtestvaluemodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "abtest",
		Description:       "ABTest blueprint is responsible for ab tests",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
