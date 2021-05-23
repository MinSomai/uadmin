package user

import (
	"github.com/uadmin/uadmin/blueprint/user/migrations"
	"github.com/uadmin/uadmin/interfaces"
)

var Blueprint = interfaces.Blueprint{
	Name: "user",
	Description: "this blueprint is about users",
	MigrationRegistry: migrations.BMigrationRegistry,
}
