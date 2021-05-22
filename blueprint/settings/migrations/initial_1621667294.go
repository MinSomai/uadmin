package migrations

type initial_1621667294 struct {
}

func (m initial_1621667294) GetName() string {
    return "initial"
}

func (m initial_1621667294) GetId() int64 {
    return 1621667294
}

func (m initial_1621667294) Up() {
}

func (m initial_1621667294) Down() {
}

func (m initial_1621667294) Deps() []string {
    return make([]string, 0)
}
