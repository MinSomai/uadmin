package http

import (
	langmodel "github.com/uadmin/uadmin/blueprint/language/models"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	"github.com/uadmin/uadmin/preloaded"
	// "github.com/uadmin/uadmin/translation"
	"net/http"
	"strconv"
)

// pageErrorHandler is handler to return 404 pages
func PageErrorHandler(w http.ResponseWriter, r *http.Request, session *sessionmodel.Session) {
	type Context struct {
		User       string
		ID         uint
		UserExists bool
		Language   langmodel.Language
		SiteName   string
		ErrMsg     string
		ErrCode    int
		RootURL    string
		Logo       string
		FavIcon    string
	}

	c := Context{}

	c.RootURL = preloaded.RootURL
	c.SiteName = preloaded.SiteName
	// @todo, redo
	// c.Language = translation.GetLanguage(r)
	c.ErrMsg = "Page Not Found"
	c.ErrCode = 404
	c.Logo = preloaded.Logo
	c.FavIcon = preloaded.FavIcon
	if r.Form.Get("err_msg") != "" {
		c.ErrMsg = r.Form.Get("err_msg")
	}
	if code, err := strconv.ParseUint(r.Form.Get("err_code"), 10, 16); err == nil {
		c.ErrCode = int(code)
	}
	if session != nil {
		user := session.User
		c.User = user.Username
		c.ID = user.ID
	}

	w.WriteHeader(c.ErrCode)
	RenderHTML(w, r, "./templates/uadmin/"+preloaded.Theme+"/404.html", c)
}

