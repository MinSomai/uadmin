package core

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/url"
	"reflect"
	"strings"
)

func NewGormAdminPage(parentPage *AdminPage, genModelI func() (interface{}, interface{}), generateForm func(modelI interface{}, ctx IAdminContext) *Form) *AdminPage {
	modelI4, _ := genModelI()
	modelName := ""
	if modelI4 != nil {
		uadminDatabase := NewUadminDatabaseWithoutConnection()
		stmt := &gorm.Statement{DB: uadminDatabase.Db}
		stmt.Parse(modelI4)
		modelName = strings.ToLower(stmt.Schema.Name)
	}
	var form *Form
	var listDisplay *ListDisplayRegistry
	var searchFieldRegistry *SearchFieldRegistry
	if modelI4 != nil {
		form = NewFormFromModelFromGinContext(&AdminContext{}, modelI4, make([]string, 0), []string{"ID"}, true, "")
		listDisplay = NewListDisplayRegistryFromGormModel(modelI4)
		searchFieldRegistry = NewSearchFieldRegistryFromGormModel(modelI4)
	}
	return &AdminPage{
		Form:           form,
		SubPages:       NewAdminPageRegistry(),
		GenerateModelI: genModelI,
		ParentPage:     parentPage,
		GetQueryset: func(adminPage *AdminPage, adminRequestParams *AdminRequestParams) IAdminFilterObjects {
			uadminDatabase := NewUadminDatabase()
			db := uadminDatabase.Db
			var paginatedQuerySet IPersistenceStorage
			var perPage int
			modelI, _ := genModelI()
			modelI1, _ := genModelI()
			modelI2, _ := genModelI()
			modelI3, _ := genModelI()
			ret := &GormAdminFilterObjects{
				InitialGormQuerySet:   NewGormPersistenceStorage(db.Model(modelI)),
				GormQuerySet:          NewGormPersistenceStorage(db.Model(modelI1)),
				PaginatedGormQuerySet: NewGormPersistenceStorage(db.Model(modelI2)),
				Model:                 modelI3,
				UadminDatabase:        uadminDatabase,
				GenerateModelI:        genModelI,
			}
			if adminRequestParams != nil && adminRequestParams.RequestURL != "" {
				url1, _ := url.Parse(adminRequestParams.RequestURL)
				queryParams, _ := url.ParseQuery(url1.RawQuery)
				for filter := range adminPage.ListFilter.Iterate() {
					filterValue := queryParams.Get(filter.URLFilteringParam)
					if filterValue != "" {
						filter.FilterQs(ret, fmt.Sprintf("%s=%s", filter.URLFilteringParam, filterValue))
					}
				}
			}
			if adminRequestParams != nil && adminRequestParams.Search != "" {
				searchFilterObjects := &GormAdminFilterObjects{
					InitialGormQuerySet:   NewGormPersistenceStorage(db),
					GormQuerySet:          NewGormPersistenceStorage(db),
					PaginatedGormQuerySet: NewGormPersistenceStorage(db),
					Model:                 modelI3,
					UadminDatabase:        uadminDatabase,
					GenerateModelI:        genModelI,
				}
				for filter := range adminPage.SearchFields.GetAll() {
					filter.Search(searchFilterObjects, adminRequestParams.Search)
				}
				ret.SetPaginatedQuerySet(ret.GetPaginatedQuerySet().Where(searchFilterObjects.GetPaginatedQuerySet().GetCurrentDB()))
				ret.SetFullQuerySet(ret.GetFullQuerySet().Where(searchFilterObjects.GetFullQuerySet().GetCurrentDB()))
			}
			if adminRequestParams != nil && adminRequestParams.Paginator.PerPage > 0 {
				perPage = adminRequestParams.Paginator.PerPage
			} else {
				perPage = adminPage.Paginator.PerPage
			}
			if adminRequestParams != nil {
				paginatedQuerySet = ret.GetPaginatedQuerySet().Offset(adminRequestParams.Paginator.Offset)
				if adminPage.Paginator.ShowLastPageOnPreviousPage {
					var countRecords int64
					ret.GetFullQuerySet().Count(&countRecords)
					if countRecords > int64(adminRequestParams.Paginator.Offset+(2*perPage)) {
						paginatedQuerySet = paginatedQuerySet.Limit(perPage)
					} else {
						paginatedQuerySet = paginatedQuerySet.Limit(int(countRecords - int64(adminRequestParams.Paginator.Offset)))
					}
				} else {
					paginatedQuerySet = paginatedQuerySet.Limit(perPage)
				}
				ret.SetPaginatedQuerySet(paginatedQuerySet)
				for listDisplay := range adminPage.ListDisplay.GetAllFields() {
					direction := listDisplay.SortBy.GetDirection()
					if len(adminRequestParams.Ordering) > 0 {
						for _, ordering := range adminRequestParams.Ordering {
							directionSort := 1
							if strings.HasPrefix(ordering, "-") {
								directionSort = -1
								ordering = ordering[1:]
							}
							if ordering == listDisplay.DisplayName {
								direction = directionSort
								listDisplay.SortBy.Sort(ret, direction)
							}
						}
					}
				}
			}
			return ret
		},
		Model:                   modelI4,
		ModelName:               modelName,
		Validators:              NewValidatorRegistry(),
		ExcludeFields:           NewFieldRegistry(),
		FieldsToShow:            NewFieldRegistry(),
		ModelActionsRegistry:    NewAdminModelActionRegistry(),
		InlineRegistry:          NewAdminPageInlineRegistry(),
		ListDisplay:             listDisplay,
		ListFilter:              &ListFilterRegistry{ListFilter: make([]*ListFilter, 0)},
		SearchFields:            searchFieldRegistry,
		Paginator:               &Paginator{PerPage: CurrentConfig.D.Uadmin.AdminPerPage, ShowLastPageOnPreviousPage: true},
		ActionsSelectionCounter: true,
		FilterOptions:           NewFilterOptionsRegistry(),
		GenerateForm:            generateForm,
	}
}

