package migrations

type initial_1621667383 struct {
}

func (m initial_1621667383) GetName() string {
    return "initial"
}

func (m initial_1621667383) GetId() int64 {
    return 1621667383
}

func (m initial_1621667383) Up() {
}

func (m initial_1621667383) Down() {
}

func (m initial_1621667383) Deps() []string {
    return make([]string, 0)
}
