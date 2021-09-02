package abtest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/abtest/migrations"
	abtestmodel "github.com/uadmin/uadmin/blueprint/abtest/models"
	"github.com/uadmin/uadmin/core"
	"strconv"
)

type Blueprint struct {
	core.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	abTestAdminPage := core.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) { return nil, nil },
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	abTestAdminPage.PageName = "AB Tests"
	abTestAdminPage.Slug = "abtest"
	abTestAdminPage.BlueprintName = "abtest"
	abTestAdminPage.Router = mainRouter

	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(abTestAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
	abtestmodelAdminPage := core.NewGormAdminPage(
		abTestAdminPage,
		func() (interface{}, interface{}) { return &abtestmodel.ABTest{}, &[]*abtestmodel.ABTest{} },
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"ContentType", "Type", "Name", "Field", "PrimaryKey", "Active", "Group", "StaticPath"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			form.ExtraStatic.ExtraJS = append(form.ExtraStatic.ExtraJS, "/static-inbuilt/uadmin/assets/js/abtestformhandler.js")
			typeField, _ := form.FieldRegistry.GetByName("Type")
			w := typeField.FieldConfig.Widget.(*core.SelectWidget)
			w.OptGroups = make(map[string][]*core.SelectOptGroup)
			w.OptGroups[""] = make([]*core.SelectOptGroup, 0)
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "unknown",
				Value:    "0",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "static",
				Value:    "1",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "model",
				Value:    "2",
			})
			typeField.FieldConfig.Widget.SetPopulate(func(renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
				a := renderContext.Model.(*abtestmodel.ABTest).Type
				return strconv.Itoa(int(a))
			})
			typeField.SetUpField = func(w core.IWidget, m interface{}, v interface{}, afo core.IAdminFilterObjects) error {
				abTestM := m.(*abtestmodel.ABTest)
				vI, _ := strconv.Atoi(v.(string))
				abTestM.Type = abtestmodel.ABTestType(vI)
				return nil
			}
			contentTypeField, _ := form.FieldRegistry.GetByName("ContentType")
			w1 := contentTypeField.FieldConfig.Widget.(*core.ContentTypeSelectorWidget)
			w1.LoadFieldsOfAllModels = true
			fieldField, _ := form.FieldRegistry.GetByName("Field")
			w2 := fieldField.FieldConfig.Widget.(*core.SelectWidget)
			w2.SetAttr("data-initialized", "false")
			w2.DontValidateForExistence = true
			return form
		},
	)
	abtestmodelAdminPage.PageName = "AB Tests"
	abtestmodelAdminPage.Slug = "abtest"
	abtestmodelAdminPage.BlueprintName = "abtest"
	abtestmodelAdminPage.Router = mainRouter
	typeListDisplay, _ := abtestmodelAdminPage.ListDisplay.GetFieldByDisplayName("Type")
	typeListDisplay.Populate = func(m interface{}) string {
		return abtestmodel.HumanizeAbTestType(m.(*abtestmodel.ABTest).Type)
	}
	contentTypeListDisplay, _ := abtestmodelAdminPage.ListDisplay.GetFieldByDisplayName("ContentType")
	contentTypeListDisplay.Populate = func(m interface{}) string {
		return m.(*abtestmodel.ABTest).ContentType.String()
	}
	abTestValueInline := core.NewAdminPageInline(
		"AB Test Values",
		core.TabularInline, func(m interface{}) (interface{}, interface{}) {
			if m != nil {
				mO := m.(*abtestmodel.ABTest)
				return &abtestmodel.ABTestValue{ABTestID: mO.ID}, &[]*abtestmodel.ABTestValue{}
			}
			return &abtestmodel.ABTestValue{}, &[]*abtestmodel.ABTestValue{}
		}, func(afo core.IAdminFilterObjects, model interface{}, rp *core.AdminRequestParams) core.IAdminFilterObjects {
			abTest := model.(*abtestmodel.ABTest)
			var db *core.UadminDatabase
			if afo == nil {
				db = core.NewUadminDatabase()
			} else {
				db = afo.(*core.AdminFilterObjects).UadminDatabase
			}
			return &core.AdminFilterObjects{
				GormQuerySet:   core.NewGormPersistenceStorage(db.Db.Model(&abtestmodel.ABTestValue{}).Where(&abtestmodel.ABTestValue{ABTestID: abTest.ID})),
				Model:          &abtestmodel.ABTestValue{},
				UadminDatabase: db,
				GenerateModelI: func() (interface{}, interface{}) {
					return &abtestmodel.ABTestValue{}, &[]*abtestmodel.ABTestValue{}
				},
			}
		},
	)
	abTestValueInline.VerboseName = "AB Test Value"
	abTestValueInline.ListDisplay.AddField(&core.ListDisplay{
		DisplayName: "Click through rate",
		MethodName:  "ClickThroughRate",
	})
	abTestValueInline.ListDisplay.AddField(&core.ListDisplay{
		DisplayName: "Preview",
		MethodName:  "PreviewFormList",
	})
	abtestmodelAdminPage.InlineRegistry.Add(abTestValueInline)
	err = abTestAdminPage.SubPages.AddAdminPage(abtestmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	core.ProjectModels.RegisterModel(func() interface{} { return &abtestmodel.ABTestValue{} })
	core.ProjectModels.RegisterModel(func() interface{} { return &abtestmodel.ABTest{} })
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "abtest",
		Description:       "ABTest blueprint is responsible for ab tests",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
