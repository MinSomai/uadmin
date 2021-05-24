package migrations

import (
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/utils"
	"github.com/uadmin/uadmin"
)

type initial_1621667392 struct {
}

func (m initial_1621667392) GetName() string {
	return "test.1621680132"
}

func (m initial_1621667392) GetId() int64 {
	return 1621667392
}

func (m initial_1621667392) Up() {
}

func (m initial_1621667392) Down() {
}

func (m initial_1621667392) Deps() []string {
	return make([]string, 0)
}

func (m initial_1621667392) IsDependentFrom(dep string) bool {
	return utils.Contains(m.Deps(), dep)
}

type migration2_1621667393 struct {
}

func (m migration2_1621667393) GetName() string {
	return "migration2.1621667393"
}

func (m migration2_1621667393) GetId() int64 {
	return 1621667393
}

func (m migration2_1621667393) Up() {
}

func (m migration2_1621667393) Down() {
}

func (m migration2_1621667393) Deps() []string {
	return []string{"test.1621680132"}
}

func (m migration2_1621667393) IsDependentFrom(dep string) bool {
	return utils.Contains(m.Deps(), dep)
}

type initialtest1_1621667392 struct {
}

func (m initialtest1_1621667392) GetName() string {
	return "test1.1621680132"
}

func (m initialtest1_1621667392) GetId() int64 {
	return 1621667392
}

func (m initialtest1_1621667392) Up() {
}

func (m initialtest1_1621667392) Down() {
}

func (m initialtest1_1621667392) Deps() []string {
	return make([]string, 0)
}

func (m initialtest1_1621667392) IsDependentFrom(dep string) bool {
	return utils.Contains(m.Deps(), dep)
}

type migration2test_1621667393 struct {
}

func (m migration2test_1621667393) GetName() string {
	return "migration2.1621667393"
}

func (m migration2test_1621667393) GetId() int64 {
	return 1621667393
}

func (m migration2test_1621667393) Up() {
}

func (m migration2test_1621667393) Down() {
}

func (m migration2test_1621667393) Deps() []string {
	return []string{"test1.1621680132", "test.1621680132"}
}

func (m migration2test_1621667393) IsDependentFrom(dep string) bool {
	return utils.Contains(m.Deps(), dep)
}

type initialblueprintconflicts_1621667392 struct {
}

func (m initialblueprintconflicts_1621667392) GetName() string {
	return "blueprint_conflicts.1621680132"
}

func (m initialblueprintconflicts_1621667392) GetId() int64 {
	return 1621667392
}

func (m initialblueprintconflicts_1621667392) Up() {
}

func (m initialblueprintconflicts_1621667392) Down() {
}

func (m initialblueprintconflicts_1621667392) Deps() []string {
	return make([]string, 0)
}

func (m initialblueprintconflicts_1621667392) IsDependentFrom(dep string) bool {
	return utils.Contains(m.Deps(), dep)
}

type migration2blueprintconflicts_1621667392 struct {
}

func (m migration2blueprintconflicts_1621667392) GetName() string {
	return "blueprint_conflicts.16216801321"
}

func (m migration2blueprintconflicts_1621667392) GetId() int64 {
	return 1621667392
}

func (m migration2blueprintconflicts_1621667392) Up() {
}

func (m migration2blueprintconflicts_1621667392) Down() {
}

func (m migration2blueprintconflicts_1621667392) Deps() []string {
	return make([]string, 0)
}

func (m migration2blueprintconflicts_1621667392) IsDependentFrom(dep string) bool {
	return utils.Contains(m.Deps(), dep)
}

var TestBlueprintMigrationRegistry *interfaces.MigrationRegistry
var Test1BlueprintMigrationRegistry *interfaces.MigrationRegistry
var TestBlueprint interfaces.Blueprint

var Test1Blueprint interfaces.Blueprint
var BlueprintWithConflictsMigrationRegistry *interfaces.MigrationRegistry
var BlueprintWithConflicts interfaces.Blueprint

var BlueprintWithNoMigrationsRegistry *interfaces.MigrationRegistry
var BlueprintWithNoMigrations interfaces.Blueprint


var ConcreteBlueprintRegistry = uadmin.BlueprintRegistry{
	RegisteredBlueprints: make(map[string]interfaces.IBlueprint),
}

func init() {
	TestBlueprintMigrationRegistry = interfaces.NewMigrationRegistry()
	TestBlueprintMigrationRegistry.AddMigration(initial_1621667392{})
	TestBlueprintMigrationRegistry.AddMigration(migration2_1621667393{})

	Test1BlueprintMigrationRegistry = interfaces.NewMigrationRegistry()
	Test1BlueprintMigrationRegistry.AddMigration(initialtest1_1621667392{})
	Test1BlueprintMigrationRegistry.AddMigration(migration2test_1621667393{})

	TestBlueprint = interfaces.Blueprint{
		Name: "user",
		Description: "this blueprint for testing",
		MigrationRegistry: TestBlueprintMigrationRegistry,
	}
	Test1Blueprint = interfaces.Blueprint{
		Name: "test1",
		Description: "this test1 blueprint for testing",
		MigrationRegistry: Test1BlueprintMigrationRegistry,
	}

	ConcreteBlueprintRegistry.Register(TestBlueprint)
	ConcreteBlueprintRegistry.Register(Test1Blueprint)

	BlueprintWithConflictsMigrationRegistry = interfaces.NewMigrationRegistry()
	BlueprintWithConflictsMigrationRegistry.AddMigration(initialblueprintconflicts_1621667392{})
	BlueprintWithConflictsMigrationRegistry.AddMigration(migration2blueprintconflicts_1621667392{})
	BlueprintWithConflicts = interfaces.Blueprint{
		Name: "blueprint_conflicts",
		Description: "blueprint with conflicts",
		MigrationRegistry: BlueprintWithConflictsMigrationRegistry,
	}

	BlueprintWithNoMigrationsRegistry = interfaces.NewMigrationRegistry()
	BlueprintWithNoMigrations = interfaces.Blueprint{
		Name: "blueprint with no migrations",
		Description: "blueprint with conflicts",
		MigrationRegistry: BlueprintWithNoMigrationsRegistry,
	}
}