package core

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"strconv"
)

type IAdminFilterObjects interface {
	WithTransaction(handler func(afo1 IAdminFilterObjects) error)
	LoadDataForModelByID(ID interface{}, model interface{})
	SaveModel(model interface{}) error
	CreateNew(model interface{}) error
	GetPaginated() <-chan *IterateAdminObjects
	IterateThroughWholeQuerySet() <-chan *IterateAdminObjects
	GetPaginatedQuerySet() IPersistenceStorage
	GetFullQuerySet() IPersistenceStorage
	SetFullQuerySet(IPersistenceStorage)
	SetPaginatedQuerySet(IPersistenceStorage)
	GetUadminDatabase() *UadminDatabase
	GetCurrentModel() interface{}
	GetInitialQuerySet() IPersistenceStorage
	SetInitialQuerySet(IPersistenceStorage)
	GenerateModelInterface() (interface{}, interface{})
	RemoveModelPermanently(model interface{}) error
}

type AdminFilterObjects struct {
	InitialGormQuerySet   IPersistenceStorage
	GormQuerySet          IPersistenceStorage
	PaginatedGormQuerySet IPersistenceStorage
	Model                 interface{}
	UadminDatabase        *UadminDatabase
	GenerateModelI        func() (interface{}, interface{})
}

type IterateAdminObjects struct {
	Model         interface{}
	ID            uint
	RenderContext *FormRenderContext
}

func (afo *AdminFilterObjects) GetPaginatedQuerySet() IPersistenceStorage {
	return afo.PaginatedGormQuerySet
}

func (afo *AdminFilterObjects) GetFullQuerySet() IPersistenceStorage {
	return afo.GormQuerySet
}

func (afo *AdminFilterObjects) SetFullQuerySet(storage IPersistenceStorage) {
	afo.GormQuerySet = storage
}

func (afo *AdminFilterObjects) GenerateModelInterface() (interface{}, interface{}) {
	return afo.GenerateModelI()
}

func (afo *AdminFilterObjects) GetInitialQuerySet() IPersistenceStorage {
	return afo.InitialGormQuerySet
}

func (afo *AdminFilterObjects) SetInitialQuerySet(storage IPersistenceStorage) {
	afo.InitialGormQuerySet = storage
}

func (afo *AdminFilterObjects) GetCurrentModel() interface{} {
	return afo.Model
}

func (afo *AdminFilterObjects) GetUadminDatabase() *UadminDatabase {
	return afo.UadminDatabase
}

func (afo *AdminFilterObjects) SetPaginatedQuerySet(storage IPersistenceStorage) {
	afo.PaginatedGormQuerySet = storage
}

func (afo *AdminFilterObjects) WithTransaction(handler func(afo1 IAdminFilterObjects) error) {
	afo.UadminDatabase.Db.Transaction(func(tx *gorm.DB) error {
		return handler(&AdminFilterObjects{UadminDatabase: &UadminDatabase{Db: tx}, GenerateModelI: afo.GenerateModelI})
	})
}

func (afo *AdminFilterObjects) LoadDataForModelByID(ID interface{}, model interface{}) {
	afo.UadminDatabase.Db.Preload(clause.Associations).First(model, ID)
}

func (afo *AdminFilterObjects) SaveModel(model interface{}) error {
	res := afo.UadminDatabase.Db.Save(model)
	return res.Error
}

func (afo *AdminFilterObjects) CreateNew(model interface{}) error {
	res := afo.UadminDatabase.Db.Model(model).Create(model)
	return res.Error
}

func (afo *AdminFilterObjects) RemoveModelPermanently(model interface{}) error {
	res := afo.UadminDatabase.Db.Unscoped().Delete(model)
	return res.Error
}

func (afo *AdminFilterObjects) GetPaginated() <-chan *IterateAdminObjects {
	chnl := make(chan *IterateAdminObjects)
	go func() {
		defer close(chnl)
		modelI, models := afo.GenerateModelI()
		modelDescription := ProjectModels.GetModelFromInterface(modelI)
		afo.PaginatedGormQuerySet.Preload(clause.Associations).Find(models)
		s := reflect.Indirect(reflect.ValueOf(models))
		for i := 0; i < s.Len(); i++ {
			model := s.Index(i).Interface()
			gormModelV := reflect.Indirect(reflect.ValueOf(model))
			ID := TransformValueForWidget(gormModelV.FieldByName(modelDescription.Statement.Schema.PrimaryFields[0].Name).Interface())
			IDN, _ := strconv.Atoi(ID.(string))
			yieldV := &IterateAdminObjects{
				Model:         model,
				ID:            uint(IDN),
				RenderContext: &FormRenderContext{Model: model},
			}
			chnl <- yieldV
		}
	}()
	return chnl
}

func (afo *AdminFilterObjects) IterateThroughWholeQuerySet() <-chan *IterateAdminObjects {
	chnl := make(chan *IterateAdminObjects)
	go func() {
		defer close(chnl)
		modelI, models := afo.GenerateModelI()
		modelDescription := ProjectModels.GetModelFromInterface(modelI)
		afo.GormQuerySet.Preload(clause.Associations).Find(models)
		s := reflect.Indirect(reflect.ValueOf(models))
		for i := 0; i < s.Len(); i++ {
			model := s.Index(i).Interface()
			gormModelV := reflect.Indirect(reflect.ValueOf(model))
			ID := TransformValueForWidget(gormModelV.FieldByName(modelDescription.Statement.Schema.PrimaryFields[0].Name).Interface())
			IDN, _ := strconv.Atoi(ID.(string))
			yieldV := &IterateAdminObjects{
				Model:         model,
				ID:            uint(IDN),
				RenderContext: &FormRenderContext{Model: model},
			}
			chnl <- yieldV
		}
	}()
	return chnl
}

