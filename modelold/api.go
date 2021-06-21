package modelold

// @todo redo
// GetListSchema returns a schema for list view
//func GetListData(a interface{}, PageLength int, r *http.Request, session *sessionmodel.Session, query string, args ...interface{}) (l *preloaded.ListData) {
	//l = &preloaded.ListData{}
	//schema, _ := GetSchema(a)
	//language := translation.GetLanguage(r)
	//translation.TranslateSchema(&schema, language.Code)
	//
	//t := reflect.TypeOf(a)
	//
	//// For Order By and Pagination
	//o := r.FormValue("o")
	//asc := !strings.HasPrefix(o, "-")
	//if !asc {
	//	o = strings.Replace(o, "-", "", 1)
	//}
	//p := r.FormValue("p")
	//page, _ := strconv.ParseInt(p, 10, 32)
	//
	//m, ok := NewModelArray(schema.ModelName, false)
	//if !ok {
	//	utils.Trail(utils.ERROR, "getListSchema.NewModelArray. No model name", schema.ModelName)
	//	return
	//}
	//
	//newModel, _ := NewModel(schema.ModelName, false)
	//iPager, isPager := newModel.Interface().(utils.AdminPager)
	//iCounter, isCounter := newModel.Interface().(utils.Counter)
	//
	//var (
	//	_query interface{}
	//	_args  []interface{}
	//)
	//
	//_query, _args = uadminhttp.GetFilter(r, session, &schema)
	//if _query.(string) != "" {
	//	if query == "" {
	//		query = _query.(string)
	//	} else {
	//		query += " AND " + _query.(string)
	//	}
	//	args = append(args, _args...)
	//}
	//if !isPager {
	//	if preloaded.OptimizeSQLQuery {
	//		database.FilterList(&schema, o, asc, int(page-1)*PageLength, PageLength, m.Addr().Interface(), query, args...)
	//	} else {
	//		database.AdminPage(o, asc, int(page-1)*PageLength, PageLength, m.Addr().Interface(), query, args...)
	//	}
	//} else {
	//	iPager.AdminPage(o, asc, int(page-1)*PageLength, PageLength, m.Addr().Interface(), query, args...)
	//}
	//if !isCounter {
	//	l.Count = database.Count(m.Interface(), query, args...)
	//} else {
	//	l.Count = iCounter.Count(m.Interface(), query, args...)
	//}
	//for i := 0; i < m.Len(); i++ {
	//	l.Rows = append(l.Rows, evaluateObject(m.Index(i).Interface(), t, &schema, language.Code, session))
	//}
	//return
//}

