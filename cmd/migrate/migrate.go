package migrate

import (
	"os"
)

func main() {
	environment := os.Getenv("environment")
	if environment == "" {
		environment = "dev"
	}
	// app1 := uadmin.NewApp(environment)
}
