package core

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"reflect"
	"strconv"
	"time"
)

type SqliteDialect struct {
	Statement *gorm.Statement
	DbType    string
}

func NewDbAdapter(db *gorm.DB, dbType string) IDbAdapter {
	return &SqliteDialect{
		DbType: dbType,
		Statement: &gorm.Statement{
			DB:      db,
			Context: context.Background(),
			Clauses: map[string]clause.Clause{},
		},
	}
}

func (d *SqliteDialect) Equals(name interface{}, args ...interface{}) {
	query := d.Statement.Quote(name) + " = ?"
	clause.Expr{SQL: query, Vars: args}.Build(d.Statement)
}

func (d *SqliteDialect) GetStringToExtractYearFromField(filterOptionField string) string {
	return fmt.Sprintf("strftime(\"%%Y\", %s)", filterOptionField)
}

func (d *SqliteDialect) GetStringToExtractMonthFromField(filterOptionField string) string {
	return fmt.Sprintf("strftime(\"%%m\", %s)", filterOptionField)
}

// @todo analyze
func (d *SqliteDialect) GetDb(alias string, dryRun bool) (*gorm.DB, error) {
	// } else if strings.ToLower(Database.Type) == "postgresql" {
	// 	if Database.Host == "" || Database.Host == "localhost" {
	// 		Database.Host = "127.0.0.1"
	// 	}
	// 	if Database.Port == 0 {
	// 		Database.Port = 5432
	// 	}
	//
	// 	if Database.User == "" {
	// 		Database.User = "root"
	// 	}
	//
	// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
	// 		Database.Host,
	// 		Database.User,
	// 		Database.Password,
	// 		Database.Name,
	// 		Database.Port,
	// 	)
	// 	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 		Logger: logger.Default.LogMode(logger.Info),
	// 	})
	//
	// 	// Check if the error is DB doesn't exist and create it
	// 	if err != nil && strings.Contains(err.Error(), "does not exist") {
	// 		err = createDB()
	//
	// 		if err == nil {
	// 			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 				Logger: logger.Default.LogMode(logger.Info),
	// 			})
	// 		}
	// 	}
	// } else if strings.ToLower(Database.Type) == "mysql" {
	// 	if Database.Host == "" || Database.Host == "localhost" {
	// 		Database.Host = "127.0.0.1"
	// 	}
	// 	if Database.Port == 0 {
	// 		Database.Port = 3306
	// 	}
	//
	// 	if Database.User == "" {
	// 		Database.User = "root"
	// 	}
	//
	// 	credential := Database.User
	//
	// 	if Database.Password != "" {
	// 		credential = fmt.Sprintf("%s:%s", Database.User, Database.Password)
	// 	}
	// 	dsn := fmt.Sprintf("%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
	// 		credential,
	// 		Database.Host,
	// 		Database.Port,
	// 		Database.Name,
	// 	)
	// 	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 		Logger: logger.Default.LogMode(logger.Info),
	// 	})
	//
	// 	// Check if the error is DB doesn't exist and create it
	// 	if err != nil && err.Error() == "Error 1049: Unknown database '"+Database.Name+"'" {
	// 		err = createDB()
	//
	// 		if err == nil {
	// 			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 				Logger: logger.Default.LogMode(logger.Info),
	// 			})
	// 		}
	// 	}
	// }
	//var (
	//	postgresDsnTemplate, _ = template.New("postgresdsn").Parse("host={{.Host}} user={{.User}} password={{.Password}} dbname={{.Name}} port=5432 sslmode=disable TimeZone=UTC")
	//)
	var aliasDatabaseSettings *DBSettings
	if alias == "default" {
		aliasDatabaseSettings = CurrentDatabaseSettings.Default
	}
	var db *gorm.DB
	var err error
	if CurrentConfig.D.Uadmin.DebugTests {
		db, err = gorm.Open(sqlite.Dialector{DriverName: "UadminSqliteDriver", DSN: aliasDatabaseSettings.Name}, &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			DryRun:                                   dryRun,
		})
	} else {
		db, err = gorm.Open(sqlite.Dialector{DriverName: "UadminSqliteDriver", DSN: aliasDatabaseSettings.Name}, &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			DryRun:                                   dryRun,
		})
	}
	return db, err
}

func (d *SqliteDialect) Exact(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s = ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) IExact(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s LIKE ? ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Contains(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s LIKE '%%' || ? || '%%' ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) IContains(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s LIKE '%%' || ? || '%%' ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) In(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s IN ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Gt(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s > ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Gte(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s >= ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Lt(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s < ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Lte(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s <= ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) StartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s LIKE ? || '%%' ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) IStartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s LIKE ? || '%%' ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) EndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s LIKE '%%' || ? ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) IEndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s LIKE '%%' || ? ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Date(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" uadmin_datetime_cast_date(%s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Year(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s BETWEEN ? AND ? ", operatorContext.TableName, field.DBName)
	year := value.(int)
	startOfTheYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfTheYear := time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, startOfTheYear, endOfTheYear)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, startOfTheYear, endOfTheYear)
	}
}

