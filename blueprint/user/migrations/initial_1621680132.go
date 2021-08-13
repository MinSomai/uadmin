package migrations

import (
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
    err = db.AutoMigrate(interfaces.UserGroup{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(interfaces.User{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(interfaces.Permission{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(interfaces.OneTimeAction{})
    if err != nil {
        return err
    }
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&interfaces.OneTimeAction{})
    contentType := &interfaces.ContentType{BlueprintName: "user", ModelName: stmt.Schema.Table}
    db.Create(contentType)
    stmt = &gorm.Statement{DB: db}
    stmt.Parse(&interfaces.User{})
    userContentType := &interfaces.ContentType{BlueprintName: "user", ModelName: stmt.Schema.Table}
    db.Create(userContentType)
    return nil
}

func (m initial_1621680132) Down(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.Migrator().DropTable(interfaces.Permission{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(interfaces.User{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(interfaces.UserGroup{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(interfaces.OneTimeAction{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(interfaces.ContentType{})
    if err != nil {
        return err
    }
    var contentType interfaces.ContentType
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&interfaces.OneTimeAction{})
    db.Model(&interfaces.ContentType{}).Where(&interfaces.ContentType{BlueprintName: "user", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Delete(&contentType)
    stmt = &gorm.Statement{DB: db}
    stmt.Parse(&interfaces.User{})
    db.Model(&interfaces.ContentType{}).Where(&interfaces.ContentType{BlueprintName: "user", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Delete(&contentType)
    return nil
}

func (m initial_1621680132) Deps() []string {
    return []string{}
}
