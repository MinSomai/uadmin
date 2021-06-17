package dialect

import (
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
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
