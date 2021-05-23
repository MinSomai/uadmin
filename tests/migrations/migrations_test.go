package migrations

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildTree(t *testing.T) {
	for blueprint := range ConcreteBlueprintRegistry.Iterate() {
		err := blueprint.GetMigrationRegistry().BuildTree()
		if err != nil {
			assert.Equal(t, false, true, err)
			return
		}
	}
}

func TestDetermineConflictsInMigrations(t *testing.T) {
	err := BlueprintWithConflicts.GetMigrationRegistry().BuildTree()
	if err != nil {
		assert.Equal(t, false, true, err)
		return
	}
	conflicts := BlueprintWithConflicts.GetMigrationRegistry().FindPotentialConflictsForBlueprint(
		BlueprintWithConflicts.GetName(),
	)
	assert.Equal(t, len(conflicts), 2)
}

func TestBuildTreeForBlueprintWithNoMigrations(t *testing.T) {
	err := BlueprintWithNoMigrations.GetMigrationRegistry().BuildTree()
	if err != nil {
		assert.Equal(t, false, true, err)
		return
	}
}
