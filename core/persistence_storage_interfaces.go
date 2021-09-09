package core

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

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
	GetCurrentDB() *gorm.DB
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

