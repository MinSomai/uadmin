package dialect

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommonDialect struct {
	Statement *gorm.Statement
	DbType    string
}

func NewCommonDialect(db *gorm.DB, db_type string) *CommonDialect {
	return &CommonDialect{
		DbType: db_type,
		Statement: &gorm.Statement{
			DB:      db,
			Context: context.Background(),
			Clauses: map[string]clause.Clause{},
		},
	}
}

func (d *CommonDialect) Equals(name interface{}, args ...interface{}) {
	query := d.Statement.Quote(name) + " = ?"
	clause.Expr{SQL: query, Vars: args}.Build(d.Statement)
}

func (d *CommonDialect) Quote(name interface{}) string {
	return d.Statement.Quote(name)
}

func (d *CommonDialect) LikeOperator() string {
	if d.DbType == "sqlite" {
		return " LIKE "
	}
	return " LIKE BINARY "
}
func (d *CommonDialect) ToString() string {
	return d.Statement.SQL.String()
}

func (d *CommonDialect) GetLastInsertId() {
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

func (d *CommonDialect) buildClauses(clause_interfaces []clause.Interface) {
	var buildNames []string
	for _, c := range clause_interfaces {
		buildNames = append(buildNames, c.Name())
		d.Statement.AddClause(c)
	}
	d.Statement.Build(buildNames...)
}

func (d *CommonDialect) QuoteTableName(tableName string) string {
	return d.Statement.Quote(tableName)
}
