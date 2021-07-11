package migrations

import (
    models2 "github.com/uadmin/uadmin/blueprint/auth/models"
    "github.com/uadmin/uadmin/interfaces"
)

type initial_1623582035 struct {
}

func (m initial_1623582035) GetName() string {
    return "auth.1623582035"
}

func (m initial_1623582035) GetId() int64 {
    return 1623582035
}

func (m initial_1623582035) Up() {
    db := interfaces.GetDB()
    err := db.AutoMigrate(models2.UserAuthToken{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623582035) Down() {
    db := interfaces.GetDB()
    err := db.Migrator().DropTable(models2.UserAuthToken{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623582035) Deps() []string {
    return []string{"user.1621680132"}
}
