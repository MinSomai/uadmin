package migrations

import (
    "github.com/uadmin/uadmin/blueprint/settings/models"
    "github.com/uadmin/uadmin/interfaces"
    "gorm.io/gorm"
)

type initial_1623082592 struct {
}

func (m initial_1623082592) GetName() string {
    return "settings.1623082592"
}

func (m initial_1623082592) GetId() int64 {
    return 1623082592
}

func (m initial_1623082592) Up(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    db.AutoMigrate(models.SettingCategory{})
    db.AutoMigrate(models.Setting{})
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&models.SettingCategory{})
    contentType := &interfaces.ContentType{BlueprintName: "settings", ModelName: stmt.Schema.Table}
    db.Create(contentType)
    stmt = &gorm.Statement{DB: db}
    stmt.Parse(&models.Setting{})
    settingContentType := &interfaces.ContentType{BlueprintName: "settings", ModelName: stmt.Schema.Table}
    db.Create(settingContentType)
    return nil
}

func (m initial_1623082592) Down(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.Migrator().DropTable(models.Setting{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(models.SettingCategory{})
    if err != nil {
        return err
    }
    var contentType interfaces.ContentType
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&models.SettingCategory{})
    db.Model(&interfaces.ContentType{}).Where(&interfaces.ContentType{BlueprintName: "settings", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Delete(&contentType)
    stmt = &gorm.Statement{DB: db}
    stmt.Parse(&models.Setting{})
    db.Model(&interfaces.ContentType{}).Where(&interfaces.ContentType{BlueprintName: "settings", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Delete(&contentType)
    return nil
}

func (m initial_1623082592) Deps() []string {
    return make([]string, 0)
}
