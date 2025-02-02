package migrations

import (
	logmodel "github.com/sergeyglazyrindev/uadmin/blueprint/logging/models"
	"github.com/sergeyglazyrindev/uadmin/core"
	"gorm.io/gorm"
)

type initial1623082882 struct {
}

func (m initial1623082882) GetName() string {
	return "logging.1623082882"
}

func (m initial1623082882) GetID() int64 {
	return 1623082882
}

func (m initial1623082882) Up(uadminDatabase *core.UadminDatabase) error {
	db := uadminDatabase.Db
	err := db.AutoMigrate(logmodel.Log{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1623082882) Down(uadminDatabase *core.UadminDatabase) error {
	db := uadminDatabase.Db
	err := db.Migrator().DropTable(logmodel.Log{})
	if err != nil {
		return err
	}
	var contentType core.ContentType
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(&logmodel.Log{})
	db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "logging", ModelName: stmt.Schema.Table}).First(&contentType)
	db.Unscoped().Where(" 1 = 1").Delete(&contentType)
	return nil
}

func (m initial1623082882) Deps() []string {
	return make([]string, 0)
}
