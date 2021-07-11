package settings

import (
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/settings/migrations"
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
		Name:              "settings",
		Description:       "Settings blueprint responsible for wide-project settings",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
