package migrations

import (
	"github.com/uadmin/uadmin/interfaces"
)

type MigrationRegistry struct {
	migrations map[string]interfaces.MigrationInterface
}

func (r MigrationRegistry) addMigration(migration interfaces.MigrationInterface) {
	r.migrations[migration.GetName()] = migration
}

func (r MigrationRegistry) FindMigrations() {
}

var BMigrationRegistry *MigrationRegistry

func init() {
    BMigrationRegistry = &MigrationRegistry{
        migrations: make(map[string]interfaces.MigrationInterface),
    }
    // placeholder to insert next migration
}