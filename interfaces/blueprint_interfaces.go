package interfaces

type TraverseMigrationResult struct {
	MigrationLeaf IMigrationLeaf
	Error error
}

type IBlueprint interface {
	GetName() string
	GetDescription() string
	GetMigrationRegistry() *MigrationRegistry
}

type IBlueprintRegistry interface {
	Iterate() <-chan IBlueprint
	GetByName(name string) IBlueprint
	Register(blueprint IBlueprint)
}

type Blueprint struct {
	Name string
	Description string
	MigrationRegistry *MigrationRegistry
}

func (b Blueprint) GetName() string {
	return b.Name
}

func (b Blueprint) GetDescription() string {
	return b.Description
}

func (b Blueprint) GetMigrationRegistry() *MigrationRegistry {
	return b.MigrationRegistry
}