package abtest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/blueprint/abtest/migrations"
	abtestmodel "github.com/uadmin/uadmin/blueprint/abtest/models"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	abTestAdminPage := admin.NewGormAdminPage(func() interface{} {return nil}, "")
	abTestAdminPage.PageName = "AB Tests"
	abTestAdminPage.Slug = "abtest"
	abTestAdminPage.BlueprintName = "abtest"
	abTestAdminPage.Router = mainRouter
	// abTestAdminPage.ListHandler = templatecontext.BuildAdminHandlerForBlueprintfunc(abTestAdminPage.PageName)
	//func (ctx *gin.Context) {
	//	type Context struct {
	//		templatecontext.AdminContext
	//		Menu     string
	//	}
	//
	//	c := &Context{}
	//	templatecontext.PopulateTemplateContextForAdminPanel(ctx, c, templatecontext.NewAdminRequestParams())
	//	menu := string(admin.CurrentDashboardAdminPanel.AdminPages.PreparePagesForTemplate(c.UserPermissionRegistry))
	//	c.Menu = menu
	//	tr := interfaces.NewTemplateRenderer("AB Tests")
	//	tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("home"), c, template.FuncMap)
	//}

	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(abTestAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
	abtestmodelAdminPage := admin.NewGormAdminPage(func() interface{} {return &abtestmodel.ABTest{}}, "abtest")
	abtestmodelAdminPage.PageName = "AB Tests"
	abtestmodelAdminPage.Slug = "abtest"
	abtestmodelAdminPage.BlueprintName = "abtest"
	abtestmodelAdminPage.Router = mainRouter
	err = abTestAdminPage.SubPages.AddAdminPage(abtestmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
	abtestvaluemodelAdminPage := admin.NewGormAdminPage(func() interface{} {return &abtestmodel.ABTestValue{}}, "abtestvalue")
	abtestvaluemodelAdminPage.PageName = "AB Test Values"
	abtestvaluemodelAdminPage.Slug = "abtestvalue"
	abtestvaluemodelAdminPage.BlueprintName = "abtest"
	abtestvaluemodelAdminPage.Router = mainRouter
	err = abTestAdminPage.SubPages.AddAdminPage(abtestvaluemodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	interfaces.ProjectModels.RegisterModel(&abtestmodel.ABTestValue{})
	interfaces.ProjectModels.RegisterModel(&abtestmodel.ABTest{})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "abtest",
		Description:       "ABTest blueprint is responsible for ab tests",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
