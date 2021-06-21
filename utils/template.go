package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/debug"
	"io/fs"
	"reflect"
	"runtime"
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
// RenderHTML creates a new template and applies a parsed template to the specified
// data object. For function, Tf is available by default and if you want to add functions
//to your template, just add them to funcs which will add them to the template with their
// original function names. If you added anonymous functions, they will be available in your
// templates as func1, func2 ...etc.
func RenderHTML(ctx *gin.Context, fsys fs.FS, path string, data interface{}, funcs ...interface{}) error {
	var err error
	var funcVal reflect.Value
	var funcName string

	funcMap := template.FuncMap{}
	for k,v := range FuncMap {
		funcMap[k] = v
	}
	for i := range funcs {
		funcVal = reflect.ValueOf(funcs[i])
		if funcVal.Type().Kind() != reflect.Func {
			Trail(WARNING, "Interface passed to RenderHTML in funcs parameter should only be a function. Got (%s) in position %d", funcVal.Type().Kind(), i)
			continue
		}

		funcName = runtime.FuncForPC(funcVal.Pointer()).Name()
		funcName = funcName[strings.LastIndex(funcName, ".")+1:]
		funcMap[funcName] = funcs[i]
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
	newT := template.New(templateNameParts[len(templateNameParts) - 1]).Funcs(funcMap)
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
