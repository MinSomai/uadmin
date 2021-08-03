package object_query_builder

import (
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/interfaces"
	"gorm.io/gorm"
	"math"
	"strconv"
	"testing"
	"time"
)

type ObjectQueryBuilderTestSuite struct {
	uadmin.UadminTestSuite
	createdUser *usermodels.User
}

func (suite *ObjectQueryBuilderTestSuite) SetupTest() {
	suite.UadminTestSuite.SetupTest()
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	db := uadminDatabase.Db
	u := usermodels.User{Username: "dsadas", Email: "ffsdfsd@example.com"}
	db.Create(&u)
	suite.createdUser = &u
}

func (suite *ObjectQueryBuilderTestSuite) TestExact() {
	operator := interfaces.ExactGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "dsadas")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "dsadaS")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestIExact() {
	operator := interfaces.IExactGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "dsadas")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "Dsadas")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestContains() {
	operator := interfaces.ContainsGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "dsad")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "Dsa")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestIContains() {
	operator := interfaces.IContainsGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "dsad")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "Dsa")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestIn() {
	operator := interfaces.InGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["ID"], []uint{suite.createdUser.ID})
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestGt() {
	operator := interfaces.GtGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["ID"], -1)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestGte() {
	operator := interfaces.GteGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["ID"], suite.createdUser.ID)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestLt() {
	operator := interfaces.LtGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["ID"], suite.createdUser.ID + 100)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestLte() {
	operator := interfaces.LteGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["ID"], suite.createdUser.ID)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestStartsWith() {
	operator := interfaces.StartsWithGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "dsad")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "Dsa")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestIStartsWith() {
	operator := interfaces.IStartsWithGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "dsad")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "Dsa")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestEndsWith() {
	operator := interfaces.EndsWithGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "das")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "Das")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestIEndsWith() {
	operator := interfaces.IEndsWithGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "das")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"], "Das")
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestRange() {
	operator := interfaces.RangeGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["ID"], []uint{suite.createdUser.ID - 1, suite.createdUser.ID + 100})
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestDate() {
	operator := interfaces.DateGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		suite.createdUser.CreatedAt.Round(0),
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		suite.createdUser.CreatedAt.Round(0).Add(-10 * 3600 * 24 * time.Second),
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["UpdatedAt"],
		suite.createdUser.CreatedAt.Round(0).Add(-10 * 3600 * 24 * time.Second),
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["UpdatedAt"],
		suite.createdUser.CreatedAt.Round(0),
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestYear() {
	operator := interfaces.YearGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		suite.createdUser.CreatedAt.Year(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestMonth() {
	operator := interfaces.MonthGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		suite.createdUser.CreatedAt.Month(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestDay() {
	operator := interfaces.DayGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		suite.createdUser.CreatedAt.Day(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestWeek() {
	operator := interfaces.WeekGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	_, isoWeek := suite.createdUser.CreatedAt.ISOWeek()
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		isoWeek,
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestWeekDay() {
	operator := interfaces.WeekDayGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		suite.createdUser.CreatedAt.Weekday(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestQuarter() {
	operator := interfaces.QuarterGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		math.Ceil(float64(suite.createdUser.CreatedAt.Month() / 3)),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestTime() {
	operator := interfaces.TimeGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		suite.createdUser.CreatedAt.Format("15:04"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestHour() {
	operator := interfaces.HourGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		suite.createdUser.CreatedAt.Hour(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestRegex() {
	operator := interfaces.RegexGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"],
		"^DSAdas$",
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"],
		"^dsaDAS1111111111$",
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestIRegex() {
	operator := interfaces.IRegexGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"],
		"^dsadas$",
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	uadminDatabase = interfaces.NewUadminDatabase()
	gormOperatorContext = interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["Username"],
		"^dsaDAS$",
	)
	gormOperatorContext.Tx.First(&u)
	uadminDatabase.Close()
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestMinute() {
	operator := interfaces.MinuteGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		suite.createdUser.CreatedAt.Minute(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestSecond() {
	operator := interfaces.SecondGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		suite.createdUser.CreatedAt.Second(),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestIsNull() {
	operator := interfaces.IsNullGormOperator{}
	var u usermodels.User
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	gormOperatorContext := interfaces.NewGormOperatorContext(uadminDatabase.Db, &usermodels.User{})
	operator.Build(
		uadminDatabase.Adapter, gormOperatorContext, gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"],
		false,
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestFilterGormModel() {
	var u usermodels.User
	createdSecond := suite.createdUser.CreatedAt.Second()
	uadminDatabase := interfaces.NewUadminDatabase()
	statement := &gorm.Statement{DB: uadminDatabase.Db}
	statement.Parse(&usermodels.User{})
	interfaces.FilterGormModel(uadminDatabase.Adapter, uadminDatabase.Db, statement.Schema, []string{"CreatedAt__isnull=false", "CreatedAt__second=" + strconv.Itoa(createdSecond)}, &usermodels.User{})
	defer uadminDatabase.Close()
	uadminDatabase.Db.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestFilterGormModelWithDependencies() {
	var p usermodels.Permission
	uadminDatabase := interfaces.NewUadminDatabase()
	statement := &gorm.Statement{DB: uadminDatabase.Db}
	statement.Parse(&usermodels.Permission{})
	contentType := interfaces.ContentType{BlueprintName: "user11", ModelName: "user11"}
	uadminDatabase.Db.Create(&contentType)
	permission := usermodels.Permission{ContentType: contentType, PermissionBits: interfaces.RevertPermBit}
	uadminDatabase.Db.Create(&permission)
	interfaces.FilterGormModel(uadminDatabase.Adapter, uadminDatabase.Db, statement.Schema, []string{"CreatedAt__isnull=false", "ContentType__ID__exact=" + strconv.Itoa(int(contentType.ID))}, &usermodels.Permission{})
	defer uadminDatabase.Close()
	uadminDatabase.Db.First(&p)
	assert.Equal(suite.T(), p.ContentTypeID, contentType.ID)
}


// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestObjectQueryBuilder(t *testing.T) {
	uadmin.Run(t, new(ObjectQueryBuilderTestSuite))
}
