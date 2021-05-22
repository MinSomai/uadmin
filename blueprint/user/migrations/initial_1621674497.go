package migrations

type initial_1621674497 struct {
}

func (m initial_1621674497) GetName() string {
    return "initial"
}

func (m initial_1621674497) GetId() int64 {
    return 1621674497
}

func (m initial_1621674497) Up() {
}

func (m initial_1621674497) Down() {
}

func (m initial_1621674497) Deps() []string {
    return make([]string, 0)
}
