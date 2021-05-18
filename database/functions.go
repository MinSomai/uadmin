package database

import (
	"database/sql"
	"fmt"
	models2 "github.com/uadmin/uadmin/blueprint/abtest/models"
	"github.com/uadmin/uadmin/colors"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/dialect"
	"github.com/uadmin/uadmin/metrics"
	"github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/utils"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// Save saves the object in the database
func Save(a interface{}) (err error) {
	encryptRecord(a)
	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		err = dialect.GetDB().Save(a).Error
		for fmt.Sprint(err) == "database is locked" {
			time.Sleep(time.Millisecond * 100)
			err = dialect.GetDB().Save(a).Error
		}
	})
	if err != nil {
		utils.Trail(utils.ERROR, "DB error in Save(%v). %s", model.GetModelName(a), err.Error())
		return err
	}
	err = customSave(a)
	if err != nil {
		utils.Trail(utils.ERROR, "DB error in customSave(%v). %s", model.GetModelName(a), err.Error())
		return err
	}
	return nil
}


func decryptArray(a interface{}) {
	model1, _ := model.NewModel(model.GetModelName(a), false)
	if schema, ok := model.GetSchema(model1); ok {
		for _, f := range schema.Fields {
			if f.Encrypt {
				// TODO: Decrypt
				allArray := reflect.ValueOf(a)
				for i := 0; i < allArray.Elem().Len(); i++ {
					encryptedValue := allArray.Elem().Index(i).FieldByName(f.Name).String()
					decryptedValue, _ := utils.Decrypt(EncryptKey, encryptedValue)
					allArray.Elem().Index(i).FieldByName(f.Name).Set(reflect.ValueOf(decryptedValue))
				}
			}
		}
	}
}

func encryptArray(a interface{}) {
	model1, _ := model.NewModel(model.GetModelName(a), false)
	if schema, ok := model.GetSchema(model1); ok {
		for _, f := range schema.Fields {
			if f.Encrypt {
				allArray := reflect.ValueOf(a)
				for i := 0; i < allArray.Elem().Len(); i++ {
					encryptedValue := allArray.Elem().Index(i).FieldByName(f.Name).String()
					decryptedValue, _ := utils.Encrypt(EncryptKey, encryptedValue)
					allArray.Elem().Index(i).FieldByName(f.Name).Set(reflect.ValueOf(decryptedValue))
				}
			}
		}
	}
}

func decryptRecord(a interface{}) {
	model1, _ := model.NewModel(model.GetModelName(a), false)
	if schema, ok := model.GetSchema(model1); ok {
		for _, f := range schema.Fields {
			if f.Encrypt {
				recordValue := reflect.ValueOf(a)
				encryptedValue := recordValue.Elem().FieldByName(f.Name).String()
				decryptedValue, _ := utils.Decrypt(EncryptKey, encryptedValue)
				recordValue.Elem().FieldByName(f.Name).Set(reflect.ValueOf(decryptedValue))
			}
		}
	}
}

func encryptRecord(a interface{}) {
	model1, _ := model.NewModel(model.GetModelName(a), false)
	if schema, ok := model.GetSchema(model1); ok {
		for _, f := range schema.Fields {
			if f.Encrypt {
				recordValue := reflect.ValueOf(a)
				encryptedValue := recordValue.Elem().FieldByName(f.Name).String()
				decryptedValue, _ := utils.Encrypt(EncryptKey, encryptedValue)
				recordValue.Elem().FieldByName(f.Name).Set(reflect.ValueOf(decryptedValue))
			}
		}
	}
}

