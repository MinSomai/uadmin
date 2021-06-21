package utils

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/debug"
	"io/fs"
	"strings"
	"text/template"
)

var FuncMap = template.FuncMap{
	"Tf": Tf,
	//"CSRF": func() string {
	//	return "dfsafsa"
	//	// return authapi.GetSession(r)
	//},
}

type IncludeContext struct {
	SiteName    string
	PageTitle string
}

type TemplateRenderer struct {
	funcMap template.FuncMap
	pageTitle string
}

func (tr *TemplateRenderer) AddFuncMap(funcName string, concreteFunc interface{}) {
	tr.funcMap[funcName] = concreteFunc
}

func (tr *TemplateRenderer) Render(ctx *gin.Context, fsys fs.FS, path string, data interface{}, funcs ...template.FuncMap) {
	Include := func(funcs1 template.FuncMap) func (templateName string) string {
		return func (templateName string) string {
			templateWriter := bytes.NewBuffer([]byte{})
			err := RenderHTMLAsString(templateWriter, fsys, config.CurrentConfig.GetPathToTemplate(templateName), data, funcs1)
			if err != nil {
				debug.Trail(debug.CRITICAL, "Error while parsing include of the template %s", templateName)
				panic(err)
			}
			return templateWriter.String()
		}
    }
	PageTitle := func() string {
		return fmt.Sprintf("%s - %s", config.CurrentConfig.D.Uadmin.SiteName, tr.pageTitle)
	}
	var funcs1 template.FuncMap
	if len(funcs) == 0 {
		funcs1 = template.FuncMap{}
		funcs1["Include"] = Include(funcs1)
		funcs1["PageTitle"] = PageTitle
	} else {
		funcs1 = funcs[0]
		funcs1["Include"] = Include(funcs1)
		funcs1["PageTitle"] = PageTitle
	}
	RenderHTML(ctx, fsys, path, data, funcs1)
}

func NewTemplateRenderer(pageTitle string) *TemplateRenderer {
	templateRenderer := TemplateRenderer{funcMap: template.FuncMap{}, pageTitle: pageTitle}
	return &templateRenderer
}


// RenderHTML creates a new template and applies a parsed template to the specified
// data object. For function, Tf is available by default and if you want to add functions
//to your template, just add them to funcs which will add them to the template with their
// original function names. If you added anonymous functions, they will be available in your
// templates as func1, func2 ...etc.
func RenderHTML(ctx *gin.Context, fsys fs.FS, path string, data interface{}, funcs ...template.FuncMap) error {
	var err error

	var funcs1 template.FuncMap
	if len(funcs) == 0 {
		funcs1 = template.FuncMap{}
	} else {
		funcs1 = funcs[0]
	}
	for k,v := range FuncMap {
		funcs1[k] = v
	}
	//// Check for ABTesting cookie
	//if cookie, err := ctx.Cookie("abt"); err != nil || cookie == nil {
	//	now := time.Now().AddDate(0, 0, 1)
	//	cookie1 := &http.Cookie{
	//		Name:    "abt",
	//		Value:   fmt.Sprint(now.Second()),
	//		Path:    "/",
	//		Expires: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
	//	}
	//	http.SetCookie(ctx.Writer, cookie1)
	//}
	templateNameParts := strings.Split(path, "/")
	newT := template.New(templateNameParts[len(templateNameParts) - 1]).Funcs(funcs1)
	newT, err = newT.ParseFS(fsys, path)
	if err != nil {
		ctx.String(500, err.Error())
		debug.Trail(debug.ERROR, "RenderHTML unable to parse %s. %s", path, err)
		return err
	}
	err = newT.Execute(ctx.Writer, data)
	if err != nil {
		// ctx.AbortWithStatus(500)
		ctx.String(500, err.Error())
		debug.Trail(debug.ERROR, "RenderHTML unable to parse %s. %s", path, err)
		return err
	}
	return nil
}

// RenderHTML creates a new template and applies a parsed template to the specified
// data object. For function, Tf is available by default and if you want to add functions
//to your template, just add them to funcs which will add them to the template with their
// original function names. If you added anonymous functions, they will be available in your
// templates as func1, func2 ...etc.
func RenderHTMLAsString(writer *bytes.Buffer, fsys fs.FS, path string, data interface{}, funcs ...template.FuncMap) error {
	var err error

	var funcs1 template.FuncMap
	if len(funcs) == 0 {
		funcs1 = template.FuncMap{}
	} else {
		funcs1 = funcs[0]
		for funcName, handler := range funcs[0] {
			funcs1[funcName] = handler
		}
	}
	for k,v := range FuncMap {
		funcs1[k] = v
	}

	//// Check for ABTesting cookie
	//if cookie, err := ctx.Cookie("abt"); err != nil || cookie == nil {
	//	now := time.Now().AddDate(0, 0, 1)
	//	cookie1 := &http.Cookie{
	//		Name:    "abt",
	//		Value:   fmt.Sprint(now.Second()),
	//		Path:    "/",
	//		Expires: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
	//	}
	//	http.SetCookie(ctx.Writer, cookie1)
	//}
	templateNameParts := strings.Split(path, "/")
	newT := template.New(templateNameParts[len(templateNameParts) - 1]).Funcs(funcs1)
	newT, err = newT.ParseFS(fsys, path)
	if err != nil {
		debug.Trail(debug.ERROR, "RenderHTML unable to parse %s. %s", path, err)
		return err
	}
	err = newT.Execute(writer, data)
	if err != nil {
		// ctx.AbortWithStatus(500)
		debug.Trail(debug.ERROR, "RenderHTML unable to parse %s. %s", path, err)
		return err
	}
	return nil
}
