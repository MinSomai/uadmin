
package migrations

type initial_1621677848 struct {
}

func (m initial_1621677848) GetName() string {
    return "initial"
}

func (m initial_1621677848) GetId() int64 {
    return 1621677848
}

func (m initial_1621677848) Up() {
}

func (m initial_1621677848) Down() {
}

func (m initial_1621677848) Deps() []string {
    return make([]string, 0)
}
