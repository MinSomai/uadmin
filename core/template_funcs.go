package core

import (
	"fmt"
	"html/template"
	"strings"
)

func add(n1 int, n2 int) int {
	return n1 + n2
}

func mul(n1 int, n2 int) int {
	return n1 * n2
}

func safe(s string) template.HTML {
	return template.HTML(s)
}

func GenerateAttrs(attrs map[string]string) template.HTML {
	attrsContent := make([]string, 0)
	for k, v := range attrs {
		attrsContent = append(attrsContent, fmt.Sprintf(" %s=\"%s\" ", template.HTMLAttr(k), template.HTML(v)))
	}
	return template.HTML(strings.Join(attrsContent, " "))
}

func attr(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}

var FuncMap = template.FuncMap{
	"Tf":            Tf,
	"add":           add,
	"mul":           mul,
	"safe":          safe,
	"attr":          attr,
	"GenerateAttrs": GenerateAttrs,
	//"CSRF": func() string {
	//	return "dfsafsa"
	//	// return authapi.GetSession(r)
	//},
}
