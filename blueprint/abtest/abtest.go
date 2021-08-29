package abtest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/abtest/migrations"
	abtestmodel "github.com/uadmin/uadmin/blueprint/abtest/models"
	"github.com/uadmin/uadmin/interfaces"
	"strconv"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	abTestAdminPage := interfaces.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) {return nil, nil},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {return nil},
	)
	abTestAdminPage.PageName = "AB Tests"
	abTestAdminPage.Slug = "abtest"
	abTestAdminPage.BlueprintName = "abtest"
	abTestAdminPage.Router = mainRouter

	err := interfaces.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(abTestAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
	abtestmodelAdminPage := interfaces.NewGormAdminPage(
		abTestAdminPage,
		func() (interface{}, interface{}) {return &abtestmodel.ABTest{}, &[]*abtestmodel.ABTest{}},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {
			fields := []string{"ContentType", "Type", "Name", "Field", "PrimaryKey", "Active", "Group", "StaticPath"}
			form := interfaces.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			form.ExtraStatic.ExtraJS = append(form.ExtraStatic.ExtraJS, "/static-inbuilt/uadmin/assets/js/abtestformhandler.js")
			typeField, _ := form.FieldRegistry.GetByName("Type")
			w := typeField.FieldConfig.Widget.(*interfaces.SelectWidget)
			w.OptGroups = make(map[string][]*interfaces.SelectOptGroup)
			w.OptGroups[""] = make([]*interfaces.SelectOptGroup, 0)
			w.OptGroups[""] = append(w.OptGroups[""], &interfaces.SelectOptGroup{
				OptLabel: "unknown",
				Value: "0",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &interfaces.SelectOptGroup{
				OptLabel: "static",
				Value: "1",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &interfaces.SelectOptGroup{
				OptLabel: "model",
				Value: "2",
			})
			typeField.FieldConfig.Widget.SetPopulate(func(m interface{}, currentField *interfaces.Field) interface{} {
				a := m.(*abtestmodel.ABTest).Type
				return strconv.Itoa(int(a))
			})
			typeField.SetUpField = func(w interfaces.IWidget, m interface{}, v interface{}, afo interfaces.IAdminFilterObjects) error {
				abTestM := m.(*abtestmodel.ABTest)
				vI, _ := strconv.Atoi(v.(string))
				abTestM.Type = abtestmodel.ABTestType(vI)
				return nil
			}
			contentTypeField, _ := form.FieldRegistry.GetByName("ContentType")
			w1 := contentTypeField.FieldConfig.Widget.(*interfaces.ContentTypeSelectorWidget)
			w1.LoadFieldsOfAllModels = true
			fieldField, _ := form.FieldRegistry.GetByName("Field")
			w2 := fieldField.FieldConfig.Widget.(*interfaces.SelectWidget)
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
	abTestValueInline := interfaces.NewAdminPageInline(
		"AB Test Values",
		interfaces.TabularInline, func(m interface{}) (interface{}, interface{}) {
			if m != nil {
				mO := m.(*abtestmodel.ABTest)
				return &abtestmodel.ABTestValue{ABTestID: mO.ID}, &[]*abtestmodel.ABTestValue{}
			}
			return &abtestmodel.ABTestValue{}, &[]*abtestmodel.ABTestValue{}
		}, func(afo interfaces.IAdminFilterObjects, model interface{}, rp *interfaces.AdminRequestParams) interfaces.IAdminFilterObjects {
			abTest := model.(*abtestmodel.ABTest)
			var db *interfaces.UadminDatabase
			if afo == nil {
				db = interfaces.NewUadminDatabase()
			} else {
				db = afo.(*interfaces.AdminFilterObjects).UadminDatabase
			}
			return &interfaces.AdminFilterObjects{
				GormQuerySet: interfaces.NewGormPersistenceStorage(db.Db.Model(&abtestmodel.ABTestValue{}).Where(&abtestmodel.ABTestValue{ABTestID: abTest.ID})),
				Model: &abtestmodel.ABTestValue{},
				UadminDatabase: db,
				GenerateModelI: func() (interface{}, interface{}) {
					return &abtestmodel.ABTestValue{}, &[]*abtestmodel.ABTestValue{}
				},
			}
		},
	)
	abTestValueInline.VerboseName = "AB Test Value"
	abTestValueInline.ListDisplay.AddField(&interfaces.ListDisplay{
		DisplayName: "Click through rate",
		MethodName: "ClickThroughRate",
	})
	abTestValueInline.ListDisplay.AddField(&interfaces.ListDisplay{
		DisplayName: "Preview",
		MethodName: "PreviewFormList",
	})
	abtestmodelAdminPage.InlineRegistry.Add(abTestValueInline)
	err = abTestAdminPage.SubPages.AddAdminPage(abtestmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &abtestmodel.ABTestValue{}})
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &abtestmodel.ABTest{}})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "abtest",
		Description:       "ABTest blueprint is responsible for ab tests",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
