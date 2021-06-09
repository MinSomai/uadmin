package logging

import (
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/logging/migrations"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(group *gin.RouterGroup) {
}

func (b Blueprint) Init(config *config.UadminConfig) {
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "logging",
		Description:       "Logging blueprint is responsible to store all actions made through admin panel",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
