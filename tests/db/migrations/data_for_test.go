package migrations

import (
	"github.com/uadmin/uadmin/core"
)

type initial_1621667392 struct {
}

func (m initial_1621667392) GetName() string {
	return "user.1621667393"
}

func (m initial_1621667392) GetId() int64 {
	return 1621667392
}

func (m initial_1621667392) Up(uadminDatabase *core.UadminDatabase) error {
	appliedMigrations = append(appliedMigrations, m.GetName())
	return nil
}

func (m initial_1621667392) Down(uadminDatabase *core.UadminDatabase) error {
	appliedMigrations = core.Remove(appliedMigrations, m.GetName())
	return nil
}

func (m initial_1621667392) Deps() []string {
	return make([]string, 0)
}

type migration2_1621667393 struct {
}

func (m migration2_1621667393) GetName() string {
	return "user.1621680132"
}

func (m migration2_1621667393) GetId() int64 {
	return 1621667393
}

func (m migration2_1621667393) Up(uadminDatabase *core.UadminDatabase) error {
	appliedMigrations = append(appliedMigrations, m.GetName())
	return nil
}

func (m migration2_1621667393) Down(uadminDatabase *core.UadminDatabase) error {
	appliedMigrations = core.Remove(appliedMigrations, m.GetName())
	return nil
}

func (m migration2_1621667393) Deps() []string {
	return []string{"user.1621667393"}
}

type initialtest1_1621667392 struct {
}

func (m initialtest1_1621667392) GetName() string {
	return "test1.1621667393"
}

func (m initialtest1_1621667392) GetId() int64 {
	return 1621667392
}

func (m initialtest1_1621667392) Up(uadminDatabase *core.UadminDatabase) error {
	appliedMigrations = append(appliedMigrations, m.GetName())
	return nil
}

func (m initialtest1_1621667392) Down(uadminDatabase *core.UadminDatabase) error {
	appliedMigrations = core.Remove(appliedMigrations, m.GetName())
	return nil
}

func (m initialtest1_1621667392) Deps() []string {
	return make([]string, 0)
}

type migration2test_1621667393 struct {
}

func (m migration2test_1621667393) GetName() string {
	return "test1.1621680132"
}

func (m migration2test_1621667393) GetId() int64 {
	return 1621667393
}

func (m migration2test_1621667393) Up(uadminDatabase *core.UadminDatabase) error {
	appliedMigrations = append(appliedMigrations, m.GetName())
	return nil
}

func (m migration2test_1621667393) Down(uadminDatabase *core.UadminDatabase) error {
	appliedMigrations = core.Remove(appliedMigrations, m.GetName())
	return nil
}

func (m migration2test_1621667393) Deps() []string {
	return []string{"test1.1621667393", "user.1621680132"}
}

type initialblueprintconflicts_1621667392 struct {
}

func (m initialblueprintconflicts_1621667392) GetName() string {
	return "user.1621680132"
}

func (m initialblueprintconflicts_1621667392) GetId() int64 {
	return 1621667392
}

func (m initialblueprintconflicts_1621667392) Up(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m initialblueprintconflicts_1621667392) Down(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m initialblueprintconflicts_1621667392) Deps() []string {
	return make([]string, 0)
}

type migration2blueprintconflicts_1621667392 struct {
}

func (m migration2blueprintconflicts_1621667392) GetName() string {
	return "user.16216801321"
}

func (m migration2blueprintconflicts_1621667392) GetId() int64 {
	return 1621667393
}

func (m migration2blueprintconflicts_1621667392) Up(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m migration2blueprintconflicts_1621667392) Down(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m migration2blueprintconflicts_1621667392) Deps() []string {
	return []string{"user.1621680132"}
}

type migration3blueprintconflicts_1621667392 struct {
}

func (m migration3blueprintconflicts_1621667392) GetName() string {
	return "user.16216801341"
}

func (m migration3blueprintconflicts_1621667392) GetId() int64 {
	return 1621667394
}

