package core

import "github.com/gin-gonic/gin"

var PopulateTemplateContextForAdminPanel func(ctx *gin.Context, context IAdminContext, adminRequestParams *AdminRequestParams)

