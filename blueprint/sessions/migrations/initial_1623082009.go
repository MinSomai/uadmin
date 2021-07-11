package migrations

import (
    "github.com/uadmin/uadmin/blueprint/sessions/models"
    "github.com/uadmin/uadmin/interfaces"
)

type initial_1623082009 struct {
}

func (m initial_1623082009) GetName() string {
    return "sessions.1623082009"
}

func (m initial_1623082009) GetId() int64 {
    return 1623082009
}

func (m initial_1623082009) Up() {
    db := interfaces.GetDB()
    err := db.AutoMigrate(models.Session{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623082009) Down() {
    db := interfaces.GetDB()
    err := db.Migrator().DropTable(models.Session{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623082009) Deps() []string {
    return []string{"user.1621680132"}
}
