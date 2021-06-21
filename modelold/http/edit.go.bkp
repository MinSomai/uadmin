package http

import (
	"encoding/json"
	"fmt"
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

func dAPIEditHandler(w http.ResponseWriter, r *http.Request, s *sessionmodel.Session) {
	urlParts := strings.Split(r.URL.Path, "/")
	modelName := urlParts[0]
	model, ok := model2.NewModel(modelName, false)
	if !ok {
		// @todo, redo
		// utils.Trail(utils.ERROR, "Couldnt return model for model name. %s", modelName)
		//utils.ReturnJSON(w, r, map[string]interface{}{
		//	"status":  "error",
		//	"err_msg": "Unknown model.",
		//})
		return
	}
	schema, _ := model2.GetSchema(model)
	tableName := schema.TableName

	// @todo, redo
	// Check CSRF
	//if utils.CheckCSRF(r) {
	//	utils.ReturnJSON(w, r, map[string]interface{}{
	//		"status":  "error",
	//		"err_msg": "Failed CSRF protection.",
	//	})
	//	return
	//}

	// Check permission
	allow := false
	if disableEditor, ok := model.Interface().(APIDisabledEditor); ok {
		allow = disableEditor.APIDisabledEdit(r)
		// This is a "Disable" method
		allow = !allow
		if !allow {
			// @todo, redo
			//utils.ReturnJSON(w, r, map[string]interface{}{
			//	"status":  "error",
			//	"err_msg": "Permission denied",
			//})
			return
		}
	}
	if publicEditor, ok := model.Interface().(APIPublicEditor); ok {
		allow = publicEditor.APIPublicEdit(r)
	}
	if !allow && s != nil {
		allow = s.User.GetAccess(modelName).Edit
	}
	if !allow {
		// @todo, redo
		//utils.ReturnJSON(w, r, map[string]interface{}{
		//	"status":  "error",
		//	"err_msg": "Permission denied",
		//})
		return
	}

	// Check if log is required
	log := preloaded.APILogEdit
	if logEditor, ok := model.Interface().(APILogEditor); ok {
		log = logEditor.APILogEdit(r)
	}

	// Get parameters
	params := getURLArgs(r)

	// remove empty file and image fields from params to avoid deleteing existing files
	// in case the request does not have a new file
	for k, v := range params {
		for _, f := range schema.Fields {
			if len(k) > 0 && k[0] == '_' && f.ColumnName == k[1:] && (f.Type == preloaded.CIMAGE || f.Type == preloaded.CFILE) {
				if v == "" && params[k+"-delete"] != "delete" {
					delete(params, k)
				}
				delete(params, k+"-delete")
				break
			}
		}
	}

	// Process Upload files
	fileList, err := dAPIUpload(w, r, &schema)
	if err != nil {
		// @todo, redo
		// utils.Trail(utils.ERROR, "dAPI Add Upload error processing. %s", err)
	}
	for k, v := range fileList {
		params["_"+k] = v
	}

	// Remove the field for file and image after adding it with an underscore
	// to params
	for k := range params {
		for _, f := range schema.Fields {
			if f.ColumnName == k && (f.Type == preloaded.CIMAGE || f.Type == preloaded.CFILE) {
				delete(params, k)
				break
			}
		}
	}

	writeMap, m2mMap := getEditMap(params, &schema, &model)

	db := dialect.GetDB("default")

	dialect1 := dialect.GetDialectForDb("default")
	sqlDialectStrings := dialect1.GetSqlDialectStrings()
	if len(urlParts) == 2 {
		// Edit multiple
		q, args := getFilters(r, params, tableName, &schema)

		modelArray, _ := model2.NewModelArray(modelName, true)
		if log {
			db.Model(model.Interface()).Where(q, args...).Scan(modelArray.Interface())
		}
		db = db.Model(model.Interface()).Where(q, args...).Updates(writeMap)
		if db.Error != nil {
			// @todo redo
			//utils.ReturnJSON(w, r, map[string]interface{}{
			//	"status":  "error",
			//	"err_msg": "Unable to update database. " + db.Error.Error(),
			//})
			return
		}

		// Process M2M
		db = dialect.GetDB("default").Begin()
		table1 := schema.ModelName
		for i := 0; i < modelArray.Elem().Len(); i++ {
			for k, v := range m2mMap {
				model, _ := model2.NewModel(k, false)
				t2Schema, _ := model2.GetSchema(model)
				table2 := t2Schema.ModelName
				// First delete exisiting records
				sql := sqlDialectStrings["deleteM2M"]
				sql = strings.Replace(sql, "{TABLE1}", table1, -1)
				sql = strings.Replace(sql, "{TABLE2}", table2, -1)
				sql = strings.Replace(sql, "{TABLE1_ID}", fmt.Sprint(database.GetID(modelArray.Elem().Index(i))), -1)
				db = db.Exec(sql)

				if v == "" {
					continue
				}

				// Now add the records
				for _, id := range strings.Split(v, ",") {
					sql = sqlDialectStrings["insertM2M"]
					sql = strings.Replace(sql, "{TABLE1}", table1, -1)
					sql = strings.Replace(sql, "{TABLE2}", table2, -1)
					sql = strings.Replace(sql, "{TABLE1_ID}", fmt.Sprint(database.GetID(modelArray.Elem().Index(i))), -1)
					sql = strings.Replace(sql, "{TABLE2_ID}", id, -1)
					db = db.Exec(sql)
				}
			}
		}
		db.Commit()

		returnDAPIJSON(w, r, map[string]interface{}{
			"status":     "ok",
			"rows_count": db.RowsAffected,
		}, params, "edit", model.Interface())
		if log {
			for i := 0; i < modelArray.Elem().Len(); i++ {
				createAPIEditLog(modelName, modelArray.Elem().Index(i).Interface(), &s.User, r)
			}
		}
	} else if len(urlParts) == 3 {
		// Edit One
		m, _ := model2.NewModel(modelName, true)
		if log {
			db.Model(model.Interface()).Where("id = ?", urlParts[2]).Scan(m.Interface())
		}
		db = db.Model(model.Interface()).Where("id = ?", urlParts[2]).Updates(writeMap)
		if db.Error != nil {
			// @todo, redo
			//utils.ReturnJSON(w, r, map[string]interface{}{
			//	"status":  "error",
			//	"err_msg": "Unable to update database. " + db.Error.Error(),
			//})
			return
		}

		// Process M2M
		db = dialect.GetDB("default").Begin()
		table1 := schema.ModelName
		for k, v := range m2mMap {
			model, _ := model2.NewModel(k, false)
			t2Schema, _ := model2.GetSchema(model)
			table2 := t2Schema.ModelName
			// First delete exisiting records
			sql := sqlDialectStrings["deleteM2M"]
			sql = strings.Replace(sql, "{TABLE1}", table1, -1)
			sql = strings.Replace(sql, "{TABLE2}", table2, -1)
			sql = strings.Replace(sql, "{TABLE1_ID}", urlParts[2], -1)
			db = db.Exec(sql)

			if v == "" {
				continue
			}

			// Now add the records
			for _, id := range strings.Split(v, ",") {
				sql = sqlDialectStrings["insertM2M"]
				sql = strings.Replace(sql, "{TABLE1}", table1, -1)
				sql = strings.Replace(sql, "{TABLE2}", table2, -1)
				sql = strings.Replace(sql, "{TABLE1_ID}", urlParts[2], -1)
				sql = strings.Replace(sql, "{TABLE2_ID}", id, -1)
				db = db.Exec(sql)
			}
		}
		db.Commit()

		if log {
			createAPIEditLog(modelName, m.Interface(), &s.User, r)
		}

		returnDAPIJSON(w, r, map[string]interface{}{
			"status":     "ok",
			"rows_count": db.RowsAffected,
		}, params, "edit", model.Interface())
	} else {
		// Error: Unknown format
		// @todo, redo
		//utils.ReturnJSON(w, r, map[string]interface{}{
		//	"status":  "error",
		//	"err_msg": "invalid format (" + r.URL.Path + ")",
		//})
		return
	}
}

func getEditMap(params map[string]string, schema *model2.ModelSchema, model *reflect.Value) (map[string]interface{}, map[string]string) {
	paramResult := map[string]interface{}{}
	m2mMap := map[string]string{}

	for k, v := range params {
		if k[0] != '_' {
			continue
		}
		k = k[1:]

		// Check M2M
		isM2M := false
		for _, f := range schema.Fields {
			if k == f.ColumnName && f.Type == preloaded.CM2M {
				m2mMap[strings.ToLower(f.TypeName)] = v
				isM2M = true
				break
			}
		}
		if isM2M {
			continue
		}

		var f *model2.F
		var isPtr = false
		for i := range schema.Fields {
			if k == schema.Fields[i].ColumnName || ((k) == schema.Fields[i].ColumnName+"_id" && schema.Fields[i].Type == preloaded.CFK) {
				f = &schema.Fields[i]
				isPtr = model.FieldByName(f.Name).Kind() == reflect.Ptr
				break
			}
		}
		if f == nil {
			continue
		}
		if v == "" && isPtr {
			paramResult[k] = nil
		} else {
			paramResult[k] = v
		}
	}

	return paramResult, m2mMap
}

func getWriteQueryFields(v string) string {
	if strings.HasPrefix(v, "_") {
		return dialect.GetDialectForDb("default").Quote(strings.TrimPrefix(v, "_"))
	}
	return ""
}

func createAPIEditLog(modelName string, m interface{}, user *usermodel.User, r *http.Request) {
	b, _ := json.Marshal(m)
	output := string(b[:len(b)-1]) + `,"_IP":"` + r.RemoteAddr + `"}`

	log := logmodel.Log{
		Username:  user.Username,
		Action:    logmodel.Action(0).Modified(),
		TableName: modelName,
		TableID:   int(database.GetID(reflect.ValueOf(m))),
		Activity:  output,
	}
	log.Save()
}

