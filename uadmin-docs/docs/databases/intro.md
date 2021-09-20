---
sidebar_position: 1
---

# Database support

Right now uadmin supports only sqlite, postgres databases, but it's easy to provide adapters for another databases, we just need to write implementation of the interface
```go
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
}
```
Please also create ci-cd job for new database type you are adding to the framework.
And don't forget to handle this database type in the core.NewDbAdapter function.  
You can get instance of the uadminDatabase using function:
```go
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
DBInstance := NewUadminDatabase()
defer DBInstance.Close()
// do whatever you want with database
```
