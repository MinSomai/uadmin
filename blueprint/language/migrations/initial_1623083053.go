package migrations

import (
    langmodel "github.com/uadmin/uadmin/blueprint/language/models"
    "github.com/uadmin/uadmin/interfaces"
    "gorm.io/gorm"
)

type initial_1623083053 struct {
}

func (m initial_1623083053) GetName() string {
    return "language.1623083053"
}

func (m initial_1623083053) GetId() int64 {
    return 1623083053
}

func (m initial_1623083053) Up() {
    db := interfaces.GetDB()
    err := db.AutoMigrate(langmodel.Language{})
    if err != nil {
        panic(err)
    }
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&langmodel.Language{})
    contentType := &interfaces.ContentType{BlueprintName: "language", ModelName: stmt.Schema.Table}
    db.Create(contentType)
}

func (m initial_1623083053) Down() {
    db := interfaces.GetDB()
    err := db.Migrator().DropTable(langmodel.Language{})
    if err != nil {
        panic(err)
    }
    var contentType interfaces.ContentType
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&langmodel.Language{})
    db.Model(&interfaces.ContentType{}).Where(&interfaces.ContentType{BlueprintName: "language", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Delete(&contentType)
}

func (m initial_1623083053) Deps() []string {
    return make([]string, 0)
}
