package migrations

import (
    models2 "github.com/uadmin/uadmin/blueprint/user/models"
    "github.com/uadmin/uadmin/interfaces"
    "gorm.io/gorm"
)

type initial_1621680132 struct {
}

func (m initial_1621680132) GetName() string {
    return "user.1621680132"
}

func (m initial_1621680132) GetId() int64 {
    return 1621680132
}

func (m initial_1621680132) Up(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.AutoMigrate(interfaces.ContentType{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(models2.UserGroup{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(models2.User{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(models2.Permission{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(models2.OneTimeAction{})
    if err != nil {
        return err
    }
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&models2.OneTimeAction{})
    contentType := &interfaces.ContentType{BlueprintName: "user", ModelName: stmt.Schema.Table}
    db.Create(contentType)
    stmt = &gorm.Statement{DB: db}
    stmt.Parse(&models2.User{})
    userContentType := &interfaces.ContentType{BlueprintName: "user", ModelName: stmt.Schema.Table}
    db.Create(userContentType)
    return nil
}

func (m initial_1621680132) Down(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.Migrator().DropTable(models2.Permission{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(models2.User{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(models2.UserGroup{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(models2.OneTimeAction{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(interfaces.ContentType{})
    if err != nil {
        return err
    }
    var contentType interfaces.ContentType
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&models2.OneTimeAction{})
    db.Model(&interfaces.ContentType{}).Where(&interfaces.ContentType{BlueprintName: "user", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Delete(&contentType)
    stmt = &gorm.Statement{DB: db}
    stmt.Parse(&models2.User{})
    db.Model(&interfaces.ContentType{}).Where(&interfaces.ContentType{BlueprintName: "user", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Delete(&contentType)
    return nil
}

func (m initial_1621680132) Deps() []string {
    return []string{}
}
