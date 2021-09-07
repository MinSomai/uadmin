package migrations

import "github.com/sergeyglazyrindev/uadmin/core"

type initial_1631027794 struct {
}

func (m initial_1631027794) GetName() string {
    return "example.1631027794"
}

func (m initial_1631027794) GetID() int64 {
    return 1631027794
}

func (m initial_1631027794) Up(uadminDatabase *core.UadminDatabase) error {
    return nil
}

func (m initial_1631027794) Down(uadminDatabase *core.UadminDatabase) error {
    return nil
}

func (m initial_1631027794) Deps() []string {
    return make([]string, 0)
}
