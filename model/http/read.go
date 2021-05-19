package http

import (
	"database/sql"
	"encoding/json"
	logmodel "github.com/uadmin/uadmin/blueprint/logging/models"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	usermodel "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/database"
	"github.com/uadmin/uadmin/dialect"
	model2 "github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/preloaded"
	"net/http"
	"reflect"
	"strings"
)

func dAPIReadHandler(w http.ResponseWriter, r *http.Request, s *sessionmodel.Session) {
	var err error
	var rowsCount int64

	urlParts := strings.Split(r.URL.Path, "/")
	modelName := urlParts[0]
	model, ok := model2.NewModel(modelName, false)
	if !ok {
		w.WriteHeader(401)
		// @todo, redo
		//utils.ReturnJSON(w, r, map[string]interface{}{
		//	"status":  "error",
		//	"err_msg": "No model found",
		//})
		return
	}
	params := getURLArgs(r)
	schema, _ := model2.GetSchema(model.Interface())

	// Check permission
	allow := false
	if disableReader, ok := model.Interface().(APIDisabledReader); ok {
		allow = disableReader.APIDisabledRead(r)
		// This is a "Disable" method
		allow = !allow
		if !allow {
			w.WriteHeader(401)
			// @todo, redo
			//utils.ReturnJSON(w, r, map[string]interface{}{
			//	"status":  "error",
			//	"err_msg": "Permission denied",
			//})
			return
		}
	}
	if publicReader, ok := model.Interface().(APIPublicReader); ok {
		allow = publicReader.APIPublicRead(r)
	}
	if !allow && s != nil {
		allow = s.User.GetAccess(modelName).Read
	}
	if !allow {
		w.WriteHeader(401)
		// @todo, redo
		//utils.ReturnJSON(w, r, map[string]interface{}{
		//	"status":  "error",
		//	"err_msg": "Permission denied",
		//})
		return
	}

	// Check if log is required
	log := preloaded.APILogRead
	if logReader, ok := model.Interface().(APILogReader); ok {
		log = logReader.APILogRead(r)
	}

	// Run prequery handler
	if preQueryReader, ok := model.Interface().(APIPreQueryReader); ok {
		if !preQueryReader.APIPreQueryRead(w, r) {
			return
		}
	}

	if len(urlParts) == 2 {
		// Read Multiple
		var m interface{}

		SQL := "SELECT {FIELDS} FROM {TABLE_NAME}"
		if val, ok := params["$distinct"]; ok && val == "1" {
			SQL = "SELECT DISTINCT {FIELDS} FROM {TABLE_NAME}"
		}

		tableName := schema.TableName
		SQL = strings.Replace(SQL, "{TABLE_NAME}", tableName, -1)

		f, customSchema := getQueryFields(r, params, tableName)
		if f != "" {
			SQL = strings.Replace(SQL, "{FIELDS}", f, -1)
		} else {
			SQL = strings.Replace(SQL, "{FIELDS}", "*", -1)
		}

		join := getQueryJoin(r, params, tableName)
		if join != "" {
			SQL += " " + join
		}

		q, args := getFilters(r, params, tableName, &schema)
		if q != "" {
			SQL += " WHERE " + q
		}

		groupBy := getQueryGroupBy(r, params)
		if groupBy != "" {
			SQL += " GROUP BY " + groupBy
		}
		order := getQueryOrder(r, params)
		if order != "" {
			SQL += " ORDER BY " + order
		}
		limit := getQueryLimit(r, params)
		if limit != "" {
			SQL += " LIMIT " + limit
		}
		offset := getQueryOffset(r, params)
		if offset != "" {
			SQL += " OFFSET " + offset
		}

		if preloaded.DebugDB {
			// @todo, redo
			//utils.Trail(utils.DEBUG, SQL)
			//utils.Trail(utils.DEBUG, "%#v", args)
		}
		var rows *sql.Rows

		if !customSchema {
			mArray, _ := model2.NewModelArray(modelName, true)
			m = mArray.Interface()
		} else {
			m = []map[string]interface{}{}
		}
		dialect1 := dialect.GetDialectForDb()
		db := dialect.GetDB().Begin()
		rows, err = dialect1.ReadRows(db, customSchema, SQL, m, args...)
		if customSchema {
			if err != nil {
				w.WriteHeader(500)
				// @todo, redo
				//utils.ReturnJSON(w, r, map[string]interface{}{
				//	"status":  "error",
				//	"err_msg": "Unable to execute SQL. " + err.Error(),
				//})
				//utils.Trail(utils.ERROR, "SQL: %v\nARGS: %v", SQL, args)
				return
			}
			m = parseCustomDBSchema(rows)
		}
		db.Commit()
		if a, ok := m.([]map[string]interface{}); ok {
			rowsCount = int64(len(a))
		} else {
			rowsCount = int64(reflect.ValueOf(m).Elem().Len())
		}
		// Preload
		if params["$preload"] == "1" {
			mList := reflect.ValueOf(m)
			for i := 0; i < mList.Elem().Len(); i++ {
				database.Preload(mList.Elem().Index(i).Addr().Interface())
			}
		}

		// Process M2M
		getQueryM2M(params, m, customSchema, modelName)

		returnDAPIJSON(w, r, map[string]interface{}{
			"status": "ok",
			"result": m,
		}, params, "read", model.Interface())
		go func() {
			if log {
				createAPIReadLog(modelName, 0, rowsCount, params, &s.User, r)
			}
		}()
		return
	} else if len(urlParts) == 3 {
		// Read One
		m, _ := model2.NewModel(modelName, true)
		database.Get(m.Interface(), "id = ?", urlParts[2])
		rowsCount = 0

		var i interface{}
		if int(database.GetID(m)) != 0 {
			i = m.Interface()
			rowsCount = 1
		}

		if params["$preload"] == "1" {
			database.Preload(m.Interface())
		}

		returnDAPIJSON(w, r, map[string]interface{}{
			"status": "ok",
			"result": i,
		}, params, "read", model.Interface())
		go func() {
			if log {
				createAPIReadLog(modelName, int(database.GetID(m)), rowsCount, map[string]string{"id": urlParts[2]}, &s.User, r)
			}
		}()
	} else {
		// Error: Unknown format
		w.WriteHeader(404)
		// @todo, redo
		//utils.ReturnJSON(w, r, map[string]interface{}{
		//	"status":  "error",
		//	"err_msg": "invalid format (" + r.URL.Path + ")",
		//})
		return
	}
}

func createAPIReadLog(modelName string, ID int, rowsCount int64, params map[string]string, user *usermodel.User, r *http.Request) {
	vals := map[string]interface{}{
		"params":     params,
		"rows_count": rowsCount,
		"_IP":        r.RemoteAddr,
	}
	output, _ := json.Marshal(vals)

	log := logmodel.Log{
		Username:  user.Username,
		Action:    logmodel.Action(0).Read(),
		TableName: modelName,
		TableID:   ID,
		Activity:  string(output),
	}
	log.Save()

}

