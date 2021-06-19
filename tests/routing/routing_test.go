package routing

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/interfaces"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type ConcreteTestSuite struct {
	suite.Suite
	app *uadmin.App
}

func (suite *ConcreteTestSuite) SetupTest() {
	app, _ := uadmin.NewTestApp()
	suite.app = app
	suite.app.BlueprintRegistry = interfaces.NewBlueprintRegistry()
	suite.app.BlueprintRegistry.Register(ConcreteBlueprint)
}

func (suite *ConcreteTestSuite) TearDownSuite() {
	err := os.Remove(suite.app.Config.D.Db.Default.Name)
	if err != nil {
		assert.Equal(suite.T(), true, false, fmt.Errorf("Couldnt remove db with name %s", suite.app.Config.D.Db.Default.Name))
	}
}

func (suite *ConcreteTestSuite) TestRouterInitialization() {
	// suite.app.Router = gin.Default()
	suite.app.InitializeRouter()
	req, _ := http.NewRequest("GET", "/user/visit", nil)
	uadmin.TestHTTPResponse(suite.T(), suite.app, req, func(w *httptest.ResponseRecorder) bool {
		return visited
	})
}

func (suite *ConcreteTestSuite) TestPingEndpoint() {
	// suite.app.Router = gin.Default()
	suite.app.InitializeRouter()
	req, _ := http.NewRequest("GET", "/ping", nil)
	uadmin.TestHTTPResponse(suite.T(), suite.app, req, func(w *httptest.ResponseRecorder) bool {
		return w.Body.String() == "{\"message\":\"pong\"}"
	})
	req1, _ := http.NewRequest("GET", "/static-inbuilt/uadmin/assets/moment.js", nil)
	uadmin.TestHTTPResponse(suite.T(), suite.app, req1, func(w *httptest.ResponseRecorder) bool {
		return w.Code == 200
	})
}


// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMigrations(t *testing.T) {
	uadmin.ClearApp()
	suite.Run(t, new(ConcreteTestSuite))
}