type GormPersistenceStorage struct {
	Db *gorm.DB
}

func NewGormPersistenceStorage(db *gorm.DB) *GormPersistenceStorage {
	return &GormPersistenceStorage{Db: db}
}

func (gps *GormPersistenceStorage) Association(column string) IPersistenceAssociation {
	return gps.Db.Association(column)
}

func (gps *GormPersistenceStorage) Model(value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Model(value)
	return gps
}

func (gps *GormPersistenceStorage) Clauses(conds ...clause.Expression) IPersistenceStorage {
	gps.Db = gps.Db.Clauses(conds...)
	return gps
}

func (gps *GormPersistenceStorage) GetCurrentDB() *gorm.DB {
	return gps.Db
}

func (gps *GormPersistenceStorage) Table(name string, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Table(name, args...)
	return gps
}

func (gps *GormPersistenceStorage) Distinct(args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Distinct(args...)
	return gps
}

func (gps *GormPersistenceStorage) Select(query interface{}, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Select(query, args...)
	return gps
}

func (gps *GormPersistenceStorage) Omit(columns ...string) IPersistenceStorage {
	gps.Db = gps.Db.Omit(columns...)
	return gps
}

func (gps *GormPersistenceStorage) Where(query interface{}, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Where(query, args...)
	return gps
}

func (gps *GormPersistenceStorage) Not(query interface{}, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Not(query, args...)
	return gps
}

func (gps *GormPersistenceStorage) Or(query interface{}, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Or(query, args...)
	return gps
}

func (gps *GormPersistenceStorage) Joins(query string, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Joins(query, args...)
	return gps
}

func (gps *GormPersistenceStorage) Group(name string) IPersistenceStorage {
	gps.Db = gps.Db.Group(name)
	return gps
}

func (gps *GormPersistenceStorage) Having(query interface{}, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Having(query, args...)
	return gps
}

func (gps *GormPersistenceStorage) Order(value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Order(value)
	return gps
}

func (gps *GormPersistenceStorage) Limit(limit int) IPersistenceStorage {
	gps.Db = gps.Db.Limit(limit)
	return gps
}

