package migrations

import (
    "github.com/uadmin/uadmin/blueprint/approval/models"
    "github.com/uadmin/uadmin/interfaces"
)

type initial_1623083268 struct {
}

func (m initial_1623083268) GetName() string {
    return "approval.1623083268"
}

func (m initial_1623083268) GetId() int64 {
    return 1623083268
}

func (m initial_1623083268) Up() {
    db := interfaces.GetDB()
    err := db.AutoMigrate(models.Approval{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623083268) Down() {
    db := interfaces.GetDB()
    err := db.Migrator().DropTable(models.Approval{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623083268) Deps() []string {
    return make([]string, 0)
}
