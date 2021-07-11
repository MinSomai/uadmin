package language

import (
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/language/migrations"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
}

func (b Blueprint) Init(config *interfaces.UadminConfig) {
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "language",
		Description:       "Language blueprint is responsible for managing languages used in the project",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
