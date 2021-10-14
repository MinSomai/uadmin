package main

import (
	proofit_example "github.com/sergeyglazyrindev/proofit-example"
	"os"
)

func main() {
	environment := os.Getenv("environment")
	if environment == "" {
		environment = "dev"
	}
	app1 := proofit_example.NewProofitApp(environment)
	app1.ExecuteCommand()
}
