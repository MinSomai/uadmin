package http

import (
	"encoding/json"
	"fmt"
	logmodel "github.com/uadmin/uadmin/blueprint/logging/models"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	usermodel "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/database"
	model2 "github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/utils"
	"net/http"
	"strconv"
	"strings"
)

// processDelete is a handler for processing deleting records from a table
func processDelete(a interface{}, w http.ResponseWriter, r *http.Request, session *sessionmodel.Session, user *usermodel.User) {

	if r.FormValue("listID") == "" || r.FormValue("listID") == "," {
		return
	}
	tempID := strings.Split(r.FormValue("listID"), ",")
	var tempIDs []uint
	modelName, ok := a.(string)
	if !ok {
		PageErrorHandler(w, r, session)
		return
	}

	if !user.GetAccess(modelName).Delete {
		return
	}

	// Check CSRF
	if utils.CheckCSRF(r) {
		PageErrorHandler(w, r, session)
		return
	}

	for _, v := range tempID {
		temp, _ := strconv.ParseUint(v, 10, 64)
		tempIDs = append(tempIDs, uint(temp))
	}

	if preloaded.LogDelete {
		for _, v := range tempIDs {
			log := logmodel.Log{}
			log.Username = user.Username
			log.Action = log.Action.Deleted()
			log.TableName = modelName
			log.TableID = int(v)

			m, ok := model2.NewModel(modelName, false)
			if !ok {
				utils.Trail(utils.ERROR, "processDelete invalid model name: %s", modelName)
			}
			database.Get(m.Addr().Interface(), "id = ?", v)
			model, _ := model2.NewModel(modelName, false)
			s, _ := model2.GetSchema(model)
			model2.GetFormData(m.Interface(), r, session, &s, user)
			jsonifyValue := map[string]string{}
			for _, ff := range s.Fields {
				jsonifyValue[ff.Name] = fmt.Sprint(ff.Value)
			}

			json1, _ := json.Marshal(jsonifyValue)
			log.Activity = string(json1)

			log.Save()
		}
	}

	m, ok := model2.NewModel(modelName, true)
	if !ok {
		PageErrorHandler(w, r, session)
		return
	}

	type Deleter interface {
		Delete(interface{}, string, ...interface{})
	}

	deleter, ok := m.Interface().(Deleter)
	if ok {
		deleter.Delete(m.Interface(), "id IN (?)", tempIDs)
	} else {
		database.DeleteList(m.Interface(), "id IN (?)", tempIDs)
	}
}

