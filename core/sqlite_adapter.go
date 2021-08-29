package core

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
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

func (d *SqliteDialect) Quote(name interface{}) string {
	return d.Statement.Quote(name)
}

func (d *SqliteDialect) LikeOperator() string {
	if d.DbType == "sqlite" {
		return " LIKE "
	}
	return " LIKE BINARY "
}
func (d *SqliteDialect) ToString() string {
	return d.Statement.SQL.String()
}

func (d *SqliteDialect) GetLastInsertId() {
	var last_insert_id_func string
	if d.DbType == "sqlite" {
		last_insert_id_func = "last_insert_rowid()"
	} else {
		last_insert_id_func = "LAST_INSERT_ID()"
	}
	clause_interfaces := []clause.Interface{clause.Select{
		Expression: clause.Expr{
			SQL: last_insert_id_func + " AS lastid",
		},
	},
	}
	d.buildClauses(clause_interfaces)
}

func (d *SqliteDialect) buildClauses(clause_interfaces []clause.Interface) {
	var buildNames []string
	for _, c := range clause_interfaces {
		buildNames = append(buildNames, c.Name())
		d.Statement.AddClause(c)
	}
	d.Statement.Build(buildNames...)
}

func (d *SqliteDialect) QuoteTableName(tableName string) string {
	return d.Statement.Quote(tableName)
}

func (d *SqliteDialect) GetStringToExtractYearFromField(filterOptionField string) string {
	return fmt.Sprintf("strftime(\"%%Y\", %s)", filterOptionField)
}

func (d *SqliteDialect) GetStringToExtractMonthFromField(filterOptionField string) string {
	return fmt.Sprintf("strftime(\"%%m\", %s)", filterOptionField)
}

// @todo analyze
func (d *SqliteDialect) Delete(db *gorm.DB, model reflect.Value, query interface{}, args ...interface{}) *gorm.DB {
	// if Database.Type == "mysql" {
	// 	db := GetDB()
	//
	// 	if log {
	// 		db.Model(model.Interface()).Where(q, args...).Scan(modelArray.Interface())
	// 	}
	//
	// 	db = db.Where(q, args...).Delete(model)
	// 	if db.Error != nil {
	// 		ReturnJSON(w, r, map[string]interface{}{
	// 			"status":  "error",
	// 			"err_msg": "Unable to execute DELETE SQL. " + db.Error.Error(),
	// 		})
	// 		return
	// 	}
	// 	rowsCount = db.RowsAffected
	// 	if log {
	// 		for i := 0; i < modelArray.Elem().Len(); i++ {
	// 			createAPIDeleteLog(modelName, modelArray.Elem().Index(i).Interface(), &s.User, r)
	// 		}
	// 	}
	//
	// } else if Database.Type == "postgres" {
	// 	db := GetDB()
	//
	// 	if log {
	// 		db.Model(model.Interface()).Where(q, args...).Scan(modelArray.Interface())
	// 	}
	//
	// 	db = db.Where(q, args...).Delete(model.Interface())
	// 	if db.Error != nil {
	// 		ReturnJSON(w, r, map[string]interface{}{
	// 			"status":  "error",
	// 			"err_msg": "Unable to execute DELETE SQL. " + db.Error.Error(),
	// 		})
	// 		return
	// 	}
	// 	rowsCount = db.RowsAffected
	// 	if log {
	// 		for i := 0; i < modelArray.Elem().Len(); i++ {
	// 			createAPIDeleteLog(modelName, modelArray.Elem().Index(i).Interface(), &s.User, r)
	// 		}
	// 	}
	//
	// } else
	db = db.Exec("PRAGMA case_sensitive_like=ON;")
	db = db.Where(query, args...).Delete(model)
	db = db.Exec("PRAGMA case_sensitive_like=OFF;")
	return db
}

