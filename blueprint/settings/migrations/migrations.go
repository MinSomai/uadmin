package migrations

import (
	"github.com/uadmin/uadmin/core"
)

var BMigrationRegistry *core.MigrationRegistry

func init() {
	BMigrationRegistry = core.NewMigrationRegistry()

	BMigrationRegistry.AddMigration(initial1623082592{})

	BMigrationRegistry.AddMigration(insertall1623263908{})
	// placeholder to insert next migration
}
