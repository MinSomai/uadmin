
package migrations

import (
    "github.com/uadmin/uadmin/utils"
)

type initial_1621680132 struct {
}

func (m initial_1621680132) GetName() string {
    return "user.initial"
}

func (m initial_1621680132) GetId() int64 {
    return 1621680132
}

func (m initial_1621680132) Up() {
}

func (m initial_1621680132) Down() {
}

func (m initial_1621680132) Deps() []string {
    return make([]string, 0)
}

func (m initial_1621680132) IsDependentFrom(dep string) bool {
    return utils.Contains(m.Deps(), dep)
}