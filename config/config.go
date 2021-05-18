package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// DBSettings !
type DBSettings struct {
	Type     string `json:"type"` // sqlite, mysql
	Name     string `json:"name"` // File/DB name
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

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
		Admin struct {
			ListenPort int `yaml:"listen_port"`
			SSL        struct {
				ListenPort int `yaml:"listen_port"`
			} `yaml:"ssl"`
		} `yaml:"admin"`
		Api struct {
			ListenPort int `yaml:"listen_port"`
			SSL        struct {
				ListenPort int `yaml:"listen_port"`
			} `yaml:"ssl"`
		} `yaml:"api"`
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
	return c
}
