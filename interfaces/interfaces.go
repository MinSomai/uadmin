package interfaces

type IBlueprintInterface struct {
	Name string
}

type CommandInterface interface {
	Proceed()
	ParseArgs()
	GetHelpText() string
}

type MigrationInterface interface {
	Up()
	Down()
	GetName() string
	GetId() int64
	Deps() []string
}