// @todo redo
// evaluateObject !
//func evaluateObject(obj interface{}, t reflect.Type, s *ModelSchema, lang string, session *sessionmodel.Session) (y []interface{}) {
//	value := reflect.ValueOf(obj)
//	for index := 0; index < len(s.Fields); index++ {
//		if s.Fields[index].IsMethod {
//			if strings.Contains(s.Fields[index].Name, "__List") {
//				in := []reflect.Value{}
//				method := value.MethodByName(s.Fields[index].Name)
//				ret := method.Call(in)
//				y = append(y, template.HTML(utils.StripHTMLScriptTag(fmt.Sprint(ret[0].Interface()))))
//			}
//			continue
//		}
//
//		field, _ := t.FieldByName(s.Fields[index].Name)
//		if strings.ToLower(string(field.Name[0])) == string(field.Name[0]) {
//			continue
//		}
//		if !s.Fields[index].ListDisplay {
//			continue
//		}
//
//		v := value.FieldByName(field.Name)
//		if s.Fields[index].Type == preloaded.CID {
//			id, ok := v.Interface().(uint)
//			if !ok {
//				utils.Trail(utils.ERROR, "evaluateObject.Interface.(uadmin.Model) ID NOT OK. %#v", v.Interface())
//			}
//			temp := template.HTML(fmt.Sprintf("<a class='clickable Row_id no-style bold' data-id='%d' href='%s%s/%d'>%s</a>", id, preloaded.RootURL, s.ModelName, id, html.EscapeString(utils.GetString(obj))))
//			y = append(y, temp)
//		} else if s.Fields[index].Type == preloaded.CNUMBER {
//			temp := v.Interface()
//			y = append(y, temp)
//
//		} else if s.Fields[index].Type == preloaded.CPROGRESSBAR {
//			tempProgressValue, _ := strconv.ParseFloat(fmt.Sprint(v.Interface()), 64) // 10
//			tempProgressColor := ""
//			maxThreshold := 0.0
//			var maxColor string
//			if len(s.Fields[index].ProgressBar) == 1 {
//				for tempThreshold, tempColor := range s.Fields[index].ProgressBar {
//					tempProgressColor = tempColor
//					maxThreshold = tempThreshold
//				}
//			} else {
//				for tempThreshold, tempColor := range s.Fields[index].ProgressBar {
//					if tempThreshold > maxThreshold {
//						maxThreshold = tempThreshold
//						maxColor = tempColor
//					}
//				}
//				currentThreshold := maxThreshold
//				tempProgressColor = maxColor
//				for tempThreshold, tempColor := range s.Fields[index].ProgressBar {
//					if tempThreshold > tempProgressValue && tempThreshold < currentThreshold {
//						tempProgressColor = tempColor
//						currentThreshold = tempThreshold
//					}
//				}
//			}
//			tempProgressWidth := tempProgressValue / maxThreshold * 100.0
//			if tempProgressValue > maxThreshold {
//				tempProgressWidth = 100.0
//				tempProgressColor = maxColor
//			}
//
//			tempColor := helper.GetRGB(tempProgressColor)
//			DarkerFactor := 0.67
//			tempDarker := []int{
//				int(float64(tempColor[0]) * DarkerFactor),
//				int(float64(tempColor[1]) * DarkerFactor),
//				int(float64(tempColor[2]) * DarkerFactor),
//			}
//
//			tempGradient1 := fmt.Sprintf("#%02x%02x%02x", tempColor[0], tempColor[1], tempColor[2])
//			tempGradient2 := fmt.Sprintf("#%02x%02x%02x", tempDarker[0], tempDarker[1], tempDarker[2])
//
//			temp := template.HTML(fmt.Sprintf("<div style='border:solid 1px;'><div style='width:%d%%; background-image:linear-gradient(%s,%s); text-align:center;'>%.2f</div></div>", int(tempProgressWidth), tempGradient1, tempGradient2, tempProgressValue))
//			y = append(y, temp)
//		} else if s.Fields[index].Type == preloaded.CMONEY {
//			temp := utils.Commaf(v.Interface())
//			y = append(y, temp)
//		} else if s.Fields[index].Type == preloaded.CLINK {
//			if fmt.Sprint(v) != "" {
//				URL := v.String()
//				if URL != "" {
//					if !strings.Contains(URL, "?") {
//						URL += "?"
//					} else {
//						URL += "&"
//					}
//					URL += "x-csrf-token=" + session.Key
//				}
//				temp := template.HTML(fmt.Sprintf("<a class='btn btn-primary' href='%s'>%s</a>", URL, s.Fields[index].Name))
//				y = append(y, temp)
//			} else {
//				temp := template.HTML("<span></span>")
//				y = append(y, temp)
//			}
//		} else if s.Fields[index].Type == preloaded.CDATE {
//			if fmt.Sprint(v.Type())[0] == '*' {
//				if v.IsNil() {
//					y = append(y, "")
//				} else {
//					v = v.Elem()
//					d, _ := v.Interface().(time.Time)
//					y = append(y, d.Format("2006-01-02 15:04:05"))
//				}
//			} else {
//				d, _ := v.Interface().(time.Time)
//				y = append(y, d.Format("2006-01-02 15:04:05"))
//			}
//
//		} else if s.Fields[index].Type == preloaded.CFK {
//			vID := value.FieldByName(field.Name + "ID")
//			cIndex, _ := vID.Interface().(uint)
//			fkFieldName := strings.ToLower(value.FieldByName(s.Fields[index].Name).Type().Name())
//			if value.FieldByName(s.Fields[index].Name).Type().Kind() == reflect.Ptr {
//				fkFieldName = strings.ToLower(value.FieldByName(s.Fields[index].Name).Type().Elem().Name())
//			}
//
//			// Fetch that record from DB
//			fkModel, _ := NewModel(fkFieldName, true)
//			database.GetStringer(fkModel.Interface(), "id = ?", cIndex)
//			temp := template.HTML(fmt.Sprintf("<a class='clickable no-style bold' href='%s%s/%d'>%s</a>", preloaded.RootURL, fkFieldName, cIndex, html.EscapeString(utils.GetString(fkModel.Interface()))))
//			y = append(y, temp)
//		} else if s.Fields[index].Type == preloaded.CBOOL {
//			var temp template.HTML
//			tempValue, _ := v.Interface().(bool)
//			if tempValue {
//				temp = template.HTML(`<i class="fa fa-check-circle" aria-hidden=TRUE style="color:green;"></i>`)
//			} else {
//				temp = template.HTML(`<i class="fa fa-times-circle" aria-hidden=TRUE style="color:red;"></i>`)
//			}
//			y = append(y, temp)
//		} else if s.Fields[index].Type == preloaded.CLIST {
//			cIndex := v.Int()
//			choiceAdded := false
//			// @todo probably return
//			//if s.Fields[index].LimitChoicesTo != nil {
//			//	s.Fields[index].Choices = s.Fields[index].LimitChoicesTo(obj, &session.User)
//			//}
//			for cCounter := 0; cCounter < len(s.Fields[index].Choices); cCounter++ {
//				if uint(cIndex) == s.Fields[index].Choices[cCounter].K {
//					y = append(y, s.Fields[index].Choices[uint(cCounter)].V)
//					choiceAdded = true
//					break
//				}
//			}
//			if !choiceAdded {
//				y = append(y, cIndex)
//			}
//		} else if s.Fields[index].Type == preloaded.CIMAGE {
//			temp := template.HTML(fmt.Sprintf(`<img class="hvr-grow pointer image_trigger" style="max-width: 50px; height: auto;" src="%s" />`, v.Interface()))
//			y = append(y, temp)
//
//		} else if s.Fields[index].Type == preloaded.CFILE {
//			if v.Interface() != "" {
//				fileLocation := v.Interface().(string)
//				fileName := path.Base(fileLocation)
//				temp := template.HTML(fmt.Sprintf(`<a href="%s">%s</a>`, v.Interface(), fileName))
//				y = append(y, temp)
//			} else {
//				y = append(y, "")
//			}
//
//		} else if s.Fields[index].Type == preloaded.CCODE {
//			temp := template.HTML(fmt.Sprintf(`<pre style="width: 200px; white-space: pre-wrap;">%s</pre>`, v.Interface()))
//			y = append(y, temp)
//		} else if s.Fields[index].Type == preloaded.CMULTILINGUAL {
//			y = append(y, translation.Translate(fmt.Sprint(v), lang, true))
//		} else if s.Fields[index].Type == preloaded.CHTML {
//			str := helper.StripTags(fmt.Sprint(v))
//			y = append(y, str)
//		} else {
//
//			y = append(y, v)
//		}
//	}
//	return
//}

