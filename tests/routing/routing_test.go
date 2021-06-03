package routing

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/tests"
	"os"
	"testing"
)

type ConcreteTestSuite struct {
	suite.Suite
	app *uadmin.App
}

func (suite *ConcreteTestSuite) SetupTest() {
	app, _ := tests.NewTestApp()
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
	assert.True(suite.T(), false)
}


// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMigrations(t *testing.T) {
	uadmin.ClearApp()
	suite.Run(t, new(ConcreteTestSuite))
}
