package migrations

import (
	"github.com/uadmin/uadmin/interfaces"
)

type MigrationRegistry struct {
	migrations map[string]interfaces.IMigration
}

func (r MigrationRegistry) AddMigration(migration interfaces.IMigration) {
	r.migrations[migration.GetName()] = migration
}

func (r MigrationRegistry) FindMigrations() {
}

var BMigrationRegistry *MigrationRegistry

func init() {
    BMigrationRegistry = &MigrationRegistry{
        migrations: make(map[string]interfaces.IMigration),
    }
    // placeholder to insert next migration
}