func customSave(m interface{}) (err error) {
	a := m
	t := reflect.TypeOf(a)
	if t.Kind() == reflect.Ptr {
		a = reflect.ValueOf(m).Elem().Interface()
		t = reflect.TypeOf(a)
	}
	value := reflect.ValueOf(a)
	dialect1 := dialect.GetDialectForDb()
	sqlDialectStrings := dialect1.GetSqlDialectStrings()
	for i := 0; i < t.NumField(); i++ {
		// Check if there is any m2m fields
		if t.Field(i).Type.Kind() == reflect.Slice && t.Field(i).Type.Elem().Kind() == reflect.Struct {
			table1 := strings.ToLower(t.Name())
			table2 := strings.ToLower(t.Field(i).Type.Elem().Name())

			// Delete existing records
			sql := sqlDialectStrings["deleteM2M"]
			sql = strings.Replace(sql, "{TABLE1}", table1, -1)
			sql = strings.Replace(sql, "{TABLE2}", table2, -1)
			sql = strings.Replace(sql, "{TABLE1_ID}", fmt.Sprint(GetID(value)), -1)

			metrics.TimeMetric("uadmin/db/duration", 1000, func() {
				err = dialect.GetDB().Exec(sql).Error
				for fmt.Sprint(err) == "database is locked" {
					time.Sleep(time.Millisecond * 100)
					err = dialect.GetDB().Exec(sql).Error
				}
			})
			if err != nil {
				utils.Trail(utils.ERROR, "Unable to delete m2m records. %s", err)
				utils.Trail(utils.ERROR, sql)
				return err
			}
			// Insert records
			for index := 0; index < value.Field(i).Len(); index++ {
				sql := sqlDialectStrings["insertM2M"]
				sql = strings.Replace(sql, "{TABLE1}", table1, -1)
				sql = strings.Replace(sql, "{TABLE2}", table2, -1)
				sql = strings.Replace(sql, "{TABLE1_ID}", fmt.Sprint(GetID(value)), -1)
				sql = strings.Replace(sql, "{TABLE2_ID}", fmt.Sprint(GetID(value.Field(i).Index(index))), -1)

				metrics.TimeMetric("uadmin/db/duration", 1000, func() {
					err = dialect.GetDB().Exec(sql).Error
					for fmt.Sprint(err) == "database is locked" {
						time.Sleep(time.Millisecond * 100)
						err = dialect.GetDB().Exec(sql).Error
					}
				})
				if err != nil {
					utils.Trail(utils.ERROR, "Unable to insert m2m records. %s", err)
					utils.Trail(utils.ERROR, sql)
					return err
				}
			}

		}
	}
	return nil
}

// EncryptKey is a key for encryption and decryption of data in the DB.
var EncryptKey = []byte{}

// GetID !
func GetID(m reflect.Value) uint {
	if m.Kind() == reflect.Ptr {
		return uint(m.Elem().FieldByName("ID").Uint())
	}
	return uint(m.FieldByName("ID").Uint())
}

var DbOK = false

// initializeDB opens the connection the DB
func initializeDB(a ...interface{}) {
	// Open the connection the the DB
	db := dialect.GetDB()

	// Migrate schema
	for i, model := range a {
		utils.Trail(utils.WORKING, "Initializing DB: [%s%d/%d%s]", colors.FGGreenB, i+1, len(a), colors.FGNormal)
		db.AutoMigrate(model)
		customMigration(model)
	}
	utils.Trail(utils.OK, "Initializing DB: [%s%d/%d%s]", colors.FGGreenB, len(a), len(a), colors.FGNormal)
}

func customMigration(a interface{}) (err error) {
	t := reflect.TypeOf(a)
	currentdialect := dialect.GetDialectForDb()
	sqlDialectStrings := currentdialect.GetSqlDialectStrings()
	for i := 0; i < t.NumField(); i++ {
		// Check if there is any m2m fields
		if t.Field(i).Type.Kind() == reflect.Slice && t.Field(i).Type.Elem().Kind() == reflect.Struct {
			table1 := strings.ToLower(t.Name())
			table2 := strings.ToLower(t.Field(i).Type.Elem().Name())

			//Check if the table is created for the m2m field
			if !dialect.GetDB().Migrator().HasTable(table1 + "_" + table2) {
				sql1 := sqlDialectStrings["createM2MTable"]
				sql1 = strings.Replace(sql1, "{TABLE1}", table1, -1)
				sql1 = strings.Replace(sql1, "{TABLE2}", table2, -1)
				err = dialect.GetDB().Exec(sql1).Error
				if err != nil {
					utils.Trail(utils.ERROR, "Unable to create M2M table. %s", err)
					utils.Trail(utils.ERROR, sql1)
					return err
				}
			}
		}
	}
	return err
}

func InitializeDbSettingsFromConfig(config *config.UadminConfig) {
	dialect.CurrentDatabaseSettings = config.D.Db.Default
}

func createDB() error {
	return fmt.Errorf("CreateDB: Unknown database type " + dialect.CurrentDatabaseSettings.Type)
}

// ClearDB clears the db object
func ClearDB() {
	dialect.Db = nil
}

