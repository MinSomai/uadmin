package proofit_example

import "github.com/sergeyglazyrindev/uadmin"

func NewProofitApp(environment string) *uadmin.App {
	app1 := uadmin.NewApp(environment, true)
	// next two lines are mandatory for uadmin to determine your blueprints and everything else.
	app1.Initialize()
	app1.InitializeRouter()
	return app1
}
