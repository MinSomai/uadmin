package interfaces

import (
	"container/list"
	"fmt"
	"sort"
	"strings"
)

type IMigrationRegistry interface {
	GetByName(migrationName string) (IMigration, error)
	AddMigration(migration IMigration)
	GetSortedMigrations() MigrationList
}

type IMigration interface {
	Up(uadminDatabase *UadminDatabase) error
	Down(uadminDatabase *UadminDatabase) error
	GetName() string
	GetId() int64
	Deps() []string
}

type IMigrationNode interface {
	IsApplied() bool
	GetMigration() IMigration
	SetItAsRoot()
	IsRoot() bool
	AddChild(node IMigrationNode)
	AddDep(node IMigrationNode)
	GetChildrenCount() int
	GetChildren() *list.List
	GetDeps() *list.List
	TraverseDeps(migrationList []string, depList MigrationDepList) MigrationDepList
	TraverseChildren(migrationList []string) []string
	IsDummy() bool
	Downgrade(uadminDatabase *UadminDatabase) error
	Apply(uadminDatabase *UadminDatabase) error
}

type IMigrationTree interface {
	GetRoot() IMigrationNode
	SetRoot(root IMigrationNode)
	GetNodeByMigrationName(migrationName string) (IMigrationNode, error)
	AddNode(node IMigrationNode) error
	TreeBuilt()
	IsTreeBuilt() bool
}

func GetBluePrintNameFromMigrationName(migrationName string) string {
	return strings.Split(migrationName, ".")[0]
}

type MigrationNode struct {
	Deps *list.List
	Node IMigration
	Children *list.List
	applied bool
	isRoot bool
	dummy bool
}

func (n MigrationNode) IsDummy() bool {
	return n.dummy
}

func (n MigrationNode) Apply(uadminDatabase *UadminDatabase) error {
	res := n.Node.Up(uadminDatabase)
	if res == nil {
		n.applied = true
	}
	return res
}

func (n MigrationNode) Downgrade(uadminDatabase *UadminDatabase) error {
	res := n.Node.Down(uadminDatabase)
	if res == nil {
		n.applied = false
	}
	return res
}

func (n MigrationNode) GetMigration() IMigration {
	return n.Node
}

func (n MigrationNode) GetChildren() *list.List {
	return n.Children
}

func (n MigrationNode) GetDeps() *list.List {
	return n.Deps
}

func (n MigrationNode) GetChildrenCount() int {
	return n.Children.Len()
}

func (n MigrationNode) IsApplied() bool {
	return n.applied
}

func (n MigrationNode) SetItAsRoot() {
	n.isRoot = true
}

func (n MigrationNode) IsRoot() bool {
	return n.isRoot
}

func (n MigrationNode) AddChild(node IMigrationNode) {
	n.Children.PushBack(node)
}

func (n MigrationNode) AddDep(node IMigrationNode) {
	n.Deps.PushBack(node)
}

func (n MigrationNode) TraverseDeps(migrationList []string, depList MigrationDepList) MigrationDepList {
	for l := n.GetDeps().Front(); l != nil; l = l.Next() {
		migration := l.Value.(IMigrationNode)
		if migration.IsDummy() {
			continue
		}
		migrationName := l.Value.(IMigrationNode).GetMigration().GetName()
		if !Contains(migrationList, migrationName) && !Contains(depList, migrationName) {
			depList = append(depList, l.Value.(IMigrationNode).GetMigration().GetName())
			depList = l.Value.(IMigrationNode).TraverseDeps(migrationList, depList)
		}
	}
	return depList
}