// All fetches all object in the database
func All(a interface{}) (err error) {
	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		err = dialect.GetDB().Find(a).Error
		for fmt.Sprint(err) == "database is locked" {
			time.Sleep(time.Millisecond * 100)
			err = dialect.GetDB().Find(a).Error
		}
	})
	if err != nil {
		utils.Trail(utils.ERROR, "DB error in All(%v). %s", model.GetModelName(a), err.Error())
		return err
	}
	decryptArray(a)
	return nil
}

// Get fetches the first record from the database matching query and args
func Get(a interface{}, query interface{}, args ...interface{}) (err error) {
	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		err = dialect.GetDB().Where(query, args...).First(a).Error
		for fmt.Sprint(err) == "database is locked" {
			time.Sleep(time.Millisecond * 100)
			err = dialect.GetDB().Where(query, args...).First(a).Error
		}
	})

	if err != nil {
		if err.Error() != "record not found" {
			utils.Trail(utils.ERROR, "DB error in Get(%s)-(%v). %s", model.GetModelName(a), a, err.Error())
		}
		return err
	}

	err = customGet(a)
	if err != nil {
		utils.Trail(utils.ERROR, "DB error in customGet(%v). %s", model.GetModelName(a), err.Error())
		return err
	}
	decryptRecord(a)
	return nil
}

// GetABTest is like Get function but implements AB testing for the results
func GetABTest(r *http.Request, a interface{}, query interface{}, args ...interface{}) (err error) {
	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		Get(a, query, args...)
	})

	// Check if there are any active A/B tests for any field in this model
	abt := models2.GetABT(r)
	modelName := model.GetModelName(a)
	models2.AbTestsMutex.Lock()
	for k, v := range models2.ModelABTests {
		if strings.HasPrefix(k, modelName+"__") && strings.HasSuffix(k, "__"+fmt.Sprint(GetID(reflect.ValueOf(a)))) {
			if len(v) != 0 {
				index := abt % len(v)
				fName := model.Schema[modelName].Fields[v[index].fname].Name

				// TODO: Support more data types
				switch model.Schema[modelName].Fields[v[index].fname].Type {
				case preloaded.CSTRING:
					reflect.ValueOf(a).Elem().FieldByName(fName).SetString(v[index].v)
				case preloaded.CIMAGE:
					reflect.ValueOf(a).Elem().FieldByName(fName).SetString(v[index].v)
				}

				// Increment impressions
				v[index].imp++
				models2.ModelABTests[k] = v
			}
		}
	}
	models2.AbTestsMutex.Unlock()
	return nil
}

// GetStringer fetches the first record from the database matching query and args
// and get only fields tagged with `stringer` tag. If no field has `stringer` tag
// then it gets all the fields
func GetStringer(a interface{}, query interface{}, args ...interface{}) (err error) {
	stringers := []string{}
	modelName := model.GetModelName(a)
	for _, f := range model.Schema[modelName].Fields {
		if f.Stringer {
			stringers = append(stringers, dialect.GetDB().Config.NamingStrategy.ColumnName("", f.Name))
		}
	}
	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		if len(stringers) == 0 {
			err = dialect.GetDB().Where(query, args...).First(a).Error
			for fmt.Sprint(err) == "database is locked" {
				time.Sleep(time.Millisecond * 100)
				err = dialect.GetDB().Where(query, args...).First(a).Error
			}
		} else {
			err = dialect.GetDB().Select(stringers).Where(query, args...).First(a).Error
			for fmt.Sprint(err) == "database is locked" {
				time.Sleep(time.Millisecond * 100)
				err = dialect.GetDB().Select(stringers).Where(query, args...).First(a).Error
			}
		}
	})
	if err != nil {
		if err.Error() != "record not found" {
			utils.Trail(utils.ERROR, "DB error in Get(%s)-(%v). %s", model.GetModelName(a), a, err.Error())
		}
		return err
	}

	//err = customGet(a)
	if err != nil {
		utils.Trail(utils.ERROR, "DB error in customGet(%v). %s", model.GetModelName(a), err.Error())
		return err
	}
	decryptRecord(a)
	return nil
}

