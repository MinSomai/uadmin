package migrations

type initial_1621667246 struct {
}

func (m initial_1621667246) GetName() string {
    return "initial"
}

func (m initial_1621667246) GetId() int64 {
    return 1621667246
}

func (m initial_1621667246) Up() {
}

func (m initial_1621667246) Down() {
}

func (m initial_1621667246) Deps() []string {
    return make([]string, 0)
}
