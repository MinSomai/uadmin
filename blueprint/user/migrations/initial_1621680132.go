package migrations

import (
	"github.com/uadmin/uadmin/core"
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

func (m initial_1621680132) Up(uadminDatabase *core.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.AutoMigrate(core.ContentType{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(core.UserGroup{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(core.User{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(core.Permission{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(core.OneTimeAction{})
    if err != nil {
        return err
    }
    return nil
}

func (m initial_1621680132) Down(uadminDatabase *core.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.Migrator().DropTable(core.Permission{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(core.User{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(core.UserGroup{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(core.OneTimeAction{})
    if err != nil {
        return err
    }
    err = db.Migrator().DropTable(core.ContentType{})
    if err != nil {
        return err
    }
    var contentType core.ContentType
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&core.OneTimeAction{})
    db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "user", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Unscoped().Delete(&contentType)
    stmt = &gorm.Statement{DB: db}
    stmt.Parse(&core.User{})
    db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "user", ModelName: stmt.Schema.Table}).First(&contentType)
    db.Unscoped().Delete(&contentType)
    return nil
}

func (m initial_1621680132) Deps() []string {
    return []string{}
}
