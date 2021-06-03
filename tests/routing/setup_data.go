package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/interfaces"
)

type BlueprintRouting struct {
	interfaces.Blueprint
}
var ConcreteBlueprint BlueprintRouting

func (b BlueprintRouting) InitRouter(group *gin.RouterGroup) {
}

func init() {
	ConcreteBlueprint = BlueprintRouting{
		interfaces.Blueprint{
			Name:              "user",
			Description:       "blueprint",
			MigrationRegistry: interfaces.NewMigrationRegistry(),
		},
	}
}