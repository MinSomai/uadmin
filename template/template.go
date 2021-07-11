package template

import (
	"github.com/uadmin/uadmin/utils"
	"text/template"
)

func add(n1 int, n2 int) int {
	return n1 + n2
}

var FuncMap = template.FuncMap{
	"Tf": utils.Tf,
	"add": add,
	//"CSRF": func() string {
	//	return "dfsafsa"
	//	// return authapi.GetSession(r)
	//},
}

