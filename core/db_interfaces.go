package core

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"text/template"
)

type DeleteRowStructure struct {
	SQL         string
	Values      []interface{}
	Explanation string
	Table       string
	Cond        string
}

type IDbAdapter interface {
	Equals(name interface{}, args ...interface{})
	Quote(name interface{}) string
	LikeOperator() string
	ToString() string
	GetLastInsertID()
	buildClauses(clauseInterfaces []clause.Interface)
	QuoteTableName(tableName string) string
	Delete(db *gorm.DB, model reflect.Value, query interface{}, args ...interface{}) *gorm.DB
	ReadRows(db *gorm.DB, customSchema bool, SQL string, m interface{}, args ...interface{}) (*sql.Rows, error)
	GetSQLDialectStrings() map[string]string
	GetDb(alias string, dryRun bool) (*gorm.DB, error)
	CreateDb() error
	Transaction(handler func()) error
	GetStringToExtractYearFromField(filterOptionField string) string
	GetStringToExtractMonthFromField(filterOptionField string) string
	Exact(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	IExact(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Contains(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	IContains(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	In(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Gt(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Gte(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Lt(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Lte(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	StartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	IStartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	EndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	IEndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Range(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Date(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Year(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Month(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Day(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Week(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	WeekDay(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Quarter(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Time(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Hour(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Minute(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Second(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	IsNull(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	Regex(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	IRegex(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool)
	BuildDeleteString(table string, cond string, values ...interface{}) *DeleteRowStructure
}

var Db *gorm.DB

type UadminDatabase struct {
	Db      *gorm.DB
	Adapter IDbAdapter
}

func (uad *UadminDatabase) Close() {
	db, _ := uad.Db.DB()
	db.Close()
}

func NewUadminDatabase(alias1 ...string) *UadminDatabase {
	var alias string
	if len(alias1) == 0 {
		alias = "default"
	} else {
		alias = alias1[0]
	}
	adapter := GetAdapterForDb(alias)
	Db, _ = adapter.GetDb(
		alias, false,
	)
	return &UadminDatabase{Db: Db, Adapter: adapter}
}

func NewUadminDatabaseWithoutConnection(alias1 ...string) *UadminDatabase {
	var alias string
	if len(alias1) == 0 {
		alias = "default"
	} else {
		alias = alias1[0]
	}
	adapter := GetAdapterForDb(alias)
	Db, _ = adapter.GetDb(
		alias, true,
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
func GetDB(alias1 ...string) *gorm.DB {
	var alias string
	if len(alias1) == 0 {
		alias = "default"
	} else {
		alias = alias1[0]
	}
	var err error

	// Check if there is a database config file
	dialect := GetAdapterForDb(alias)
	Db, err = dialect.GetDb(
		alias, false,
	)
	if err != nil {
		Trail(ERROR, "unable to connect to DB. %s", err)
		Db.Error = fmt.Errorf("unable to connect to DB. %s", err)
	}
	return Db
}

func GetAdapterForDb(alias1 ...string) IDbAdapter {
	var databaseConfig *DBSettings
	var alias string
	if len(alias1) == 0 {
		alias = "default"
	} else {
		alias = alias1[0]
	}
	if alias == "default" {
		databaseConfig = CurrentDatabaseSettings.Default
	}
	return NewDbAdapter(Db, databaseConfig.Type)
}