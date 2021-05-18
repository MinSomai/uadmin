package http

import (
	"encoding/json"
	langmodel "github.com/uadmin/uadmin/blueprint/language/models"
	logmodel "github.com/uadmin/uadmin/blueprint/logging/models"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	"github.com/uadmin/uadmin/dialect"
	model2 "github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/translation"
	"github.com/uadmin/uadmin/utils"
	"net/http"
	"strings"
)

func dAPISchemaHandler(w http.ResponseWriter, r *http.Request, s *sessionmodel.Session) {
	urlParts := strings.Split(r.URL.Path, "/")
	model, _ := model2.NewModel(urlParts[0], false)
	modelName := dialect.GetDB().Config.NamingStrategy.TableName(model.Type().Name())
	params := getURLArgs(r)

	// Check permission
	allow := false
	if disableSchemer, ok := model.Interface().(APIDisabledSchemer); ok {
		allow = disableSchemer.APIDisabledSchema(r)
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
	if publicSchemer, ok := model.Interface().(APIPublicSchemer); ok {
		allow = publicSchemer.APIPublicSchema(r)
	}
	if !allow && s != nil {
		allow = s.User.GetAccess(modelName).Read
	}
	if !allow {
		utils.ReturnJSON(w, r, map[string]interface{}{
			"status":  "error",
			"err_msg": "Permission denied",
		})
		return
	}

	schema, _ := model2.GetSchema(model)

	// Get Language
	lang := r.URL.Query().Get("language")
	if lang == "" {
		if langC, err := r.Cookie("language"); err != nil || (langC != nil && langC.Value == "") {
			lang = langmodel.GetDefaultLanguage().Code
		} else {
			lang = langC.Value
		}
	}

	// Translation
	translation.TranslateSchema(&schema, lang)

	if r.URL.Query().Get("$choices") == "1" {
		// Load Choices for FK
		for i := range schema.Fields {
			if schema.Fields[i].Type == preloaded.CFK || schema.Fields[i].Type == preloaded.CM2M {
				choices := utils.GetChoices(schema.Fields[i].TypeName)
				schema.Fields[i].Choices = choices
			}
		}
	}

	returnDAPIJSON(w, r, map[string]interface{}{
		"status": "ok",
		"result": schema,
	}, params, "schema", model.Interface())

	go func() {
		// Check if log is required
		log := preloaded.APILogSchema
		if logSchemer, ok := model.Interface().(APILogSchemer); ok {
			log = logSchemer.APILogSchema(r)
		}

		if log {
			user := ""
			if s != nil {
				user = s.User.Username
			}
			activity, _ := json.Marshal(map[string]interface{}{
				"_IP": r.RemoteAddr,
			})
			log := logmodel.Log{
				Username:  user,
				Action:    logmodel.Action(0).GetSchema(),
				TableName: modelName,
				Activity:  string(activity),
			}
			log.Save()
		}
	}()
}
