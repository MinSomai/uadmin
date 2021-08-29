package interfaces

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"math"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type AdminRequestPaginator struct {
	PerPage int
	Offset int
}

type AdminRequestParams struct {
	CreateSession bool
	GenerateCSRFToken bool
	NeedAllLanguages bool
	Paginator *AdminRequestPaginator
	RequestURL string
	Search string
	Ordering []string
}

func (arp *AdminRequestParams) GetOrdering() string {
	return strings.Join(arp.Ordering, ",")
}

func NewAdminRequestParams() *AdminRequestParams {
	return &AdminRequestParams{
		CreateSession: true,
		GenerateCSRFToken: true,
		NeedAllLanguages: false,
		Paginator: &AdminRequestPaginator{},
	}
}

func NewAdminRequestParamsFromGinContext(ctx *gin.Context) *AdminRequestParams {
	ret := &AdminRequestParams{
		CreateSession: true,
		GenerateCSRFToken: true,
		NeedAllLanguages: false,
		Paginator: &AdminRequestPaginator{},
	}
	if ctx.Query("perpage") != "" {
		perPage, _ := strconv.Atoi(ctx.Query("perpage"))
		ret.Paginator.PerPage = perPage
 	} else {
		ret.Paginator.PerPage = CurrentConfig.D.Uadmin.AdminPerPage
	}
	if ctx.Query("offset") != "" {
		offset, _ := strconv.Atoi(ctx.Query("offset"))
		ret.Paginator.Offset = offset
	}
	if ctx.Query("p") != "" {
		page, _ := strconv.Atoi(ctx.Query("p"))
		if page > 1 {
			ret.Paginator.Offset = (page - 1) * ret.Paginator.PerPage
		}
	}
	ret.RequestURL = ctx.Request.URL.String()
	ret.Search = ctx.Query("search")
	orderingParts := strings.Split(ctx.Query("initialOrder"), ",")
	currentOrder := ctx.Query("o")
	currentOrderNameWithoutDirection := currentOrder
	if strings.HasPrefix(currentOrderNameWithoutDirection, "-") {
		currentOrderNameWithoutDirection = currentOrderNameWithoutDirection[1:]
	}
	foundNewOrder := false
	for i, part := range orderingParts {
		if strings.HasPrefix(part, "-") {
			part = part[1:]
		}
		if part == currentOrderNameWithoutDirection {
			orderingParts[i] = currentOrder
			foundNewOrder = true
		}
	}
	if !foundNewOrder {
		orderingParts = append(orderingParts, currentOrder)
	}
	finalOrderingParts := make([]string, 0)
	for _, part := range orderingParts {
		if part != "" {
			finalOrderingParts = append(finalOrderingParts, part)
		}
	}
	ret.Ordering = finalOrderingParts
	return ret
}

type AdminActionPlacement struct {
	DisplayOnEditPage bool
	//DisplayToTheTop bool
	//DisplayToTheBottom bool
	//DisplayToTheRight bool
	//DisplayToTheLeft bool
	ShowOnTheListPage bool
}

type IAdminModelActionInterface interface {
}

type AdminModelAction struct {
	IAdminModelActionInterface
	ActionName string
	Description string
	ShowFutureChanges bool
	RedirectToRootModelPage bool
	Placement *AdminActionPlacement
	PermName CustomPermission
	Handler func (adminPage *AdminPage, afo IAdminFilterObjects, ctx *gin.Context) (bool, int64)
	IsDisabled func (afo IAdminFilterObjects, ctx *gin.Context) bool
	SlugifiedActionName string
	RequestMethod string
	RequiresExtraSteps bool
}

func prepareAdminModelActionName(adminModelAction string) string {
	slugifiedAdminModelAction := AsciiRegex.ReplaceAllLiteralString(adminModelAction, "")
	slugifiedAdminModelAction = strings.Replace(strings.ToLower(slugifiedAdminModelAction), " ", "_", -1)
	slugifiedAdminModelAction = strings.Replace(strings.ToLower(slugifiedAdminModelAction), ".", "_", -1)
	return slugifiedAdminModelAction
}

