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
	abTestAdminPage := admin.NewGormAdminPage(nil, func() (interface{}, interface{}) {return nil, nil}, "")
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
	abtestmodelAdminPage := admin.NewGormAdminPage(abTestAdminPage, func() (interface{}, interface{}) {return &abtestmodel.ABTest{}, &[]*abtestmodel.ABTest{}}, "abtest")
	abtestmodelAdminPage.PageName = "AB Tests"
	abtestmodelAdminPage.Slug = "abtest"
	abtestmodelAdminPage.BlueprintName = "abtest"
	abtestmodelAdminPage.Router = mainRouter
	adminContext := &interfaces.AdminContext{}
	abTestForm := interfaces.NewFormFromModelFromGinContext(adminContext, &abtestmodel.ABTest{}, make([]string, 0), []string{}, true, "")
	abtestmodelAdminPage.Form = abTestForm
	nameField, _ := abTestForm.FieldRegistry.GetByName("Name")
	nameListDisplay := interfaces.NewListDisplay(nameField)
	nameListDisplay.Ordering = 2
	abtestmodelAdminPage.ListDisplay.AddField(nameListDisplay)
	typeField, _ := abTestForm.FieldRegistry.GetByName("Type")
	typeListDisplay := interfaces.NewListDisplay(typeField)
	typeListDisplay.Populate = func(m interface{}) string {
		return abtestmodel.HumanizeAbTestType(m.(*abtestmodel.ABTest).Type)
	}
	typeListDisplay.Ordering = 3
	abtestmodelAdminPage.ListDisplay.AddField(typeListDisplay)
	staticPathField, _ := abTestForm.FieldRegistry.GetByName("StaticPath")
	staticPathListDisplay := interfaces.NewListDisplay(staticPathField)
	staticPathListDisplay.Ordering = 4
	abtestmodelAdminPage.ListDisplay.AddField(staticPathListDisplay)
	modelNameField, _ := abTestForm.FieldRegistry.GetByName("ModelName")
	modelNameListDisplay := interfaces.NewListDisplay(modelNameField)
	modelNameListDisplay.Ordering = 5
	abtestmodelAdminPage.ListDisplay.AddField(modelNameListDisplay)
	fieldField, _ := abTestForm.FieldRegistry.GetByName("Field")
	fieldListDisplay := interfaces.NewListDisplay(fieldField)
	fieldListDisplay.Ordering = 6
	abtestmodelAdminPage.ListDisplay.AddField(fieldListDisplay)
	primaryKeyField, _ := abTestForm.FieldRegistry.GetByName("PrimaryKey")
	primaryKeyListDisplay := interfaces.NewListDisplay(primaryKeyField)
	primaryKeyListDisplay.Ordering = 7
	abtestmodelAdminPage.ListDisplay.AddField(primaryKeyListDisplay)
	activeField, _ := abTestForm.FieldRegistry.GetByName("Active")
	activeListDisplay := interfaces.NewListDisplay(activeField)
	activeListDisplay.Ordering = 8
	abtestmodelAdminPage.ListDisplay.AddField(activeListDisplay)
	groupField, _ := abTestForm.FieldRegistry.GetByName("Group")
	groupListDisplay := interfaces.NewListDisplay(groupField)
	groupListDisplay.Ordering = 9
	abtestmodelAdminPage.ListDisplay.AddField(groupListDisplay)
	resetABTestField, _ := abTestForm.FieldRegistry.GetByName("ResetABTest")
	resetABTestListDisplay := interfaces.NewListDisplay(resetABTestField)
	resetABTestListDisplay.Ordering = 10
	abtestmodelAdminPage.ListDisplay.AddField(resetABTestListDisplay)
	err = abTestAdminPage.SubPages.AddAdminPage(abtestmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &abtestmodel.ABTestValue{}})
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &abtestmodel.ABTest{}})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "abtest",
		Description:       "ABTest blueprint is responsible for ab tests",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
