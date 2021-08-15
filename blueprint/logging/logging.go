package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/blueprint/logging/migrations"
	logmodel "github.com/uadmin/uadmin/blueprint/logging/models"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	logAdminPage := admin.NewGormAdminPage(nil, func() (interface{}, interface{}) {return nil, nil}, "")
	logAdminPage.PageName = "Logs"
	logAdminPage.Slug = "log"
	logAdminPage.BlueprintName = "logging"
	logAdminPage.Router = mainRouter
	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(logAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing log blueprint: %s", err))
	}
	logmodelAdminPage := admin.NewGormAdminPage(logAdminPage, func() (interface{}, interface{}) {return &logmodel.Log{}, &[]*logmodel.Log{}}, "log")
	logmodelAdminPage.PageName = "Logs"
	logmodelAdminPage.Slug = "log"
	logmodelAdminPage.BlueprintName = "logging"
	logmodelAdminPage.Router = mainRouter
	adminContext := &interfaces.AdminContext{}
	logForm := interfaces.NewFormFromModelFromGinContext(adminContext, &logmodel.Log{}, make([]string, 0), []string{}, true, "")
	logmodelAdminPage.Form = logForm
	actionField, _ := logForm.FieldRegistry.GetByName("Action")
	actionListDisplay := interfaces.NewListDisplay(actionField)
	actionListDisplay.Ordering = 2
	actionListDisplay.Populate = func(m interface{}) string {
		return logmodel.HumanizeAction(m.(*logmodel.Log).Action)
	}
	logmodelAdminPage.ListDisplay.AddField(actionListDisplay)
	usernameField, _ := logForm.FieldRegistry.GetByName("Username")
	usernameListDisplay := interfaces.NewListDisplay(usernameField)
	usernameListDisplay.Ordering = 3
	logmodelAdminPage.ListDisplay.AddField(usernameListDisplay)
	tableNameField, _ := logForm.FieldRegistry.GetByName("TableName")
	tableNameListDisplay := interfaces.NewListDisplay(tableNameField)
	logmodelAdminPage.ListDisplay.AddField(tableNameListDisplay)
	tableNameListDisplay.Ordering = 4
	activityField, _ := logForm.FieldRegistry.GetByName("Activity")
	activityListDisplay := interfaces.NewListDisplay(activityField)
	logmodelAdminPage.ListDisplay.AddField(activityListDisplay)
	activityListDisplay.Ordering = 5
	rollbackField, _ := logForm.FieldRegistry.GetByName("RollBack")
	rollbackListDisplay := interfaces.NewListDisplay(rollbackField)
	logmodelAdminPage.ListDisplay.AddField(rollbackListDisplay)
	rollbackListDisplay.Ordering = 6
	createdAtField, _ := logForm.FieldRegistry.GetByName("CreatedAt")
	createdAtListDisplay := interfaces.NewListDisplay(createdAtField)
	createdAtListDisplay.Ordering = 7
	logmodelAdminPage.ListDisplay.AddField(createdAtListDisplay)
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