func NewAdminModelAction(actionName string, placement *AdminActionPlacement) *AdminModelAction {
	return &AdminModelAction{
		RedirectToRootModelPage: true,
		ActionName: actionName,
		Placement: placement,
		SlugifiedActionName: prepareAdminModelActionName(actionName),
		RequestMethod: "POST",
	}
}

type IPersistenceMigrator interface {
	AutoMigrate(dst ...interface{}) error
	CurrentDatabase() string
	FullDataTypeOf(*schema.Field) clause.Expr
	CreateTable(dst ...interface{}) error
	DropTable(dst ...interface{}) error
	HasTable(dst interface{}) bool
	RenameTable(oldName interface{}, newName interface{}) error
	AddColumn(dst interface{}, field string) error
	DropColumn(dst interface{}, field string) error
	AlterColumn(dst interface{}, field string) error
	MigrateColumn(dst interface{}, field *schema.Field, columnType gorm.ColumnType) error
	HasColumn(dst interface{}, field string) bool
	RenameColumn(dst interface{}, oldName string, field string) error
	ColumnTypes(dst interface{}) ([]gorm.ColumnType, error)
	CreateView(name string, option gorm.ViewOption) error
	DropView(name string) error
	CreateConstraint(dst interface{}, name string) error
	DropConstraint(dst interface{}, name string) error
	HasConstraint(dst interface{}, name string) bool
	CreateIndex(dst interface{}, name string) error
	DropIndex(dst interface{}, name string) error
	HasIndex(dst interface{}, name string) bool
	RenameIndex(dst interface{}, oldName string, newName string) error
}

type IPersistenceAssociation interface {
	Find(out interface{}, conds ...interface{}) error
	Append(values ...interface{}) error
	Replace(values ...interface{}) error
	Delete(values ...interface{}) error
	Clear() error
	Count() (count int64)
}

type IPersistenceIterateRow interface {
	Scan(dest ...interface{}) error
	Err() error
}

type IPersistenceIterateRows interface {
	Next() bool
	NextResultSet() bool
	Err() error
	Columns() ([]string, error)
	ColumnTypes() ([]*sql.ColumnType, error)
	Scan(dest ...interface{}) error
	Close() error
}

