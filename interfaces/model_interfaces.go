package interfaces

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"sort"
	"strings"
)

type ModelDescription struct {
	Model interface{}
	Statement *gorm.Statement
	GenerateModelI func() interface{}
}

type ProjectModelRegistry struct {
	models map[string]*ModelDescription
}

func (pmr *ProjectModelRegistry) RegisterModel(generateModelI func() interface{}) {
	model := generateModelI()
	uadminDatabase := NewUadminDatabaseWithoutConnection()
	statement := &gorm.Statement{DB: uadminDatabase.Db}
	statement.Parse(model)
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	modelName := v.Type().Name()
	pmr.models[modelName] = &ModelDescription{Model: model, Statement: statement, GenerateModelI: generateModelI}
}

func (pmr *ProjectModelRegistry) Iterate() <- chan *ModelDescription {
	chnl := make(chan *ModelDescription)
	go func() {
		defer close(chnl)
		for _, modelDescription := range pmr.models {
			chnl <- modelDescription
		}
	}()
	return chnl
}

func (pmr *ProjectModelRegistry) GetModelByName(modelName string) *ModelDescription {
	model, exists := pmr.models[modelName]
	if !exists {
		panic(fmt.Errorf("no model with name %s registered in the project", modelName))
	}
	return model
}

func (pmr *ProjectModelRegistry) GetModelFromInterface(model interface{}) *ModelDescription {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	modelName := v.Type().Name()
	modelI, _ := pmr.models[modelName]
	return modelI
}

var ProjectModels *ProjectModelRegistry

func init() {
	ProjectModels = &ProjectModelRegistry{
		models: make(map[string]*ModelDescription),
	}
}

func ClearProjectModels() {
	ProjectModels = &ProjectModelRegistry{
		models: make(map[string]*ModelDescription),
	}
}

type RemovalTreeNode struct {
	Model interface{}
	ModelDescription *ModelDescription
	Next []*RemovalTreeNode
	Prev []*RemovalTreeNode
	RawSQL []*DeleteRowStructure
	Visited bool
	Level int
}

type RemovalTreeNodeStringified struct {
	Explanation string
	Level int
}

type RemovalOrderList []*RemovalTreeNode

func TraverseRemovalTreeNode(nodeToVisit *RemovalTreeNode, removalOrderList *RemovalOrderList) {
	*removalOrderList = append(*removalOrderList, nodeToVisit)
	nodeToVisit.Visited = false
	for _, nextToRemove := range nodeToVisit.Next {
		nextToRemove.Visited = false
		*removalOrderList = append(*removalOrderList, nextToRemove)
		TraverseRemovalTreeNode(nextToRemove, removalOrderList)
	}
}

func (rtn *RemovalTreeNode) RemoveFromDatabase(uadminDatabase *UadminDatabase) error {
	var removalOrder RemovalOrderList
	TraverseRemovalTreeNode(rtn, &removalOrder)
	sort.Slice(removalOrder, func(i, j int) bool {
		return removalOrder[i].Level > removalOrder[j].Level
	})
	for _, removalTreeNode := range removalOrder {
		if len(removalTreeNode.RawSQL) > 0 {
			for _, rawSQL := range removalTreeNode.RawSQL {
				res := uadminDatabase.Db.Exec(rawSQL.SQL, rawSQL.Values...)
				if res.Error != nil {
					return res.Error
				}
			}
		}
		res1 := uadminDatabase.Db.Unscoped().Delete(removalTreeNode.Model)
		if res1.Error != nil {
			return res1.Error
		}
	}
	return nil
}

func (rtn *RemovalTreeNode) BuildDeletionTreeStringified(uadminDatabase *UadminDatabase) []*RemovalTreeNodeStringified {
	var removalTreeStringified []*RemovalTreeNodeStringified
	var removalOrder RemovalOrderList
	TraverseRemovalTreeNode(rtn, &removalOrder)
	//sort.Slice(removalOrder, func(i, j int) bool {
	//	return removalOrder[i].Level > removalOrder[j].Level
	//})
	removalTreeStringified = make([]*RemovalTreeNodeStringified, 0)
	for _, removalTreeNode := range removalOrder {
		uadminDatabase.Db.Unscoped().First(removalTreeNode.Model)
		if len(removalTreeNode.RawSQL) > 0 {
			for _, rawSQL := range removalTreeNode.RawSQL {
				removalTreeStringified = append(removalTreeStringified, &RemovalTreeNodeStringified{
					Explanation: fmt.Sprintf("Association with %s", rawSQL.Table),
					Level: removalTreeNode.Level,
				})
			}
		}
		gormModelV := reflect.Indirect(reflect.ValueOf(removalTreeNode.Model))
		Idv := TransformValueForWidget(gormModelV.FieldByName(removalTreeNode.ModelDescription.Statement.Schema.PrimaryFields[0].Name).Interface())
		modelAdminPage := CurrentAdminPageRegistry.GetByModelName(removalTreeNode.ModelDescription.Statement.Schema.Name)
		if modelAdminPage != nil {
			url := fmt.Sprintf("%s/%s/%s", CurrentConfig.D.Uadmin.RootAdminURL, modelAdminPage.ParentPage.Slug, modelAdminPage.Slug)
			removalTreeStringified = append(removalTreeStringified, &RemovalTreeNodeStringified{
				Explanation: fmt.Sprintf("<a target='_blank' href='%s/%s'>%s</a>", url, Idv, reflect.ValueOf(removalTreeNode.Model).MethodByName("String").Call([]reflect.Value{})[0]),
				Level: removalTreeNode.Level,
			})
		} else {
			removalTreeStringified = append(removalTreeStringified, &RemovalTreeNodeStringified{
				Explanation: fmt.Sprintf("%s", reflect.ValueOf(removalTreeNode.Model).MethodByName("String").Call([]reflect.Value{})[0]),
				Level: removalTreeNode.Level,
			})
		}
	}
	return removalTreeStringified
}

