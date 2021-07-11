package interfaces

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"reflect"
	"text/template"
)

type DbDialect interface {
	Equals(name interface{}, args ...interface{})
	Quote(name interface{}) string
	LikeOperator() string
	ToString() string
	GetLastInsertId()
	buildClauses(clause_interfaces []clause.Interface)
	QuoteTableName(tableName string) string
	Delete(db *gorm.DB, model reflect.Value, query interface{}, args ...interface{}) *gorm.DB
	ReadRows(db *gorm.DB, customSchema bool, SQL string, m interface{}, args ...interface{}) (*sql.Rows, error)
	GetSqlDialectStrings() map[string]string
	GetDb(alias_ ...string) (*gorm.DB, error)
	CreateDb() error
	Transaction(handler func()) error
}

var Db *gorm.DB

type UadminDatabase struct {
	db *gorm.DB
	dialect DbDialect
}
type Database struct {
	config    *UadminConfig
	databases map[string]*UadminDatabase
}

var (
	postgresDsnTemplate, _ = template.New("postgresdsn").Parse("host={{.Host}} user={{.User}} password={{.Password}} dbname={{.Name}} port=5432 sslmode=disable TimeZone=UTC")
)

func NewDatabase(config *UadminConfig) *Database {
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
		//var tplBytes bytes.Buffer
		//databaseConfig, _ := reflections.GetField(d.config.D.Db, strings.Title(alias))
		//err := postgresDsnTemplate.Execute(&tplBytes, databaseConfig)
		//if err != nil {
		//	panic(err)
		//}
		//databaseOpened, err := gorm.Open(postgres.Open(tplBytes.String()), &gorm.Config{})
		//if err != nil {
		//	panic(err)
		//}
		databaseOpened := GetDB(alias)
		d.databases[alias] = &UadminDatabase{
			db: databaseOpened,
			dialect: NewDbDialect(databaseOpened, d.config.D.Db.Default.Type),
		}
	}
	return d.databases[alias].db
}

type DatabaseSettings struct {
	Default *DBSettings
}
var CurrentDatabaseSettings *DatabaseSettings

// GetDB returns a pointer to the DB
func GetDB(alias_ ...string) *gorm.DB {
	var alias string
	if len(alias_) == 0 {
		alias = "default"
	} else {
		alias = alias_[0]
	}
	if Db != nil{
		return Db
	}
	var err error

	// Check if there is a database config file
	if CurrentDatabaseSettings == nil {
		buf, err := ioutil.ReadFile(".database")
		if err == nil {
			err = json.Unmarshal(buf, CurrentDatabaseSettings)
			if err != nil {
				Trail(WARNING, ".database file is not a valid json file. %s", err)
			}
		}
	}
	dialect := GetDialectForDb(alias)
	Db, err = dialect.GetDb(
		alias,
	)
	if err != nil {
		Trail(ERROR, "unable to connect to DB. %s", err)
		Db.Error = fmt.Errorf("unable to connect to DB. %s", err)
	}
	return Db
}



func GetDialectForDb(alias_ ...string) DbDialect {
	var databaseConfig *DBSettings
	var alias string
	if len(alias_) == 0 {
		alias = "default"
	} else {
		alias = alias_[0]
	}
	if alias == "default" {
		databaseConfig = CurrentDatabaseSettings.Default
	}
	return NewDbDialect(Db, databaseConfig.Type)
}

