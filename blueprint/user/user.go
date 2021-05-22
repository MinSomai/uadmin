package user

import (
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/blueprint/user/migrations"
)

type Blueprint struct {
	interfaces.IBlueprint
}

func (b Blueprint) GetName() string {
	return "user"
}

func (b Blueprint) GetMigrationRegistry() interfaces.IMigrationRegistry {
	return interfaces.IMigrationRegistry(migrations.BMigrationRegistry)
}