type IPersistenceStorage interface {
	Association(column string) IPersistenceAssociation
	Model(value interface{}) IPersistenceStorage
	Clauses(conds ...clause.Expression) IPersistenceStorage
	Table(name string, args ...interface{}) IPersistenceStorage
	Distinct(args ...interface{}) IPersistenceStorage
	Select(query interface{}, args ...interface{}) IPersistenceStorage
	Omit(columns ...string) IPersistenceStorage
	Where(query interface{}, args ...interface{}) IPersistenceStorage
	Not(query interface{}, args ...interface{}) IPersistenceStorage
	Or(query interface{}, args ...interface{}) IPersistenceStorage
	Joins(query string, args ...interface{}) IPersistenceStorage
	Group(name string) IPersistenceStorage
	Having(query interface{}, args ...interface{}) IPersistenceStorage
	Order(value interface{}) IPersistenceStorage
	Limit(limit int) IPersistenceStorage
	Offset(offset int) IPersistenceStorage
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) IPersistenceStorage
	Preload(query string, args ...interface{}) IPersistenceStorage
	Attrs(attrs ...interface{}) IPersistenceStorage
	Assign(attrs ...interface{}) IPersistenceStorage
	Unscoped() IPersistenceStorage
	Raw(sql string, values ...interface{}) IPersistenceStorage
	Migrator() IPersistenceMigrator
	AutoMigrate(dst ...interface{}) error
	Session(config *gorm.Session) IPersistenceStorage
	WithContext(ctx context.Context) IPersistenceStorage
	Debug() IPersistenceStorage
	Set(key string, value interface{}) IPersistenceStorage
	Get(key string) (interface{}, bool)
	InstanceSet(key string, value interface{}) IPersistenceStorage
	InstanceGet(key string) (interface{}, bool)
	AddError(err error) error
	DB() (*sql.DB, error)
	SetupJoinTable(model interface{}, field string, joinTable interface{}) error
	Use(plugin gorm.Plugin) error
	Create(value interface{}) IPersistenceStorage
	CreateInBatches(value interface{}, batchSize int) IPersistenceStorage
	Save(value interface{}) IPersistenceStorage
	First(dest interface{}, conds ...interface{}) IPersistenceStorage
	Take(dest interface{}, conds ...interface{}) IPersistenceStorage
	Last(dest interface{}, conds ...interface{}) IPersistenceStorage
	Find(dest interface{}, conds ...interface{}) IPersistenceStorage
	FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) IPersistenceStorage
	FirstOrInit(dest interface{}, conds ...interface{}) IPersistenceStorage
	FirstOrCreate(dest interface{}, conds ...interface{}) IPersistenceStorage
	Update(column string, value interface{}) IPersistenceStorage
	Updates(values interface{}) IPersistenceStorage
	UpdateColumn(column string, value interface{}) IPersistenceStorage
	UpdateColumns(values interface{}) IPersistenceStorage
	Delete(value interface{}, conds ...interface{}) IPersistenceStorage
	Count(count *int64) IPersistenceStorage
	Row() IPersistenceIterateRow
	Rows() (IPersistenceIterateRows, error)
	Scan(dest interface{}) IPersistenceStorage
	Pluck(column string, dest interface{}) IPersistenceStorage
	ScanRows(rows IPersistenceIterateRows, dest interface{}) error
	Transaction(fc func(*gorm.DB) error, opts ...*sql.TxOptions) (err error)
	Begin(opts ...*sql.TxOptions) IPersistenceStorage
	Commit() IPersistenceStorage
	Rollback() IPersistenceStorage
	SavePoint(name string) IPersistenceStorage
	RollbackTo(name string) IPersistenceStorage
	Exec(sql string, values ...interface{}) IPersistenceStorage
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

type IAdminFilterObjects interface {
	WithTransaction(handler func(afo1 IAdminFilterObjects) error)
	LoadDataForModelById(Id interface{}, model interface{})
	SaveModel(model interface{}) error
	CreateNew(model interface{}) error
	GetPaginated() <- chan *IterateAdminObjects
	IterateThroughWholeQuerySet() <- chan *IterateAdminObjects
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
	InitialGormQuerySet IPersistenceStorage
	GormQuerySet IPersistenceStorage
	PaginatedGormQuerySet IPersistenceStorage
	Model interface{}
	UadminDatabase *UadminDatabase
	GenerateModelI func() (interface{}, interface{})
}

type IterateAdminObjects struct {
	Model interface {}
	Id uint
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

func (afo *AdminFilterObjects) LoadDataForModelById(Id interface{}, model interface{}) {
	modelI, _ := afo.GenerateModelI()
	afo.UadminDatabase.Db.Model(modelI).Preload(clause.Associations).First(model, Id)
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

func (afo *AdminFilterObjects) GetPaginated() <- chan *IterateAdminObjects {
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
			Id := TransformValueForWidget(gormModelV.FieldByName(modelDescription.Statement.Schema.PrimaryFields[0].Name).Interface())
			IdN, _ := strconv.Atoi(Id.(string))
			yieldV := &IterateAdminObjects{
				Model: model,
				Id: uint(IdN),
				RenderContext: &FormRenderContext{Model: model},
			}
			chnl <- yieldV
		}
	}()
	return chnl
}

