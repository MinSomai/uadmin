package interfaces

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/config"
	"sort"
)

type TraverseMigrationResult struct {
	Node  IMigrationNode
	Error error
}

type IBlueprint interface {
	GetName() string
	GetDescription() string
	GetMigrationRegistry() IMigrationRegistry
	InitRouter(group *gin.RouterGroup)
	Init(config *config.UadminConfig)
}

type IBlueprintRegistry interface {
	Iterate() <-chan IBlueprint
	GetByName(name string) (IBlueprint, error)
	Register(blueprint IBlueprint)
	GetMigrationTree() IMigrationTree
	TraverseMigrations() <- chan *TraverseMigrationResult
	TraverseMigrationsDownTo(downToMigration string) <- chan *TraverseMigrationResult
	InitializeRouting(router *gin.Engine)
	Initialize(config *config.UadminConfig)
}

type Blueprint struct {
	Name string
	Description string
	MigrationRegistry IMigrationRegistry
}

func (b Blueprint) GetName() string {
	return b.Name
}

func (b Blueprint) InitRouter(group *gin.RouterGroup) {
	panic(fmt.Errorf("has to be redefined in concrete blueprint"))
}

func (b Blueprint) GetDescription() string {
	return b.Description
}

func (b Blueprint) Init(config *config.UadminConfig) {

}

func (b Blueprint) GetMigrationRegistry() IMigrationRegistry {
	return b.MigrationRegistry
}

type BlueprintRegistry struct {
	RegisteredBlueprints map[string]IBlueprint
	MigrationTree        IMigrationTree
}

func (r BlueprintRegistry) Iterate() <-chan IBlueprint {
	chnl := make(chan IBlueprint)
	go func() {
		for _, blueprint := range r.RegisteredBlueprints {
			chnl <- blueprint
		}
		// Ensure that at the end of the loop we close the channel!
		close(chnl)
	}()
	return chnl
}

func (r BlueprintRegistry) GetByName(name string) (IBlueprint, error) {
	blueprint, ok := r.RegisteredBlueprints[name]
	var err error
	if !ok {
		err = fmt.Errorf("Couldn't find blueprint with name %s", name)
	}
	return blueprint, err
}

func (r BlueprintRegistry) Register(blueprint IBlueprint) {
	r.RegisteredBlueprints[blueprint.GetName()] = blueprint
}

func (r BlueprintRegistry) GetMigrationTree() IMigrationTree {
	return r.MigrationTree
}

func (r BlueprintRegistry) buildMigrationTree(chnl chan *TraverseMigrationResult) bool {
	if r.MigrationTree.IsTreeBuilt() {
		return true
	}
	var err error
	_, err = r.GetByName("user")
	if err != nil {
		res := &TraverseMigrationResult{
			Node:  nil,
			Error: err,
		}
		chnl <- res
		return false
	}
	var currentMigration IMigration
	var tmpMigration IMigration
	var previousNode IMigrationNode
	var currentNode IMigrationNode
	var tmpNode IMigrationNode
	var blueprintName string
	var blueprintTmp IBlueprint
	for blueprint := range r.Iterate() {
		numberOfMigrationsWithNoDeps := mapset.NewSet()
		for i, migration := range blueprint.GetMigrationRegistry().GetSortedMigrations() {
			currentMigration = migration
			if i == 0 {
				if numberOfMigrationsWithNoDeps.Cardinality() >= 1 {
					numberOfMigrationsWithNoDeps.Add(currentMigration.GetName())
					res := &TraverseMigrationResult{
						Node:  nil,
						Error: fmt.Errorf("Found two migrations with no deps : %v", numberOfMigrationsWithNoDeps),
					}
					chnl <- res
					return false
				}
				currentNode, err = r.MigrationTree.GetNodeByMigrationName(currentMigration.GetName())
				if currentNode == nil {
					currentNode = NewMigrationNode(
						r.MigrationTree.GetRoot(), currentMigration, nil,
					)
				}
				err = r.MigrationTree.AddNode(currentNode)
				if err != nil {
					res := &TraverseMigrationResult{
						Node:  nil,
						Error: err,
					}
					chnl <- res
					return false
				}
				r.MigrationTree.GetRoot().AddChild(
					currentNode,
				)
				numberOfMigrationsWithNoDeps.Add(currentMigration.GetName())
			}
			if i != 0 {
				currentNode, err = r.MigrationTree.GetNodeByMigrationName(currentMigration.GetName())
				if currentNode == nil {
					currentNode = NewMigrationNode(
						previousNode, currentMigration, nil,
					)
				}
				// previousNode.AddChild(currentNode)
				err = r.MigrationTree.AddNode(currentNode)
				if err != nil {
					res := &TraverseMigrationResult{
						Node:  nil,
						Error: err,
					}
					chnl <- res
					return false
				}
			}
			for _, dep := range currentMigration.Deps() {
				blueprintName = GetBluePrintNameFromMigrationName(dep)
				if blueprintName == blueprint.GetName() {
					tmpMigration, err = blueprint.GetMigrationRegistry().GetByName(dep)
					if err != nil {
						res := &TraverseMigrationResult{
							Node:  nil,
							Error: err,
						}
						chnl <- res
						return false
					}
				} else {
					blueprintTmp, err = r.GetByName(blueprintName)
					if err != nil {
						res := &TraverseMigrationResult{
							Node:  nil,
							Error: err,
						}
						chnl <- res
						return false
					}
					tmpMigration, err = blueprintTmp.GetMigrationRegistry().GetByName(dep)
					if err != nil {
						res := &TraverseMigrationResult{
							Node:  nil,
							Error: err,
						}
						chnl <- res
						return false
					}
				}
				tmpNode, err = r.MigrationTree.GetNodeByMigrationName(tmpMigration.GetName())
				if tmpNode == nil {
					tmpNode = NewMigrationNode(
						previousNode, tmpMigration, nil,
					)
				}
				currentNode.AddDep(tmpNode)
				tmpNode.AddChild(currentNode)
			}
			previousNode = currentNode
		}
	}
	for blueprint := range r.Iterate() {
		numberOfMigrationsWithNoDescendants := mapset.NewSet()
		for _, migration := range blueprint.GetMigrationRegistry().GetSortedMigrations() {
			currentNode, err = r.MigrationTree.GetNodeByMigrationName(migration.GetName())
			if err != nil {

			}
			if currentNode.GetChildrenCount() == 0 {
				numberOfMigrationsWithNoDescendants.Add(migration.GetName())
			}
		}
		if numberOfMigrationsWithNoDescendants.Cardinality() >= 2 {
			res := &TraverseMigrationResult{
				Node:  nil,
				Error: fmt.Errorf("Found two or more migrations with no children from the same blueprint: %v", numberOfMigrationsWithNoDescendants),
			}
			chnl <- res
			return false
		}
	}
	return true
}