// @todo analyze
func (d *SqliteDialect) ReadRows(db *gorm.DB, customSchema bool, SQL string, m interface{}, args ...interface{}) (*sql.Rows, error) {
	// if Database.Type == "mysql" {
	// 	db := GetDB()
	// 	if !customSchema {
	// 		db.Raw(SQL, args...).Scan(m)
	// 	} else {
	// 		rows, err = db.Raw(SQL, args...).Rows()
	// 		if err != nil {
	// 			w.WriteHeader(500)
	// 			ReturnJSON(w, r, map[string]interface{}{
	// 				"status":  "error",
	// 				"err_msg": "Unable to execute SQL. " + err.Error(),
	// 			})
	// 			Trail(ERROR, "SQL: %v\nARGS: %v", SQL, args)
	// 			return
	// 		}
	// 		m = parseCustomDBSchema(rows)
	// 	}
	// 	if a, ok := m.([]map[string]interface{}); ok {
	// 		rowsCount = int64(len(a))
	// 	} else {
	// 		rowsCount = int64(reflect.ValueOf(m).Elem().Len())
	// 	}
	// } else if Database.Type == "postgresql" {
	// 	db := GetDB()
	// 	if !customSchema {
	// 		db.Raw(SQL, args...).Scan(m)
	// 	} else {
	// 		rows, err = db.Raw(SQL, args...).Rows()
	// 		if err != nil {
	// 			w.WriteHeader(500)
	// 			ReturnJSON(w, r, map[string]interface{}{
	// 				"status":  "error",
	// 				"err_msg": "Unable to execute SQL. " + err.Error(),
	// 			})
	// 			Trail(ERROR, "SQL: %v\nARGS: %v", SQL, args)
	// 			return
	// 		}
	// 		m = parseCustomDBSchema(rows)
	// 	}
	// 	if a, ok := m.([]map[string]interface{}); ok {
	// 		rowsCount = int64(len(a))
	// 	} else {
	// 		rowsCount = int64(reflect.ValueOf(m).Elem().Len())
	// 	}
	db.Exec("PRAGMA case_sensitive_like=ON;")
	var rows *sql.Rows
	var err error
	if !customSchema {
		db.Raw(SQL, args...).Scan(m)
	} else {
		rows, err = db.Raw(SQL, args...).Rows()
	}
	db.Exec("PRAGMA case_sensitive_like=OFF;")
	return rows, err
}

