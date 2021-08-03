package uadmin

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/utils"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"testing"
	"time"
)

type UadminTestSuite struct {
	suite.Suite
	App *App
}

func (suite *UadminTestSuite) SetupTest() {
	app := NewFullAppForTests()
	suite.App = app
}

func (suite *UadminTestSuite) TearDownSuite() {
	ClearTestApp()
}

func ClearTestApp() {
	appForTests = nil
	appInstance = nil
}

func failOnPanic(t *testing.T) {
	r := recover()
	if r != nil {
		t.Errorf("test panicked: %v\n%s", r, debug.Stack())
		t.FailNow()
	}
}

func newSuiteInformation() *suite.SuiteInformation {
	testStats := make(map[string]*suite.TestInformation)

	return &suite.SuiteInformation{
		TestStats: testStats,
	}
}

var allTestsFilter = func(_, _ string) (bool, error) { return true, nil }

func methodFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^Test", name); !ok {
		return false, nil
	}
	return regexp.MatchString("", name)
}

func startStats(s *suite.SuiteInformation, testName string) {
	s.TestStats[testName] = &suite.TestInformation{
		TestName: testName,
		Start:    time.Now(),
	}
}

func endStats(s *suite.SuiteInformation, testName string, passed bool) {
	s.TestStats[testName].End = time.Now()
	s.TestStats[testName].Passed = passed
}

func Run(t *testing.T, currentsuite suite.TestingSuite) {
	defer failOnPanic(t)

	currentsuite.SetT(t)

	var suiteSetupDone bool

	var stats *suite.SuiteInformation
	if _, ok := currentsuite.(suite.WithStats); ok {
		stats = newSuiteInformation()
	}

	tests := []testing.InternalTest{}
	methodFinder := reflect.TypeOf(currentsuite)
	suiteName := methodFinder.Elem().Name()

	for i := 0; i < methodFinder.NumMethod(); i++ {
		method := methodFinder.Method(i)

		ok, err := methodFilter(method.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "testify: invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}

		if !ok {
			continue
		}

		if !suiteSetupDone {
			if stats != nil {
				stats.Start = time.Now()
			}

			if setupAllSuite, ok := currentsuite.(suite.SetupAllSuite); ok {
				setupAllSuite.SetupSuite()
			}

			suiteSetupDone = true
		}

		test := testing.InternalTest{
			Name: method.Name,
			F: func(t *testing.T) {
				parentT := currentsuite.T()
				currentsuite.SetT(t)
				defer failOnPanic(t)
				defer func() {
					if stats != nil {
						passed := !t.Failed()
						endStats(stats, method.Name, passed)
					}

					if afterTestSuite, ok := currentsuite.(suite.AfterTest); ok {
						afterTestSuite.AfterTest(suiteName, method.Name)
					}

					if tearDownTestSuite, ok := currentsuite.(suite.TearDownTestSuite); ok {
						tearDownTestSuite.TearDownTest()
					}

					currentsuite.SetT(parentT)
				}()

				if setupTestSuite, ok := currentsuite.(suite.SetupTestSuite); ok {
					setupTestSuite.SetupTest()
				}
				if beforeTestSuite, ok := currentsuite.(suite.BeforeTest); ok {
					beforeTestSuite.BeforeTest(methodFinder.Elem().Name(), method.Name)
				}

				if stats != nil {
					startStats(stats, method.Name)
				}
				appForTests.Config.InTests = true
				utils.SentEmailsDuringTests.ClearTestEmails()
				if appForTests.Config.D.Db.Default.Type == "sqlite" {
					upCommand := MigrateCommand{}
					upCommand.Proceed("up", make([]string, 0))
					method.Func.Call([]reflect.Value{reflect.ValueOf(currentsuite)})
					downCommand := MigrateCommand{}
					downCommand.Proceed("down", make([]string, 0))
				} else {
					uadminDatabase := interfaces.NewUadminDatabase()
					uadminDatabase.Db.Transaction(func(tx *gorm.DB) error {
						method.Func.Call([]reflect.Value{reflect.ValueOf(currentsuite)})
						// return nil will commit the whole transaction
						return fmt.Errorf("dont commit")
					})
					uadminDatabase.Close()
				}
			},
		}
		tests = append(tests, test)
	}
	if suiteSetupDone {
		defer func() {
			if tearDownAllSuite, ok := currentsuite.(suite.TearDownAllSuite); ok {
				tearDownAllSuite.TearDownSuite()
			}

			if suiteWithStats, measureStats := currentsuite.(suite.WithStats); measureStats {
				stats.End = time.Now()
				suiteWithStats.HandleStats(suiteName, stats)
			}
		}()
	}

	runTests(t, tests)
}

type runner interface {
	Run(name string, f func(t *testing.T)) bool
}

func runTests(t testing.TB, tests []testing.InternalTest) {
	if len(tests) == 0 {
		t.Log("warning: no tests to run")
		return
	}

	r, ok := t.(runner)
	if !ok { // backwards compatibility with Go 1.6 and below
		if !testing.RunTests(allTestsFilter, tests) {
			t.Fail()
		}
		return
	}

	for _, test := range tests {
		r.Run(test.Name, test.F)
	}
}

func NewTestApp() *App {
	a := App{}
	a.DashboardAdminPanel = admin.NewDashboardAdminPanel()
	admin.CurrentDashboardAdminPanel = a.DashboardAdminPanel
	a.Config = interfaces.NewConfig("configs/" + "test" + ".yaml")
	a.CommandRegistry = &CommandRegistry{
		Actions: make(map[string]interfaces.ICommand),
	}
	a.BlueprintRegistry = interfaces.NewBlueprintRegistry()
	a.Database = interfaces.NewDatabase(a.Config)
	a.Router = gin.Default()
	a.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	a.RegisterBaseCommands()
	interfaces.CurrentDatabaseSettings = &interfaces.DatabaseSettings{
		Default: a.Config.D.Db.Default,
	}
	StoreCurrentApp(&a)
	return &a
}

// Helper function to process a request and test its response
func TestHTTPResponse(t *testing.T, app *App, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	app.Router.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

var appForTests *App

func NewFullAppForTests() *App {
	if appForTests != nil {
		return appForTests
	}
	a := NewApp("test")
	appForTests = a
	// appForTests.DashboardAdminPanel.RegisterHttpHandlers(a.Router)
	StoreCurrentApp(a)
	return a
}
//type UadminTestSuite struct {
//	suite.Suite
//}
//
//func (suite *UadminTestSuite) SetupTest() {
//	db := dialect.GetDB()
//	db = db.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")
//	if db.Error != nil {
//		assert.Equal(suite.T(), true, false, "Couldnt setup isolation level for db")
//	}
//	db = db.Exec("BEGIN")
//	if db.Error != nil {
//		assert.Equal(suite.T(), true, false, "Couldnt start transaction")
//	}
//}
//
//func (suite *UadminTestSuite) TearDownSuite() {
//	db := dialect.GetDB()
//	db = db.Exec("ROLLBACK")
//	if db.Error != nil {
//		assert.Equal(suite.T(), true, false, "Couldnt rollback transaction")
//	}
//}
