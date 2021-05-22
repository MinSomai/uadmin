package migrations

type initial_1621667357 struct {
}

func (m initial_1621667357) GetName() string {
    return "initial"
}

func (m initial_1621667357) GetId() int64 {
    return 1621667357
}

func (m initial_1621667357) Up() {
}

func (m initial_1621667357) Down() {
}

func (m initial_1621667357) Deps() []string {
    return make([]string, 0)
}
