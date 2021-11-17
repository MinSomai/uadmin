package example

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/uadmin/core"
	"github.com/sergeyglazyrindev/uadmin_example/blueprint/example/migrations"
	"github.com/sergeyglazyrindev/uadmin_example/blueprint/example/models"
)

type Blueprint struct {
	core.Blueprint
}

func (b Blueprint) InitRouter(app core.IApp, group *gin.RouterGroup) {
	// initialize administrator page for this blueprint.
	todosAdminPage := core.NewGormAdminPage(
		nil,
		nil,
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	todosAdminPage.PageName = "Example"
	todosAdminPage.Slug = "example"
	todosAdminPage.BlueprintName = "example"
	todosAdminPage.Router = app.GetRouter()
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(todosAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing blueprint: %s", err))
	}
	// initialize administrator page for your specific model.
	todosModelAdminPage := core.NewGormAdminPage(
		todosAdminPage,
		&models.Todo{},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			// define fields that you want to have in your admin panel
			fields := []string{"TaskAlias", "TaskDescription"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			return form
		},
	)
	todosModelAdminPage.PageName = "Todos"
	todosModelAdminPage.Slug = "todo"
	todosModelAdminPage.BlueprintName = "example"
	todosModelAdminPage.Router = app.GetRouter()
	err = todosAdminPage.SubPages.AddAdminPage(todosModelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing blueprint: %s", err))
	}
}

func (b Blueprint) InitApp(app core.IApp) {
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &models.Todo{}, &[]*models.Todo{} })
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "example",
		Description:       "blueprint for uadmin example",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
