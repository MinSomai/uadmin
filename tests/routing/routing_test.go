package routing

import (
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin/core"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ConcreteTestSuite struct {
	uadmin.TestSuite
}

func (suite *ConcreteTestSuite) SetupTest() {
	//app := uadmin.NewTestApp()
	//suite.app = app
	suite.TestSuite.SetupTest()
	suite.App.BlueprintRegistry = core.NewBlueprintRegistry()
	suite.App.BlueprintRegistry.Register(ConcreteBlueprint)
}

func (suite *ConcreteTestSuite) TearDownSuite() {
	uadmin.ClearTestApp()
}

func (suite *ConcreteTestSuite) TestRouterInitialization() {
	// suite.app.Router = gin.Default()
	routergroup := suite.App.Router.Group("/" + "user")
	ConcreteBlueprint.InitRouter(suite.App.Router, routergroup)
	req, _ := http.NewRequest("GET", "/user/visit", nil)
	uadmin.TestHTTPResponse(suite.T(), suite.App, req, func(w *httptest.ResponseRecorder) bool {
		return visited
	})
}

func (suite *ConcreteTestSuite) TestPingEndpoint() {
	// suite.app.Router = gin.Default()
	req, _ := http.NewRequest("GET", "/ping", nil)
	uadmin.TestHTTPResponse(suite.T(), suite.App, req, func(w *httptest.ResponseRecorder) bool {
		return w.Body.String() == "{\"message\":\"pong\"}\n"
	})
	req1, _ := http.NewRequest("GET", "/static-inbuilt/uadmin/assets/moment.js", nil)
	uadmin.TestHTTPResponse(suite.T(), suite.App, req1, func(w *httptest.ResponseRecorder) bool {
		return w.Code == 200
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRouting(t *testing.T) {
	uadmin.ClearApp()
	uadmin.RunTests(t, new(ConcreteTestSuite))
}
