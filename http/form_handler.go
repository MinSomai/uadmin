package http

//// formHandler handles form view requests to render forms and process POST requests to edit
//// the form content. It also handles delete requests for inlines of the form.
//func formHandler(w http.ResponseWriter, r *http.Request, session *sessionmodel.Session) {
//	r.ParseMultipartForm(32 << 20)
//	type Context struct {
//		User            string
//		ID              uint
//		Schema          model.ModelSchema
//		SaveAndContinue bool
//		IsUpdated       bool
//		//Demo            bool
//		CanUpdate bool
//		SiteName  string
//		Language  langmodel.Language
//		Direction string
//		RootURL   string
//		ReadOnlyF string
//		CSRF      string
//		Logo      string
//		FavIcon   string
//	}
//	var err error
//	c := Context{}
//
//	c.RootURL = preloaded.RootURL
//	c.Language = translation.GetLanguage(r)
//	c.User = session.User.Username
//	c.SiteName = preloaded.SiteName
//	c.CSRF = authapi.GetSession(r)
//	c.Logo = preloaded.Logo
//	c.FavIcon = preloaded.FavIcon
//	user := session.User
//
//	URLPath := strings.Split(strings.TrimPrefix(r.URL.Path, preloaded.RootURL), "/")
//
//	ModelName := URLPath[0]
//	ModelID, _ := strconv.ParseUint(URLPath[1], 10, 64)
//	ID := uint(ModelID)
//	_ = ID
//
//	m, ok := model.NewModel(ModelName, false)
//	if !ok {
//		PageErrorHandler(w, r, session)
//		return
//	}
//
//	// Check user permissions
//	perm := user.GetAccess(ModelName)
//	if !perm.Read {
//		PageErrorHandler(w, r, session)
//		return
//	}
//	c.CanUpdate = perm.Add || perm.Edit
//
//	c.Schema, _ = model.GetSchema(m.Interface())
//
//	// Filter inlines that the user does not have permission to
//	inlinesList := []*model.ModelSchema{}
//	for i := range c.Schema.Inlines {
//		if user.GetAccess(c.Schema.Inlines[i].ModelName).Read {
//			inlinesList = append(inlinesList, c.Schema.Inlines[i])
//		}
//	}
//	c.Schema.Inlines = inlinesList
//
//	r.Form.Set("ModelID", fmt.Sprint(ModelID))
//	InlineModelName := ""
//	if r.FormValue("listModelName") != "" {
//		InlineModelName = strings.ToLower(r.FormValue("listModelName"))
//	}
//
//	if r.Method == preloaded.CPOST {
//		// Check CSRF
//		if utils.CheckCSRF(r) {
//			PageErrorHandler(w, r, session)
//			return
//		}
//		if r.FormValue("delete") == "delete" {
//			if InlineModelName != "" {
//				processDelete(InlineModelName, w, r, session, &user)
//			}
//			c.IsUpdated = true
//			http.Redirect(w, r, fmt.Sprint(preloaded.RootURL+r.URL.Path), http.StatusSeeOther)
//		} else {
//			// Process the form and check for validation errors
//			// @todo, redo
//			//m = utils.ProcessForm(ModelName, w, r, session, &c.Schema)
//			//m = m.Elem()
//			if r.FormValue("new_url") != "" {
//				r.URL, err = url.Parse(r.FormValue("new_url"))
//				if err != nil {
//					utils.Trail(utils.ERROR, "formHandler unable to parse new_url(%s). %s", r.FormValue("new_url"), err)
//					return
//				}
//			}
//		}
//	}
//
//	if r.FormValue("new_url") == "" {
//		if preloaded.OptimizeSQLQuery {
//			database.GetForm(m.Addr().Interface(), &c.Schema, "id = ?", ModelID)
//		} else {
//			database.Get(m.Addr().Interface(), "id = ?", ModelID)
//		}
//	}
//
//	// Return 404 incase the ID doens't exist in the DB and its not in new form
//	if URLPath[1] != "new" {
//		if database.GetID(m) == 0 {
//			PageErrorHandler(w, r, session)
//			return
//		}
//	}
//
//	// Check if Save and Continue
//	inlines := model.GetInlines()
//	c.SaveAndContinue = (URLPath[1] == "new" && len(inlines[ModelName]) > 0 && r.URL.Query().Get("return_url") == "")
//
//	// Disable fk for inline form
//	if r.URL.Query().Get("return_url") != "" {
//		for k := range r.URL.Query() {
//			if c.Schema.FieldByName(k).Type == preloaded.CFK {
//				c.ReadOnlyF = c.Schema.FieldByName(k).Name
//			}
//		}
//	}
//
//	// @todo, maybe return
//	// Process User Custom Schema Logic
//	//if c.Schema.FormModifier != nil {
//	//	c.Schema.FormModifier(&c.Schema, m.Addr().Interface(), &user)
//	//}
//
//	// Add data to Schema
//	// @todo probably, return
//	// model.GetFormData(m.Interface(), r, session, &c.Schema, &user)
//	translation.TranslateSchema(&c.Schema, c.Language.Code)
//
//	RenderHTML(w, r, "./templates/uadmin/"+c.Schema.GetFormTheme()+"/form.html", c)
//
//	// Store Read Log in a separate go routine
//	if preloaded.LogRead {
//		go func() {
//			if ModelID > 0 {
//				log := &logmodel.Log{}
//				log.ParseRecord(m, m.Type().Name(), uint(ModelID), &session.User, log.Action.Read(), r)
//				log.Save()
//			}
//		}()
//	}
//}
//

