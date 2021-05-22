package migrations

type initial_1621667330 struct {
}

func (m initial_1621667330) GetName() string {
    return "initial"
}

func (m initial_1621667330) GetId() int64 {
    return 1621667330
}

func (m initial_1621667330) Up() {
}

func (m initial_1621667330) Down() {
}

func (m initial_1621667330) Deps() []string {
    return make([]string, 0)
}
