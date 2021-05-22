package migrations

type MigrationName struct {
}

func (m MigrationName) GetName() string {
    return "concreteMigrationName"
}

func (m MigrationName) GetId() int64 {
    return concreteMigrationId
}

func (m MigrationName) Up() {
}

func (m MigrationName) Down() {
}

func (m MigrationName) Deps() []string {
    return dependencyId
}