func (gps *GormPersistenceStorage) Offset(offset int) IPersistenceStorage {
	gps.Db = gps.Db.Offset(offset)
	return gps
}

func (gps *GormPersistenceStorage) Scopes(funcs ...func(*gorm.DB) *gorm.DB) IPersistenceStorage {
	gps.Db = gps.Db.Scopes(funcs...)
	return gps
}

func (gps *GormPersistenceStorage) Preload(query string, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Preload(query, args...)
	return gps
}

func (gps *GormPersistenceStorage) Attrs(attrs ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Attrs(attrs...)
	return gps
}

func (gps *GormPersistenceStorage) Assign(attrs ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Assign(attrs...)
	return gps
}

func (gps *GormPersistenceStorage) Unscoped() IPersistenceStorage {
	gps.Db = gps.Db.Unscoped()
	return gps
}

func (gps *GormPersistenceStorage) Raw(sql string, values ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Raw(sql, values...)
	return gps
}

func (gps *GormPersistenceStorage) Migrator() IPersistenceMigrator {
	return gps.Db.Migrator()
}

func (gps *GormPersistenceStorage) AutoMigrate(dst ...interface{}) error {

	return gps.Db.AutoMigrate(dst...)
}

func (gps *GormPersistenceStorage) Session(config *gorm.Session) IPersistenceStorage {
	gps.Db = gps.Db.Session(config)
	return gps
}

func (gps *GormPersistenceStorage) WithContext(ctx context.Context) IPersistenceStorage {
	gps.Db = gps.Db.WithContext(ctx)
	return gps
}

func (gps *GormPersistenceStorage) Debug() IPersistenceStorage {
	gps.Db = gps.Db.Debug()
	return gps
}

func (gps *GormPersistenceStorage) Set(key string, value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Set(key, value)
	return gps
}

func (gps *GormPersistenceStorage) Get(key string) (interface{}, bool) {
	return gps.Db.Get(key)
}

func (gps *GormPersistenceStorage) InstanceSet(key string, value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.InstanceSet(key, value)
	return gps
}

func (gps *GormPersistenceStorage) InstanceGet(key string) (interface{}, bool) {
	return gps.Db.InstanceGet(key)
}

func (gps *GormPersistenceStorage) AddError(err error) error {
	return gps.Db.AddError(err)
}

func (gps *GormPersistenceStorage) DB() (*sql.DB, error) {
	return gps.Db.DB()
}

func (gps *GormPersistenceStorage) SetupJoinTable(model interface{}, field string, joinTable interface{}) error {
	return gps.Db.SetupJoinTable(model, field, joinTable)
}

func (gps *GormPersistenceStorage) Use(plugin gorm.Plugin) error {
	return gps.Db.Use(plugin)
}

func (gps *GormPersistenceStorage) Create(value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Create(value)
	return gps
}

func (gps *GormPersistenceStorage) CreateInBatches(value interface{}, batchSize int) IPersistenceStorage {
	gps.Db = gps.Db.CreateInBatches(value, batchSize)
	return gps
}

func (gps *GormPersistenceStorage) Save(value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Save(value)
	return gps
}

func (gps *GormPersistenceStorage) First(dest interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.First(dest, conds...)
	return gps
}

func (gps *GormPersistenceStorage) Take(dest interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Take(dest, conds...)
	return gps
}

func (gps *GormPersistenceStorage) Last(dest interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Last(dest, conds...)
	return gps
}

func (gps *GormPersistenceStorage) Find(dest interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Find(dest, conds...)
	return gps
}

func (gps *GormPersistenceStorage) FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) IPersistenceStorage {
	gps.Db = gps.Db.FindInBatches(dest, batchSize, fc)
	return gps
}

func (gps *GormPersistenceStorage) FirstOrInit(dest interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.FirstOrInit(dest, conds...)
	return gps
}

func (gps *GormPersistenceStorage) FirstOrCreate(dest interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.FirstOrCreate(dest, conds...)
	return gps
}

func (gps *GormPersistenceStorage) Update(column string, value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Update(column, value)
	return gps
}