// GetForm fetches the first record from the database matching query and args
// where it selects only visible fields in the form based on given schema
func GetForm(a interface{}, s *model.ModelSchema, query interface{}, args ...interface{}) (err error) {
	// get a list of visible fields
	columnList := []string{}
	dialect1 := dialect.GetDialectForDb()
	m2mList := []string{}
	for _, f := range s.Fields {
		if !f.Hidden {
			if f.Type == preloaded.CM2M {
				m2mList = append(m2mList, f.ColumnName)
			} else if f.Type == preloaded.CFK {
				columnList = append(columnList, dialect1.Quote(f.ColumnName+"_id"))
				// } else if f.IsMethod {
			} else {
				columnList = append(columnList, dialect1.Quote(f.ColumnName))
			}
		}
	}

	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		err = dialect.GetDB().Select(columnList).Where(query, args...).First(a).Error
		for fmt.Sprint(err) == "database is locked" {
			time.Sleep(time.Millisecond * 100)
			err = dialect.GetDB().Select(columnList).Where(query, args...).First(a).Error
		}
	})

	if err != nil {
		if err.Error() != "record not found" {
			utils.Trail(utils.ERROR, "DB error in Get(%s)-(%v). %s", model.GetModelName(a), a, err.Error())
		}
		return err
	}

	err = customGet(a, m2mList...)
	if err != nil {
		utils.Trail(utils.ERROR, "DB error in customGet(%v). %s", model.GetModelName(a), err.Error())
		return err
	}
	decryptRecord(a)
	return nil
}

func customGet(m interface{}, m2m ...string) (err error) {
	a := m
	t := reflect.TypeOf(a)
	var ignore bool
	dialect1 := dialect.GetDialectForDb()
	sqlDialectStrings := dialect1.GetSqlDialectStrings()
	if t.Kind() == reflect.Ptr {
		a = reflect.ValueOf(m).Elem().Interface()
		t = reflect.TypeOf(a)
	}
	value := reflect.ValueOf(a)
	for i := 0; i < t.NumField(); i++ {
		ignore = false
		if len(m2m) != 0 {
			ignore = true
			for _, fName := range m2m {
				if fName == t.Field(i).Name {
					ignore = false
					break
				}
			}
		}
		if ignore {
			continue
		}

		// Skip private fields
		if t.Field(i).Anonymous || t.Field(i).Name[0:1] != strings.ToUpper(t.Field(i).Name[0:1]) {
			continue
		}

		// Check if there is any m2m fields
		if t.Field(i).Type.Kind() == reflect.Slice && t.Field(i).Type.Elem().Kind() == reflect.Struct {
			table1 := strings.ToLower(t.Name())
			table2 := strings.ToLower(t.Field(i).Type.Elem().Name())

			sqlSelect := sqlDialectStrings["selectM2M"]
			sqlSelect = strings.Replace(sqlSelect, "{TABLE1}", table1, -1)
			sqlSelect = strings.Replace(sqlSelect, "{TABLE2}", table2, -1)
			sqlSelect = strings.Replace(sqlSelect, "{TABLE1_ID}", fmt.Sprint(GetID(value)), -1)

			var rows *sql.Rows
			rows, err = dialect.GetDB().Raw(sqlSelect).Rows()
			if err != nil {
				utils.Trail(utils.ERROR, "Unable to get m2m records. %s", err)
				utils.Trail(utils.ERROR, sqlSelect)
				return err
			}
			defer rows.Close()
			var fkID uint
			tmpDst := reflect.New(reflect.SliceOf(t.Field(i).Type.Elem())).Elem()
			for rows.Next() {
				rows.Scan(&fkID)
				tempModel := reflect.New(t.Field(i).Type.Elem()).Elem()
				Get(tempModel.Addr().Interface(), "id = ?", fkID)
				tmpDst = reflect.Append(tmpDst, tempModel)
			}
			reflect.ValueOf(m).Elem().Field(i).Set(tmpDst)
		}
	}
	return nil
}

// Filter fetches records from the database
func Filter(a interface{}, query interface{}, args ...interface{}) (err error) {
	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		err = dialect.GetDB().Where(query, args...).Find(a).Error
		for fmt.Sprint(err) == "database is locked" {
			time.Sleep(time.Millisecond * 100)
			err = dialect.GetDB().Where(query, args...).Find(a).Error
		}
	})

	if err != nil {
		utils.Trail(utils.ERROR, "DB error in Filter(%v). %s\n", model.GetModelName(a), err.Error())
		return err
	}
	decryptArray(a)
	return nil
}

