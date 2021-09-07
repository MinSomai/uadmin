package migrations

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin/core"
	"os"
	"strings"
	"testing"
	"time"
)

type MigrationTestSuite struct {
	suite.Suite
	app *uadmin.App
	db  *core.UadminDatabase
}

func (suite *MigrationTestSuite) SetupTest() {
	appliedMigrations = make([]string, 0)
	app := uadmin.NewTestApp()
	suite.app = app
	db := core.NewUadminDatabase()
	defer db.Close()
	db.Db.Exec("DROP TABLE migrations")
	db.Db.AutoMigrate(uadmin.Migration{})
	suite.app.BlueprintRegistry = core.NewBlueprintRegistry()
	suite.app.BlueprintRegistry.Register(TestBlueprint)
	suite.app.BlueprintRegistry.Register(Test1Blueprint)
}

func (suite *MigrationTestSuite) TearDownSuite() {
	err := os.Remove(suite.app.Config.D.Db.Default.Name)
	if err != nil {
		assert.Equal(suite.T(), true, false, fmt.Errorf("Couldnt remove db with name %s", suite.app.Config.D.Db.Default.Name))
	}
	uadmin.ClearTestApp()
}

func (suite *MigrationTestSuite) TestUpgradeDatabase() {
	suite.app.BlueprintRegistry.Register(TestBlueprint)
	suite.app.BlueprintRegistry.Register(Test1Blueprint)
	suite.app.TriggerCommandExecution("migrate", "up", make([]string, 0))
	appliedMigrationsExpected := mapset.NewSet()
	appliedMigrationsExpected.Add("user.1621667393")
	appliedMigrationsExpected.Add("user.1621680132")
	appliedMigrationsExpected.Add("test1.1621667393")
	appliedMigrationsExpected.Add("test1.1621680132")
	appliedMigrationsActual := mapset.NewSet()
	for _, migrationName := range appliedMigrations {
		appliedMigrationsActual.Add(migrationName)
	}
	assert.Equal(suite.T(), appliedMigrationsExpected, appliedMigrationsActual)
	var appliedMigrationsDb []uadmin.Migration
	db := core.NewUadminDatabase()
	defer db.Close()
	db.Db.Find(&appliedMigrationsDb)
	assert.Equal(suite.T(), 4, len(appliedMigrationsDb))
}

func (suite *MigrationTestSuite) TestDowngradeDatabase() {
	suite.app.BlueprintRegistry.Register(TestBlueprint)
	suite.app.BlueprintRegistry.Register(Test1Blueprint)
	appliedMigrationsExpected := mapset.NewSet()
	appliedMigrationsExpected.Add("user.1621667393")
	appliedMigrationsExpected.Add("user.1621680132")
	appliedMigrationsExpected.Add("test1.1621667393")
	appliedMigrationsExpected.Add("test1.1621680132")
	db := core.NewUadminDatabase()
	defer db.Close()
	db.Db.Create(
		&uadmin.Migration{MigrationName: "user.1621667393", AppliedAt: time.Now()},
	)
	db.Db.Create(
		&uadmin.Migration{MigrationName: "user.1621680132", AppliedAt: time.Now()},
	)
	db.Db.Create(
		&uadmin.Migration{MigrationName: "test1.1621667393", AppliedAt: time.Now()},
	)
	db.Db.Create(
		&uadmin.Migration{MigrationName: "test1.1621680132", AppliedAt: time.Now()},
	)
	var appliedMigrationsDb []uadmin.Migration
	db.Db.Find(&appliedMigrationsDb)
	assert.Equal(suite.T(), 4, len(appliedMigrationsDb))
	suite.app.TriggerCommandExecution("migrate", "down", []string{""})
	appliedMigrationsDb = make([]uadmin.Migration, 0)
	db.Db.Find(&appliedMigrationsDb)
	assert.Equal(suite.T(), 0, len(appliedMigrationsDb))
}

func (suite *MigrationTestSuite) TestTraverseDatabaseForUpgrade() {
	concreteBlueprintRegistry := core.NewBlueprintRegistry()
	concreteBlueprintRegistry.Register(TestBlueprint)
	concreteBlueprintRegistry.Register(Test1Blueprint)
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	for res := range concreteBlueprintRegistry.TraverseMigrations() {
		res.Node.Apply(uadminDatabase)
	}
	appliedMigrationsExpected := mapset.NewSet()
	appliedMigrationsExpected.Add("user.1621667393")
	appliedMigrationsExpected.Add("user.1621680132")
	appliedMigrationsExpected.Add("test1.1621667393")
	appliedMigrationsExpected.Add("test1.1621680132")
	appliedMigrationsActual := mapset.NewSet()
	for _, migrationName := range appliedMigrations {
		appliedMigrationsActual.Add(migrationName)
	}
	assert.Equal(suite.T(), appliedMigrationsExpected, appliedMigrationsActual)
}

