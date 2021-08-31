package migrations

import "github.com/uadmin/uadmin/core"

type adding_use_1623259185 struct {
}

func (m adding_use_1623259185) GetName() string {
	return "user.1623259185"
}

func (m adding_use_1623259185) GetId() int64 {
	return 1623259185
}

func (m adding_use_1623259185) Up(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m adding_use_1623259185) Down(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m adding_use_1623259185) Deps() []string {
	return []string{"user.1621680132"}
}
