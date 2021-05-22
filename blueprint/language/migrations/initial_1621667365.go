package migrations

type initial_1621667365 struct {
}

func (m initial_1621667365) GetName() string {
    return "initial"
}

func (m initial_1621667365) GetId() int64 {
    return 1621667365
}

func (m initial_1621667365) Up() {
}

func (m initial_1621667365) Down() {
}

func (m initial_1621667365) Deps() []string {
    return make([]string, 0)
}