func (afo *AdminFilterObjects) IterateThroughWholeQuerySet() <- chan *IterateAdminObjects {
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
			Id := TransformValueForWidget(gormModelV.FieldByName(modelDescription.Statement.Schema.PrimaryFields[0].Name).Interface())
			IdN, _ := strconv.Atoi(Id.(string))
			yieldV := &IterateAdminObjects{
				Model: model,
				Id: uint(IdN),
				RenderContext: &FormRenderContext{Model: model},
			}
			chnl <- yieldV
		}
	}()
	return chnl
}

type ISortInterface interface {
	Order(afo IAdminFilterObjects)
}

type ISortBy interface {
	Sort (afo IAdminFilterObjects, direction int)
}

type SortBy struct {
	ISortBy
	Direction int // -1 descending order, 1 ascending order
	Field *Field
}

func (sb *SortBy) Sort(afo IAdminFilterObjects, direction int) {
	sortBy := sb.Field.DBName
	if direction == -1 {
		sortBy += " desc"
	}
	afo.SetPaginatedQuerySet(afo.GetPaginatedQuerySet().Order(sortBy))
}

type ListDisplayRegistry struct {
	ListDisplayFields map[string]*ListDisplay
	MaxOrdering int
	Prefix string
	Placement string
}

func (ldr *ListDisplayRegistry) GetFieldsCount() int {
	return len(ldr.ListDisplayFields)
}

func (ldr *ListDisplayRegistry) SetPrefix(prefix string) {
	ldr.Prefix = prefix
	for _, ld := range ldr.ListDisplayFields {
		ld.SetPrefix(prefix)
	}
}

func (ldr *ListDisplayRegistry) ClearAllFields() {
	ldr.MaxOrdering = 0
	ldr.ListDisplayFields = make(map[string]*ListDisplay)
}

func (ldr *ListDisplayRegistry) IsThereAnyEditable() bool {
	for ld := range ldr.GetAllFields() {
		if ld.IsEditable {
			return true
		}
	}
	return false
}
func (ldr *ListDisplayRegistry) AddField(ld *ListDisplay) {
	ldr.ListDisplayFields[ld.DisplayName] = ld
	ldr.MaxOrdering = int(math.Max(float64(ldr.MaxOrdering + 1), float64(ld.Ordering + 1)))
	ld.Ordering = ldr.MaxOrdering
}

func (ldr *ListDisplayRegistry) BuildFormForListEditable(adminContext IAdminContext, ID uint, model interface{}) *FormListEditable {
	return NewFormListEditableFromListDisplayRegistry(adminContext, ldr.Prefix, ID, model, ldr)
}

func (ldr *ListDisplayRegistry) BuildListEditableFormForNewModel(adminContext IAdminContext, ID string, model interface{}) *FormListEditable {
	return NewFormListEditableForNewModelFromListDisplayRegistry(adminContext, ldr.Prefix, ID, model, ldr)
}

func (ldr *ListDisplayRegistry) GetAllFields() <- chan *ListDisplay {
	chnl := make(chan *ListDisplay)
	go func() {
		defer close(chnl)
		dFields := make([]*ListDisplay, 0)
		for _, dField := range ldr.ListDisplayFields {
			dFields = append(dFields, dField)
		}
		sort.Slice(dFields, func(i, j int) bool {
			if dFields[i].Ordering == dFields[j].Ordering {
				return dFields[i].DisplayName < dFields[j].DisplayName
			}
			return dFields[i].Ordering < dFields[j].Ordering
		})
		for _, dField := range dFields {
			chnl <- dField
		}
	}()
	return chnl
}

func (ldr *ListDisplayRegistry) GetFieldByDisplayName(displayName string) (*ListDisplay, error) {
	listField, exists := ldr.ListDisplayFields[displayName]
	if !exists {
		return nil, fmt.Errorf("found no display field with name %s", displayName)
	}
	return listField, nil
}

type AdminModelActionRegistry struct {
	AdminModelActions map[string]*AdminModelAction
}

func (amar *AdminModelActionRegistry) AddModelAction(ma *AdminModelAction) {
	amar.AdminModelActions[ma.SlugifiedActionName] = ma
}