// @todo redo
//func GetFormData(a interface{}, r *http.Request, session *sessionmodel.Session, s *ModelSchema, user *usermodel.User) {
//	// This holds the formatted value of the field
//	var value interface{}
//	var f *F
//	var err error
//
//	// Get the type of model
//	t := reflect.TypeOf(a)
//
//	// Get the value of the model
//	modelValue := reflect.ValueOf(a)
//
//	// Get the primary key
//	newForm := r.FormValue("ModelID") == "0"
//	ModelID64, _ := strconv.ParseUint(r.FormValue("ModelID"), 10, 64)
//
//	// Loop over the fields of the model and get schema information
//	for index := 0; index < t.NumField(); index++ {
//		// Read field value
//		fieldValue := modelValue.Field(index)
//
//		// Get the field from schema
//		f = &F{}
//		fName := t.Field(index).Name
//		if t.Field(index).Anonymous {
//			fName = "ID"
//		}
//		for i := range s.Fields {
//			if s.Fields[i].Name == fName {
//				f = &s.Fields[i]
//				break
//			}
//		}
//		if f.Hidden || !f.FormDisplay {
//			continue
//		}
//		if f.Type == preloaded.CFK {
//			fieldValue = modelValue.FieldByName(f.Name + "ID")
//		}
//
//		// Check if the field was not found in schema
//		if f.Name == "" {
//			continue
//		}
//
//		// For new records:
//		// Overide field value with any values passed in request
//		// If not available check if there is a default value for the field
//		//if newForm {
//		if r.FormValue(t.Field(index).Name) != "" {
//			if f.Type == preloaded.CNUMBER || f.Type == preloaded.CLIST {
//				fValue, _ := strconv.ParseInt(r.FormValue(t.Field(index).Name), 10, 64)
//				fieldValue = reflect.ValueOf(fValue)
//			} else if f.Type == preloaded.CDATE {
//				var tm time.Time
//				if t.Field(index).Type.Kind() == reflect.Ptr {
//					var ptm *time.Time
//					if r.FormValue(t.Field(index).Name) == "" {
//						fieldValue = reflect.ValueOf(ptm)
//					} else {
//						tm, err = time.Parse("2006-01-02 15:04", r.FormValue(t.Field(index).Name))
//						if err != nil {
//							utils.Trail(utils.ERROR, "getFormData unable to parse time format(%s). %s", r.FormValue(t.Field(index).Name), err)
//						}
//						fieldValue = reflect.ValueOf(&tm)
//					}
//				} else {
//					tm, err = time.Parse("2006-01-02 15:04", r.FormValue(t.Field(index).Name))
//					if err != nil {
//						utils.Trail(utils.ERROR, "getFormData unable to parse time format(%s). %s", r.FormValue(t.Field(index).Name), err)
//					}
//					fieldValue = reflect.ValueOf(tm)
//				}
//			} else if f.Type == preloaded.CBOOL {
//				fieldValue = reflect.ValueOf(r.FormValue(t.Field(index).Name) == "on")
//			} else {
//				fieldValue = reflect.ValueOf(r.FormValue(t.Field(index).Name))
//			}
//		} else if f.DefaultValue != "" && newForm {
//			DefaultValue := f.DefaultValue
//			DefaultValue = strings.Replace(DefaultValue, "{NOW}", time.Now().Format("2006-01-02 15:04:05"), -1)
//			fieldValue = reflect.ValueOf(DefaultValue)
//		}
//
//		// Check for approval
//		if f.Approval {
//			// Check if there is an approval record
//			approvals := []approvalmodel.Approval{}
//			database.Filter(&approvals, "model_name = ? AND column_name = ? AND model_pk = ?", strings.ToLower(t.Name()), f.Name, ModelID64)
//			if len(approvals) != 0 {
//
//				// Get the last approval
//				lastA := approvals[len(approvals)-1]
//				f.ApprovalAction = lastA.ApprovalAction
//				f.NewValue = lastA.NewValueDescription
//				f.ChangedBy = lastA.ChangedBy
//				f.ChangeDate = &lastA.ChangeDate
//				f.ApprovalDate = lastA.ApprovalDate
//				f.ApprovalBy = lastA.ApprovalBy
//				f.ApprovalID = lastA.ID
//				f.OldValue = lastA.OldValue
//
//				// Remove required if the field has a pending approval
//				f.Required = f.Required && (f.ApprovalAction != ApprovalAction(0))
//			}
//		}
//
//		// Check the data type
//		if f.Type == preloaded.CID {
//			m, ok := fieldValue.Interface().(Model)
//			if !ok {
//				utils.Trail(utils.ERROR, "Unable tp parse value of ID for %s.%s (%#v)", t.Name(), f.Name, fieldValue.Interface())
//			}
//			value = m.ID
//		} else if f.Type == preloaded.CNUMBER {
//			if f.Format != "" {
//				value = fmt.Sprintf(f.Format, fieldValue.Interface())
//			} else {
//				value = fieldValue.Interface()
//			}
//		} else if f.Type == preloaded.CFK {
//			// Get selected items's ID
//			fkValue, _ := strconv.ParseUint(fmt.Sprint(fieldValue.Interface()), 10, 64)
//			value = fkValue
//
//			if f.LimitChoicesTo == nil {
//				fkType := t.Field(index).Type.Name()
//				if t.Field(index).Type.Kind() == reflect.Ptr {
//					fkType = t.Field(index).Type.Elem().Name()
//				}
//				fkList, _ := NewModelArray(strings.ToLower(fkType), false)
//				database.All(fkList.Addr().Interface())
//
//				// Build choices
//				f.Choices = []Choice{
//					{
//						K:        0,
//						V:        "-",
//						Selected: uint(fkValue) == 0,
//					},
//				}
//
//				for i := 0; i < fkList.Len(); i++ {
//					f.Choices = append(f.Choices, Choice{
//						K:        database.GetID(fkList.Index(i)),
//						V:        utils.GetString(fkList.Index(i).Interface()),
//						Selected: uint(fkValue) == database.GetID(fkList.Index(i)),
//					})
//				}
//			} else {
//				f.Choices = f.LimitChoicesTo(a, &session.User)
//
//				for i := 0; i < len(f.Choices); i++ {
//					f.Choices[i].Selected = uint(fkValue) == f.Choices[i].K
//				}
//			}
//
//		} else if f.Type == preloaded.CM2M {
//			if fmt.Sprint(reflect.TypeOf(fieldValue.Interface())) == "string" {
//				continue
//			}
//			fKType := reflect.TypeOf(fieldValue.Interface()).Elem()
//			m, ok := NewModelArray(strings.ToLower(fKType.Name()), false)
//
//			if !ok {
//				utils.Trail(utils.ERROR, "GetListSchema.NewModelArray. No model name (%s)", s.ModelName)
//			}
//			if f.LimitChoicesTo == nil {
//				database.All(m.Addr().Interface())
//				f.Choices = []Choice{}
//				for i := 0; i < m.Len(); i++ {
//					item := m.Index(i).Interface()
//					id := database.GetID(m.Index(i))
//					// if id == myID {
//					// 	continue
//					// }
//					f.Choices = append(f.Choices, Choice{
//						K: id,
//						V: utils.GetString(item),
//					})
//
//				}
//			} else {
//				f.Choices = f.LimitChoicesTo(a, &session.User)
//			}
//
//			for i := 0; i < fieldValue.Len(); i++ {
//				for counter, val := range f.Choices {
//					itemID := database.GetID(fieldValue.Index(i))
//					if val.K == itemID {
//						f.Choices[counter].Selected = true
//					}
//				}
//			}
//		} else if f.Type == preloaded.CDATE {
//			if newForm && t.Field(index).Type.Kind() != reflect.Ptr {
//				value = time.Now().Format("2006-01-02 15:04:05")
//			} else {
//				var d *time.Time
//				// If the date is not a pointer to date make it a pointer
//				if t.Field(index).Type.Kind() != reflect.Ptr {
//					tempD, _ := fieldValue.Interface().(time.Time)
//					d = &tempD
//				} else {
//					d, _ = fieldValue.Interface().(*time.Time)
//				}
//				if d == nil {
//					value = ""
//				} else {
//					value = d.Format("2006-01-02 15:04:05") //2006-01-02 15:04:05
//				}
//			}
//		} else if f.Type == preloaded.CBOOL {
//			d, ok := fieldValue.Interface().(bool)
//			if !ok {
//				utils.Trail(utils.ERROR, "Unable to parse bool value for %s.%s (%#v)", t.Name(), f.Name, fieldValue.Interface())
//			}
//			value = d
//		} else if f.Type == preloaded.CLIST {
//			value = fieldValue.Int()
//			if f.LimitChoicesTo != nil {
//				f.Choices = append([]Choice{{"-", 0, false}}, f.LimitChoicesTo(a, &session.User)...)
//			}
//			for i := range f.Choices {
//				f.Choices[i].Selected = f.Choices[i].K == uint(fieldValue.Int())
//			}
//		} else if f.Type == preloaded.CMULTILINGUAL {
//			value = fieldValue.Interface()
//			for i := range langmodel.GetActiveLanguages() {
//				f.Translations[i].Value = translation.Translate(fmt.Sprint(value), langmodel.ActiveLangs[i].Code, false)
//				if f.ChangedBy != "" {
//					f.Translations[i].NewValue = translation.Translate(fmt.Sprint(f.NewValue), langmodel.ActiveLangs[i].Code, false)
//					f.Translations[i].OldValue = translation.Translate(fmt.Sprint(f.OldValue), langmodel.ActiveLangs[i].Code, false)
//				}
//			}
//		} else if f.Type == preloaded.CLINK {
//			URL := fieldValue.String()
//			if URL != "" {
//				if strings.Contains(URL, "?") {
//					URL += "&"
//				} else {
//					URL += "?"
//				}
//				URL += "x-csrf-token=" + session.Key
//			}
//			value = URL
//		} else {
//			value = fieldValue.Interface()
//		}
//		f.Value = value
//
//	}
//
//	// Get data from method fields
//	for index := 0; index < t.NumMethod(); index++ {
//		// Check if the method should be included in the field list
//		if strings.Contains(t.Method(index).Name, "__Form") {
//			if strings.ToLower(string(t.Method(index).Name[0])) == string(t.Method(index).Name[0]) {
//				continue
//			}
//
//			in := []reflect.Value{}
//			ret := modelValue.Method(index).Call(in)
//			s.FieldByName(t.Method(index).Name).Value = template.HTML(utils.StripHTMLScriptTag(fmt.Sprint(ret[0].Interface())))
//		}
//	}
//
//	inlineData := []preloaded.ListData{}
//	if uint(ModelID64) != 0 {
//		for _, inlineS := range s.Inlines {
//			inlineModel, _ := NewModel(strings.ToLower(inlineS.ModelName), false)
//			//inlineQ := fmt.Sprintf("%s = %d", foreignKeys[s.ModelName][strings.ToLower(inlineS.ModelName)], ModelID64)
//			//r.Form.Set("inline_id", inlineQ)
//
//			// Check if there the inline has a ListModifier
//			query := ""
//			args := []interface{}{}
//			if inlineS.ListModifier != nil {
//				query, args = inlineS.ListModifier(inlineS, user)
//			}
//			// Add the fk for the inline
//			if query != "" {
//				query += " AND "
//			}
//			query += fmt.Sprintf("%s = ?", foreignKeys[s.ModelName][strings.ToLower(inlineS.ModelName)])
//			args = append(args, ModelID64)
//
//			rows := GetListData(inlineModel.Interface(), preloaded.PageLength, r, session, query, args...)
//			r.Form.Del("inline_id")
//			if rows.Count == 0 {
//				rows.Rows = [][]interface{}{}
//			}
//			inlineData = append(inlineData, *rows)
//		}
//	}
//	s.InlinesData = inlineData
//	s.ModelID = uint(ModelID64)
//}

// @todo, redo
//// GetFieldsAPI returns a list of fields in a model
//func GetFieldsAPI(w http.ResponseWriter, r *http.Request, session *sessionmodel.Session) {
//	modelName := r.FormValue("m")
//
//	response := []string{}
//	s := ModelSchema{}
//	for _, v := range Schema {
//		if v.ModelName == modelName {
//			s = v
//			break
//		}
//	}
//
//	for _, f := range s.Fields {
//		response = append(response, f.Name)
//	}
//	utils.ReturnJSON(w, r, response)
//}
//