func (m migration3blueprintconflicts_1621667392) Up(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m migration3blueprintconflicts_1621667392) Down(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m migration3blueprintconflicts_1621667392) Deps() []string {
	return []string{"user.16216801321"}
}

type migration4blueprintconflicts_1621667392 struct {
}

func (m migration4blueprintconflicts_1621667392) GetName() string {
	return "user.16216801381"
}

func (m migration4blueprintconflicts_1621667392) GetId() int64 {
	return 1621667395
}

func (m migration4blueprintconflicts_1621667392) Up(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m migration4blueprintconflicts_1621667392) Down(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m migration4blueprintconflicts_1621667392) Deps() []string {
	return []string{"user.16216801321"}
}

type nodeps1_1621667392 struct {
}

func (m nodeps1_1621667392) GetName() string {
	return "test1.1621680132"
}

func (m nodeps1_1621667392) GetId() int64 {
	return 1621667392
}

func (m nodeps1_1621667392) Up(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m nodeps1_1621667392) Down(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m nodeps1_1621667392) Deps() []string {
	return make([]string, 0)
}

type nodeps2_1621667392 struct {
}

func (m nodeps2_1621667392) GetName() string {
	return "test1.16216801321"
}

func (m nodeps2_1621667392) GetId() int64 {
	return 1621667392
}

func (m nodeps2_1621667392) Up(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m nodeps2_1621667392) Down(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m nodeps2_1621667392) Deps() []string {
	return make([]string, 0)
}

type loopedmigration1_1621667392 struct {
}

func (m loopedmigration1_1621667392) GetName() string {
	return "user.1621680132"
}

func (m loopedmigration1_1621667392) GetId() int64 {
	return 1621667392
}

func (m loopedmigration1_1621667392) Up(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m loopedmigration1_1621667392) Down(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m loopedmigration1_1621667392) Deps() []string {
	return []string{"user1.16216801321"}
}

type loopedmigration2_1621667392 struct {
}

func (m loopedmigration2_1621667392) GetName() string {
	return "user1.16216801321"
}

func (m loopedmigration2_1621667392) GetId() int64 {
	return 1621667392
}

func (m loopedmigration2_1621667392) Up(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m loopedmigration2_1621667392) Down(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m loopedmigration2_1621667392) Deps() []string {
	return []string{"user.1621680132"}
}

type samenamemigration1_1621667392 struct {
}

func (m samenamemigration1_1621667392) GetName() string {
	return "user.1621680132"
}

func (m samenamemigration1_1621667392) GetId() int64 {
	return 1621667392
}

func (m samenamemigration1_1621667392) Up(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m samenamemigration1_1621667392) Down(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m samenamemigration1_1621667392) Deps() []string {
	return make([]string, 0)
}

type samenamemigration2_1621667392 struct {
}

func (m samenamemigration2_1621667392) GetName() string {
	return "user.1621680132"
}

func (m samenamemigration2_1621667392) GetId() int64 {
	return 1621667392
}

func (m samenamemigration2_1621667392) Up(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m samenamemigration2_1621667392) Down(uadminDatabase *core.UadminDatabase) error {
	return nil
}

func (m samenamemigration2_1621667392) Deps() []string {
	return make([]string, 0)
}

var TestBlueprintMigrationRegistry *core.MigrationRegistry
var Test1BlueprintMigrationRegistry *core.MigrationRegistry
var TestBlueprint core.Blueprint

var Test1Blueprint core.Blueprint
var BlueprintWithConflictsMigrationRegistry *core.MigrationRegistry
var BlueprintWithConflicts core.Blueprint

var BlueprintWithNoMigrationsRegistry *core.MigrationRegistry
var BlueprintWithNoMigrations core.Blueprint

var BlueprintWithTwoSameDeps core.Blueprint

var Blueprint1WithLoopedMigrations core.Blueprint
var Blueprint2WithLoopedMigrations core.Blueprint

