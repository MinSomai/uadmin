package settings

import (
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/settings/migrations"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(group *gin.RouterGroup) {
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "settings",
		Description:       "Settings blueprint responsible for wide-project settings",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
