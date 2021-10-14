package proofit_example

import (
	"embed"
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin/core"
)

func NewProofitApp(environment string) *uadmin.App {
	app1 := uadmin.NewApp(environment, true)
	// next two lines are mandatory for uadmin to determine your blueprints and everything else.
	app1.Initialize()
	app1.InitializeRouter()
	core.CurrentConfig.OverridenTemplatesFS = &templatesRoot
	return app1
}

//go:embed templates
var templatesRoot embed.FS
