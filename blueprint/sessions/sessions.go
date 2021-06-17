package sessions

import (
	"github.com/gin-gonic/gin"
	interfaces2 "github.com/uadmin/uadmin/blueprint/sessions/interfaces"
	"github.com/uadmin/uadmin/blueprint/sessions/migrations"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
	SessionAdapterRegistry *interfaces2.SessionProviderRegistry
}

func (b Blueprint) InitRouter(group *gin.RouterGroup) {
}

func (b Blueprint) Init(config *config.UadminConfig) {
	b.SessionAdapterRegistry.RegisterNewAdapter(&interfaces2.DbSession{}, true)
}

var ConcreteBlueprint = Blueprint{
	Blueprint: interfaces.Blueprint{
		Name:              "sessions",
		Description:       "Sessions blueprint responsible to keep session data in database",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
	SessionAdapterRegistry: interfaces2.NewSessionRegistry(),
}
