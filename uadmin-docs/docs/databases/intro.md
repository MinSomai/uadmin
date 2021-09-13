---
sidebar_position: 1
---

# Database support

Right now uadmin supports only sqlite database, but it's easy to provide adapters for another databases, we just need to write implementation of the interface
```go
type IDbAdapter interface {
	Equals(name interface{}, args ...interface{})
	GetDb(alias string, dryRun bool) (*gorm.DB, error)
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
```
and don't forget to handle this database type in the core.NewDbAdapter function.  
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
