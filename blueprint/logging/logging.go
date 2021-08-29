package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/logging/migrations"
	logmodel "github.com/uadmin/uadmin/blueprint/logging/models"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	logAdminPage := interfaces.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) {return nil, nil},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {return nil},
	)
	logAdminPage.PageName = "Logs"
	logAdminPage.Slug = "log"
	logAdminPage.BlueprintName = "logging"
	logAdminPage.Router = mainRouter
	err := interfaces.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(logAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing log blueprint: %s", err))
	}
	logmodelAdminPage := interfaces.NewGormAdminPage(
		logAdminPage,
		func() (interface{}, interface{}) {return &logmodel.Log{}, &[]*logmodel.Log{}},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {
			fields := []string{"Username", "Action", "Activity", "CreatedAt", "ContentType", "ModelPK"}
			form := interfaces.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			form.ExtraStatic.ExtraJS = append(form.ExtraStatic.ExtraJS, "/static-inbuilt/uadmin/assets/highlight.js/highlight.pack.js")
			form.ExtraStatic.ExtraJS = append(form.ExtraStatic.ExtraJS, "/static-inbuilt/uadmin/assets/js/initialize.highlight.js")
			form.ExtraStatic.ExtraCSS = append(form.ExtraStatic.ExtraCSS, "/static-inbuilt/uadmin/assets/highlight.js/styles/default.css")
			actionField, _ := form.FieldRegistry.GetByName("Action")
			actionField.FieldConfig.Widget.SetPopulate(func(m interface{}, currentField *interfaces.Field) interface{} {
				a := m.(*logmodel.Log).Action
				return logmodel.HumanizeAction(a)
			})
			actionField.Ordering = 0
			usernameField, _ := form.FieldRegistry.GetByName("Username")
			usernameField.Ordering = 0
			activityField, _ := form.FieldRegistry.GetByName("Activity")
			activityField.FieldConfig.Widget.SetTemplateName("admin/widgets/textareajson")
			activityField.Ordering = 9
			createdAtField, _ := form.FieldRegistry.GetByName("CreatedAt")
			createdAtField.FieldConfig.Widget.SetReadonly(true)
			createdAtField.Ordering = 10
			return form
		},
	)
	logmodelAdminPage.PageName = "Logs"
	logmodelAdminPage.Slug = "log"
	logmodelAdminPage.BlueprintName = "logging"
	logmodelAdminPage.Router = mainRouter
	contentTypeListDisplay, _ := logmodelAdminPage.ListDisplay.GetFieldByDisplayName("ContentType")
	contentTypeListDisplay.Populate = func(m interface{}) string {
		return m.(*logmodel.Log).ContentType.String()
	}
	actionListDisplay, _ := logmodelAdminPage.ListDisplay.GetFieldByDisplayName("Action")
	actionListDisplay.Populate = func(m interface{}) string {
		return logmodel.HumanizeAction(m.(*logmodel.Log).Action)
	}
	logmodelAdminPage.NoPermissionToAddNew = true
	err = logAdminPage.SubPages.AddAdminPage(logmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing log blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	interfaces.ProjectModels.RegisterModel(func() interface{} {return &logmodel.Log{}})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "logging",
		Description:       "Logging blueprint is responsible to store all actions made through admin panel",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
