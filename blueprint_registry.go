package uadmin

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/uadmin/uadmin/interfaces"
)

type BlueprintRegistry struct {
	RegisteredBlueprints map[string]interfaces.IBlueprint
}

func (r BlueprintRegistry) Iterate() <-chan interfaces.IBlueprint {
	chnl := make(chan interfaces.IBlueprint)
	go func() {
		for _, blueprint := range r.RegisteredBlueprints {
			chnl <- blueprint
		}
		// Ensure that at the end of the loop we close the channel!
		close(chnl)
	}()
	return chnl
}

func (r BlueprintRegistry) GetByName(name string) (interfaces.IBlueprint, error) {
	blueprint, ok := r.RegisteredBlueprints[name]
	var err error
	if !ok {
		err = fmt.Errorf("Couldn't find blueprint with name %s", name)
	}
	return blueprint, err
}

func (r BlueprintRegistry) Register(blueprint interfaces.IBlueprint) {
	r.RegisteredBlueprints[blueprint.GetName()] = blueprint
}

func (r BlueprintRegistry) traverseMigrations() <- chan *interfaces.TraverseMigrationResult {
	chnl := make(chan *interfaces.TraverseMigrationResult)
	go func() {
		appliedMigrations := ensureDatabaseIsReadyForMigrationsAndReadAllApplied()
		appliedMigrationNames := mapset.NewSet()
		for _, migration := range appliedMigrations {
			appliedMigrationNames.Add(migration.MigrationName)
		}
		for blueprint := range appInstance.BlueprintRegistry.Iterate() {
			err := blueprint.GetMigrationRegistry().BuildTree(blueprint.GetName())
			if err != nil {
				res := &interfaces.TraverseMigrationResult{
					MigrationLeaf: interfaces.IMigrationLeaf(nil),
					Error: err,
				}
				chnl <- res
				close(chnl)
				return
			}
		}
		close(chnl)
		//db := appInstance.Database.ConnectTo("default")
		//var migrationName string
		//userBlueprint, err := appInstance.BlueprintRegistry.GetByName("user")
		//if err != nil {
		//	res := &interfaces.TraverseMigrationResult{
		//		MigrationLeaf: nil,
		//		Error: err,
		//	}
		//	chnl <- res
		//	close(chnl)
		//	return
		//}
		//for migrationToApply := range userBlueprint.GetMigrationRegistry().MigrationTree.TraverseInOrder() {
		//	migrationName = migrationToApply.GetMigration().GetName()
		//	// check for some stupid mistake in traversal algo
		//	if appliedMigrationNames.Contains(migrationName) {
		//		continue
		//	}
		//	migrationToApply.GetMigration().Up()
		//	appliedMigrationNames.Add(migrationName)
		//	db = db.Create(&Migration{
		//		MigrationName: migrationName,
		//		AppliedAt: time.Now(),
		//	})
		//	if db.Error != nil {
		//		return db.Error
		//	}
		//}
		//for blueprint := range appInstance.BlueprintRegistry.Iterate() {
		//	if blueprint.GetName() == "user" {
		//		continue
		//	}
		//	for migrationToApply := range blueprint.GetMigrationRegistry().MigrationTree.TraverseInOrder() {
		//		migrationName = migrationToApply.GetMigration().GetName()
		//		// check for some stupid mistake in traversal algo
		//		if appliedMigrationNames.Contains(migrationName) {
		//			continue
		//		}
		//		migrationToApply.GetMigration().Up()
		//		appliedMigrationNames.Add(migrationName)
		//		db = db.Create(&Migration{
		//			MigrationName: migrationName,
		//			AppliedAt: time.Now(),
		//		})
		//		if db.Error != nil {
		//			return db.Error
		//		}
		//	}
		//	//if err != nil {
		//	//	return err
		//	//}
		//}
	}()
	return chnl
}