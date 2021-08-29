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

func (m initial_1623083268) Up(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.AutoMigrate(models.Approval{})
    if err != nil {
        return err
    }
    return nil
}

func (m initial_1623083268) Down(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.Migrator().DropTable(models.Approval{})
    if err != nil {
        return err
    }
    var contentType interfaces.ContentType
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&models.Approval{})
    db.Model(&interfaces.ContentType{}).Where(&interfaces.ContentType{BlueprintName: "approval", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Unscoped().Delete(&contentType)
    return nil
}

func (m initial_1623083268) Deps() []string {
    return make([]string, 0)
}
