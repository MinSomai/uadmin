package uadmin

import (
	"github.com/uadmin/uadmin/interfaces"
	"testing"
	"github.com/stretchr/testify/assert"
)

type MigrationRegistry struct {
	migrations map[string]interfaces.IMigration
}

func (r MigrationRegistry) addMigration(migration interfaces.IMigration) {
	r.migrations[migration.GetName()] = migration
}

func (r MigrationRegistry) FindMigrations() {
}

type initialb1_1621667392 struct {
}

func (m initialb1_1621667392) GetName() string {
	return "initial"
}

func (m initialb1_1621667392) GetId() int64 {
	return 1621667392
}

func (m initialb1_1621667392) Up() {
}

func (m initialb1_1621667392) Down() {
}

func (m initialb1_1621667392) Deps() []string {
	return make([]string, 0)
}


var BMigrationRegistry *MigrationRegistry

// TestFormHandler is a unit testing function for formHandler() function
func TestMigrate(t *testing.T) {
	BMigrationRegistry = &MigrationRegistry{
		migrations: make(map[string]interfaces.IMigration),
	}
	BMigrationRegistry.addMigration(initialb1_1621667392{})
	assert.Equal(t, false, true)
}
