package migrations

import "github.com/sergeyglazyrindev/uadmin/core"

type addinguse1623259185 struct {
}

func (m addinguse1623259185) GetName() string {
	return "user.1623259185"
}

func (m addinguse1623259185) GetID() int64 {
	return 1623259185
}

func (m addinguse1623259185) Up(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m addinguse1623259185) Down(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m addinguse1623259185) Deps() []string {
	return []string{"user.1621680132"}
}