func (r BlueprintRegistry) TraverseMigrations() <- chan *TraverseMigrationResult {
	chnl := make(chan *TraverseMigrationResult)
	go func() {
		wasTreeBuilt := r.buildMigrationTree(chnl)
		if !wasTreeBuilt {
			close(chnl)
			return
		}
		r.MigrationTree.TreeBuilt()
		applyMigrationsInOrder := make([]string, 0)
		allBlueprintRoots := r.MigrationTree.GetRoot().GetChildren()
		for l := allBlueprintRoots.Front(); l != nil; l = l.Next() {
			applyMigrationsInOrder = append(applyMigrationsInOrder, l.Value.(IMigrationNode).GetMigration().GetName())
			migrationDepList := l.Value.(IMigrationNode).TraverseDeps(applyMigrationsInOrder, make(MigrationDepList, 0))
			sort.Reverse(migrationDepList)
			for _, m := range migrationDepList {
				applyMigrationsInOrder = append(applyMigrationsInOrder, m)
			}
			applyMigrationsInOrder = l.Value.(IMigrationNode).TraverseChildren(applyMigrationsInOrder)
		}
		for _, migrationName := range applyMigrationsInOrder {
			node, err := r.MigrationTree.GetNodeByMigrationName(migrationName)
			if err != nil {
				res := &TraverseMigrationResult{
					Node:  nil,
					Error: fmt.Errorf("Not found migration node with name : %s", migrationName),
				}
				chnl <- res
				close(chnl)
				return
			}
			res := &TraverseMigrationResult{
				Node:  node,
				Error: nil,
			}
			chnl <- res
		}
		close(chnl)
	}()
	return chnl
}

func (r BlueprintRegistry) InitializeRouting(router *gin.Engine) {
	for blueprint := range r.Iterate() {
		routergroup := router.Group("/" + blueprint.GetName())
		blueprint.InitRouter(routergroup)
	}
	router.GET( "/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func (r BlueprintRegistry) Initialize(config *config.UadminConfig) {
	for blueprint := range r.Iterate() {
		blueprint.Init(config)
	}
}

func (r BlueprintRegistry) TraverseMigrationsDownTo(downToMigration string) <- chan *TraverseMigrationResult {
	chnl := make(chan *TraverseMigrationResult)
	go func() {
		wasTreeBuilt := r.buildMigrationTree(chnl)
		if !wasTreeBuilt {
			close(chnl)
			return
		}
		applyMigrationsInOrder := make([]string, 0)
		allBlueprintRoots := r.MigrationTree.GetRoot().GetChildren()
		for l := allBlueprintRoots.Front(); l != nil; l = l.Next() {
			applyMigrationsInOrder = append(applyMigrationsInOrder, l.Value.(IMigrationNode).GetMigration().GetName())
			migrationDepList := l.Value.(IMigrationNode).TraverseDeps(applyMigrationsInOrder, make(MigrationDepList, 0))
			sort.Reverse(migrationDepList)
			for _, m := range migrationDepList {
				applyMigrationsInOrder = append(applyMigrationsInOrder, m)
			}
			applyMigrationsInOrder = l.Value.(IMigrationNode).TraverseChildren(applyMigrationsInOrder)
		}
		downgradeMigrationsInOrder := make(MigrationDepList, 0)
		var foundMigration = false
		for _, migrationName := range applyMigrationsInOrder {
			if len(downToMigration) > 0 && migrationName == downToMigration {
				foundMigration = true
			}
			if len(downToMigration) > 0 && !foundMigration {
				continue
			}
			downgradeMigrationsInOrder = append(downgradeMigrationsInOrder, migrationName)
		}
		sort.Reverse(downgradeMigrationsInOrder)
		for _, migrationName := range downgradeMigrationsInOrder {
			node, err := r.MigrationTree.GetNodeByMigrationName(migrationName)
			if err != nil {
				res := &TraverseMigrationResult{
					Node:  nil,
					Error: fmt.Errorf("Not found migration node with name : %s", migrationName),
				}
				chnl <- res
				close(chnl)
				return
			}
			res := &TraverseMigrationResult{
				Node:  node,
				Error: nil,
			}
			chnl <- res
		}
		close(chnl)
	}()
	return chnl
}

func NewBlueprintRegistry() IBlueprintRegistry {
	return &BlueprintRegistry{
		RegisteredBlueprints: make(map[string]IBlueprint),
		MigrationTree: NewMigrationTree(),
	}
}