// Preload fills the data from foreign keys into structs. You can pass in preload alist of fields
// to be preloaded. If nothing is passed, every foreign key is preloaded
func Preload(a interface{}, preload ...string) (err error) {
	modelName := strings.ToLower(reflect.TypeOf(a).Elem().Name())
	model1, _ := model.NewModel(modelName, false)
	if len(preload) == 0 {
		if schema, ok := model.GetSchema(model1); ok {
			for _, f := range schema.Fields {
				if f.Type == "fk" {
					preload = append(preload, f.Name)
				}
			}
		} else {
			utils.Trail(utils.ERROR, "DB.Preload No model named %s", modelName)
			return fmt.Errorf("DB.Preload No model named %s", modelName)
		}
	}
	value := reflect.ValueOf(a).Elem()
	for _, p := range preload {
		fkType := value.FieldByName(p).Type().Name()
		if value.FieldByName(p).Type().Kind() == reflect.Ptr {
			fkType = value.FieldByName(p).Type().Elem().Name()
		}
		fieldStruct, _ := model.NewModel(strings.ToLower(fkType), true)
		metrics.TimeMetric("uadmin/db/duration", 1000, func() {
			err = dialect.GetDB().Where("id = ?", value.FieldByName(p+"ID").Interface()).First(fieldStruct.Interface()).Error
			for fmt.Sprint(err) == "database is locked" {
				time.Sleep(time.Millisecond * 100)
				err = dialect.GetDB().Where("id = ?", value.FieldByName(p+"ID").Interface()).First(fieldStruct.Interface()).Error
			}
		})

		//		err = Get(fieldStruct.Interface(), "id = ?", value.FieldByName(p+"ID").Interface())
		if err != nil && err.Error() != "record not found" {
			utils.Trail(utils.ERROR, "DB error in Preload(%s).%s %s\n", modelName, p, err.Error())
			return err
		}
		if GetID(fieldStruct) != 0 {
			if value.FieldByName(p).Type().Kind() == reflect.Ptr {
				value.FieldByName(p).Set(fieldStruct)
			} else {
				value.FieldByName(p).Set(fieldStruct.Elem())
			}
		}
	}
	return customGet(a)
}

// Delete records from database
func Delete(a interface{}) (err error) {
	// Sanity Check for ID = 0
	if GetID(reflect.ValueOf(a)) == 0 {
		return nil
	}
	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		err = dialect.GetDB().Delete(a).Error
		for fmt.Sprint(err) == "database is locked" {
			time.Sleep(time.Millisecond * 100)
			err = dialect.GetDB().Delete(a).Error
		}
	})

	if err != nil {
		utils.Trail(utils.ERROR, "DB error in Delete(%v). %s\n", model.GetModelName(a), err.Error())
		return err
	}
	return nil
}

// DeleteList deletes multiple records from database
func DeleteList(a interface{}, query interface{}, args ...interface{}) (err error) {
	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		err = dialect.GetDB().Where(query, args...).Delete(a).Error
		for fmt.Sprint(err) == "database is locked" {
			time.Sleep(time.Millisecond * 100)
			err = dialect.GetDB().Where(query, args...).Delete(a).Error
		}
	})

	if err != nil {
		utils.Trail(utils.ERROR, "DB error in DeleteList(%v). %s\n", model.GetModelName(a), err.Error())
		return err
	}
	return nil
}

// FilterBuilder changes a map filter into a query
func FilterBuilder(params map[string]interface{}) (query string, args []interface{}) {
	keys := []string{}
	for key, value := range params {
		keys = append(keys, key)
		args = append(args, value)
	}
	query = strings.Join(keys, " AND ")
	return
}

// AdminPage !
func AdminPage(order string, asc bool, offset int, limit int, a interface{}, query interface{}, args ...interface{}) (err error) {
	dialect1 := dialect.GetDialectForDb()
	if order != "" {
		orderby := " desc"
		if asc {
			orderby = " asc"
		}
		order = dialect1.Quote(order)
		orderby += " "
		order += orderby
	} else {
		order = "id desc"
	}
	if limit > 0 {
		metrics.TimeMetric("uadmin/db/duration", 1000, func() {
			err = dialect.GetDB().Where(query, args...).Order(order).Offset(offset).Limit(limit).Find(a).Error
			for fmt.Sprint(err) == "database is locked" {
				time.Sleep(time.Millisecond * 100)
				err = dialect.GetDB().Where(query, args...).Order(order).Offset(offset).Limit(limit).Find(a).Error
			}
		})

		if err != nil {
			utils.Trail(utils.ERROR, "DB error in AdminPage(%v). %s\n", model.GetModelName(a), err.Error())
			return err
		}
		decryptArray(a)
		return nil
	}
	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		err = dialect.GetDB().Where(query, args...).Order(order).Find(a).Error
		for fmt.Sprint(err) == "database is locked" {
			time.Sleep(time.Millisecond * 100)
			err = dialect.GetDB().Where(query, args...).Order(order).Find(a).Error
		}
	})

	if err != nil {
		utils.Trail(utils.ERROR, "DB error in AdminPage(%v). %s\n", model.GetModelName(a), err.Error())
		return err
	}
	decryptArray(a)
	return nil
}