func (suite *MigrationTestSuite) TestTraverseDatabaseForDowngrade() {
	concreteBlueprintRegistry := core.NewBlueprintRegistry()
	concreteBlueprintRegistry.Register(TestBlueprint)
	concreteBlueprintRegistry.Register(Test1Blueprint)
	toDowngradeMigrationsExpected := mapset.NewSet()
	toDowngradeMigrationsExpected.Add("user.1621667393")
	toDowngradeMigrationsExpected.Add("user.1621680132")
	toDowngradeMigrationsExpected.Add("test1.1621667393")
	toDowngradeMigrationsExpected.Add("test1.1621680132")
	downgradedMigrationsActual := mapset.NewSet()
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	for res := range concreteBlueprintRegistry.TraverseMigrationsDownTo("") {
		res.Node.Downgrade(uadminDatabase)
		downgradedMigrationsActual.Add(res.Node.GetMigration().GetName())
	}
	assert.Equal(suite.T(), toDowngradeMigrationsExpected, downgradedMigrationsActual)
}

func (suite *MigrationTestSuite) TestBuildTreeForBlueprintWithNoMigrations() {
	blueprintRegistry := core.NewBlueprintRegistry()
	blueprintRegistry.Register(BlueprintWithNoMigrations)
	for res := range blueprintRegistry.TraverseMigrations() {
		if res.Error != nil {
			assert.Equal(suite.T(), false, true, res.Error)
			return
		}
	}
}

func (suite *MigrationTestSuite) TestBuildTreeWithNoUserBlueprint() {
	blueprintRegistry := core.NewBlueprintRegistry()
	blueprintRegistry.Register(Test1Blueprint)
	for res := range blueprintRegistry.TraverseMigrations() {
		assert.Equal(suite.T(), res.Error, fmt.Errorf("Couldn't find blueprint with name user"))
		return
	}
	assert.True(suite.T(), false)
}

func (suite *MigrationTestSuite) TestBuildTreeWithTwoNoDepsMigrationsFromtheSameBlueprint() {
	blueprintRegistry := core.NewBlueprintRegistry()
	blueprintRegistry.Register(BlueprintWithTwoSameDeps)
	for res := range blueprintRegistry.TraverseMigrations() {
		assert.True(suite.T(), strings.Contains(res.Error.Error(), "Found two or more migrations with no children from the same blueprint"))
		return
	}
	assert.True(suite.T(), false)
}

func (suite *MigrationTestSuite) TestBuildTreeWithTwoNoChildMigrationsFromtheSameBlueprint() {
	blueprintRegistry := core.NewBlueprintRegistry()
	blueprintRegistry.Register(BlueprintWithConflicts)
	for res := range blueprintRegistry.TraverseMigrations() {
		assert.True(suite.T(), strings.Contains(res.Error.Error(), "Found two or more migrations with no children from the same blueprint"))
		return
	}
	assert.True(suite.T(), false)
}

func (suite *MigrationTestSuite) TestBuildTreeWithLoop() {
	blueprintRegistry := core.NewBlueprintRegistry()
	blueprintRegistry.Register(Blueprint1WithLoopedMigrations)
	blueprintRegistry.Register(Blueprint2WithLoopedMigrations)
	for range blueprintRegistry.TraverseMigrations() {
	}
}

//func (suite *MigrationTestSuite) TestBuildTreeWithTwoSameMigrationNames() {
//	blueprintRegistry := interfaces.NewBlueprintRegistry()
//	blueprintRegistry.Register(Blueprint1WithSameMigrationNames)
//	blueprintRegistry.Register(Blueprint2WithSameMigrationNames)
//	for res := range blueprintRegistry.TraverseMigrations() {
//		assert.True(suite.T(), strings.Contains(res.Error.Error(), "has been added to tree before"))
//	}
//}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSqliteMigrations(t *testing.T) {
	uadmin.ClearApp()
	suite.Run(t, new(MigrationTestSuite))
}
