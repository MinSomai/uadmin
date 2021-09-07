package main

import (
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin_example/blueprint/example"
	"os"
)

func main() {
	environment := os.Getenv("environment")
	if environment == "" {
		environment = "dev"
	}
	app1 := uadmin.NewApp(environment, true)
	// please configure everything you need for your project here.
	app1.BlueprintRegistry.Register(example.ConcreteBlueprint)
	// next two lines are mandatory for uadmin to determine your blueprints and everything else.
	app1.Initialize()
	app1.InitializeRouter()
	// and here is the command handler that will be used for any functionality you may need.
	app1.ExecuteCommand()
}
