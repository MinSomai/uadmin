package interfaces

type IMigrationRegistry interface {
	FindMigrations() <-chan IMigration
}

type IBlueprint interface {
	GetName() string
	GetMigrationRegistry() IMigrationRegistry
}

type IBlueprintRegistry interface {
	Iterate() <-chan IBlueprint
	GetByName(name string) IBlueprint
	Register(blueprint IBlueprint)
}

type ICommand interface {
	Proceed()
	ParseArgs()
	GetHelpText() string
}

type IMigration interface {
	Up()
	Down()
	GetName() string
	GetId() int64
	Deps() []string
}