func (gps *GormPersistenceStorage) Updates(values interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Updates(values)
	return gps
}

func (gps *GormPersistenceStorage) UpdateColumn(column string, value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.UpdateColumn(column, value)
	return gps
}

func (gps *GormPersistenceStorage) UpdateColumns(values interface{}) IPersistenceStorage {
	gps.Db = gps.Db.UpdateColumns(values)
	return gps
}

func (gps *GormPersistenceStorage) Delete(value interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Delete(value, conds...)
	return gps
}

func (gps *GormPersistenceStorage) Count(count *int64) IPersistenceStorage {
	gps.Db = gps.Db.Count(count)
	return gps
}

func (gps *GormPersistenceStorage) Row() IPersistenceIterateRow {
	return gps.Db.Row()
}

func (gps *GormPersistenceStorage) Rows() (IPersistenceIterateRows, error) {
	return gps.Db.Rows()
}

func (gps *GormPersistenceStorage) Scan(dest interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Scan(dest)
	return gps
}

func (gps *GormPersistenceStorage) Pluck(column string, dest interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Pluck(column, dest)
	return gps
}

func (gps *GormPersistenceStorage) ScanRows(rows IPersistenceIterateRows, dest interface{}) error {
	return gps.Db.ScanRows(rows.(*sql.Rows), dest)
}