// FilterList fetches the all record from the database matching query and args
// where it selects only visible fields in the form based on given schema
func FilterList(s *model.ModelSchema, order string, asc bool, offset int, limit int, a interface{}, query interface{}, args ...interface{}) (err error) {
	dialect1 := dialect.GetDialectForDb()
	// get a list of visible fields
	columnList := []string{}
	for _, f := range s.Fields {
		if f.ListDisplay {
			if f.Type == preloaded.CFK {
				columnList = append(columnList, dialect1.Quote(dialect.GetDB().Config.NamingStrategy.ColumnName("", f.Name)+"_id"))
			} else if f.Type == preloaded.CM2M {
			} else if f.IsMethod {
			} else {
				columnList = append(columnList, dialect1.Quote(dialect.GetDB().Config.NamingStrategy.ColumnName("", f.Name)))
			}
		}
	}
	if order != "" {
		orderby := " desc"
		if asc {
			orderby = " asc"
		}
		order = dialect1.Quote(order)
		orderby += " "
		order += orderby
	} else {
		order = "id desc"
	}
	if limit > 0 {
		metrics.TimeMetric("uadmin/db/duration", 1000, func() {
			err = dialect.GetDB().Select(columnList).Where(query, args...).Order(order).Offset(offset).Limit(limit).Find(a).Error
			for fmt.Sprint(err) == "database is locked" {
				time.Sleep(time.Millisecond * 100)
				err = dialect.GetDB().Select(columnList).Where(query, args...).Order(order).Offset(offset).Limit(limit).Find(a).Error
			}
		})

		if err != nil {
			utils.Trail(utils.ERROR, "DB error in FilterList(%v) query:%s, args(%#v). %s\n", model.GetModelName(a), query, args, err.Error())
			return err
		}
		decryptArray(a)
		return nil
	}
	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		err = dialect.GetDB().Select(columnList).Where(query, args...).Order(order).Find(a).Error
		for fmt.Sprint(err) == "database is locked" {
			time.Sleep(time.Millisecond * 100)
			err = dialect.GetDB().Select(columnList).Where(query, args...).Order(order).Find(a).Error
		}
	})

	if err != nil {
		utils.Trail(utils.ERROR, "DB error in FilterList(%v) query:%s, args(%#v). %s\n", model.GetModelName(a), query, args, err.Error())
		return err
	}
	decryptArray(a)
	return nil
}

// Count return the count of records in a table based on a filter
func Count(a interface{}, query interface{}, args ...interface{}) int {
	var count int64
	var err error
	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		err = dialect.GetDB().Model(a).Where(query, args...).Count(&count).Error
		for fmt.Sprint(err) == "database is locked" {
			time.Sleep(time.Millisecond * 100)
			err = dialect.GetDB().Model(a).Where(query, args...).Count(&count).Error
		}
	})

	if err != nil {
		utils.Trail(utils.ERROR, "DB error in Count(%v). %s\n", model.GetModelName(a), err.Error())
	}
	return int(count)
}

// Update !
func Update(a interface{}, fieldName string, value interface{}, query string, args ...interface{}) (err error) {
	metrics.TimeMetric("uadmin/db/duration", 1000, func() {
		err = dialect.GetDB().Model(a).Where(query, args...).Update(fieldName, value).Error
		for fmt.Sprint(err) == "database is locked" {
			time.Sleep(time.Millisecond * 100)
			err = dialect.GetDB().Model(a).Where(query, args...).Update(fieldName, value).Error
		}
	})

	if err != nil {
		utils.Trail(utils.ERROR, "DB error in Update(%v). %s\n", model.GetModelName(a), err.Error())
	}
	return nil
}
