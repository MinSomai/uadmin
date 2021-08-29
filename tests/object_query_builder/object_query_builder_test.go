package object_query_builder

import (
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/core"
	"gorm.io/gorm"
	"math"
	"strconv"
	"testing"
	"time"
)

type ObjectQueryBuilderTestSuite struct {
	uadmin.UadminTestSuite
	createdUser *core.User
}

func (suite *ObjectQueryBuilderTestSuite) SetupTest() {
	suite.UadminTestSuite.SetupTest()
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	db := uadminDatabase.Db
	u := core.User{Username: "dsadas", Email: "ffsdfsd@example.com"}
	db.Create(&u)
	suite.createdUser = &u
}

func (suite *ObjectQueryBuilderTestSuite) TestExact() {
	operator := core.ExactGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsadas")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsadaS")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestIExact() {
	operator := core.IExactGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsadas")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Dsadas")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestContains() {
	operator := core.ContainsGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsad")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Dsa")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestIContains() {
	operator := core.IContainsGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsad")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Dsa")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestIn() {
	operator := core.InGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["ID"]}, []uint{suite.createdUser.ID})
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestGt() {
	operator := core.GtGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["ID"]}, -1)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestGte() {
	operator := core.GteGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["ID"]}, suite.createdUser.ID)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestLt() {
	operator := core.LtGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["ID"]}, suite.createdUser.ID + 100)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestLte() {
	operator := core.LteGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["ID"]}, suite.createdUser.ID)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestStartsWith() {
	operator := core.StartsWithGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsad")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Dsa")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestIStartsWith() {
	operator := core.IStartsWithGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsad")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Dsa")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestEndsWith() {
	operator := core.EndsWithGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "das")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Das")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestIEndsWith() {
	operator := core.IEndsWithGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "das")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Das")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestRange() {
	operator := core.RangeGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["ID"]}, []uint{suite.createdUser.ID - 1, suite.createdUser.ID + 100})
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestDate() {
	operator := core.DateGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},	suite.createdUser.CreatedAt.Round(0),
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Round(0).Add(-10 * 3600 * 24 * time.Second),
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["UpdatedAt"]},
		suite.createdUser.CreatedAt.Round(0).Add(-10 * 3600 * 24 * time.Second),
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["UpdatedAt"]},
		suite.createdUser.CreatedAt.Round(0),
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestYear() {
	operator := core.YearGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Year(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestMonth() {
	operator := core.MonthGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Month(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestDay() {
	operator := core.DayGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Day(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestWeek() {
	operator := core.WeekGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	_, isoWeek := suite.createdUser.CreatedAt.ISOWeek()
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		isoWeek,
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestWeekDay() {
	operator := core.WeekDayGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Weekday(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestQuarter() {
	operator := core.QuarterGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		math.Ceil(float64(suite.createdUser.CreatedAt.Month() / 3)),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestTime() {
	operator := core.TimeGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Format("15:04"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestHour() {
	operator := core.HourGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Hour(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestRegex() {
	operator := core.RegexGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]},
		"^DSAdas$",
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]},
		"^dsaDAS1111111111$",
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestIRegex() {
	operator := core.IRegexGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]},
		"^dsadas$",
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = core.NewUadminDatabase()
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]},
		"^dsaDAS$",
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestMinute() {
	operator := core.MinuteGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Minute(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestSecond() {
	operator := core.SecondGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Second(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestIsNull() {
	operator := core.IsNullGormOperator{}
	var u core.User
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(uadminDatabase.Db), &core.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		false,
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestFilterGormModel() {
	var u core.User
	createdSecond := suite.createdUser.CreatedAt.Second()
	uadminDatabase := core.NewUadminDatabase()
	statement := &gorm.Statement{DB: uadminDatabase.Db}
	statement.Parse(&core.User{})
	core.FilterGormModel(uadminDatabase.Adapter, core.NewGormPersistenceStorage(uadminDatabase.Db), statement.Schema, []string{"CreatedAt__isnull=false", "CreatedAt__second=" + strconv.Itoa(createdSecond)}, &core.User{})
	defer uadminDatabase.Close()
	uadminDatabase.Db.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestFilterGormModelWithDependencies() {
	var p core.Permission
	uadminDatabase := core.NewUadminDatabase()
	statement := &gorm.Statement{DB: uadminDatabase.Db}
	statement.Parse(&core.Permission{})
	contentType := core.ContentType{BlueprintName: "user11", ModelName: "user11"}
	uadminDatabase.Db.Create(&contentType)
	permission := core.Permission{ContentType: contentType, PermissionBits: core.RevertPermBit}
	uadminDatabase.Db.Create(&permission)
	core.FilterGormModel(uadminDatabase.Adapter, core.NewGormPersistenceStorage(uadminDatabase.Db), statement.Schema, []string{"CreatedAt__isnull=false", "ContentType__ID__exact=" + strconv.Itoa(int(contentType.ID))}, &core.Permission{})
	defer uadminDatabase.Close()
	uadminDatabase.Db.First(&p)
	assert.Equal(suite.T(), p.ContentTypeID, contentType.ID)
}


// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestObjectQueryBuilder(t *testing.T) {
	uadmin.Run(t, new(ObjectQueryBuilderTestSuite))
}
