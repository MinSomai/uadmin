package migrations

import (
    "github.com/uadmin/uadmin/blueprint/menu/models"
    "github.com/uadmin/uadmin/interfaces"
)

type initial_1623081544 struct {
}

func (m initial_1623081544) GetName() string {
    return "menu.1623081544"
}

func (m initial_1623081544) GetId() int64 {
    return 1623081544
}

func (m initial_1623081544) Up() {
    db := interfaces.GetDB()
    err := db.AutoMigrate(models.DashboardMenu{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623081544) Down() {
    db := interfaces.GetDB()
    err := db.Migrator().DropTable(models.DashboardMenu{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623081544) Deps() []string {
    return make([]string, 0)
}
