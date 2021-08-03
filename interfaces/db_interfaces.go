package interfaces

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
	"text/template"
)

type IDbAdapter interface {
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
	GetStringToExtractYearFromField(filterOptionField string) string
	GetStringToExtractMonthFromField(filterOptionField string) string
	Exact(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	IExact(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Contains(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	IContains(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	In(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Gt(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Gte(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Lt(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Lte(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	StartsWith(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	IStartsWith(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	EndsWith(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	IEndsWith(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Range(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Date(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Year(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Month(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Day(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Week(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	WeekDay(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Quarter(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Time(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Hour(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Minute(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Second(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	IsNull(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	Regex(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
	IRegex(operatorContext *GormOperatorContext, field *schema.Field, value interface{})
}

var Db *gorm.DB

type UadminDatabase struct {
	Db *gorm.DB
	Adapter IDbAdapter
}

func (uad *UadminDatabase) Close() {
	db, _ := uad.Db.DB()
	db.Close()
}

func NewUadminDatabase(alias_ ...string) *UadminDatabase {
	var alias string
	if len(alias_) == 0 {
		alias = "default"
	} else {
		alias = alias_[0]
	}
	adapter := GetAdapterForDb(alias)
	Db, _ = adapter.GetDb(
		alias,
	)
	return &UadminDatabase{Db: Db, Adapter: adapter}
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
	//d.databases[alias] = &UadminDatabase{
	//	db: databaseOpened,
	//	dialect: NewDbAdapter(databaseOpened, d.config.D.Db.Default.Type),
	//}
	// return d.databases[alias].db
	return GetDB(alias)
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
	var err error

	// Check if there is a database config file
	dialect := GetAdapterForDb(alias)
	Db, err = dialect.GetDb(
		alias,
	)
	if err != nil {
		Trail(ERROR, "unable to connect to DB. %s", err)
		Db.Error = fmt.Errorf("unable to connect to DB. %s", err)
	}
	return Db
}



func GetAdapterForDb(alias_ ...string) IDbAdapter {
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
	return NewDbAdapter(Db, databaseConfig.Type)
}

