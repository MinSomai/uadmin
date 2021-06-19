package approval

import (
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/approval/migrations"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
}

func (b Blueprint) Init(config *config.UadminConfig) {
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "approval",
		Description:       "Approval blueprint is responsible for approving things in the project",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
