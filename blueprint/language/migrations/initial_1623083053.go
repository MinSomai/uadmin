package migrations

import (
    langmodel "github.com/uadmin/uadmin/blueprint/language/models"
    "github.com/uadmin/uadmin/interfaces"
)

type initial_1623083053 struct {
}

func (m initial_1623083053) GetName() string {
    return "language.1623083053"
}

func (m initial_1623083053) GetId() int64 {
    return 1623083053
}

func (m initial_1623083053) Up() {
    db := interfaces.GetDB()
    err := db.AutoMigrate(langmodel.Language{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623083053) Down() {
    db := interfaces.GetDB()
    err := db.Migrator().DropTable(langmodel.Language{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623083053) Deps() []string {
    return make([]string, 0)
}
