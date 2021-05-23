package interfaces

import (
	"container/list"
	"sort"
	"strings"
)

type IMigrationRegistry interface {
	FindMigrations() <-chan IMigration
}

type IMigration interface {
	Up()
	Down()
	GetName() string
	GetId() int64
	Deps() []string
	IsDependentFrom(dep string) bool
}

type IMigrationLeaf interface {
	GetAncestors() <-chan IMigrationLeaf
	GetDescendants() <-chan IMigrationLeaf
	AddAncestor(IMigrationLeaf)
	AddDescendant(IMigrationLeaf)
	GetMigration() IMigration
	DoesLeafHaveDescendants() bool
}

type IMigrationTree interface {
	GetRoot() IMigrationLeaf
	SetRoot(root *MigrationLeaf)
	Traverse() <-chan MigrationLeaf
	AddLeaf(IMigrationLeaf)
}

type MigrationLeaf struct {
	Migration IMigration
	IsRoot bool
	Ancestors *list.List
	Descendants *list.List
}

func getBluePrintNameFromMigrationName(migrationName string) string {
	return strings.Split(migrationName, ".")[0]
}

func (l MigrationLeaf) DoesLeafHaveDescendants() bool {
	return l.Descendants.Len() != 0
}

func (l MigrationLeaf) GetAncestors() <-chan IMigrationLeaf {
	chnl := make(chan IMigrationLeaf)
	go func() {
		for e := l.Ancestors.Front(); e != nil; e = e.Next() {
			chnl <- e.Value.(IMigrationLeaf)
		}
		close(chnl)
	}()
	return chnl
}

func (l MigrationLeaf) GetMigration() IMigration{
	return l.Migration
}

func (l MigrationLeaf) GetDescendants() <-chan IMigrationLeaf {
	chnl := make(chan IMigrationLeaf)
	go func() {
		for e := l.Descendants.Front(); e != nil; e = e.Next() {
			chnl <- e.Value.(IMigrationLeaf)
		}
		close(chnl)
	}()
	return chnl
}

func (l MigrationLeaf) AddAncestor(migrationLeaf IMigrationLeaf) {
	l.Ancestors.PushBack(migrationLeaf)
}

func (l MigrationLeaf) AddDescendant(migrationLeaf IMigrationLeaf) {
	l.Descendants.PushBack(migrationLeaf)
}

func NewMigrationLeaf(migration IMigration) *MigrationLeaf {
	return &MigrationLeaf{
		Migration: migration,
		Ancestors: list.New(),
		Descendants: list.New(),
	}
}

type MigrationTree struct {
	Root *MigrationLeaf
	Leafs *list.List
}

func (t MigrationTree) AddLeaf(leaf IMigrationLeaf) {
	t.Leafs.PushBack(leaf)
}

func (t MigrationTree) GetRoot() IMigrationLeaf {
	return t.Root
}

func (t MigrationTree) SetRoot(root *MigrationLeaf) {
	root.IsRoot = true
	*t.Root = *root
}

type MigrationList []IMigration

func (m MigrationList) Len() int { return len(m) }
func (m MigrationList) Less(i, j int) bool {
	return m[i].GetId() < m[j].GetId()
}
func (m MigrationList) Swap(i, j int){ m[i], m[j] = m[j], m[i] }

type MigrationRegistry struct {
	Migrations map[string]IMigration
	MigrationTree MigrationTree
}

func (r MigrationRegistry) AddMigration(migration IMigration) {
	r.Migrations[migration.GetName()] = migration
}

func (t MigrationTree) traverse() <-chan IMigrationLeaf {
	chnl := make(chan IMigrationLeaf)
	go func() {
		if t.Leafs.Len() == 0 {
			close(chnl)
			return
		}
		var currentLeaf = t.Leafs.Front()
		for true {
			if currentLeaf == nil {
				break
			}
			chnl <- currentLeaf.Value.(IMigrationLeaf)
			currentLeaf = currentLeaf.Next()
		}
		close(chnl)
	}()
	return chnl
}

func (t MigrationTree) findPotentialConflictsForBlueprint(blueprintName string) []string {
	var conflicts []string
	for mLeaf := range t.traverse() {
		migrationBlueprintName := getBluePrintNameFromMigrationName(mLeaf.GetMigration().GetName())
		if migrationBlueprintName != blueprintName {
			continue
		}
		if !mLeaf.DoesLeafHaveDescendants() {
			conflicts = append(conflicts, mLeaf.GetMigration().GetName())
		}
	}
	if len(conflicts) > 1 {
		return conflicts
	}
	return make([]string, 0)
}

func (r MigrationRegistry) BuildTree() error {
	if len(r.Migrations) == 0 {
		return nil
	}
	sortedMigrations := r.GetSortedMigrations()
	r.MigrationTree.SetRoot(NewMigrationLeaf(sortedMigrations[0]))
	currentLeaf := r.MigrationTree.GetRoot()
	r.MigrationTree.AddLeaf(currentLeaf)
	for _, migration := range sortedMigrations[1:] {
		migrationLeaf := NewMigrationLeaf(migration)
		r.MigrationTree.AddLeaf(migrationLeaf)
		if migration.IsDependentFrom(currentLeaf.GetMigration().GetName()) {
			currentLeaf.AddDescendant(migrationLeaf)
		}
		if currentLeaf.GetMigration().IsDependentFrom(migration.GetName()) {
			currentLeaf.AddAncestor(migrationLeaf)
		}
		currentLeaf = migrationLeaf
	}
	return nil
}

func (r MigrationRegistry) FindPotentialConflictsForBlueprint(blueprintName string) []string {
	return r.MigrationTree.findPotentialConflictsForBlueprint(blueprintName)
}


func (r MigrationRegistry) GetSortedMigrations() MigrationList {
	sortedMigrations := make(MigrationList, len(r.Migrations))
	i := 0
	for _, migration := range r.Migrations {
		sortedMigrations[i] = migration
		i += 1
	}
	sort.Sort(sortedMigrations)
	return sortedMigrations
}

func (r MigrationRegistry) FindMigrations() <-chan *IMigration{
	chnl := make(chan *IMigration)
	go func() {
		close(chnl)
	}()
	return chnl
}

func NewMigrationRegistry() *MigrationRegistry {
	return &MigrationRegistry{
		Migrations: make(map[string]IMigration),
		MigrationTree: MigrationTree{
			Root: &MigrationLeaf{
				Ancestors: list.New(),
				Descendants: list.New(),
			},
			Leafs: list.New(),
		},
	}
}