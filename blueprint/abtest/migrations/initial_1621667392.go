package migrations

type initial_1621667392 struct {
}

func (m initial_1621667392) GetName() string {
    return "initial"
}

func (m initial_1621667392) GetId() int64 {
    return 1621667392
}

func (m initial_1621667392) Up() {
}

func (m initial_1621667392) Down() {
}

func (m initial_1621667392) Deps() []string {
    return make([]string, 0)
}
