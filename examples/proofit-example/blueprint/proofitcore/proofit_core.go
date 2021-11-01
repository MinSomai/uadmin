package proofitcore

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/proofit-example/blueprint/proofitcore/migrations"
	models2 "github.com/sergeyglazyrindev/proofit-example/blueprint/proofitcore/models"
	"github.com/sergeyglazyrindev/uadmin/core"
	"gorm.io/gorm/schema"
	"strconv"
)

type Blueprint struct {
	core.Blueprint
}

func (b Blueprint) InitRouter(app core.IApp, group *gin.RouterGroup) {
	// initialize administrator page for this blueprint.
	proofItAdminPage := core.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) { return nil, nil },
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	proofItAdminPage.PageName = "Proof It Core"
	proofItAdminPage.Slug = "proofit"
	proofItAdminPage.BlueprintName = "proofit_core"
	proofItAdminPage.Router = app.GetRouter()
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(proofItAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing blueprint: %s", err))
	}
	// initialize administrator page for your specific model.
	proofItAppAdminPage := core.NewGormAdminPage(
		proofItAdminPage,
		func() (interface{}, interface{}) {
			return &models2.ProofItApp{}, &[]*models2.ProofItApp{}
		},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			// define fields that you want to have in your admin panel
			model := modelI.(*models2.ProofItApp)
			fields := []string{"Name", "PublicKey", "SecretKey"}
			if model.ID == "" {
				fields = []string{"Name"}
			}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			languagesWidget := &core.SelectWidget{}
			languagesWidget.Multiple = true
			languagesWidget.InitializeAttrs()
			languagesWidget.RenderForAdmin()
			languagesWidget.SetName("Languages")
			languagesWidget.SetRequired()
			languagesWidget.OptGroups = make(map[string][]*core.SelectOptGroup)
			languagesWidget.OptGroups[""] = make([]*core.SelectOptGroup, 0)
			allLanguages := core.GetActiveLanguages()
			for _, lang := range allLanguages {
				languagesWidget.OptGroups[""] = append(languagesWidget.OptGroups[""], &core.SelectOptGroup{
					OptLabel: lang.Name,
					Value:    strconv.Itoa(int(lang.ID)),
				})

			}
			languagesWidget.Populate = func(renderContext *core.FormRenderContext, f *core.Field) interface{} {
				m := renderContext.Model.(*models2.ProofItApp)
				retV := make([]string, 0)
				for _, lang := range m.Languages {
					retV = append(retV, strconv.Itoa(int(lang.ID)))
				}
				return retV
			}
			field := &core.Field{
				Field: schema.Field{
					Name:         "Languages",
					DBName:       "",
					DefaultValue: "",
				},
				UadminFieldType: core.TextUadminFieldType,
				FieldConfig:     &core.FieldConfig{Widget: languagesWidget},
				Required:        true,
				DisplayName:     "Languages",
			}
			field.SetUpField = func(w core.IWidget, modelI interface{}, v interface{}, afo core.IAdminFilterObjects) error {
				m := modelI.(*models2.ProofItApp)
				lang := &core.Language{}
				realV := v.([]string)
				m.Languages = make([]core.Language, 0)
				for _, v1 := range realV {
					afo.GetDB().First(lang, v1)
					m.Languages = append(m.Languages, *lang)
					lang = &core.Language{}
				}
				return nil
			}
			form.FieldRegistry.AddField(field)
			return form
		},
	)
	proofItAppAdminPage.PageName = "Proof IT App"
	proofItAppAdminPage.Slug = "proofitapp"
	proofItAppAdminPage.BlueprintName = "proofit_core"
	proofItAppAdminPage.Router = app.GetRouter()
	err = proofItAdminPage.SubPages.AddAdminPage(proofItAppAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing blueprint: %s", err))
	}
	scienceCategoryAdminPage := core.NewGormAdminPage(
		proofItAdminPage,
		func() (interface{}, interface{}) {
			return &models2.ScienceCategory{}, &[]*models2.ScienceCategory{}
		},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			// define fields that you want to have in your admin panel
			fields := []string{"Name", "Icon"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			iconField, _ := form.FieldRegistry.GetByName("Icon")
			s3Storage := core.NewAWSS3Storage("test", &core.AWSConfig{
				S3: &core.AWSS3Config{
					Region:    "eu-central-1",
					AccessKey: "{{YOUR_ACCESS_TOKEN}}",
					SecretKey: "{{YOUR_SECRET_TOKEN}}",
				},
			})
			s3Storage.(*core.AWSS3Storage).Bucket = "testimagestorageforpippackage"
			s3Storage.(*core.AWSS3Storage).Domain = "https://testimagestorageforpippackage.s3.eu-central-1.amazonaws.com"
			iconField.FieldConfig.Widget.(*core.FileWidget).Storage = s3Storage
			return form
		},
	)
	core.UadminFormCongirurableOptionInstance.AddFieldFormOptions(&core.FieldFormOptions{
		WidgetType: "image",
		Name:       "ScienceCategoryPhotoFormOptions",
		WidgetPopulate: func(renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
			fsStorage := core.NewFsStorage()
			photo := renderContext.Model.(*models2.ScienceCategory).Icon
			if photo == "" {
				return ""
			}
			return fmt.Sprintf("%s%s", fsStorage.GetUploadURL(), photo)
		},
	})
	scienceCategoryAdminPage.PageName = "Science category"
	scienceCategoryAdminPage.Slug = "science-category"
	scienceCategoryAdminPage.BlueprintName = "proofit_core"
	scienceCategoryAdminPage.Router = app.GetRouter()
	scienceCategoryExpertsInline := core.NewAdminPageInline(
		"Experts",
		core.TabularInline, func(m interface{}) (interface{}, interface{}) {
			if m != nil {
				mO := m.(*models2.ScienceCategory)
				return &models2.ExpertScienceCategory{ScienceCategoryID: mO.ID}, &[]*models2.ExpertScienceCategory{}
			}
			return &models2.ExpertScienceCategory{}, &[]*models2.ExpertScienceCategory{}
		}, func(adminContext core.IAdminContext, afo core.IAdminFilterObjects, model interface{}) core.IAdminFilterObjects {
			scienceCategory := model.(*models2.ScienceCategory)
			var db *core.UadminDatabase
			if afo == nil {
				db = core.NewUadminDatabase()
			} else {
				db = afo.(*core.GormAdminFilterObjects).UadminDatabase
			}
			return &core.GormAdminFilterObjects{
				GormQuerySet:   core.NewGormPersistenceStorage(db.Db.Model(&models2.ExpertScienceCategory{ScienceCategoryID: scienceCategory.ID}).Where(&models2.ExpertScienceCategory{ScienceCategoryID: scienceCategory.ID})),
				Model:          &models2.ExpertScienceCategory{},
				UadminDatabase: db,
				GenerateModelI: func() (interface{}, interface{}) {
					return &models2.ExpertScienceCategory{}, &[]*models2.ExpertScienceCategory{}
				},
			}
		},
	)
	scienceCategoryExpertsInline.VerboseName = "Category experts"
	ldExpertField, _ := scienceCategoryExpertsInline.ListDisplay.GetFieldByDisplayName("Expert")
	widget := ldExpertField.Field.FieldConfig.Widget.(*core.ForeignKeyWidget)
	widget.GenerateModelInterface = func() (interface{}, interface{}) {
		return &models2.Expert{}, &[]*models2.Expert{}
	}
	ldExpertField.Field.FieldConfig.Widget = widget
	scienceCategoryExpertsInline.ListDisplay.RemoveFieldByName("ScienceCategory")
	// add custom fields to abTestValue inline
	scienceCategoryAdminPage.InlineRegistry.Add(scienceCategoryExpertsInline)
	scienceCategoryLocalizedInline := core.NewAdminPageInline(
		"Localizations",
		core.TabularInline, func(m interface{}) (interface{}, interface{}) {
			if m != nil {
				mO := m.(*models2.ScienceCategory)
				return &models2.ScienceCategoryLocalized{ScienceCategoryID: mO.ID}, &[]*models2.ScienceCategoryLocalized{}
			}
			return &models2.ScienceCategoryLocalized{}, &[]*models2.ScienceCategoryLocalized{}
		}, func(adminContext core.IAdminContext, afo core.IAdminFilterObjects, model interface{}) core.IAdminFilterObjects {
			scienceCategory := model.(*models2.ScienceCategory)
			var db *core.UadminDatabase
			if afo == nil {
				db = core.NewUadminDatabase()
			} else {
				db = afo.(*core.GormAdminFilterObjects).UadminDatabase
			}
			return &core.GormAdminFilterObjects{
				GormQuerySet:   core.NewGormPersistenceStorage(db.Db.Model(&models2.ScienceCategoryLocalized{ScienceCategoryID: scienceCategory.ID}).Where(&models2.ScienceCategoryLocalized{ScienceCategoryID: scienceCategory.ID})),
				Model:          &models2.ScienceCategoryLocalized{},
				UadminDatabase: db,
				GenerateModelI: func() (interface{}, interface{}) {
					return &models2.ScienceCategoryLocalized{}, &[]*models2.ScienceCategoryLocalized{}
				},
			}
		},
	)
	scienceCategoryLocalizedInline.VerboseName = "Localizations"
	languageCodeField, _ := scienceCategoryLocalizedInline.ListDisplay.GetFieldByDisplayName("LanguageCode")
	w := languageCodeField.Field.FieldConfig.Widget.(*core.SelectWidget)
	allLanguages := core.GetActiveLanguages()
	w.OptGroups = make(map[string][]*core.SelectOptGroup)
	w.OptGroups[""] = make([]*core.SelectOptGroup, 0)
	for _, lang := range allLanguages {
		w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
			OptLabel: lang.Name,
			Value:    lang.Code,
		})
	}
	scienceCategoryAdminPage.InlineRegistry.Add(scienceCategoryLocalizedInline)
	err = proofItAdminPage.SubPages.AddAdminPage(scienceCategoryAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing blueprint: %s", err))
	}
	scienceTermAdminPage := core.NewGormAdminPage(
		proofItAdminPage,
		func() (interface{}, interface{}) {
			return &models2.ScienceTerm{}, &[]*models2.ScienceTerm{}
		},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			// define fields that you want to have in your admin panel
			fields := []string{"Type", "Alias", "Discussion", "ScienceCategory"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			discussionField, _ := form.FieldRegistry.GetByName("Discussion")
			discussionField.FieldConfig.Widget.(*core.ForeignKeyWidget).GenerateModelInterface = func() (interface{}, interface{}) {
				return &models2.Discussion{}, &[]*models2.Discussion{}
			}
			scienceCategoryField, _ := form.FieldRegistry.GetByName("ScienceCategory")
			scienceCategoryField.FieldConfig.Widget.(*core.ForeignKeyWidget).GenerateModelInterface = func() (interface{}, interface{}) {
				return &models2.ScienceCategory{}, &[]*models2.ScienceCategory{}
			}
			typeField, _ := form.FieldRegistry.GetByName("Type")
			w := typeField.FieldConfig.Widget.(*core.SelectWidget)
			w.OptGroups = make(map[string][]*core.SelectOptGroup)
			w.OptGroups[""] = make([]*core.SelectOptGroup, 0)
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "axiom",
				Value:    "1",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "law",
				Value:    "2",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "deduction",
				Value:    "3",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "induction",
				Value:    "4",
			})
			typeField.FieldConfig.Widget.SetPopulate(func(renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
				a := renderContext.Model.(*models2.ScienceTerm).Type
				return strconv.Itoa(int(a))
			})
			typeField.SetUpField = func(w core.IWidget, m interface{}, v interface{}, afo core.IAdminFilterObjects) error {
				scienceTerm := m.(*models2.ScienceTerm)
				vI, _ := strconv.Atoi(v.(string))
				scienceTerm.Type = models2.ScienceTermType(vI)
				return nil
			}
			return form
		},
	)
	scienceTermAdminPage.PageName = "Science term"
	scienceTermAdminPage.Slug = "science-term"
	scienceTermAdminPage.BlueprintName = "proofit_core"
	scienceTermAdminPage.Router = app.GetRouter()
	scienceTermTypeListDisplay, _ := scienceTermAdminPage.ListDisplay.GetFieldByDisplayName("Type")
	scienceTermTypeListDisplay.Populate = func(m interface{}) string {
		return models2.HumanizeScienceTermType(m.(*models2.ScienceTerm).Type)
	}
	scienceTermLocalizedInline := core.NewAdminPageInline(
		"Localizations",
		core.TabularInline, func(m interface{}) (interface{}, interface{}) {
			if m != nil {
				mO := m.(*models2.ScienceTerm)
				return &models2.ScienceTermLocalized{ScienceTermID: mO.ID}, &[]*models2.ScienceTermLocalized{}
			}
			return &models2.ScienceTermLocalized{}, &[]*models2.ScienceTermLocalized{}
		}, func(adminContext core.IAdminContext, afo core.IAdminFilterObjects, model interface{}) core.IAdminFilterObjects {
			scienceTerm := model.(*models2.ScienceTerm)
			var db *core.UadminDatabase
			if afo == nil {
				db = core.NewUadminDatabase()
			} else {
				db = afo.(*core.GormAdminFilterObjects).UadminDatabase
			}
			return &core.GormAdminFilterObjects{
				GormQuerySet:   core.NewGormPersistenceStorage(db.Db.Model(&models2.ScienceTermLocalized{ScienceTermID: scienceTerm.ID}).Where(&models2.ScienceTermLocalized{ScienceTermID: scienceTerm.ID})),
				Model:          &models2.ScienceTermLocalized{},
				UadminDatabase: db,
				GenerateModelI: func() (interface{}, interface{}) {
					return &models2.ScienceTermLocalized{}, &[]*models2.ScienceTermLocalized{}
				},
			}
		},
	)
	scienceTermLocalizedInline.VerboseName = "Localizations"
	languageCodeField, _ = scienceTermLocalizedInline.ListDisplay.GetFieldByDisplayName("LanguageCode")
	w = languageCodeField.Field.FieldConfig.Widget.(*core.SelectWidget)
	w.OptGroups = make(map[string][]*core.SelectOptGroup)
	w.OptGroups[""] = make([]*core.SelectOptGroup, 0)
	for _, lang := range allLanguages {
		w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
			OptLabel: lang.Name,
			Value:    lang.Code,
		})
	}
	scienceTermAdminPage.InlineRegistry.Add(scienceTermLocalizedInline)
	err = proofItAdminPage.SubPages.AddAdminPage(scienceTermAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing blueprint: %s", err))
	}
	expertAdminPage := core.NewGormAdminPage(
		proofItAdminPage,
		func() (interface{}, interface{}) {
			return &models2.Expert{}, &[]*models2.Expert{}
		},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			// define fields that you want to have in your admin panel
			fields := []string{"User"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			userField, _ := form.FieldRegistry.GetByName("User")
			userFieldWidget := userField.FieldConfig.Widget.(*core.FkLinkWidget)
			userFieldWidget.Context = "edit"
			languagesWidget := &core.SelectWidget{}
			languagesWidget.Multiple = true
			languagesWidget.RenderForAdmin()
			languagesWidget.InitializeAttrs()
			languagesWidget.SetName("Languages")
			languagesWidget.SetRequired()
			languagesWidget.OptGroups = make(map[string][]*core.SelectOptGroup)
			languagesWidget.OptGroups[""] = make([]*core.SelectOptGroup, 0)
			allLanguages := core.GetActiveLanguages()
			for _, lang := range allLanguages {
				languagesWidget.OptGroups[""] = append(languagesWidget.OptGroups[""], &core.SelectOptGroup{
					OptLabel: lang.Name,
					Value:    strconv.Itoa(int(lang.ID)),
				})

			}
			languagesWidget.Populate = func(renderContext *core.FormRenderContext, f *core.Field) interface{} {
				m := renderContext.Model.(*models2.Expert)
				retV := make([]string, 0)
				for _, lang := range m.Languages {
					retV = append(retV, strconv.Itoa(int(lang.ID)))
				}
				return retV
			}
			field := &core.Field{
				Field: schema.Field{
					Name:         "Languages",
					DBName:       "",
					DefaultValue: "",
				},
				UadminFieldType: core.TextUadminFieldType,
				FieldConfig:     &core.FieldConfig{Widget: languagesWidget},
				Required:        true,
				DisplayName:     "Languages",
			}
			field.SetUpField = func(w core.IWidget, modelI interface{}, v interface{}, afo core.IAdminFilterObjects) error {
				m := modelI.(*models2.Expert)
				lang := &core.Language{}
				realV := v.([]string)
				m.Languages = make([]core.Language, 0)
				for _, v1 := range realV {
					afo.GetDB().First(lang, v1)
					m.Languages = append(m.Languages, *lang)
					lang = &core.Language{}
				}
				return nil
			}
			form.FieldRegistry.AddField(field)
			//userFieldWidget.GetQuerySet = func(formRenderContext *core.FormRenderContext) core.IPersistenceStorage {
			//	uadminDatabase := core.NewUadminDatabase()
			//	return core.NewGormPersistenceStorage(uadminDatabase.Db)
			//}
			//userFieldWidget.GenerateModelInterface = func() (interface{}, interface{}) {
			//	return &core.User{}, &[]*core.User{}
			//}
			//modelDescription := core.ProjectModels.GetModelFromInterface(&core.User{})
			//userAdminPage := core.CurrentDashboardAdminPanel.AdminPages.GetByModelName(modelDescription.Statement.Schema.Name)
			//userFieldWidget.AddNewLink = userAdminPage.GenerateLinkToAddNewModel()
			// form.Debug = true
			return form
		},
	)
	expertAdminPage.PageName = "Expert"
	expertAdminPage.Slug = "expert"
	expertAdminPage.BlueprintName = "proofit_core"
	expertAdminPage.Router = app.GetRouter()
	err = proofItAdminPage.SubPages.AddAdminPage(expertAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing blueprint: %s", err))
	}
	discussionAdminPage := core.NewGormAdminPage(
		proofItAdminPage,
		func() (interface{}, interface{}) {
			return &models2.Discussion{}, &[]*models2.Discussion{}
		},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			// define fields that you want to have in your admin panel
			fields := []string{"Author"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			authorField, _ := form.FieldRegistry.GetByName("Author")
			authorField.FieldConfig.Widget.(*core.ForeignKeyWidget).GenerateModelInterface = func() (interface{}, interface{}) {
				return &models2.Expert{}, &[]*models2.Expert{}
			}
			return form
		},
	)
	discussionAdminPage.PreloadData = func(afo core.IAdminFilterObjects) {
		afo.SetPaginatedQuerySet(afo.GetPaginatedQuerySet().Preload("Author.User"))
	}
	discussionAdminPage.PageName = "Discussion"
	discussionAdminPage.Slug = "discussion"
	discussionAdminPage.BlueprintName = "proofit_core"
	discussionAdminPage.Router = app.GetRouter()
	discussionLocalizedInline := core.NewAdminPageInline(
		"Localizations",
		core.TabularInline, func(m interface{}) (interface{}, interface{}) {
			if m != nil {
				mO := m.(*models2.Discussion)
				return &models2.DiscussionLocalized{DiscussionID: mO.ID}, &[]*models2.DiscussionLocalized{}
			}
			return &models2.DiscussionLocalized{}, &[]*models2.DiscussionLocalized{}
		}, func(adminContext core.IAdminContext, afo core.IAdminFilterObjects, model interface{}) core.IAdminFilterObjects {
			discussion := model.(*models2.Discussion)
			var db *core.UadminDatabase
			if afo == nil {
				db = core.NewUadminDatabase()
			} else {
				db = afo.(*core.GormAdminFilterObjects).UadminDatabase
			}
			return &core.GormAdminFilterObjects{
				GormQuerySet:   core.NewGormPersistenceStorage(db.Db.Model(&models2.DiscussionLocalized{DiscussionID: discussion.ID}).Where(&models2.DiscussionLocalized{DiscussionID: discussion.ID})),
				Model:          &models2.DiscussionLocalized{},
				UadminDatabase: db,
				GenerateModelI: func() (interface{}, interface{}) {
					return &models2.DiscussionLocalized{}, &[]*models2.DiscussionLocalized{}
				},
			}
		},
	)
	discussionLocalizedInline.VerboseName = "Localizations"
	languageCodeField, _ = discussionLocalizedInline.ListDisplay.GetFieldByDisplayName("LanguageCode")
	w = languageCodeField.Field.FieldConfig.Widget.(*core.SelectWidget)
	w.OptGroups = make(map[string][]*core.SelectOptGroup)
	w.OptGroups[""] = make([]*core.SelectOptGroup, 0)
	for _, lang := range allLanguages {
		w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
			OptLabel: lang.Name,
			Value:    lang.Code,
		})
	}
	discussionAdminPage.InlineRegistry.Add(discussionLocalizedInline)
	discussionReviewInline := core.NewAdminPageInline(
		"Discussion Reviews",
		core.TabularInline, func(m interface{}) (interface{}, interface{}) {
			if m != nil {
				mO := m.(*models2.Discussion)
				return &models2.DiscussionReview{DiscussionID: mO.ID}, &[]*models2.DiscussionReview{}
			}
			return &models2.DiscussionReview{}, &[]*models2.DiscussionReview{}
		}, func(adminContext core.IAdminContext, afo core.IAdminFilterObjects, model interface{}) core.IAdminFilterObjects {
			discussion := model.(*models2.Discussion)
			var db *core.UadminDatabase
			if afo == nil {
				db = core.NewUadminDatabase()
			} else {
				db = afo.(*core.GormAdminFilterObjects).UadminDatabase
			}
			return &core.GormAdminFilterObjects{
				GormQuerySet:   core.NewGormPersistenceStorage(db.Db.Model(&models2.DiscussionReview{DiscussionID: discussion.ID}).Where(&models2.DiscussionReview{DiscussionID: discussion.ID})),
				Model:          &models2.DiscussionReview{},
				UadminDatabase: db,
				GenerateModelI: func() (interface{}, interface{}) {
					return &models2.DiscussionReview{}, &[]*models2.DiscussionReview{}
				},
			}
		},
	)
	discussionReviewInline.VerboseName = "Discussion Reviews"
	scienceTermField1, _ := discussionReviewInline.ListDisplay.GetFieldByDisplayName("ScienceTerm")
	scienceTermField1.Field.FieldConfig.Widget.(*core.ForeignKeyWidget).GenerateModelInterface = func() (interface{}, interface{}) {
		return &models2.ScienceTerm{}, &[]*models2.ScienceTerm{}
	}
	reviewAuthorField1, _ := discussionReviewInline.ListDisplay.GetFieldByDisplayName("Author")
	reviewAuthorField1.Field.FieldConfig.Widget.(*core.ForeignKeyWidget).GenerateModelInterface = func() (interface{}, interface{}) {
		return &models2.Expert{}, &[]*models2.Expert{}
	}
	discussionReviewReasonField, _ := discussionReviewInline.ListDisplay.GetFieldByDisplayName("Reason")
	discussionReviewReasonField.Populate = func(m interface{}) string {
		return models2.HumanizeReasonType(m.(*models2.DiscussionReview).Reason)
	}
	discussionReviewReasonField.Field.FieldConfig.Widget.SetPopulate(func(renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
		return models2.HumanizeReasonType(renderContext.Model.(*models2.DiscussionReview).Reason)
	})
	discussionReviewReasonField.Field.SetUpField = func(w core.IWidget, m interface{}, v interface{}, afo core.IAdminFilterObjects) error {
		dR := m.(*models2.DiscussionReview)
		vI, _ := strconv.Atoi(v.(string))
		dR.Reason = models2.ReasonType(vI)
		return nil
	}
	w = discussionReviewReasonField.Field.FieldConfig.Widget.(*core.SelectWidget)
	w.OptGroups = make(map[string][]*core.SelectOptGroup)
	w.OptGroups[""] = make([]*core.SelectOptGroup, 0)
	w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
		OptLabel: "conforms",
		Value:    "2",
	})
	w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
		OptLabel: "contradicts",
		Value:    "1",
	})
	discussionReviewReasonField.Field.FieldConfig.Widget.SetPopulate(func(renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
		a := renderContext.Model.(*models2.DiscussionReview).Reason
		return strconv.Itoa(int(a))
	})
	discussionReviewReasonField.Field.SetUpField = func(w core.IWidget, m interface{}, v interface{}, afo core.IAdminFilterObjects) error {
		m1 := m.(*models2.DiscussionReview)
		vI, _ := strconv.Atoi(v.(string))
		m1.Reason = models2.ReasonType(vI)
		return nil
	}
	// add custom fields to abTestValue inline
	discussionAdminPage.InlineRegistry.Add(discussionReviewInline)
	// initialize inline for abtest, it shows all abtest values that belong to the current abtest object
	expertCategoriesInline := core.NewAdminPageInline(
		"Expert Categories",
		core.TabularInline, func(m interface{}) (interface{}, interface{}) {
			if m != nil {
				mO := m.(*models2.Expert)
				return &models2.ExpertScienceCategory{ExpertID: mO.ID}, &[]*models2.ExpertScienceCategory{}
			}
			return &models2.ExpertScienceCategory{}, &[]*models2.ExpertScienceCategory{}
		}, func(adminContext core.IAdminContext, afo core.IAdminFilterObjects, model interface{}) core.IAdminFilterObjects {
			expert := model.(*models2.Expert)
			var db *core.UadminDatabase
			if afo == nil {
				db = core.NewUadminDatabase()
			} else {
				db = afo.(*core.GormAdminFilterObjects).UadminDatabase
			}
			return &core.GormAdminFilterObjects{
				GormQuerySet:   core.NewGormPersistenceStorage(db.Db.Model(&models2.ExpertScienceCategory{ExpertID: expert.ID}).Where(&models2.ExpertScienceCategory{ExpertID: expert.ID})),
				Model:          &models2.ExpertScienceCategory{},
				UadminDatabase: db,
				GenerateModelI: func() (interface{}, interface{}) {
					return &models2.ExpertScienceCategory{}, &[]*models2.ExpertScienceCategory{}
				},
			}
		},
	)
	expertCategoriesInline.VerboseName = "Expert Categories"
	ldScienceCategory, _ := expertCategoriesInline.ListDisplay.GetFieldByDisplayName("ScienceCategory")
	widget = ldScienceCategory.Field.FieldConfig.Widget.(*core.ForeignKeyWidget)
	widget.GenerateModelInterface = func() (interface{}, interface{}) {
		return &models2.ScienceCategory{}, &[]*models2.ScienceCategory{}
	}
	ldScienceCategory.Field.FieldConfig.Widget = widget
	expertCategoriesInline.ListDisplay.RemoveFieldByName("Expert")
	// add custom fields to abTestValue inline
	expertAdminPage.InlineRegistry.Add(expertCategoriesInline)
	err = proofItAdminPage.SubPages.AddAdminPage(discussionAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing blueprint: %s", err))
	}
	discussionCommentAdminPage := core.NewGormAdminPage(
		proofItAdminPage,
		func() (interface{}, interface{}) {
			return &models2.DiscussionComment{}, &[]*models2.DiscussionComment{}
		},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			// define fields that you want to have in your admin panel
			fields := []string{"Discussion", "Parent", "Author"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			discussionField, _ := form.FieldRegistry.GetByName("Discussion")
			discussionField.FieldConfig.Widget.(*core.ForeignKeyWidget).GenerateModelInterface = func() (interface{}, interface{}) {
				return &models2.Discussion{}, &[]*models2.Discussion{}
			}
			parentField, _ := form.FieldRegistry.GetByName("Parent")
			parentField.FieldConfig.Widget.(*core.ForeignKeyWidget).GenerateModelInterface = func() (interface{}, interface{}) {
				return &models2.DiscussionComment{}, &[]*models2.DiscussionComment{}
			}
			authorField, _ := form.FieldRegistry.GetByName("Author")
			authorField.FieldConfig.Widget.(*core.ForeignKeyWidget).GenerateModelInterface = func() (interface{}, interface{}) {
				return &models2.Expert{}, &[]*models2.Expert{}
			}
			return form
		},
	)
	discussionCommentAdminPage.PageName = "Discussion Comment"
	discussionCommentAdminPage.Slug = "discussion-comment"
	discussionCommentAdminPage.BlueprintName = "proofit_core"
	discussionCommentAdminPage.Router = app.GetRouter()
	discussionCommentLocalizedInline := core.NewAdminPageInline(
		"Localizations",
		core.TabularInline, func(m interface{}) (interface{}, interface{}) {
			if m != nil {
				mO := m.(*models2.DiscussionComment)
				return &models2.DiscussionCommentLocalized{DiscussionCommentID: mO.ID}, &[]*models2.DiscussionCommentLocalized{}
			}
			return &models2.DiscussionCommentLocalized{}, &[]*models2.DiscussionCommentLocalized{}
		}, func(adminContext core.IAdminContext, afo core.IAdminFilterObjects, model interface{}) core.IAdminFilterObjects {
			discussionComment := model.(*models2.DiscussionComment)
			var db *core.UadminDatabase
			if afo == nil {
				db = core.NewUadminDatabase()
			} else {
				db = afo.(*core.GormAdminFilterObjects).UadminDatabase
			}
			return &core.GormAdminFilterObjects{
				GormQuerySet:   core.NewGormPersistenceStorage(db.Db.Model(&models2.DiscussionCommentLocalized{DiscussionCommentID: discussionComment.ID}).Where(&models2.DiscussionCommentLocalized{DiscussionCommentID: discussionComment.ID})),
				Model:          &models2.DiscussionCommentLocalized{},
				UadminDatabase: db,
				GenerateModelI: func() (interface{}, interface{}) {
					return &models2.DiscussionCommentLocalized{}, &[]*models2.DiscussionCommentLocalized{}
				},
			}
		},
	)
	discussionCommentLocalizedInline.VerboseName = "Localizations"
	languageCodeField, _ = discussionCommentLocalizedInline.ListDisplay.GetFieldByDisplayName("LanguageCode")
	w = languageCodeField.Field.FieldConfig.Widget.(*core.SelectWidget)
	w.OptGroups = make(map[string][]*core.SelectOptGroup)
	w.OptGroups[""] = make([]*core.SelectOptGroup, 0)
	for _, lang := range allLanguages {
		w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
			OptLabel: lang.Name,
			Value:    lang.Code,
		})
	}
	discussionCommentAdminPage.InlineRegistry.Add(discussionCommentLocalizedInline)
	discussionCommentReviewInline := core.NewAdminPageInline(
		"Discussion Comment Reviews",
		core.TabularInline, func(m interface{}) (interface{}, interface{}) {
			if m != nil {
				mO := m.(*models2.DiscussionComment)
				return &models2.DiscussionCommentReview{DiscussionCommentID: mO.ID}, &[]*models2.DiscussionCommentReview{}
			}
			return &models2.DiscussionCommentReview{}, &[]*models2.DiscussionCommentReview{}
		}, func(adminContext core.IAdminContext, afo core.IAdminFilterObjects, model interface{}) core.IAdminFilterObjects {
			discussionComment := model.(*models2.DiscussionComment)
			var db *core.UadminDatabase
			if afo == nil {
				db = core.NewUadminDatabase()
			} else {
				db = afo.(*core.GormAdminFilterObjects).UadminDatabase
			}
			return &core.GormAdminFilterObjects{
				GormQuerySet:   core.NewGormPersistenceStorage(db.Db.Model(&models2.DiscussionCommentReview{DiscussionCommentID: discussionComment.ID}).Where(&models2.DiscussionCommentReview{DiscussionCommentID: discussionComment.ID})),
				Model:          &models2.DiscussionCommentReview{},
				UadminDatabase: db,
				GenerateModelI: func() (interface{}, interface{}) {
					return &models2.DiscussionCommentReview{}, &[]*models2.DiscussionCommentReview{}
				},
			}
		},
	)
	discussionCommentReviewInline.VerboseName = "Discussion Comment Reviews"
	scienceTermField2, _ := discussionCommentReviewInline.ListDisplay.GetFieldByDisplayName("ScienceTerm")
	scienceTermField2.Field.FieldConfig.Widget.(*core.ForeignKeyWidget).GenerateModelInterface = func() (interface{}, interface{}) {
		return &models2.ScienceTerm{}, &[]*models2.ScienceTerm{}
	}
	reviewAuthorField2, _ := discussionCommentReviewInline.ListDisplay.GetFieldByDisplayName("Author")
	reviewAuthorField2.Field.FieldConfig.Widget.(*core.ForeignKeyWidget).GenerateModelInterface = func() (interface{}, interface{}) {
		return &models2.Expert{}, &[]*models2.Expert{}
	}
	discussionReviewReasonField, _ = discussionCommentReviewInline.ListDisplay.GetFieldByDisplayName("Reason")
	discussionReviewReasonField.Populate = func(m interface{}) string {
		return models2.HumanizeReasonType(m.(*models2.DiscussionCommentReview).Reason)
	}
	discussionReviewReasonField.Field.FieldConfig.Widget.SetPopulate(func(renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
		return models2.HumanizeReasonType(renderContext.Model.(*models2.DiscussionCommentReview).Reason)
	})
	discussionReviewReasonField.Field.SetUpField = func(w core.IWidget, m interface{}, v interface{}, afo core.IAdminFilterObjects) error {
		dR := m.(*models2.DiscussionCommentReview)
		vI, _ := strconv.Atoi(v.(string))
		dR.Reason = models2.ReasonType(vI)
		return nil
	}
	w = discussionReviewReasonField.Field.FieldConfig.Widget.(*core.SelectWidget)
	w.OptGroups = make(map[string][]*core.SelectOptGroup)
	w.OptGroups[""] = make([]*core.SelectOptGroup, 0)
	w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
		OptLabel: "conforms",
		Value:    "2",
	})
	w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
		OptLabel: "contradicts",
		Value:    "1",
	})
	discussionReviewReasonField.Field.FieldConfig.Widget.SetPopulate(func(renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
		a := renderContext.Model.(*models2.DiscussionCommentReview).Reason
		return strconv.Itoa(int(a))
	})
	discussionReviewReasonField.Field.SetUpField = func(w core.IWidget, m interface{}, v interface{}, afo core.IAdminFilterObjects) error {
		m1 := m.(*models2.DiscussionCommentReview)
		vI, _ := strconv.Atoi(v.(string))
		m1.Reason = models2.ReasonType(vI)
		return nil
	}
	// add custom fields to abTestValue inline
	discussionCommentAdminPage.InlineRegistry.Add(discussionCommentReviewInline)
	discussionCommentAdminPage.PreloadData = func(afo core.IAdminFilterObjects) {
		afo.SetPaginatedQuerySet(afo.GetPaginatedQuerySet().Preload("Author.User").Preload("Discussion.Author.User"))
	}
	err = proofItAdminPage.SubPages.AddAdminPage(discussionCommentAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing blueprint: %s", err))
	}
}

func (b Blueprint) InitApp(app core.IApp) {
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) {
		return &models2.ExpertScienceCategory{}, &[]*models2.ExpertScienceCategory{}
	})
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) {
		return &models2.DiscussionReview{}, &[]*models2.DiscussionReview{}
	})
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) {
		return &models2.DiscussionCommentReview{}, &[]*models2.DiscussionCommentReview{}
	})
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) {
		return &models2.ScienceCategoryLocalized{}, &[]*models2.ScienceCategoryLocalized{}
	})
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) {
		return &models2.DiscussionLocalized{}, &[]*models2.DiscussionLocalized{}
	})
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) {
		return &models2.ScienceTermLocalized{}, &[]*models2.ScienceTermLocalized{}
	})
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) {
		return &models2.DiscussionCommentLocalized{}, &[]*models2.DiscussionCommentLocalized{}
	})
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "proofit_core",
		Description:       "Proofit-Core",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
