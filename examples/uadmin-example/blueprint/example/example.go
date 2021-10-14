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

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	// initialize administrator page for this blueprint.
	todosAdminPage := core.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) { return nil, nil },
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	todosAdminPage.PageName = "Example"
	todosAdminPage.Slug = "example"
	todosAdminPage.BlueprintName = "example"
	todosAdminPage.Router = mainRouter
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(todosAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing blueprint: %s", err))
	}
	// initialize administrator page for your specific model.
	todosModelAdminPage := core.NewGormAdminPage(
		todosAdminPage,
		func() (interface{}, interface{}) {
			return &models.Todo{}, &[]*models.Todo{}
		},
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
	todosModelAdminPage.Router = mainRouter
	err = todosAdminPage.SubPages.AddAdminPage(todosModelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	// if you want to use admin page for your model, please make sure you registered your model in project models.
	core.ProjectModels.RegisterModel(func() interface{} { return &models.Todo{} })
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "example",
		Description:       "blueprint for uadmin example",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
