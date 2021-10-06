package core

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

type IterateAdminObjects struct {
	Model         interface{}
	ID            uint
	RenderContext *FormRenderContext
}
