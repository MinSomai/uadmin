package migrations

import (
    "github.com/uadmin/uadmin/blueprint/approval/models"
    "github.com/uadmin/uadmin/interfaces"
    "gorm.io/gorm"
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
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&models.Approval{})
    contentType := &interfaces.ContentType{BlueprintName: "approval", ModelName: stmt.Schema.Table}
    db.Create(contentType)
}

func (m initial_1623083268) Down() {
    db := interfaces.GetDB()
    err := db.Migrator().DropTable(models.Approval{})
    if err != nil {
        panic(err)
    }
    var contentType interfaces.ContentType
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&models.Approval{})
    db.Model(&interfaces.ContentType{}).Where(&interfaces.ContentType{BlueprintName: "approval", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Delete(&contentType)
}

func (m initial_1623083268) Deps() []string {
    return make([]string, 0)
}
