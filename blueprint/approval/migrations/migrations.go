package migrations

import (
	"github.com/uadmin/uadmin/interfaces"
)

type MigrationRegistry struct {
	migrations map[string]interfaces.IMigration
}

func (r MigrationRegistry) addMigration(migration interfaces.IMigration) {
	r.migrations[migration.GetName()] = migration
}

func (r MigrationRegistry) FindMigrations() <-chan interfaces.IMigration{
	chnl := make(chan interfaces.IMigration)
	go func() {
		close(chnl)
	}()
	return chnl
}

var BMigrationRegistry *MigrationRegistry

func init() {
    BMigrationRegistry = &MigrationRegistry{
        migrations: make(map[string]interfaces.IMigration),
    }
    BMigrationRegistry.addMigration(initial_1621667383{})

    // placeholder to insert next migration
}