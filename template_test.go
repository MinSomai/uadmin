package uadmin

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

func TestEmbeddingTemplates(t *testing.T) {
	app := NewApp("test")
	t1, _ := template.ParseFS(app.Config.TemplatesFS, "templates/test.html")
	templateBuffer := &bytes.Buffer{}
	t1.Execute(templateBuffer, struct {
		Title string
	}{Title: "test"})
	assert.Contains(t, templateBuffer.String(), "test")
}
