package migrations

import (
	"github.com/uadmin/uadmin/interfaces"
)

var BMigrationRegistry *interfaces.MigrationRegistry

func init() {
    BMigrationRegistry = interfaces.NewMigrationRegistry()
    
    BMigrationRegistry.AddMigration(initial_1623081544{})
    
    BMigrationRegistry.AddMigration(Adddashbo_1623217408{})
    // placeholder to insert next migration
}
