package http

import (
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	"github.com/uadmin/uadmin/model"
	"net/http"
)

func dAPIAllModelsHandler(w http.ResponseWriter, r *http.Request, session *sessionmodel.Session) {
	response := []interface{}{}
	for _, v := range model.ModelList {
		response = append(response, model.Schema[model.GetModelName(v)])
	}
	// @todo, redo
	//utils.ReturnJSON(w, r, map[string]interface{}{
	//	"status": "ok",
	//	"result": response,
	//})
}

