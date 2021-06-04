package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/interfaces"
)

type BlueprintRouting struct {
	interfaces.Blueprint
}
var ConcreteBlueprint BlueprintRouting
var visited = false

func (b BlueprintRouting) InitRouter(group *gin.RouterGroup) {
	group.GET("/visit", func(c *gin.Context) {
		visited = true
	})
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