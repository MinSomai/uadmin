package http

import (
	"encoding/json"
	logmodel "github.com/uadmin/uadmin/blueprint/logging/models"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	usermodel "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/database"
	dialect2 "github.com/uadmin/uadmin/dialect"
	model2 "github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/utils"
	"net/http"
	"reflect"
	"strings"
)

func dAPIDeleteHandler(w http.ResponseWriter, r *http.Request, s *sessionmodel.Session) {
	var rowsCount int64
	urlParts := strings.Split(r.URL.Path, "/")
	modelName := urlParts[0]
	model, ok := model2.NewModel(modelName, false)
	if !ok {
		utils.Trail(utils.ERROR, "Couldnt return model for model name. %s", modelName)
		utils.ReturnJSON(w, r, map[string]interface{}{
			"status":  "error",
			"err_msg": "Unknown model.",
		})
		return
	}

	schema, _ := model2.GetSchema(model)
	tableName := schema.TableName
	params := getURLArgs(r)

	// Check CSRF
	if utils.CheckCSRF(r) {
		utils.ReturnJSON(w, r, map[string]interface{}{
			"status":  "error",
			"err_msg": "Failed CSRF protection.",
		})
		return
	}
	dialect := dialect2.GetDialectForDb()
	// Check permission
	allow := false
	if disableDeleter, ok := model.Interface().(APIDisabledDeleter); ok {
		allow = disableDeleter.APIDisabledDelete(r)
		// This is a "Disable" method
		allow = !allow
		if !allow {
			utils.ReturnJSON(w, r, map[string]interface{}{
				"status":  "error",
				"err_msg": "Permission denied",
			})
			return
		}
	}
	if publicDeleter, ok := model.Interface().(APIPublicDeleter); ok {
		allow = publicDeleter.APIPublicDelete(r)
	}
	if !allow && s != nil {
		allow = s.User.GetAccess(modelName).Delete
	}
	if !allow {
		utils.ReturnJSON(w, r, map[string]interface{}{
			"status":  "error",
			"err_msg": "Permission denied",
		})
		return
	}

	// Check if log is required
	log := preloaded.APILogDelete
	if logDeleter, ok := model.Interface().(APILogDeleter); ok {
		log = logDeleter.APILogDelete(r)
	}

	if len(urlParts) == 2 {
		// Delete Multiple
		q, args := getFilters(r, params, tableName, &schema)

		modelArray, _ := model2.NewModelArray(modelName, true)

		// Block Delete All
		if q == "deleted_at IS NULL" {
			utils.ReturnJSON(w, r, map[string]interface{}{
				"status":  "error",
				"err_msg": "Delete all is blocked",
			})
			return
		}

		db := dialect2.GetDB().Begin()
		if log {
			db.Model(model.Interface()).Where(q, args...).Scan(modelArray.Interface())
		}
		dialect.Delete(db, model, q, args...)
		db.Commit()
		if db.Error != nil {
			utils.ReturnJSON(w, r, map[string]interface{}{
				"status":  "error",
				"err_msg": "Unable to COMMIT SQL. " + db.Error.Error(),
			})
			return
		}
		rowsCount = db.RowsAffected
		if log {
			for i := 0; i < modelArray.Elem().Len(); i++ {
				createAPIDeleteLog(modelName, modelArray.Elem().Index(i).Interface(), &s.User, r)
			}
		}
		returnDAPIJSON(w, r, map[string]interface{}{
			"status":     "ok",
			"rows_count": rowsCount,
		}, params, "delete", model.Interface())
	} else if len(urlParts) == 3 {
		// Delete One
		m, _ := model2.NewModel(modelName, true)

		db := dialect2.GetDB()
		if log {
			db.Model(model.Interface()).Where("id = ?", urlParts[2]).Scan(m.Interface())
		}
		db = db.Where("id = ?", urlParts[2]).Delete(model.Interface())
		if db.Error != nil {
			utils.ReturnJSON(w, r, map[string]interface{}{
				"status":  "error",
				"err_msg": "Unable to execute DELETE SQL. " + db.Error.Error(),
			})
			return
		}

		if log {
			createAPIDeleteLog(modelName, m.Interface(), &s.User, r)
		}

		returnDAPIJSON(w, r, map[string]interface{}{
			"status":     "ok",
			"rows_count": db.RowsAffected,
		}, params, "delete", model.Interface())
	} else {
		// Error: Unknown format
		utils.ReturnJSON(w, r, map[string]interface{}{
			"status":  "error",
			"err_msg": "invalid format (" + r.URL.Path + ")",
		})
		return
	}
}

func createAPIDeleteLog(modelName string, m interface{}, user *usermodel.User, r *http.Request) {
	b, _ := json.Marshal(m)
	output := string(b[:len(b)-1]) + `,"_IP":"` + r.RemoteAddr + `"}`

	log := logmodel.Log{
		Username:  user.Username,
		Action:    logmodel.Action(0).Deleted(),
		TableName: modelName,
		TableID:   int(database.GetID(reflect.ValueOf(m))),
		Activity:  output,
	}
	log.Save()
}

