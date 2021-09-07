package example

import (
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/uadmin_example/blueprint/example/migrations"
	"github.com/sergeyglazyrindev/uadmin/core"
)

type Blueprint struct {
	core.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
}

func (b Blueprint) Init() {
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "example",
		Description:       "blueprint for uadmin example",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
