---
sidebar_position: 1
---

# Admin filter objects

Admin filter objects implements core.IAdminFilterObjects interface.

```go
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
```
Right now uadmin supports only objects that stored in database and we use gorm to interact with database.
But later on we want to provide implementations for the objects stored in NoSQL, like Mongo, etc