func BuildRemovalTree(uadminDatabase *UadminDatabase, model interface{}, level ...int) *RemovalTreeNode {
	var realLevel int
	if len(level) == 0 {
		realLevel = 1
	} else {
		realLevel = level[0] + 1
	}
	modelInfo := ProjectModels.GetModelFromInterface(model)
	removalTreeNode := &RemovalTreeNode{
		Model: model,
		Next: make([]*RemovalTreeNode, 0),
		Prev: make([]*RemovalTreeNode, 0),
		RawSQL: make([]*DeleteRowStructure, 0),
		ModelDescription: modelInfo,
		Level: realLevel,
	}
	for modelDescription := range ProjectModels.Iterate() {
		for _, relationShip := range modelDescription.Statement.Schema.Relationships.Relations {
			if relationShip.Type == "many_to_many" {
				uadminDatabase.Db.Model(model).Preload(relationShip.Name)
				for _, reference := range relationShip.References {
					if reference.PrimaryKey.Schema.Table == modelInfo.Statement.Table {
						gormModelV := reflect.Indirect(reflect.ValueOf(model))
						cond := fmt.Sprintf(
							"%s = ?",
							reference.ForeignKey.DBName,
						)
						deleteSQL := uadminDatabase.Adapter.BuildDeleteString(
							reference.ForeignKey.Schema.Table,
							cond,
							TransformValueForWidget(gormModelV.FieldByName(modelInfo.Statement.Schema.PrimaryFields[0].Name).Interface()),
						)
						deleteSQL.Table = reference.ForeignKey.Schema.Table
						removalTreeNode.RawSQL = append(removalTreeNode.RawSQL, deleteSQL)
					}
				}
			}
			if relationShip.Type == "belongs_to" {
				relationsString := []string{}
				foundRelation := false
				primaryKeyName := ""
				primaryStructField := ""
				for _, reference := range relationShip.References {
					if reference.PrimaryKey.Schema.Table == modelInfo.Statement.Table {
						foundRelation = true
						primaryKeyName = reference.PrimaryKey.DBName
						primaryStructField = reference.ForeignKey.Name
					}
					relationsString = append(
						relationsString,
						fmt.Sprintf(
							"%s.%s = %s.%s",
							modelDescription.Statement.Table, reference.ForeignKey.DBName, modelInfo.Statement.Table,
							reference.PrimaryKey.DBName,
						),
					)
				}
				if !foundRelation {
					continue
				}
				db := uadminDatabase.Db.Model(modelDescription.GenerateModelI())
				if relationShip.Field.NotNull {
					db = db.Joins(
						fmt.Sprintf(
							"INNER JOIN %s on %s",
							modelInfo.Statement.Table, strings.Join(relationsString, " AND "),
						),
					)
				} else {
					db = db.Joins(
						fmt.Sprintf(
							"LEFT JOIN %s on %s",
							modelInfo.Statement.Table, strings.Join(relationsString, " AND "),
						),
					)
				}
				gormModelV := reflect.Indirect(reflect.ValueOf(model))
				Idv := TransformValueForWidget(gormModelV.FieldByName(modelInfo.Statement.Schema.PrimaryFields[0].Name).Interface())

				rows, _ := db.Unscoped().Preload(primaryStructField).Where(fmt.Sprintf("%s.%s = ?", modelInfo.Statement.Table, primaryKeyName), Idv).Rows()
				for rows.Next() {
					newModel1 := modelDescription.GenerateModelI()
					uadminDatabase.Db.ScanRows(rows, newModel1)
					removalTreeNode.Next = append(removalTreeNode.Next, BuildRemovalTree(uadminDatabase, newModel1, realLevel))
				}
				rows.Close()
			}
		}
	}
	return removalTreeNode
}
