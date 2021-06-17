package auth

import (
	"github.com/gin-gonic/gin"
	interfaces3 "github.com/uadmin/uadmin/blueprint/auth/interfaces"
	"github.com/uadmin/uadmin/blueprint/auth/migrations"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
	AuthAdapterRegistry *interfaces3.AuthProviderRegistry
}

func (b Blueprint) InitRouter(group *gin.RouterGroup) {
	for adapter := range b.AuthAdapterRegistry.Iterate() {
		adapterGroup := group.Group("/" + adapter.GetName())
		adapterGroup.POST("/signin/", adapter.Signin)
		adapterGroup.POST("/signup/", adapter.Signup)
		adapterGroup.POST("/logout/", adapter.Logout)
		adapterGroup.GET("/status/", adapter.IsAuthenticated)
	}
}

func (b Blueprint) Init(config *config.UadminConfig) {
	b.AuthAdapterRegistry.RegisterNewAdapter(&interfaces3.DirectAuthProvider{})
	b.AuthAdapterRegistry.RegisterNewAdapter(&interfaces3.TokenAuthProvider{})
}

var ConcreteBlueprint = Blueprint{
	Blueprint: interfaces.Blueprint{
		Name:              "auth",
		Description:       "blueprint for auth functionality",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
	AuthAdapterRegistry: interfaces3.NewAuthProviderRegistry(),
}
