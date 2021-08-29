package migrations

import (
    abtestmodel "github.com/uadmin/uadmin/blueprint/abtest/models"
	"github.com/uadmin/uadmin/core"
    "gorm.io/gorm"
)

type initial_1623083395 struct {
}

func (m initial_1623083395) GetName() string {
    return "abtest.1623083395"
}

func (m initial_1623083395) GetId() int64 {
    return 1623083395
}

func (m initial_1623083395) Up(uadminDatabase *core.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.AutoMigrate(abtestmodel.ABTest{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(abtestmodel.ABTestValue{})
    if err != nil {
        return err
    }
    return nil
}

func (m initial_1623083395) Down(uadminDatabase *core.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.Migrator().DropTable(abtestmodel.ABTestValue{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(abtestmodel.ABTest{})
    if err != nil {
        return err
    }
    var contentType core.ContentType
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&abtestmodel.ABTestValue{})
    db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "abtest", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Unscoped().Delete(&contentType)
    stmt = &gorm.Statement{DB: db}
    stmt.Parse(&abtestmodel.ABTest{})
    db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "abtest", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Unscoped().Delete(&contentType)
    return nil
}

func (m initial_1623083395) Deps() []string {
    return make([]string, 0)
}