func (amar *AdminModelActionRegistry) IsThereAnyActions() bool {
	return len(amar.AdminModelActions) > 0
}

func (amar *AdminModelActionRegistry) GetAllModelActions() <- chan *AdminModelAction {
	chnl := make(chan *AdminModelAction)
	go func() {
		defer close(chnl)
		mActions := make([]*AdminModelAction, 0)
		for _, mAction := range amar.AdminModelActions {
			mActions = append(mActions, mAction)
		}
		sort.Slice(mActions, func(i, j int) bool {
			return mActions[i].ActionName < mActions[j].ActionName
		})
		for _, mAction := range mActions {
			chnl <- mAction
		}
	}()
	return chnl
}

func (amar *AdminModelActionRegistry) GetModelActionByName(actionName string) (*AdminModelAction, error) {
	mAction, exists := amar.AdminModelActions[actionName]
	if !exists {
		return nil, fmt.Errorf("found no model action with name %s", actionName)
	}
	return mAction, nil
}

type IListDisplayInterface interface {
	GetValue (m interface{}) string
}

type ListDisplay struct {
	IListDisplayInterface
	DisplayName string
	Field *Field
	ChangeLink bool
	Ordering int
	SortBy *SortBy
	Populate func (m interface{}) string
	MethodName string
	IsEditable bool
	Prefix string
}

func (ld *ListDisplay) SetPrefix(prefix string) {
	ld.Prefix = prefix
}

func (ld *ListDisplay) GetOrderingName(initialOrdering []string) string {
	for _, part := range initialOrdering {
		negativeOrdering := false
		if strings.HasPrefix(part, "-") {
			part = part[1:]
			negativeOrdering = true
		}
		if part == ld.DisplayName {
			if negativeOrdering {
				return ld.DisplayName
			}
			return "-" + ld.DisplayName
		}
	}

	return ld.DisplayName
}

func (ld *ListDisplay) IsEligibleForOrdering() bool {
	return ld.SortBy != nil
}

func (ld *ListDisplay) GetValue(m interface{}, forExportP ...bool) string {
	forExport := false
	if len(forExportP) > 0 {
		forExport = forExportP[0]
	}
	if ld.MethodName != "" {
		values := reflect.ValueOf(m).MethodByName(ld.MethodName).Call([]reflect.Value{})
		return values[0].String()
	}
	if ld.Populate != nil {
		return ld.Populate(m)
	}
	if ld.Field.FieldConfig.Widget.GetPopulate() != nil {
		return TransformValueForListDisplay(ld.Field.FieldConfig.Widget.GetPopulate()(m, ld.Field))
	}
	if ld.Field.FieldConfig.Widget.IsValueConfigured() {
		return TransformValueForListDisplay(ld.Field.FieldConfig.Widget.GetValue())
	}
	gormModelV := reflect.Indirect(reflect.ValueOf(m))
	if reflect.ValueOf(m).IsZero() || gormModelV.IsZero() { // || gormModelV.FieldByName(ld.Field.Name).IsZero()
		return ""
	}
	return TransformValueForListDisplay(gormModelV.FieldByName(ld.Field.Name).Interface(), forExport)
}

func NewListDisplay(field *Field) *ListDisplay {
	displayName := ""
	if field != nil {
		displayName = field.DisplayName
	}
	return &ListDisplay{
		DisplayName: displayName, Field: field, ChangeLink: true,
		SortBy: &SortBy{Field: field, Direction: 1},
	}
}

type IListFilterInterface interface {
	FilterQs (afo IAdminFilterObjects, filterString string)
}

type ListFilter struct {
	IListFilterInterface
	Title string
	UrlFilteringParam string
	OptionsToShow []*FieldChoice
	FetchOptions func(m interface{}) []*FieldChoice
	CustomFilterQs func(afo IAdminFilterObjects, filterString string)
	Template string
	Ordering int
}

