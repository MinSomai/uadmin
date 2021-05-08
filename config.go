package uadmin

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
)

// Info from config file
type UadminConfig struct {
	D struct {
		Test string `yaml:"test"`
		Db   struct {
			Default *DBSettings
		} `yaml: "db"`
		Auth struct {
			JWT_SECRET_TOKEN string `yaml:"jwt_secret_token"`
		} `yaml: "auth"`
	}
}

// Reads info from config file
func NewConfig(file string) *UadminConfig {
	_, err := os.Stat(file)
	if err != nil {
		log.Fatal("Config file is missing: ", file)
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	c := new(UadminConfig)
	err = yaml.Unmarshal([]byte(content), &c.D)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	spew.Dump(c.D.Db.Default)
	return c
}
