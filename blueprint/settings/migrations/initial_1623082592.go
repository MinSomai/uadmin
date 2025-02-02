package migrations

import (
	"github.com/sergeyglazyrindev/uadmin/blueprint/settings/models"
	"github.com/sergeyglazyrindev/uadmin/core"
	"gorm.io/gorm"
)

type initial1623082592 struct {
}

func (m initial1623082592) GetName() string {
	return "settings.1623082592"
}

func (m initial1623082592) GetID() int64 {
	return 1623082592
}

func (m initial1623082592) Up(uadminDatabase *core.UadminDatabase) error {
	db := uadminDatabase.Db
	db.AutoMigrate(models.SettingCategory{})
	db.AutoMigrate(models.Setting{})
	return nil
}

func (m initial1623082592) Down(uadminDatabase *core.UadminDatabase) error {
	db := uadminDatabase.Db
	err := db.Migrator().DropTable(models.Setting{})
	if err != nil {
		return err
	}
	err = db.Migrator().DropTable(models.SettingCategory{})
	if err != nil {
		return err
	}
	var contentType core.ContentType
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(&models.SettingCategory{})
	db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "settings", ModelName: stmt.Schema.Table}).First(&contentType)
	db.Unscoped().Where(" 1 = 1").Delete(&contentType)
	stmt = &gorm.Statement{DB: db}
	stmt.Parse(&models.Setting{})
	db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "settings", ModelName: stmt.Schema.Table}).First(&contentType)
	db.Unscoped().Where(" 1 = 1").Delete(&contentType)
	return nil
}

func (m initial1623082592) Deps() []string {
	return make([]string, 0)
}