func (lf *ListFilter) FilterQs (afo IAdminFilterObjects, filterString string) {
	if lf.CustomFilterQs != nil {
		lf.CustomFilterQs(afo, filterString)
	} else {
		statement := &gorm.Statement{DB: afo.GetUadminDatabase().Db}
		statement.Parse(afo.GetCurrentModel())
		schema1 := statement.Schema
		operatorContext := FilterGormModel(afo.GetUadminDatabase().Adapter, afo.GetFullQuerySet(), schema1, []string{filterString}, afo.GetCurrentModel())
		afo.SetFullQuerySet(operatorContext.Tx)
		operatorContext = FilterGormModel(afo.GetUadminDatabase().Adapter, afo.GetPaginatedQuerySet(), schema1, []string{filterString}, afo.GetCurrentModel())
		afo.SetPaginatedQuerySet(operatorContext.Tx)
	}
}

func (lf *ListFilter) IsItActive (fullUrl *url.URL) bool {
	return strings.Contains(fullUrl.String(), lf.UrlFilteringParam)
}

func (lf *ListFilter) GetURLToClearFilter(fullUrl *url.URL) string {
	clonedUrl := CloneNetUrl(fullUrl)
	qs := clonedUrl.Query()
	qs.Del(lf.UrlFilteringParam)
	clonedUrl.RawQuery = qs.Encode()
	return clonedUrl.String()
}

func (lf *ListFilter) IsThatOptionActive(option *FieldChoice, fullUrl *url.URL) bool {
	qs := fullUrl.Query()
	value := qs.Get(lf.UrlFilteringParam)
	if value != "" {
		optionValue := TransformValueForListDisplay(option.Value)
		if optionValue == value {
			return true
		}
	}
	return false
}

func (lf *ListFilter) GetURLForOption(option *FieldChoice, fullUrl *url.URL) string {
	clonedUrl := CloneNetUrl(fullUrl)
	qs := clonedUrl.Query()
	qs.Set(lf.UrlFilteringParam, TransformValueForListDisplay(option.Value))
	clonedUrl.RawQuery = qs.Encode()
	return clonedUrl.String()
}

type ListFilterRegistry struct {
	ListFilter []*ListFilter
}

type ListFilterList [] *ListFilter

func (apl ListFilterList) Len() int { return len(apl) }
func (apl ListFilterList) Less(i, j int) bool {
	return apl[i].Ordering < apl[j].Ordering
}
func (apl ListFilterList) Swap(i, j int){ apl[i], apl[j] = apl[j], apl[i] }



func (lfr *ListFilterRegistry) Iterate() <- chan *ListFilter {
	chnl := make(chan *ListFilter)
	go func() {
		lfList := make(ListFilterList, 0)
		defer close(chnl)
		for _, lF := range lfr.ListFilter {
			lfList = append(lfList, lF)
		}
		sort.Sort(lfList)
		for _, lf := range lfList {
			chnl <- lf
		}
	}()
	return chnl
}

func (lfr *ListFilterRegistry) IsEmpty() bool {
	return !(len(lfr.ListFilter) > 0)
}

func (lfr *ListFilterRegistry) Add(lf *ListFilter) {
	lfr.ListFilter = append(lfr.ListFilter, lf)
}

type ISearchFieldInterface interface {
	Search (afo IAdminFilterObjects, searchString string)
}

type SearchField struct {
	ISearchFieldInterface
	Field *Field
	CustomSearch func(afo IAdminFilterObjects, searchString string)
}