var Blueprint1WithSameMigrationNames core.Blueprint
var Blueprint2WithSameMigrationNames core.Blueprint

var appliedMigrations = make([]string, 0)

func init() {
	TestBlueprintMigrationRegistry = core.NewMigrationRegistry()
	TestBlueprintMigrationRegistry.AddMigration(initial_1621667392{})
	TestBlueprintMigrationRegistry.AddMigration(migration2_1621667393{})

	Test1BlueprintMigrationRegistry = core.NewMigrationRegistry()
	Test1BlueprintMigrationRegistry.AddMigration(initialtest1_1621667392{})
	Test1BlueprintMigrationRegistry.AddMigration(migration2test_1621667393{})

	TestBlueprint = core.Blueprint{
		Name:              "user",
		Description:       "this blueprint for testing",
		MigrationRegistry: TestBlueprintMigrationRegistry,
	}
	Test1Blueprint = core.Blueprint{
		Name:              "test1",
		Description:       "this test1 blueprint for testing",
		MigrationRegistry: Test1BlueprintMigrationRegistry,
	}

	BlueprintWithTwoSameDeps = core.Blueprint{
		Name:              "user",
		Description:       "this blueprint for testing",
		MigrationRegistry: core.NewMigrationRegistry(),
	}
	BlueprintWithTwoSameDeps.GetMigrationRegistry().AddMigration(nodeps1_1621667392{})
	BlueprintWithTwoSameDeps.GetMigrationRegistry().AddMigration(nodeps2_1621667392{})

	BlueprintWithConflictsMigrationRegistry = core.NewMigrationRegistry()
	BlueprintWithConflictsMigrationRegistry.AddMigration(initialblueprintconflicts_1621667392{})
	BlueprintWithConflictsMigrationRegistry.AddMigration(migration2blueprintconflicts_1621667392{})
	BlueprintWithConflictsMigrationRegistry.AddMigration(migration3blueprintconflicts_1621667392{})
	BlueprintWithConflictsMigrationRegistry.AddMigration(migration4blueprintconflicts_1621667392{})
	BlueprintWithConflicts = core.Blueprint{
		Name:              "user",
		Description:       "blueprint with conflicts",
		MigrationRegistry: BlueprintWithConflictsMigrationRegistry,
	}

	BlueprintWithNoMigrationsRegistry = core.NewMigrationRegistry()
	BlueprintWithNoMigrations = core.Blueprint{
		Name:              "user",
		Description:       "blueprint with no migrations",
		MigrationRegistry: BlueprintWithNoMigrationsRegistry,
	}

	Blueprint1WithLoopedMigrations = core.Blueprint{
		Name:              "user",
		Description:       "blueprint with looped migrations 1",
		MigrationRegistry: core.NewMigrationRegistry(),
	}
	Blueprint1WithLoopedMigrations.GetMigrationRegistry().AddMigration(loopedmigration1_1621667392{})
	Blueprint2WithLoopedMigrations = core.Blueprint{
		Name:              "user1",
		Description:       "blueprint with looped migrations 2",
		MigrationRegistry: core.NewMigrationRegistry(),
	}
	Blueprint2WithLoopedMigrations.GetMigrationRegistry().AddMigration(loopedmigration2_1621667392{})

	Blueprint1WithSameMigrationNames = core.Blueprint{
		Name:              "user",
		Description:       "blueprint with same migration names",
		MigrationRegistry: core.NewMigrationRegistry(),
	}
	Blueprint2WithSameMigrationNames = core.Blueprint{
		Name:              "user1",
		Description:       "blueprint with same migration names",
		MigrationRegistry: core.NewMigrationRegistry(),
	}
	Blueprint1WithSameMigrationNames.GetMigrationRegistry().AddMigration(samenamemigration1_1621667392{})
	Blueprint2WithSameMigrationNames.GetMigrationRegistry().AddMigration(samenamemigration2_1621667392{})
}
