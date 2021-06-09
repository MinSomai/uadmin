package migrations

import (
    logmodel "github.com/uadmin/uadmin/blueprint/logging/models"
    "github.com/uadmin/uadmin/dialect"
)

type initial_1623082882 struct {
}

func (m initial_1623082882) GetName() string {
    return "logging.1623082882"
}

func (m initial_1623082882) GetId() int64 {
    return 1623082882
}

func (m initial_1623082882) Up() {
    db := dialect.GetDB()
    err := db.AutoMigrate(logmodel.Log{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623082882) Down() {
    db := dialect.GetDB()
    err := db.Migrator().DropTable(logmodel.Log{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623082882) Deps() []string {
    return make([]string, 0)
}
