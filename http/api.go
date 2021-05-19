package http

import (
	"encoding/json"
	// authapi "github.com/uadmin/uadmin/blueprint/auth/api"
	//userapi "github.com/uadmin/uadmin/blueprint/user/api"
	//model2 "github.com/uadmin/uadmin/model"
	//modelhttp "github.com/uadmin/uadmin/model/http"
	"github.com/uadmin/uadmin/preloaded"
	"net/http"
	"strings"
)

// apiHandler !
func apiHandler(w http.ResponseWriter, r *http.Request) {
	// @todo, probably return
	//session := authapi.IsAuthenticated(r)
	Path := strings.TrimPrefix(r.URL.Path, preloaded.RootURL+"api")
	// Handle requests for dAPI
	if strings.HasPrefix(Path, "/d/") || Path == "/d" {
		// @todo, probably return
		// modelhttp.DAPIHandler(w, r, session)
		return
	}

	// @todo, probably return
	// For all other APIs, if the user is not authenticated
	// then send them to login page
	//if session == nil {
	//	userapi.LoginHandler(w, r)
	//	return
	//}

	if strings.HasPrefix(Path, "/upload_image") {
		// @todo, probably return
		// UploadImageHandler(w, r, session)
		return
	}
	if strings.HasPrefix(Path, "/search") {
		// TODO: Move to separate file
		// @todo, probably return
		// modelName := r.FormValue("m")
		//model, ok := model2.NewModel(modelName, false)
		//if !ok {
		//	PageErrorHandler(w, r, session)
		//	return
		//}
		// s, _ := model2.GetSchema(model)

		//query := ""
		//args := []interface{}{}
		//if s.ListModifier != nil {
		//	query, args = s.ListModifier(&s, &session.User)
		//}

		// ld := model2.GetListData(model.Interface(), preloaded.PageLength, r, session, query, args...)

		type Context struct {
			List      [][]string `json:"list"`
			PageCount int        `json:"page_count"`
		}
		context := Context{
			List: [][]string{},
		}

		// @todo, probably return
		//for i := range ld.Rows {
		//	context.List = append(context.List, []string{})
		//	for j := range ld.Rows[i] {
		//		switch ld.Rows[i][j].(type) {
		//		case template.HTML:
		//			context.List[i] = append(context.List[i], fmt.Sprint(ld.Rows[i][j]))
		//		default:
		//			context.List[i] = append(context.List[i], html.EscapeString(fmt.Sprint(ld.Rows[i][j])))
		//		}
		//	}
		//}
		// context.PageCount = utils.PaginationHandler(ld.Count, preloaded.PageLength)

		bytes, _ := json.Marshal(context)
		w.Write(bytes)
		return
	}
	if strings.HasPrefix(Path, "/get_models") {
		// @todo, probably return
		// GetModelsAPI(w, r, session)
		return
	}
	if strings.HasPrefix(Path, "/get_fields") {
		// @todo, probably return
		// model2.GetFieldsAPI(w, r, session)
		return
	}
}

