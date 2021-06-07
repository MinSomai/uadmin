package sessions

import (
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/sessions/migrations"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(group *gin.RouterGroup) {
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "sessions",
		Description:       "Sessions blueprint responsible to keep session data in database",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
