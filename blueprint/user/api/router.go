package api

import (
	"github.com/uadmin/uadmin/utils"

	"github.com/gin-gonic/gin"
)

func InitializeRouter(r *gin.Engine) {
	utils.InitializeRouter(r, "users", apiHandlers)
}