// @todo analyze
func (d *SqliteDialect) GetSqlDialectStrings() map[string]string {
	// var sqlDialect = map[string]map[string]string{
	// 	"mysql": {
	// 		"createM2MTable": "CREATE TABLE `{TABLE1}_{TABLE2}` (`table1_id` int(10) unsigned NOT NULL, `table2_id` int(10) unsigned NOT NULL, PRIMARY KEY (`table1_id`,`table2_id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;",
	// 		"selectM2M":      "SELECT `table2_id` FROM `{TABLE1}_{TABLE2}` WHERE table1_id={TABLE1_ID};",
	// 		"deleteM2M":      "DELETE FROM `{TABLE1}_{TABLE2}` WHERE `table1_id`={TABLE1_ID};",
	// 		"insertM2M":      "INSERT INTO `{TABLE1}_{TABLE2}` VALUES ({TABLE1_ID}, {TABLE2_ID});",
	// 	},
	// 	"postgresql": {
	// 		"createM2MTable": "CREATE TABLE \"{TABLE1}_{TABLE2}\" (\"table1_id\" int(10) unsigned NOT NULL, \"table2_id\" int(10) unsigned NOT NULL, PRIMARY KEY (\"table1_id\",\"table2_id\"));",
	// 		"selectM2M":      "SELECT \"table2_id\" FROM \"{TABLE1}_{TABLE2}\" WHERE table1_id={TABLE1_ID};",
	// 		"deleteM2M":      "DELETE FROM \"{TABLE1}_{TABLE2}\" WHERE \"table1_id\"={TABLE1_ID};",
	// 		"insertM2M":      "INSERT INTO \"{TABLE1}_{TABLE2}\" VALUES ({TABLE1_ID}, {TABLE2_ID});",
	// 	},
	// 	"sqlite": {
	// 		//"createM2MTable": "CREATE TABLE `{TABLE1}_{TABLE2}` (`{TABLE1}_id`	INTEGER NOT NULL,`{TABLE2}_id` INTEGER NOT NULL, PRIMARY KEY(`{TABLE1}_id`,`{TABLE2}_id`));",
	// 		"createM2MTable": "CREATE TABLE `{TABLE1}_{TABLE2}` (`table1_id`	INTEGER NOT NULL,`table2_id` INTEGER NOT NULL, PRIMARY KEY(`table1_id`,`table2_id`));",
	// 		"selectM2M": "SELECT `table2_id` FROM `{TABLE1}_{TABLE2}` WHERE table1_id={TABLE1_ID};",
	// 		"deleteM2M": "DELETE FROM `{TABLE1}_{TABLE2}` WHERE `table1_id`={TABLE1_ID};",
	// 		"insertM2M": "INSERT INTO `{TABLE1}_{TABLE2}` VALUES ({TABLE1_ID}, {TABLE2_ID});",
	// 	},
	// }
	return map[string]string{
		//"createM2MTable": "CREATE TABLE `{TABLE1}_{TABLE2}` (`{TABLE1}_id`	INTEGER NOT NULL,`{TABLE2}_id` INTEGER NOT NULL, PRIMARY KEY(`{TABLE1}_id`,`{TABLE2}_id`));",
		"createM2MTable": "CREATE TABLE `{TABLE1}_{TABLE2}` (`table1_id`	INTEGER NOT NULL,`table2_id` INTEGER NOT NULL, PRIMARY KEY(`table1_id`,`table2_id`));",
		"selectM2M": "SELECT `table2_id` FROM `{TABLE1}_{TABLE2}` WHERE table1_id={TABLE1_ID};",
		"deleteM2M": "DELETE FROM `{TABLE1}_{TABLE2}` WHERE `table1_id`={TABLE1_ID};",
		"insertM2M": "INSERT INTO `{TABLE1}_{TABLE2}` VALUES ({TABLE1_ID}, {TABLE2_ID});",
	}
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
	var aliasDatabaseSettings *DBSettings
	if alias == "default" {
		aliasDatabaseSettings = CurrentDatabaseSettings.Default
	}
	var db *gorm.DB
	var err error
	if CurrentConfig.D.Uadmin.DebugTests {
		db, err = gorm.Open(sqlite.Dialector{DriverName: "UadminSqliteDriver", DSN: aliasDatabaseSettings.Name}, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: true,
			DryRun: dryRun,
		})
	} else {
		db, err = gorm.Open(sqlite.Dialector{DriverName: "UadminSqliteDriver", DSN: aliasDatabaseSettings.Name}, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: true,
			DryRun: dryRun,
		})
	}
	//var db *gorm.DB
	//var err error
	//
	//if d.DbType == "sqlite" {
	//	if name == "" {
	//		name = "uadmin.db"
	//	}
	//
	//}
	return db, err
}

// @todo analyze
func (d *SqliteDialect) Transaction(handler func()) error {
	return nil
}

// @todo analyze
func (d *SqliteDialect) CreateDb() error {
	// if Database.Type == "mysql" {
	// 	credential := Database.User
	//
	// 	if Database.Password != "" {
	// 		credential = fmt.Sprintf("%s:%s", Database.User, Database.Password)
	// 	}
	//
	// 	dsn := fmt.Sprintf("%s@(%s:%d)/?charset=utf8&parseTime=True&loc=Local",
	// 		credential,
	// 		Database.Host,
	// 		Database.Port,
	// 	)
	// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 		Logger: logger.Default.LogMode(logger.Info),
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	Trail(INFO, "Database doens't exist, creating a new database")
	// 	db = db.Exec("CREATE SCHEMA `" + Database.Name + "` DEFAULT CHARACTER SET utf8 COLLATE utf8_bin")
	//
	// 	if db.Error != nil {
	// 		return fmt.Errorf(db.Error.Error())
	// 	}
	//
	// 	return nil
	// } else if Database.Type == "postgresql" {
	// 	// credential := Database.User
	//
	// 	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=disable TimeZone=UTC",
	// 		Database.Host,
	// 		Database.User,
	// 		Database.Password,
	// 		Database.Port,
	// 	)
	// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 		Logger: logger.Default.LogMode(logger.Info),
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	Trail(INFO, "Database doens't exist, creating a new database")
	// 	db = db.Exec("CREATE DATABASE \"" + Database.Name + "\";")
	//
	// 	if db.Error != nil {
	// 		return fmt.Errorf(db.Error.Error())
	// 	}
	//
	// 	return nil
	//
	// }
	var err error
	return err
}

