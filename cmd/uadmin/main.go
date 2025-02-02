package main

import (
	"github.com/sergeyglazyrindev/uadmin"
	"os"
)

func main() {
	environment := os.Getenv("environment")
	if environment == "" {
		environment = "dev"
	}
	app1 := uadmin.NewApp(environment)
	app1.ExecuteCommand()
}
