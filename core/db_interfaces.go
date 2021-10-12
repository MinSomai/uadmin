package core

import (
	"fmt"
	"gorm.io/gorm"
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
	GetDb(alias string, dryRun bool) (*gorm.DB, error)
	GetStringToExtractYearFromField(filterOptionField string) string
	GetStringToExtractMonthFromField(filterOptionField string) string
	Exact(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	IExact(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Contains(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	IContains(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	In(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Gt(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Gte(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Lt(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Lte(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	StartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	IStartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	EndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	IEndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Range(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Date(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Year(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Month(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Day(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Week(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	WeekDay(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Quarter(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Time(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Hour(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Minute(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Second(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	IsNull(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	Regex(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	IRegex(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder *SQLConditionBuilder)
	BuildDeleteString(table string, cond string, values ...interface{}) *DeleteRowStructure
	SetIsolationLevelForTests(db *gorm.DB)
	Close(db *gorm.DB)
	ClearTestDatabase()
	SetTimeZone(db *gorm.DB, timezone string)
	InitializeDatabaseForTests(databaseSettings *DBSettings)
	StartDBShell(databaseSettings *DBSettings) error
	GetLastError() error
}

var Db *gorm.DB

type UadminDatabase struct {
	Db      *gorm.DB
	Adapter IDbAdapter
}

func (uad *UadminDatabase) Close() {
	uad.Adapter.Close(uad.Db)
}

func (uad *UadminDatabase) ForcefullyClose() {
	db1, _ := uad.Db.DB()
	db1.Close()
}

var UadminTestDatabase *UadminDatabase

func NewUadminDatabase(alias1 ...string) *UadminDatabase {
	if CurrentConfig.InTests && UadminTestDatabase != nil {
		return UadminTestDatabase
	}
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
	if CurrentConfig.InTests && UadminTestDatabase != nil {
		return UadminTestDatabase
	}
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
	return GetDB(alias)
}

type DatabaseSettings struct {
	Default *DBSettings
	Slave   *DBSettings
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
	} else {
		databaseConfig = CurrentDatabaseSettings.Slave
	}
	return NewDbAdapter(Db, databaseConfig.Type)
}

type DbAdapterRegistry struct {
	dbTypeToAdapter map[string]func(db *gorm.DB) IDbAdapter
}

func (dar *DbAdapterRegistry) RegisterAdapter(dbType string, createAdapterHandler func(db *gorm.DB) IDbAdapter) {
	dar.dbTypeToAdapter[dbType] = createAdapterHandler
}

var GlobalDbAdapterRegistry *DbAdapterRegistry

func InitializeGlobalAdapterRegistry() {
	if GlobalDbAdapterRegistry == nil {
		GlobalDbAdapterRegistry = &DbAdapterRegistry{
			dbTypeToAdapter: make(map[string]func(db *gorm.DB) IDbAdapter),
		}
	}
}

func NewDbAdapter(db *gorm.DB, dbType string) IDbAdapter {
	adapter, ok := GlobalDbAdapterRegistry.dbTypeToAdapter[dbType]
	if !ok {
		panic("no adapter " + dbType + " has been registered")
	}
	return adapter(db)
}