func (sf *SearchField) Search(afo IAdminFilterObjects, searchString string) {
	if sf.CustomSearch != nil {
		sf.CustomSearch(afo, searchString)
	} else {
		operator := ExactGormOperator{}
		gormOperatorContext := NewGormOperatorContext(afo.GetFullQuerySet(), afo.GetCurrentModel())
		operator.Build(afo.GetUadminDatabase().Adapter, gormOperatorContext, sf.Field, searchString)
		afo.SetFullQuerySet(gormOperatorContext.Tx)
		gormOperatorContext = NewGormOperatorContext(afo.GetPaginatedQuerySet(), afo.GetCurrentModel())
		operator.Build(afo.GetUadminDatabase().Adapter, gormOperatorContext, sf.Field, searchString)
		afo.SetPaginatedQuerySet(gormOperatorContext.Tx)
	}
}

type PaginationType string

var LimitPaginationType PaginationType = "limit"
var CursorPaginationType PaginationType = "cursor"

type IPaginationInterface interface {
	Paginate (afo IAdminFilterObjects)
}

type Paginator struct {
	IPaginationInterface
	PerPage int
	AllowEmptyFirstPage bool
	ShowLastPageOnPreviousPage bool
	Count int
	NumPages int
	Offset int
	Template string
	PaginationType PaginationType
}

func (p *Paginator) Paginate(afo IAdminFilterObjects) {

}

type DisplayFilterOption struct {
	FilterField string
	FilterValue string
	DisplayAs string
}

type FilterOption struct {
	FieldName string
	FetchOptions func(afo IAdminFilterObjects) []*DisplayFilterOption
}

type FilterOptionsRegistry struct {
	FilterOption []*FilterOption
}

func (for1 *FilterOptionsRegistry) AddFilterOption(fo *FilterOption) {
	for1.FilterOption = append(for1.FilterOption, fo)
}

func (for1 *FilterOptionsRegistry) GetAll() <- chan *FilterOption {
	chnl := make(chan *FilterOption)
	go func() {
		defer close(chnl)
		for _, fo := range for1.FilterOption {
			chnl <- fo
		}
	}()
	return chnl
}

func NewFilterOptionsRegistry() *FilterOptionsRegistry {
	return &FilterOptionsRegistry{FilterOption: make([]*FilterOption, 0)}
}

func NewFilterOption() *FilterOption {
	return &FilterOption{}
}

func FetchOptionsFromGormModelFromDateTimeField(afo IAdminFilterObjects, filterOptionField string) []*DisplayFilterOption {
	ret := make([]*DisplayFilterOption, 0)
	uadminDatabase := NewUadminDatabase()
	defer uadminDatabase.Close()
	filterString := uadminDatabase.Adapter.GetStringToExtractYearFromField(filterOptionField)
	rows, _ := afo.GetInitialQuerySet().Select(filterString + " as year, count(*) as total").Group(filterString).Rows()
	var filterValue uint
	var filterCount uint
	for rows.Next() {
		rows.Scan(&filterValue, &filterCount)
		filterString := strconv.Itoa(int(filterValue))
		ret = append(ret, &DisplayFilterOption{
			FilterField: filterOptionField,
			FilterValue: filterString,
			DisplayAs: filterString,
		})
	}
	if len(ret) < 2 {
		ret = make([]*DisplayFilterOption, 0)
		filterString := uadminDatabase.Adapter.GetStringToExtractMonthFromField(filterOptionField)
		rows, _ := afo.GetInitialQuerySet().Select(filterString + " as month, count(*) as total").Group(filterString).Rows()
		var filterValue uint
		var filterCount uint
		for rows.Next() {
			rows.Scan(&filterValue, &filterCount)
			filterString := strconv.Itoa(int(filterValue))
			filteredMonth, _ := strconv.Atoi(filterString)
			ret = append(ret, &DisplayFilterOption{
				FilterField: filterOptionField,
				FilterValue: filterString,
				DisplayAs: time.Month(filteredMonth).String(),
			})
		}
	}
	return ret
}

var CurrentAdminPageRegistry *AdminPageRegistry

type AdminBreadcrumb struct {
	Name string
	Url string
	IsActive bool
	Icon string
}

type AdminBreadCrumbsRegistry struct {
	BreadCrumbs []*AdminBreadcrumb
}

