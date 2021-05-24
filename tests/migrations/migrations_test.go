package migrations

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/dialect"
	"gorm.io/gorm"
	"os"
	"testing"
)

type MigrationTestSuite struct {
	suite.Suite
	db *gorm.DB
	app *uadmin.App
}

func (suite *MigrationTestSuite) SetupTest() {
	app := uadmin.NewApp("test")
	suite.app = app
	dialect.CurrentDatabaseSettings = &dialect.DatabaseSettings{
		Default: app.Config.D.Db.Default,
	}
	dialectdb := dialect.NewDbDialect(app.Database.ConnectTo("default"), app.Config.D.Db.Default.Type)
	db, err := dialectdb.GetDb(
		"default",
	)
	suite.db = db
	if err != nil {
		assert.Equal(suite.T(), true, false, fmt.Errorf("Couldnt initialize db %s", err))
	}
	//db = db.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")
	//if db.Error != nil {
	//	assert.Equal(suite.T(), true, false, "Couldnt setup isolation level for db")
	//}
	//db = db.Exec("BEGIN")
	//if db.Error != nil {
	//	assert.Equal(suite.T(), true, false, "Couldnt start transaction")
	//}
}

func (suite *MigrationTestSuite) TearDownSuite() {
	err := os.Remove(suite.app.Config.D.Db.Default.Name)
	if err != nil {
		assert.Equal(suite.T(), true, false, fmt.Errorf("Couldnt remove db with name %s", suite.db.Name()))
	}
	//db := uadmin.GetDB()
	//db = db.Exec("ROLLBACK")
	//if db.Error != nil {
	//	assert.Equal(suite.T(), true, false, "Couldnt rollback transaction")
	//}
}


func (suite *MigrationTestSuite) TestUpgradeDatabase() {
	suite.app.BlueprintRegistry = &ConcreteBlueprintRegistry
	migrateCommand := new(uadmin.MigrateCommand)
	err := migrateCommand.Proceed("up", make([]string, 0))
	if err != nil {
		assert.Equal(suite.T(), false, true, err)
	}
	assert.Equal(suite.T(), false, true, "dsadsa")

}

func (suite *MigrationTestSuite) TestDetermineConflictsInMigrations() {
	err := BlueprintWithConflicts.GetMigrationRegistry().BuildTree(BlueprintWithConflicts.GetName())
	if err == nil {
		assert.Equal(suite.T(), false, true, `
			Had to find conflicts but didnt. Please check that we properly built tree and determined conflicts
			in migrations
		`)
		return
	}
	conflicts := BlueprintWithConflicts.GetMigrationRegistry().FindPotentialConflictsForBlueprint(
		BlueprintWithConflicts.GetName(),
	)
	assert.Equal(suite.T(), len(conflicts), 2)
}

func (suite *MigrationTestSuite) TestBuildTreeForBlueprintWithNoMigrations() {
	err := BlueprintWithNoMigrations.GetMigrationRegistry().BuildTree(BlueprintWithNoMigrations.GetName())
	if err != nil {
		assert.Equal(suite.T(), false, true, err)
		return
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMigrations(t *testing.T) {
	suite.Run(t, new(MigrationTestSuite))
}