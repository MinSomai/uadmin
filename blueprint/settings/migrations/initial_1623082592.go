
package migrations

import (
    "github.com/uadmin/uadmin/blueprint/settings/models"
    "github.com/uadmin/uadmin/dialect"
)

type initial_1623082592 struct {
}

func (m initial_1623082592) GetName() string {
    return "settings.1623082592"
}

func (m initial_1623082592) GetId() int64 {
    return 1623082592
}

func (m initial_1623082592) Up() {
    db := dialect.GetDB("default")
    db.AutoMigrate(models.SettingCategory{})
    db.AutoMigrate(models.Setting{})
}

func (m initial_1623082592) Down() {
    db := dialect.GetDB("default")
    err := db.Migrator().DropTable(models.Setting{})
    if err != nil {
        panic(err)
    }
    err = db.Migrator().DropTable(models.SettingCategory{})
    if err != nil {
        panic(err)
    }
}

func (m initial_1623082592) Deps() []string {
    return make([]string, 0)
}