func (abcr *AdminBreadCrumbsRegistry) AddBreadCrumb(breadcrumb *AdminBreadcrumb) {
	if len(abcr.BreadCrumbs) != 0 {
		abcr.BreadCrumbs[0].IsActive = false
	} else {
		breadcrumb.IsActive = true
	}
	abcr.BreadCrumbs = append(abcr.BreadCrumbs, breadcrumb)
}

func (abcr *AdminBreadCrumbsRegistry) GetAll() <- chan *AdminBreadcrumb {
	chnl := make(chan *AdminBreadcrumb)
	go func() {
		defer close(chnl)
		for _, adminBreadcrumb := range abcr.BreadCrumbs {
			chnl <- adminBreadcrumb
		}
	}()
	return chnl
}

func NewAdminBreadCrumbsRegistry() *AdminBreadCrumbsRegistry {
	ret := &AdminBreadCrumbsRegistry{BreadCrumbs: make([]*AdminBreadcrumb, 0)}
	return ret
}

func NewListDisplayRegistry() *ListDisplayRegistry {
	ret := &ListDisplayRegistry{
		ListDisplayFields: make(map[string]*ListDisplay),
	}
	return ret
}

func NewListDisplayRegistryFromGormModelForInlines(modelI interface{}) *ListDisplayRegistry {
	ret := &ListDisplayRegistry{
		ListDisplayFields: make(map[string]*ListDisplay),
	}
	uadminDatabase := NewUadminDatabaseWithoutConnection()
	stmt := &gorm.Statement{DB: uadminDatabase.Db}
	stmt.Parse(modelI)
	gormModelV := reflect.Indirect(reflect.ValueOf(modelI))
	for _, field := range stmt.Schema.Fields {
		uadminTag := field.Tag.Get("uadmin")
		if !strings.Contains(uadminTag, "inline") && field.Name != "ID" {
			continue
		}
		uadminField := NewUadminFieldFromGormField(gormModelV, field, nil, true)
		ld := NewListDisplay(uadminField)
		if field.Name != "ID" {
			ld.IsEditable = true
		}
		ret.AddField(ld)
	}
	return ret
}

func NewListDisplayRegistryFromGormModel(modelI interface{}) *ListDisplayRegistry {
	if modelI == nil {
		return nil
	}
	ret := &ListDisplayRegistry{
		ListDisplayFields: make(map[string]*ListDisplay),
	}
	uadminDatabase := NewUadminDatabaseWithoutConnection()
	stmt := &gorm.Statement{DB: uadminDatabase.Db}
	stmt.Parse(modelI)
	gormModelV := reflect.Indirect(reflect.ValueOf(modelI))
	for _, field := range stmt.Schema.Fields {
		uadminTag := field.Tag.Get("uadmin")
		if !strings.Contains(uadminTag, "list") && field.Name != "ID" {
			continue
		}
		uadminField := NewUadminFieldFromGormField(gormModelV, field, nil, true)
		ret.AddField(NewListDisplay(uadminField))
	}
	return ret
}

func NewSearchFieldRegistryFromGormModel(modelI interface{}) *SearchFieldRegistry {
	if modelI == nil {
		return nil
	}
	ret := &SearchFieldRegistry{Fields: make([]*SearchField, 0)}
	uadminDatabase := NewUadminDatabaseWithoutConnection()
	stmt := &gorm.Statement{DB: uadminDatabase.Db}
	stmt.Parse(modelI)
	gormModelV := reflect.Indirect(reflect.ValueOf(modelI))
	for _, field := range stmt.Schema.Fields {
		uadminTag := field.Tag.Get("uadmin")
		if !strings.Contains(uadminTag, "search") && field.Name != "ID" {
			continue
		}
		uadminField := NewUadminFieldFromGormField(gormModelV, field, nil, true)
		searchField := &SearchField{
			Field: uadminField,
		}
		ret.AddField(searchField)
	}
	return ret
}