package migrations

import (
	"github.com/uadmin/uadmin/core"
)

type create_all_1623263607 struct {
}

func (m create_all_1623263607) GetName() string {
    return "language.1623263607"
}

func (m create_all_1623263607) GetId() int64 {
    return 1623263607
}

func (m create_all_1623263607) Up(uadminDatabase *core.UadminDatabase) error {
    langs := [][]string{
        {"English", "English", "en"},
    }
    db := uadminDatabase.Db
    tx := db
    for _, lang := range langs {
        l := core.Language{
            EnglishName: lang[0],
            Name:        lang[1],
            Code:        lang[2],
            Active:      false,
        }
        if l.Code == "en" {
            l.AvailableInGui = true
            l.Active = true
            l.Default = true
        }
        tx.Create(&l)
    }
    return nil
}

func (m create_all_1623263607) Down(uadminDatabase *core.UadminDatabase) error {
    db := uadminDatabase.Db
    db.Unscoped().Where("1 = 1").Delete(&core.Language{Code: "en"})
    return nil
}

func (m create_all_1623263607) Deps() []string {
    return []string{"language.1623083053"}
}
