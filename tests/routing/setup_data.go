package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/core"
)

type BlueprintRouting struct {
	core.Blueprint
}
var ConcreteBlueprint BlueprintRouting
var visited = false

func (b BlueprintRouting) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	group.GET("/visit", func(c *gin.Context) {
		visited = true
	})
}

func init() {
	ConcreteBlueprint = BlueprintRouting{
		core.Blueprint{
			Name:              "user",
			Description:       "blueprint",
			MigrationRegistry: core.NewMigrationRegistry(),
		},
	}
}