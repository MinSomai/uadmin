package core

import (
	"github.com/gin-gonic/gin"
)

type IApp interface {
	GetConfig() *UadminConfig
	GetDatabase() *Database
	GetRouter() *gin.Engine
	GetCommandRegistry() *CommandRegistry
	GetBlueprintRegistry() IBlueprintRegistry
	GetDashboardAdminPanel() *DashboardAdminPanel
}