func (d *SqliteDialect) Month(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" uadmin_datetime_extract('month', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Day(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" uadmin_datetime_extract('day', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Week(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" uadmin_datetime_extract('week', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) WeekDay(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" uadmin_datetime_extract('week_day', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Quarter(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" uadmin_datetime_extract('quarter', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Hour(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" uadmin_datetime_extract('hour', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Minute(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" uadmin_datetime_extract('minute', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Second(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" uadmin_datetime_extract('second', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Regex(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s REGEX ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) IRegex(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" %s.%s REGEX '(?i)' || ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) Time(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	query := fmt.Sprintf(" uadmin_datetime_cast_time(%s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, args)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, args)
	}
}

func (d *SqliteDialect) IsNull(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	isTruthyValue := IsTruthyValue(value)
	isNull := " IS NULL "
	if !isTruthyValue {
		isNull = " IS NOT NULL "
	}
	query := fmt.Sprintf(" %s.%s %s ", operatorContext.TableName, field.DBName, isNull)
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query)
	}
}

func (d *SqliteDialect) Range(operatorContext *GormOperatorContext, field *Field, value interface{}, forSearching bool) {
	s := reflect.ValueOf(value)
	var f interface{}
	var second interface{}
	for i := 0; i < s.Len(); i++ {
		if i == 0 {
			f = s.Index(i).Interface()
		} else if i == 1 {
			second = s.Index(i).Interface()
			break
		}
	}
	query := fmt.Sprintf(" %s.%s BETWEEN ? AND ? ", operatorContext.TableName, field.DBName)
	if forSearching {
		operatorContext.Tx = operatorContext.Tx.Or(query, f, second)
	} else {
		operatorContext.Tx = operatorContext.Tx.Where(query, f, second)
	}
}

func (d *SqliteDialect) BuildDeleteString(table string, cond string, values ...interface{}) *DeleteRowStructure {
	deleteRowStructure := &DeleteRowStructure{SQL: fmt.Sprintf("DELETE FROM %s WHERE %s", table, cond), Values: values}
	return deleteRowStructure
}

func sqliteUadminDatetimeParse(dt string, tzName string, connTzname string) *time.Time {
	if dt == "" {
		return nil
	}
	ret, _ := time.Parse("2006-01-_2 15:04:00-05", dt)
	return &ret
}

func sqliteUadminDatetimeCastDate(dt string, tzName string, connTzname string) string {
	dtTmp := sqliteUadminDatetimeParse(dt, tzName, connTzname)
	if dtTmp == nil {
		return ""
	}
	res := dtTmp.Round(0).Format(time.RFC3339)
	return res
}

func sqliteUadminDatetimeCastTime(dt string, tzName string, connTzname string) string {
	dtTmp := sqliteUadminDatetimeParse(dt, tzName, connTzname)
	if dtTmp == nil {
		return ""
	}
	res := dtTmp.Format("15:04")
	return res
}

func sqliteUadminDatetimeExtract(extract string, dt string, tzName string, connTzname string) string {
	dtTmp := sqliteUadminDatetimeParse(dt, tzName, connTzname)
	if dtTmp == nil {
		return ""
	}
	if extract == "month" {
		return strconv.Itoa(int(dtTmp.Month()))
	}
	if extract == "quarter" {
		return strconv.Itoa(int(math.Ceil(float64(dtTmp.Month() / 3))))
	}
	if extract == "day" {
		return strconv.Itoa(dtTmp.Day())
	}
	if extract == "hour" {
		return strconv.Itoa(dtTmp.Hour())
	}
	if extract == "minute" {
		return strconv.Itoa(dtTmp.Minute())
	}
	if extract == "second" {
		return strconv.Itoa(dtTmp.Second())
	}
	if extract == "week" {
		_, isoWeek := dtTmp.ISOWeek()
		return strconv.Itoa(isoWeek)
	}
	if extract == "week_day" {
		return strconv.Itoa(int(dtTmp.Weekday()))
	}
	return ""
}

func init() {
	sql.Register("UadminSqliteDriver", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			var err error
			if err = conn.RegisterFunc("uadmin_datetime_cast_date", sqliteUadminDatetimeCastDate, true); err != nil {
				return err
			}
			if err = conn.RegisterFunc("uadmin_datetime_extract", sqliteUadminDatetimeExtract, true); err != nil {
				return err
			}
			if err = conn.RegisterFunc("uadmin_datetime_time", sqliteUadminDatetimeCastTime, true); err != nil {
				return err
			}
			for operator := range ProjectGormOperatorRegistry.GetAll() {
				err = operator.RegisterDbHandlers(conn)
				if err != nil {
					return err
				}
			}
			//func sqlite_uadmin_regex(re_pattern string, re_string string) bool {
			//	regex := regexp.MustCompile(re_pattern)
			//	return regex.Find([]byte(re_string)) != nil
			//}
			//if err := conn.RegisterFunc("uadmin_regex", sqlite_uadmin_regex, true); err != nil {
			//	return err
			//}
			//if err := conn.RegisterFunc("uadmin_datetime_cast_year", sqlite_uadmin_datetime_cast_year, true); err != nil {
			//	return err
			//}
			return nil
		},
	})
}
