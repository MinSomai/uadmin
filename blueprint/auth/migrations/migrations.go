package migrations

import (
	"github.com/uadmin/uadmin/interfaces"
)

var BMigrationRegistry *interfaces.MigrationRegistry

func init() {
    BMigrationRegistry = interfaces.NewMigrationRegistry()
    
    BMigrationRegistry.AddMigration(initial_1623582035{})
    // placeholder to insert next migration
}
