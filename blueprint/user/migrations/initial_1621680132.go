
package migrations

import (
    models2 "github.com/uadmin/uadmin/blueprint/user/models"
    "github.com/uadmin/uadmin/dialect"
)

type initial_1621680132 struct {
}

func (m initial_1621680132) GetName() string {
    return "user.initial"
}

func (m initial_1621680132) GetId() int64 {
    return 1621680132
}

func (m initial_1621680132) Up() {
    db := dialect.GetDB("dialect")
    err := db.AutoMigrate(models2.UserGroup{})
    if err != nil {
        panic(err)
    }
    err = db.AutoMigrate(models2.GroupPermission{})
    if err != nil {
        panic(err)
    }
    err = db.AutoMigrate(models2.User{})
    if err != nil {
        panic(err)
    }
    err = db.AutoMigrate(models2.UserPermission{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1621680132) Down() {
    db := dialect.GetDB("default")
    err := db.Migrator().DropTable(models2.UserPermission{})
    if err != nil {
        panic(err)
    }
    err = db.Migrator().DropTable(models2.User{})
    if err != nil {
        panic(err)
    }
    err = db.Migrator().DropTable(models2.GroupPermission{})
    if err != nil {
        panic(err)
    }
    err = db.Migrator().DropTable(models2.UserGroup{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1621680132) Deps() []string {
    return []string{"menu.1623081544"}
}
