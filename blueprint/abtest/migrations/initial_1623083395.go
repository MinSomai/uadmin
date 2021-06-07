package migrations

import (
    abtestmodel "github.com/uadmin/uadmin/blueprint/abtest/models"
    "github.com/uadmin/uadmin/dialect"
)

type initial_1623083395 struct {
}

func (m initial_1623083395) GetName() string {
    return "abtest.1623083395"
}

func (m initial_1623083395) GetId() int64 {
    return 1623083395
}

func (m initial_1623083395) Up() {
    db := dialect.GetDB("dialect")
    err := db.AutoMigrate(abtestmodel.ABTest{})
    if err != nil {
        panic(err)
    }
    err = db.AutoMigrate(abtestmodel.ABTestValue{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623083395) Down() {
    db := dialect.GetDB("dialect")
    err := db.Migrator().DropTable(abtestmodel.ABTestValue{})
    if err != nil {
        panic(err)
    }
    err = db.Migrator().DropTable(abtestmodel.ABTest{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623083395) Deps() []string {
    return make([]string, 0)
}
