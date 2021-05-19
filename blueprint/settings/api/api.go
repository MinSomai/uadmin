package api

import (
	langmodel "github.com/uadmin/uadmin/blueprint/language/models"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	imageapi "github.com/uadmin/uadmin/blueprint/image/api"
	settingmodel "github.com/uadmin/uadmin/blueprint/settings/models"
	"github.com/uadmin/uadmin/database"
	"github.com/uadmin/uadmin/dialect"
	"github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/preloaded"
	// "github.com/uadmin/uadmin/translation"
	"net/http"
	"strings"
)

func SettingsHandler(w http.ResponseWriter, r *http.Request, session *sessionmodel.Session) {
	r.ParseMultipartForm(32 << 20)
	type SCat struct {
		ID       uint
		Name     string
		Icon     string
		Settings []settingmodel.Setting
	}
	type Context struct {
		User     string
		SiteName string
		Language langmodel.Language
		RootURL  string
		SCat     []SCat
		Logo     string
		FavIcon  string
	}
	tMap := map[settingmodel.DataType]string{
		settingmodel.DataType(0).File():  preloaded.CFILE,
		settingmodel.DataType(0).Image(): preloaded.CIMAGE,
	}

	if session == nil {
		// @todo, redo
		// uadminhttp.PageErrorHandler(w, r, session)
		return
	}

	// Check if the user has permission to settings models
	perm := session.User.GetAccess("setting")
	if !perm.Read {
		// @todo, redo
		// uadminhttp.PageErrorHandler(w, r, session)
		return
	}

	settings := []settingmodel.Setting{}
	database.All(&settings)
	if r.Method == preloaded.CPOST {
		if !perm.Edit {
			// @todo, redo
			// uadminhttp.PageErrorHandler(w, r, session)
			return
		}
		//var tempSet Setting
		tx := dialect.GetDB().Begin()

		for _, s := range settings {
			v, ok := r.Form[s.Code]
			if s.DataType == s.DataType.Image() || s.DataType == s.DataType.File() {
				sParts := strings.SplitN(s.Code, ".", 2)

				// Process Files and Images
				_, _, err := r.FormFile(s.Code)
				if err != nil {
					continue
				}

				schema, _ := model.GetSchema(s)
				schema.FieldByName(sParts[1])

				f := model.F{Name: s.Code, Type: tMap[s.DataType], UploadTo: "/static/settings/"}

				val := imageapi.ProcessUpload(r, &f, "setting", session, &schema)
				if val == "" {
					continue
				}
				s.ParseFormValue([]string{val})
			} else if s.DataType == s.DataType.Boolean() {
				if ok {
					s.Value = "1"
				} else {
					s.Value = "0"
				}
			} else {
				s.ParseFormValue(v)
			}
			s.ApplyValue()
			tx.Save(&s)
		}
		tx.Commit()
	}

	c := Context{}

	c.RootURL = preloaded.RootURL
	// @todo, redo
	// c.Language = translation.GetLanguage(r)
	c.SiteName = preloaded.SiteName
	c.User = session.User.Username
	c.SCat = []SCat{}
	c.Logo = preloaded.Logo
	c.FavIcon = preloaded.FavIcon

	catList := []settingmodel.SettingCategory{}
	database.All(&catList)

	for _, cat := range catList {
		c.SCat = append(c.SCat, SCat{
			ID:       cat.ID,
			Name:     cat.Name,
			Icon:     cat.Icon,
			Settings: []settingmodel.Setting{},
		})
		database.Filter(&c.SCat[len(c.SCat)-1].Settings, "category_id = ?", cat.ID)
	}

	// @todo, redo
	// uadminhttp.RenderHTML(w, r, "./templates/uadmin/"+preloaded.Theme+"/setting.html", c)
}
