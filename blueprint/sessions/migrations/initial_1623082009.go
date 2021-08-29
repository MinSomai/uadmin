package migrations

import (
	"github.com/uadmin/uadmin/core"
)

type initial_1623082009 struct {
}

func (m initial_1623082009) GetName() string {
    return "sessions.1623082009"
}

func (m initial_1623082009) GetId() int64 {
    return 1623082009
}

func (m initial_1623082009) Up(uadminDatabase *core.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.AutoMigrate(core.Session{})
    if err != nil {
        return err
    }
    return nil
}

func (m initial_1623082009) Down(uadminDatabase *core.UadminDatabase) error {
    db := uadminDatabase.Db
    err := db.Migrator().DropTable(core.Session{})
    if err != nil {
        return err
    }
    return nil
}

func (m initial_1623082009) Deps() []string {
    return []string{"user.1621680132"}
}
