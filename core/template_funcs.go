package core

import "html/template"

func add(n1 int, n2 int) int {
	return n1 + n2
}

func mul(n1 int, n2 int) int {
	return n1 * n2
}

func safe(s string) template.HTML {
	return template.HTML(s)
}

func attr(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}

var FuncMap = template.FuncMap{
	"Tf":  Tf,
	"add": add,
	"mul": mul,
	"safe": safe,
	"attr": attr,
	//"CSRF": func() string {
	//	return "dfsafsa"
	//	// return authapi.GetSession(r)
	//},
}
