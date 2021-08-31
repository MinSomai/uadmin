package migrations

import (
	"github.com/uadmin/uadmin/core"
)

var BMigrationRegistry *core.MigrationRegistry

func init() {
	BMigrationRegistry = core.NewMigrationRegistry()

	BMigrationRegistry.AddMigration(initial_1623082592{})

	BMigrationRegistry.AddMigration(insert_all_1623263908{})
	// placeholder to insert next migration
}
