package migrations

import (
    abtestmodel "github.com/uadmin/uadmin/blueprint/abtest/models"
    "github.com/uadmin/uadmin/interfaces"
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

func (m initial_1623083395) Up(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.AutoMigrate(abtestmodel.ABTest{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(abtestmodel.ABTestValue{})
    if err != nil {
        return err
    }
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&abtestmodel.ABTest{})
    contentType := &interfaces.ContentType{BlueprintName: "abtest", ModelName: stmt.Schema.Table}
    db.Create(contentType)
    stmt = &gorm.Statement{DB: db}
    stmt.Parse(&abtestmodel.ABTestValue{})
    contentType = &interfaces.ContentType{BlueprintName: "abtest", ModelName: stmt.Schema.Table}
    db.Create(contentType)
    return nil
}

func (m initial_1623083395) Down(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.Migrator().DropTable(abtestmodel.ABTestValue{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(abtestmodel.ABTest{})
    if err != nil {
        return err
    }
    var contentType interfaces.ContentType
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&abtestmodel.ABTestValue{})
    db.Model(&interfaces.ContentType{}).Where(&interfaces.ContentType{BlueprintName: "abtest", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Delete(&contentType)
    stmt = &gorm.Statement{DB: db}
    stmt.Parse(&abtestmodel.ABTest{})
    db.Model(&interfaces.ContentType{}).Where(&interfaces.ContentType{BlueprintName: "abtest", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Delete(&contentType)
    return nil
}

func (m initial_1623083395) Deps() []string {
    return make([]string, 0)
}
