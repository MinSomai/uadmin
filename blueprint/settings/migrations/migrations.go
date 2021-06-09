package migrations

import (
	"github.com/uadmin/uadmin/interfaces"
)

var BMigrationRegistry *interfaces.MigrationRegistry

func init() {
    BMigrationRegistry = interfaces.NewMigrationRegistry()
    
    BMigrationRegistry.AddMigration(initial_1623082592{})
    
    BMigrationRegistry.AddMigration(insert_all_1623263908{})
    // placeholder to insert next migration
}
