package core

import "text/template"

func add(n1 int, n2 int) int {
	return n1 + n2
}

func mul(n1 int, n2 int) int {
	return n1 * n2
}

var FuncMap = template.FuncMap{
	"Tf": Tf,
	"add": add,
	"mul": mul,
	//"CSRF": func() string {
	//	return "dfsafsa"
	//	// return authapi.GetSession(r)
	//},
}
