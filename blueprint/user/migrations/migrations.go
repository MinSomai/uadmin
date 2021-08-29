package migrations

import (
	"github.com/uadmin/uadmin/core"
)

var BMigrationRegistry *core.MigrationRegistry

func init() {
    BMigrationRegistry = core.NewMigrationRegistry()
    
    BMigrationRegistry.AddMigration(initial_1621680132{})
    
    BMigrationRegistry.AddMigration(adding_use_1623259185{})

    // placeholder to insert next migration
}
