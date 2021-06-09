package migrations

import (
	"github.com/uadmin/uadmin/interfaces"
)

var BMigrationRegistry *interfaces.MigrationRegistry

func init() {
    BMigrationRegistry = interfaces.NewMigrationRegistry()
    
    BMigrationRegistry.AddMigration(initial_1623083053{})
    
    BMigrationRegistry.AddMigration(create_all_1623263607{})
    // placeholder to insert next migration
}
