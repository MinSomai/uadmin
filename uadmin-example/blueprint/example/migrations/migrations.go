package migrations

import (
	"github.com/sergeyglazyrindev/uadmin/core"
)

var BMigrationRegistry *core.MigrationRegistry

func init() {
    BMigrationRegistry = core.NewMigrationRegistry()
    
    BMigrationRegistry.AddMigration(initial_1631027794{})
    // placeholder to insert next migration
}
