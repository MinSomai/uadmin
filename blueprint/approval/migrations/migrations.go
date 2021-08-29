package migrations

import (
	"github.com/uadmin/uadmin/core"
)

var BMigrationRegistry *core.MigrationRegistry

func init() {
    BMigrationRegistry = core.NewMigrationRegistry()
    
    BMigrationRegistry.AddMigration(initial_1623083268{})
    // placeholder to insert next migration
}
