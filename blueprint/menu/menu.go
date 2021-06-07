
package menu

import (
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/menu/migrations"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(group *gin.RouterGroup) {
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "menu",
		Description:       "Dashboard menu blueprint responsible for menu in admin panel",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
