package migrations

type initial_1621667305 struct {
}

func (m initial_1621667305) GetName() string {
    return "initial"
}

func (m initial_1621667305) GetId() int64 {
    return 1621667305
}

func (m initial_1621667305) Up() {
}

func (m initial_1621667305) Down() {
}

func (m initial_1621667305) Deps() []string {
    return make([]string, 0)
}
