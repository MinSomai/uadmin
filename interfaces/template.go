package interfaces

import (
	"github.com/gin-gonic/gin"
	"io/fs"
	"text/template"
)

type ITemplateRenderer interface {
	AddFuncMap(funcName string, concreteFunc interface{})
	Render(ctx *gin.Context, fsys fs.FS, path string, data interface{}, funcs ...template.FuncMap)
}
