package database

import (
	config2 "github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/dialect"
	"text/template"

	"bytes"
	"strings"

	"github.com/oleiade/reflections"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UadminDatabase struct {
	db *gorm.DB
	dialect dialect.DbDialect
}
type Database struct {
	config    *config2.UadminConfig
	databases map[string]*UadminDatabase
}

var (
	postgresDsnTemplate, _ = template.New("postgresdsn").Parse("host={{.Host}} user={{.User}} password={{.Password}} dbname={{.Name}} port=5432 sslmode=disable TimeZone=UTC")
)

func NewDatabase(config *config2.UadminConfig) *Database {
	database := Database{}
	database.config = config
	database.databases = make(map[string]*UadminDatabase)
	return &database
}

func (d Database) ConnectTo(alias string) *gorm.DB {
	if alias == "" {
		alias = "default"
	}
	database := d.databases[alias]
	if database == nil {
		var tplBytes bytes.Buffer
		databaseConfig, _ := reflections.GetField(d.config.D.Db, strings.Title(alias))
		err := postgresDsnTemplate.Execute(&tplBytes, databaseConfig)
		if err != nil {
			panic(err)
		}
		databaseOpened, err := gorm.Open(postgres.Open(tplBytes.String()), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		d.databases[alias] = &UadminDatabase{
			db: databaseOpened,
			dialect: dialect.NewDbDialect(databaseOpened, d.config.D.Db.Default.Type),
		}
	}
	return d.databases[alias].db
}