func (n MigrationNode) TraverseChildren(migrationList []string) []string {
	for l := n.GetChildren().Front(); l != nil; l = l.Next() {
		migration := l.Value.(IMigrationNode)
		if migration.IsDummy() {
			continue
		}
		migrationName := l.Value.(IMigrationNode).GetMigration().GetName()
		if !Contains(migrationList, migrationName) {
			migrationList = append(migrationList, l.Value.(IMigrationNode).GetMigration().GetName())
			migrationDepList := n.TraverseDeps(migrationList, make(MigrationDepList, 0))
			sort.Reverse(migrationDepList)
			for _, m := range migrationDepList {
				migrationList = append(migrationList, m)
			}
			migrationList = n.TraverseChildren(migrationList)
		}
	}
	return migrationList
}

func NewMigrationNode(dep IMigrationNode, node IMigration, child IMigrationNode) IMigrationNode {
	depsList := list.New()
	if dep != nil {
		depsList.PushBack(dep)
	}
	childrenList := list.New()
	if child != nil {
		childrenList.PushBack(child)
	}
	return &MigrationNode{
		Deps: depsList,
		Node: node,
		Children: childrenList,
		applied: false,
		dummy: false,
		isRoot: false,
	}
}

func NewMigrationRootNode() IMigrationNode {
	return &MigrationNode{
		Deps:     list.New(),
		Node:     nil,
		Children: list.New(),
		applied:  false,
		dummy: true,
		isRoot: true,
	}
}

type MigrationTree struct {
	Root  IMigrationNode
	nodes map[string]IMigrationNode
	treeBuilt *bool
}

func (t MigrationTree) TreeBuilt() {
	*t.treeBuilt = true
}

func (t MigrationTree) IsTreeBuilt() bool {
	return *t.treeBuilt
}

func (t MigrationTree) GetNodeByMigrationName(migrationName string) (IMigrationNode, error){
	node, ok := t.nodes[migrationName]
	if ok {
		return node, nil
	} else {
		return nil, fmt.Errorf("No node with name %s has been found", migrationName)
	}
}

func (t MigrationTree) AddNode(node IMigrationNode) error {
	_, ok := t.nodes[node.GetMigration().GetName()]
	if ok {
		return fmt.Errorf("Migration with name %s has been added to tree before", node.GetMigration().GetName())
	}
	t.nodes[node.GetMigration().GetName()] = node
	return nil
}

func (t MigrationTree) GetRoot() IMigrationNode {
	return t.Root
}

func (t MigrationTree) SetRoot(root IMigrationNode) {
	root.SetItAsRoot()
	t.Root = root
}

type MigrationList []IMigration

func (m MigrationList) Len() int { return len(m) }
func (m MigrationList) Less(i, j int) bool {
	return m[i].GetId() < m[j].GetId()
}
func (m MigrationList) Swap(i, j int){ m[i], m[j] = m[j], m[i] }

type MigrationDepList []string

func (m MigrationDepList) Len() int { return len(m) }
func (m MigrationDepList) Less(i, j int) bool {
	return i < j
}
func (m MigrationDepList) Swap(i, j int){ m[i], m[j] = m[j], m[i] }

type MigrationRegistry struct {
	migrations map[string]IMigration
}

func (r MigrationRegistry) AddMigration(migration IMigration) {
	r.migrations[migration.GetName()] = migration
}

func (r MigrationRegistry) GetByName(migrationName string) (IMigration, error) {
	migration, ok := r.migrations[migrationName]
	if ok {
		return migration, nil
	} else {
		return nil, fmt.Errorf("No migration with name %s exists", migrationName)
	}
}

func (r MigrationRegistry) GetSortedMigrations() MigrationList {
	sortedMigrations := make(MigrationList, len(r.migrations))
	i := 0
	for _, migration := range r.migrations {
		sortedMigrations[i] = migration
		i += 1
	}
	sort.Sort(sortedMigrations)
	return sortedMigrations
}

func NewMigrationRegistry() *MigrationRegistry {
	return &MigrationRegistry{
		migrations: make(map[string]IMigration),
	}
}

func NewMigrationTree() IMigrationTree {
	var builtTree bool
	return &MigrationTree{
		Root:  NewMigrationRootNode(),
		nodes: make(map[string]IMigrationNode),
		treeBuilt: &builtTree,
	}
}
