package approval

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/blueprint/approval/migrations"
	"github.com/uadmin/uadmin/blueprint/approval/models"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	approvalAdminPage := admin.NewGormAdminPage(nil, func() (interface{}, interface{}) {return nil, make([]interface{}, 0)}, "")
	approvalAdminPage.PageName = "Approvals"
	approvalAdminPage.Slug = "approval"
	approvalAdminPage.BlueprintName = "approval"
	approvalAdminPage.Router = mainRouter
	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(approvalAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing approval blueprint: %s", err))
	}
	approvalmodelAdminPage := admin.NewGormAdminPage(approvalAdminPage, func() (interface{}, interface{}) {return &models.Approval{}, &[]*models.Approval{}}, "approval")
	approvalmodelAdminPage.PageName = "Approval"
	approvalmodelAdminPage.Slug = "approval"
	approvalmodelAdminPage.BlueprintName = "approval"
	approvalmodelAdminPage.Router = mainRouter
	adminContext := &interfaces.AdminContext{}
	approvalForm := interfaces.NewFormFromModelFromGinContext(adminContext, &models.Approval{}, make([]string, 0), []string{}, true, "")
	approvalmodelAdminPage.Form = approvalForm
	IDField, _ := approvalForm.FieldRegistry.GetByName("ID")
	IDListDisplay := interfaces.NewListDisplay(IDField)
	IDListDisplay.Ordering = 1
	approvalmodelAdminPage.ListDisplay.AddField(IDListDisplay)
	approvalActionField, _ := approvalForm.FieldRegistry.GetByName("ApprovalAction")
	approvalActionListDisplay := interfaces.NewListDisplay(approvalActionField)
	approvalmodelAdminPage.ListDisplay.AddField(approvalActionListDisplay)
	approvalActionListDisplay.Populate = func(m interface{}) string {
		return models.HumanizeApprovalAction(m.(*models.Approval).ApprovalAction)
	}
	approvalActionListDisplay.Ordering = 2
	approvalByField, _ := approvalForm.FieldRegistry.GetByName("ApprovalBy")
	approvalByListDisplay := interfaces.NewListDisplay(approvalByField)
	approvalmodelAdminPage.ListDisplay.AddField(approvalByListDisplay)
	approvalByListDisplay.Ordering = 3
	approvalDateField, _ := approvalForm.FieldRegistry.GetByName("ApprovalDate")
	approvalDateListDisplay := interfaces.NewListDisplay(approvalDateField)
	approvalmodelAdminPage.ListDisplay.AddField(approvalDateListDisplay)
	approvalDateListDisplay.Ordering = 4
	modelNameField, _ := approvalForm.FieldRegistry.GetByName("ModelName")
	modelNameListDisplay := interfaces.NewListDisplay(modelNameField)
	modelNameListDisplay.Ordering = 5
	approvalmodelAdminPage.ListDisplay.AddField(modelNameListDisplay)
	modelPKField, _ := approvalForm.FieldRegistry.GetByName("ModelPK")
	modelPKListDisplay := interfaces.NewListDisplay(modelPKField)
	approvalmodelAdminPage.ListDisplay.AddField(modelPKListDisplay)
	modelPKListDisplay.Ordering = 6
	columnNameField, _ := approvalForm.FieldRegistry.GetByName("ColumnName")
	columnNameListDisplay := interfaces.NewListDisplay(columnNameField)
	columnNameListDisplay.Ordering = 7
	approvalmodelAdminPage.ListDisplay.AddField(columnNameListDisplay)
	oldValueField, _ := approvalForm.FieldRegistry.GetByName("OldValue")
	oldValueListDisplay := interfaces.NewListDisplay(oldValueField)
	oldValueListDisplay.Ordering = 8
	approvalmodelAdminPage.ListDisplay.AddField(oldValueListDisplay)
	newValueField, _ := approvalForm.FieldRegistry.GetByName("NewValue")
	newValueListDisplay := interfaces.NewListDisplay(newValueField)
	newValueListDisplay.Ordering = 9
	approvalmodelAdminPage.ListDisplay.AddField(newValueListDisplay)
	newValueDescriptionField, _ := approvalForm.FieldRegistry.GetByName("NewValueDescription")
	newValueDescriptionFieldListDisplay := interfaces.NewListDisplay(newValueDescriptionField)
	newValueDescriptionFieldListDisplay.Ordering = 10
	approvalmodelAdminPage.ListDisplay.AddField(newValueDescriptionFieldListDisplay)
	changedByField, _ := approvalForm.FieldRegistry.GetByName("ChangedBy")
	changedByListDisplay := interfaces.NewListDisplay(changedByField)
	changedByListDisplay.Ordering = 11
	approvalmodelAdminPage.ListDisplay.AddField(changedByListDisplay)
	changeDateField, _ := approvalForm.FieldRegistry.GetByName("ChangeDate")
	changeDateListDisplay := interfaces.NewListDisplay(changeDateField)
	changeDateListDisplay.Ordering = 12
	approvalmodelAdminPage.ListDisplay.AddField(changeDateListDisplay)
	viewRecordField, _ := approvalForm.FieldRegistry.GetByName("ViewRecord")
	viewRecordListDisplay := interfaces.NewListDisplay(viewRecordField)
	viewRecordListDisplay.Ordering = 13
	approvalmodelAdminPage.ListDisplay.AddField(viewRecordListDisplay)
	err = approvalAdminPage.SubPages.AddAdminPage(approvalmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing approval blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &models.Approval{}})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "approval",
		Description:       "Approval blueprint is responsible for approving things in the project",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
