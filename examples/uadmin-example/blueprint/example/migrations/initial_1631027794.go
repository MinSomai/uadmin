package migrations

import (
	"github.com/sergeyglazyrindev/uadmin/core"
	"github.com/sergeyglazyrindev/uadmin_example/blueprint/example/models"
)

type initial1631027794 struct {
}

func (m initial1631027794) GetName() string {
	return "example.1631027794"
}

func (m initial1631027794) GetID() int64 {
	return 1631027794
}

func (m initial1631027794) Up(uadminDatabase *core.UadminDatabase) error {
	// create table for your data
	uadminDatabase.Db.AutoMigrate(&models.Todo{})
	return nil
}

func (m initial1631027794) Down(uadminDatabase *core.UadminDatabase) error {
	// remove table
	db := uadminDatabase.Db
	err := db.Migrator().DropTable(models.Todo{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1631027794) Deps() []string {
	return make([]string, 0)
}
