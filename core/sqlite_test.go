package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

// User !
type UserTest struct {
	gorm.Model
	Username     string
	FirstName    string
	LastName     string
	Password     string
	Email        string
	Active       bool
	Admin        bool
	Photo        string
	//Language     []Language `gorm:"many2many:user_languages" listExclude:"true"`
	LastLogin   *time.Time
	ExpiresOn   *time.Time
	OTPRequired bool
	OTPSeed     string
}

func GetDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("not able to open database: %s", "test.db"))
	}
	// Initialize system models
	modelList := []interface{}{
		UserTest{},
	}
	// Migrate schema
	for _, model := range modelList {
		db.AutoMigrate(model)
	}
	return db
}

func TestSqlite(t *testing.T) {
	db := GetDb()
	sql_dialect := NewDbAdapter(db, "sqlite")
	sql_dialect.Equals("admin", true)
	assert.Equal(t, sql_dialect.ToString(), "`admin` = ?")
	sql_dialect = NewDbAdapter(db, "sqlite")
	sql_dialect.GetLastInsertId()
	assert.Equal(t, sql_dialect.ToString(), "SELECT last_insert_rowid() AS lastid")
	sql_dialect = NewDbAdapter(db, "sqlite")
	quoted_field := sql_dialect.Quote("test")
	assert.Equal(t, quoted_field, "`test`")
	quoted_field = sql_dialect.LikeOperator()
	assert.Equal(t, quoted_field, " LIKE ")
	quoted_table_name := sql_dialect.QuoteTableName("TestModelA")
	assert.Equal(t, quoted_table_name, "`TestModelA`")
	quoted_table_name = sql_dialect.QuoteTableName("testmodela")
	assert.Equal(t, quoted_table_name, "`testmodela`")
}

func TestSqliteFunctional(t *testing.T) {
	//db := GetDb().Begin()
	//db = db.Exec("INSERT INTO users (`username`) VALUES (\"test\")")
	//last_ids := []int{}
	//sql_dialect := NewDbDialect(db, "sqlite")
	//sql_dialect.GetLastInsertId()
	//db = db.Raw(sql_dialect.ToString())
	//db = db.Pluck("lastid", &last_ids)
	//db.Commit()
	// assert.Equal(t, 1, len(last_ids))
}

func TestSqlite_uadmin_datetime_cast_date(t *testing.T) {
	time1 := "2005-07-29 09:56:00-05"
	dt := sqlite_uadmin_datetime_parse(time1, "UTC", "UTC")
	assert.Equal(t, dt.Year(), 2005)
	dt = sqlite_uadmin_datetime_parse("", "UTC", "UTC")
	assert.Nil(t, dt)
}