package main

import (
	"github.com/uadmin/uadmin/app"
	"os"
)

func main() {
	environment := os.Getenv("environment")
	if environment == "" {
		environment = "dev"
	}
	app1 := app.NewApp(environment)
	app1.StartAdmin()
}
