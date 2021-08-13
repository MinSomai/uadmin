package migrations

import (
    "github.com/uadmin/uadmin/interfaces"
)

type initial_1623582035 struct {
}

func (m initial_1623582035) GetName() string {
    return "auth.1623582035"
}

func (m initial_1623582035) GetId() int64 {
    return 1623582035
}

func (m initial_1623582035) Up(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.AutoMigrate(interfaces.UserAuthToken{})
    if err != nil {
        return err
    }
    return nil
}

func (m initial_1623582035) Down(uadminDatabase *interfaces.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.Migrator().DropTable(interfaces.UserAuthToken{})
    if err != nil {
        return err
    }
    return nil
}

func (m initial_1623582035) Deps() []string {
    return []string{"user.1621680132"}
}