func (gps *GormPersistenceStorage) Transaction(fc func(*gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	return gps.Db.Transaction(fc, opts...)
}

func (gps *GormPersistenceStorage) Begin(opts ...*sql.TxOptions) IPersistenceStorage {
	gps.Db = gps.Db.Begin(opts...)
	return gps
}

func (gps *GormPersistenceStorage) Commit() IPersistenceStorage {
	gps.Db = gps.Db.Commit()
	return gps
}

func (gps *GormPersistenceStorage) Rollback() IPersistenceStorage {
	gps.Db = gps.Db.Rollback()
	return gps
}

func (gps *GormPersistenceStorage) SavePoint(name string) IPersistenceStorage {
	gps.Db = gps.Db.SavePoint(name)
	return gps
}

func (gps *GormPersistenceStorage) RollbackTo(name string) IPersistenceStorage {
	gps.Db = gps.Db.RollbackTo(name)
	return gps
}

func (gps *GormPersistenceStorage) Exec(sql string, values ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Exec(sql, values)
	return gps
}

type GormAdminFilterObjects struct {
	InitialGormQuerySet   IPersistenceStorage
	GormQuerySet          IPersistenceStorage
	PaginatedGormQuerySet IPersistenceStorage
	Model                 interface{}
	UadminDatabase        *UadminDatabase
	GenerateModelI        func() (interface{}, interface{})
}

func (afo *GormAdminFilterObjects) FilterQs(filterString string) {
	statement := &gorm.Statement{DB: afo.GetUadminDatabase().Db}
	statement.Parse(afo.GetCurrentModel())
	schema1 := statement.Schema
	operatorContext := FilterGormModel(afo.GetUadminDatabase().Adapter, afo.GetFullQuerySet(), schema1, []string{filterString}, afo.GetCurrentModel())
	afo.SetFullQuerySet(operatorContext.Tx)
	operatorContext = FilterGormModel(afo.GetUadminDatabase().Adapter, afo.GetPaginatedQuerySet(), schema1, []string{filterString}, afo.GetCurrentModel())
	afo.SetPaginatedQuerySet(operatorContext.Tx)

}

func (afo *GormAdminFilterObjects) Search(field *Field, searchString string) {
	operator := IContainsGormOperator{}
	gormOperatorContext := NewGormOperatorContext(afo.GetFullQuerySet(), afo.GetCurrentModel())
	operator.Build(afo.GetUadminDatabase().Adapter, gormOperatorContext, field, searchString, &SQLConditionBuilder{Type: "or"})
	afo.SetFullQuerySet(gormOperatorContext.Tx)
	gormOperatorContext = NewGormOperatorContext(afo.GetPaginatedQuerySet(), afo.GetCurrentModel())
	operator.Build(afo.GetUadminDatabase().Adapter, gormOperatorContext, field, searchString, &SQLConditionBuilder{Type: "or"})
	afo.SetPaginatedQuerySet(gormOperatorContext.Tx)
}

func (afo *GormAdminFilterObjects) GetPaginatedQuerySet() IPersistenceStorage {
	return afo.PaginatedGormQuerySet
}

func (afo *GormAdminFilterObjects) GetFullQuerySet() IPersistenceStorage {
	return afo.GormQuerySet
}

func (afo *GormAdminFilterObjects) SetFullQuerySet(storage IPersistenceStorage) {
	afo.GormQuerySet = storage
}

func (afo *GormAdminFilterObjects) GenerateModelInterface() (interface{}, interface{}) {
	return afo.GenerateModelI()
}

func (afo *GormAdminFilterObjects) GetInitialQuerySet() IPersistenceStorage {
	return afo.InitialGormQuerySet
}

func (afo *GormAdminFilterObjects) SetInitialQuerySet(storage IPersistenceStorage) {
	afo.InitialGormQuerySet = storage
}

func (afo *GormAdminFilterObjects) GetCurrentModel() interface{} {
	return afo.Model
}

func (afo *GormAdminFilterObjects) GetUadminDatabase() *UadminDatabase {
	return afo.UadminDatabase
}

func (afo *GormAdminFilterObjects) SetPaginatedQuerySet(storage IPersistenceStorage) {
	afo.PaginatedGormQuerySet = storage
}

func (afo *GormAdminFilterObjects) GetDB() IPersistenceStorage {
	return NewGormPersistenceStorage(afo.UadminDatabase.Db)
}

func (afo *GormAdminFilterObjects) WithTransaction(handler func(afo1 IAdminFilterObjects) error) {
	afo.UadminDatabase.Db.Transaction(func(tx *gorm.DB) error {
		return handler(&GormAdminFilterObjects{UadminDatabase: &UadminDatabase{Db: tx, Adapter: afo.UadminDatabase.Adapter}, GenerateModelI: afo.GenerateModelI})
	})
}

func (afo *GormAdminFilterObjects) LoadDataForModelByID(ID interface{}, model interface{}) {
	afo.UadminDatabase.Db.Preload(clause.Associations).First(model, ID)
}

func (afo *GormAdminFilterObjects) SaveModel(model interface{}) error {
	res := afo.UadminDatabase.Db.Save(model)
	return res.Error
}

func (afo *GormAdminFilterObjects) CreateNew(model interface{}) error {
	res := afo.UadminDatabase.Db.Model(model).Create(model)
	return res.Error
}

func (afo *GormAdminFilterObjects) FilterByMultipleIds(field *Field, realObjectIds []string) {
	afo.SetFullQuerySet(afo.GetFullQuerySet().Where(fmt.Sprintf("%s IN ?", field.DBName), realObjectIds))
}

func (afo *GormAdminFilterObjects) RemoveModelPermanently(model interface{}) error {
	res := afo.UadminDatabase.Db.Unscoped().Delete(model)
	return res.Error
}

func (afo *GormAdminFilterObjects) SortBy(field *Field, direction int) {
	sortBy := field.DBName
	if direction == -1 {
		sortBy += " desc"
	}
	afo.SetPaginatedQuerySet(afo.GetPaginatedQuerySet().Order(sortBy))
}

func (afo *GormAdminFilterObjects) GetPaginated() <-chan *IterateAdminObjects {
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
			yieldV := &IterateAdminObjects{
				Model:         model,
				ID:            ID.(string),
				RenderContext: &FormRenderContext{Model: model},
			}
			chnl <- yieldV
		}
	}()
	return chnl
}

func (afo *GormAdminFilterObjects) IterateThroughWholeQuerySet() <-chan *IterateAdminObjects {
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
			yieldV := &IterateAdminObjects{
				Model:         model,
				ID:            ID.(string),
				RenderContext: &FormRenderContext{Model: model},
			}
			chnl <- yieldV
		}
	}()
	return chnl
}
