package migrations

import (
    logmodel "github.com/uadmin/uadmin/blueprint/logging/models"
    "github.com/uadmin/uadmin/interfaces"
    "gorm.io/gorm"
)

type initial_1623082882 struct {
}

func (m initial_1623082882) GetName() string {
    return "logging.1623082882"
}

func (m initial_1623082882) GetId() int64 {
    return 1623082882
}

func (m initial_1623082882) Up(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.AutoMigrate(logmodel.Log{})
    if err != nil {
        return err
    }
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&logmodel.Log{})
    contentType := &interfaces.ContentType{BlueprintName: "logging", ModelName: stmt.Schema.Table}
    db.Create(contentType)
    return nil
}

func (m initial_1623082882) Down(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.Migrator().DropTable(logmodel.Log{})
    if err != nil {
        return err
    }
    var contentType interfaces.ContentType
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&logmodel.Log{})
    db.Model(&interfaces.ContentType{}).Where(&interfaces.ContentType{BlueprintName: "logging", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Delete(&contentType)
    return nil
}

func (m initial_1623082882) Deps() []string {
    return make([]string, 0)
}
