package migrations

import (
	"github.com/uadmin/uadmin/interfaces"
)

var BMigrationRegistry *interfaces.MigrationRegistry

func init() {
    BMigrationRegistry = interfaces.NewMigrationRegistry()
    
    BMigrationRegistry.AddMigration(initial_1621680132{})
    
    BMigrationRegistry.AddMigration(adding_use_1623259185{})

    // placeholder to insert next migration
}
