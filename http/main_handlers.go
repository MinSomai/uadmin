package http

import (
	"fmt"
	langmodel "github.com/uadmin/uadmin/blueprint/language/models"
	menumodel "github.com/uadmin/uadmin/blueprint/menu/models"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	"github.com/uadmin/uadmin/utils"
	"github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/translation"
	"net/http"
	"strings"
)

// listHandler !
func listHandler(w http.ResponseWriter, r *http.Request, session *sessionmodel.Session) {
	r.ParseMultipartForm(32 << 20)

	type Context struct {
		User           string
		Pagination     int
		Data           *preloaded.ListData
		Schema         model.ModelSchema
		IsUpdated      bool
		CanAdd         bool
		CanDelete      bool
		HasAccess      bool
		SiteName       string
		Language       langmodel.Language
		RootURL        string
		HasCategorical bool
		Searchable     bool
		CSRF           string
		Logo           string
		FavIcon        string
	}

	c := Context{}
	c.RootURL = preloaded.RootURL
	c.SiteName = preloaded.SiteName
	// @todo, redo
	// c.Language = translation.GetLanguage(r)
	c.User = session.User.Username
	c.CSRF = session.Key
	c.Logo = preloaded.Logo
	c.FavIcon = preloaded.FavIcon
	user := session.User

	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/")
	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	ModelName := r.URL.Path

	// Check permissions
	perm := user.GetAccess(ModelName)
	if !perm.Read {
		PageErrorHandler(w, r, session)
		return
	}
	c.HasAccess = perm.Read
	c.CanAdd = perm.Add
	c.CanDelete = perm.Delete

	// Initialize the schema
	m, ok := model.NewModel(ModelName, false)

	// Return 404 if it is an unknown model
	if !ok {
		PageErrorHandler(w, r, session)
		return
	}

	// Process delete
	if r.Method == preloaded.CPOST {
		if r.FormValue("delete") == "delete" {
			processDelete(ModelName, w, r, session, &user)
			c.IsUpdated = true
			http.Redirect(w, r, fmt.Sprint(preloaded.RootURL+r.URL.Path), http.StatusSeeOther)
		}
	}

	// Get the schema for the model
	c.Schema, _ = model.GetSchema(m.Interface())
	for i := range c.Schema.Fields {
		if c.Schema.Fields[i].CategoricalFilter {
			c.HasCategorical = true
		}
		// @todo, probably return
		//if c.Schema.Fields[i].Filter && c.Schema.Fields[i].Type == preloaded.CFK {
		//	c.Schema.Fields[i].Choices = utils.GetChoices(strings.ToLower(c.Schema.Fields[i].TypeName))
		//}
		if c.Schema.Fields[i].Searchable {
			c.Searchable = true
		}
	}

	// func (*ModelSchema, *User) (string, []interface{})
	// @todo, redo
	//query := ""
	//args := []interface{}{}
	//if c.Schema.ListModifier != nil {
	//	query, args = c.Schema.ListModifier(&c.Schema, &user)
	//}

	// @todo, probably return
	// c.Data = model.GetListData(m.Interface(), preloaded.PageLength, r, session, query, args...)
	c.Pagination = utils.PaginationHandler(c.Data.Count, preloaded.PageLength)

	RenderHTML(w, r, "./templates/uadmin/"+c.Schema.GetListTheme()+"/list.html", c)
}

// homeHandler !
func homeHandler(w http.ResponseWriter, r *http.Request, session *sessionmodel.Session) {
	type Context struct {
		User     string
		Demo     bool
		Menu     []menumodel.DashboardMenu
		SiteName string
		Language langmodel.Language
		RootURL  string
		Logo     string
		FavIcon  string
	}

	c := Context{}

	c.RootURL = preloaded.RootURL
	// @todo, redo
	// c.Language = translation.GetLanguage(r)
	c.SiteName = preloaded.SiteName
	c.User = session.User.Username
	c.Logo = preloaded.Logo
	c.FavIcon = preloaded.FavIcon

	c.Menu = session.User.GetDashboardMenu()
	for i := range c.Menu {
		c.Menu[i].MenuName = translation.Translate(c.Menu[i].MenuName, c.Language.Code, true)
	}

	RenderHTML(w, r, "./templates/uadmin/"+preloaded.Theme+"/home.html", c)
}

// GetModelsAPI returns a list of models
func GetModelsAPI(w http.ResponseWriter, r *http.Request, session *sessionmodel.Session) {
	response := []string{}
	for _, v := range model.ModelList {
		response = append(response, model.GetModelName(v))
	}
	utils.ReturnJSON(w, r, response)
}