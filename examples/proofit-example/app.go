package proofit_example

import (
	"embed"
	proofitcore2 "github.com/sergeyglazyrindev/proofit-example/blueprint/proofitcore"
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin/blueprint/abtest"
	"github.com/sergeyglazyrindev/uadmin/blueprint/approval"
	logblueprint "github.com/sergeyglazyrindev/uadmin/blueprint/logging"
	settingsblueprint "github.com/sergeyglazyrindev/uadmin/blueprint/settings"
	"github.com/sergeyglazyrindev/uadmin/core"
)

func NewProofitApp(environment string) *uadmin.App {
	app1 := uadmin.NewApp(environment, true)
	// next two lines are mandatory for uadmin to determine your blueprints and everything else.
	app1.BlueprintRegistry.Register(proofitcore2.ConcreteBlueprint)
	app1.BlueprintRegistry.DeRegister(abtest.ConcreteBlueprint)
	app1.BlueprintRegistry.DeRegister(approval.ConcreteBlueprint)
	app1.BlueprintRegistry.DeRegister(logblueprint.ConcreteBlueprint)
	app1.BlueprintRegistry.DeRegister(settingsblueprint.ConcreteBlueprint)
	app1.RegisterCommand("generate-fake-data", &CreateFakedDataCommand{})
	app1.Initialize()
	app1.InitializeRouter()
	core.CurrentConfig.OverridenTemplatesFS = &templatesRoot
	currentApp = app1
	return app1
}

//go:embed templates
var templatesRoot embed.FS

var currentApp *uadmin.App