func (d *SqliteDialect) Exact(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s = ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) IExact(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s LIKE ? ESCAPE '\\' ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Contains(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s LIKE %%?%% ESCAPE '\\' ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) IContains(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s LIKE %%?%% ESCAPE '\\' ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) In(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s IN ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Gt(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s > ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Gte(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s >= ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Lt(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s < ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Lte(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s <= ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) StartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s LIKE ?%% ESCAPE '\\' ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) IStartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s LIKE ?%% ESCAPE '\\' ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) EndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s LIKE %%? ESCAPE '\\' ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) IEndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s LIKE %%? ESCAPE '\\' ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Date(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" uadmin_datetime_cast_date(%s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Year(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	year := value.(int)
	startOfTheYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfTheYear := time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s BETWEEN ? AND ? ", operatorContext.TableName, field.DBName), startOfTheYear, endOfTheYear)
}

func (d *SqliteDialect) Month(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" uadmin_datetime_extract('month', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Day(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" uadmin_datetime_extract('day', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Week(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" uadmin_datetime_extract('week', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) WeekDay(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" uadmin_datetime_extract('week_day', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Quarter(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" uadmin_datetime_extract('quarter', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Hour(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" uadmin_datetime_extract('hour', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Minute(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" uadmin_datetime_extract('minute', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Second(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" uadmin_datetime_extract('second', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Regex(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s REGEX ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) IRegex(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s REGEX '(?i)' || ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) Time(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" uadmin_datetime_cast_time(%s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName), value)
}

func (d *SqliteDialect) IsNull(operatorContext *GormOperatorContext, field *Field, value interface{}) {
	isTruthyValue := IsTruthyValue(value)
	isNull := " IS NULL "
	if !isTruthyValue {
		isNull = " IS NOT NULL "
	}
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s %s ", operatorContext.TableName, field.DBName, isNull))
}

func (d *SqliteDialect) Range(operatorContext *GormOperatorContext, field *Field, value interface{}) {
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
	operatorContext.Tx = operatorContext.Tx.Where(fmt.Sprintf(" %s.%s BETWEEN ? AND ? ", operatorContext.TableName, field.DBName), f, second)
}

func (d *SqliteDialect) BuildDeleteString(table string, cond string, values ...interface{}) *DeleteRowStructure {
	deleteRowStructure := &DeleteRowStructure{SQL: fmt.Sprintf("DELETE FROM %s WHERE %s", table, cond), Values: values}
	return deleteRowStructure
}

func sqlite_uadmin_datetime_parse(dt string, tz_name string, conn_tzname string) *time.Time {
	if dt == "" {
		return nil
	}
	ret, _ := time.Parse("2006-01-_2 15:04:00-05", dt)
	return &ret
}

func sqlite_uadmin_datetime_cast_date(dt string, tz_name string, conn_tzname string) string {
	dtTmp := sqlite_uadmin_datetime_parse(dt, tz_name, conn_tzname)
	if dtTmp == nil {
		return ""
	}
	res := dtTmp.Round(0).Format(time.RFC3339)
	return res
}

func sqlite_uadmin_datetime_cast_time(dt string, tz_name string, conn_tzname string) string {
	dtTmp := sqlite_uadmin_datetime_parse(dt, tz_name, conn_tzname)
	if dtTmp == nil {
		return ""
	}
	res := dtTmp.Format("15:04")
	return res
}

//func sqlite_uadmin_regex(re_pattern string, re_string string) bool {
//	regex := regexp.MustCompile(re_pattern)
//	return regex.Find([]byte(re_string)) != nil
//}

func sqlite_uadmin_datetime_extract(extract string, dt string, tz_name string, conn_tzname string) string {
	dtTmp := sqlite_uadmin_datetime_parse(dt, tz_name, conn_tzname)
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
			if err = conn.RegisterFunc("uadmin_datetime_cast_date", sqlite_uadmin_datetime_cast_date, true); err != nil {
				return err
			}
			if err = conn.RegisterFunc("uadmin_datetime_extract", sqlite_uadmin_datetime_extract, true); err != nil {
				return err
			}
			if err = conn.RegisterFunc("uadmin_datetime_time", sqlite_uadmin_datetime_cast_time, true); err != nil {
				return err
			}
			for operator := range ProjectGormOperatorRegistry.GetAll() {
				err = operator.RegisterDbHandlers(conn)
				if err != nil {
					return err
				}
			}
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