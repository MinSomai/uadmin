package interfaces

type IMigrationRegistry interface {
	FindMigrations() <-chan IMigration
}

type IBlueprint interface {
	GetName() string
	GetDescription() string
	GetMigrationRegistry() IMigrationRegistry
}

type IBlueprintRegistry interface {
	Iterate() <-chan IBlueprint
	GetByName(name string) IBlueprint
	Register(blueprint IBlueprint)
}

type ICommand interface {
	Proceed(subaction string, args []string)
	GetHelpText() string
}

type IMigration interface {
	Up()
	Down()
	GetName() string
	GetId() int64
	Deps() []string
}
