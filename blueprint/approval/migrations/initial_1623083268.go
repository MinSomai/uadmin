package migrations

import (
    "github.com/uadmin/uadmin/blueprint/approval/models"
	"github.com/uadmin/uadmin/core"
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

func (m initial_1623083268) Up(uadminDatabase *core.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.AutoMigrate(models.Approval{})
    if err != nil {
        return err
    }
    return nil
}

func (m initial_1623083268) Down(uadminDatabase *core.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.Migrator().DropTable(models.Approval{})
    if err != nil {
        return err
    }
    var contentType core.ContentType
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&models.Approval{})
    db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "approval", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Unscoped().Delete(&contentType)
    return nil
}

func (m initial_1623083268) Deps() []string {
    return make([]string, 0)
}
