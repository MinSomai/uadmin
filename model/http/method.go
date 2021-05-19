package http

import (
	"fmt"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	"github.com/uadmin/uadmin/database"
	model2 "github.com/uadmin/uadmin/model"
	"net/http"
	"reflect"
	"strings"
)

func dAPIMethodHandler(w http.ResponseWriter, r *http.Request, s *sessionmodel.Session) {
	urlParts := strings.Split(r.URL.Path, "/")
	modelName := urlParts[0]
	model, _ := model2.NewModel(modelName, true)

	params := getURLArgs(r)

	if len(urlParts) < 4 {
		w.WriteHeader(400)
		// @todo, redo
		// utils.ReturnJSON(w, r, map[string]interface{}{
		//	"status":  "error",
		//	"err_msg": "Bad request, URL format should be api/d/model/method/{METHOD_NAME}/{ID}",
		//})
		return
	}

	// @todo, redo
	//if utils.CheckCSRF(r) {
		//utils.ReturnJSON(w, r, map[string]interface{}{
		//	"status":  "error",
		//	"err_msg": "Failed CSRF protection.",
		//})
	//	return
	//}

	f := model.MethodByName(urlParts[2])
	if !f.IsValid() {
		f = model.Elem().MethodByName(urlParts[2])
	}

	if !f.IsValid() {
		w.WriteHeader(404)
		// @todo, redo
		//utils.ReturnJSON(w, r, map[string]interface{}{
		//	"status":  "error",
		//	"err_msg": "Method (" + urlParts[2] + ") doesn't exist.",
		//})
		return
	}

	database.Get(model.Interface(), "id = ?", urlParts[3])
	if database.GetID(model) == 0 {
		w.WriteHeader(404)
		// @todo, redo
		//utils.ReturnJSON(w, r, map[string]interface{}{
		//	"status":  "error",
		//	"err_msg": "ID doesn't exist (" + urlParts[3] + ").",
		//})
		return
	}

	ret := model.MethodByName(urlParts[2]).Call([]reflect.Value{})

	// Return if the method has a return value
	if len(ret) != 0 {
		returnDAPIJSON(w, r, map[string]interface{}{
			"status": "ok",
			"value":  fmt.Sprint(ret[0]),
		}, params, "method", model.Interface())
	}
}

