package approval

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/approval/migrations"
	"github.com/uadmin/uadmin/blueprint/approval/models"
	"github.com/uadmin/uadmin/interfaces"
	"strconv"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	approvalAdminPage := interfaces.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) {return nil, make([]interface{}, 0)},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {return nil},
	)
	approvalAdminPage.PageName = "Approvals"
	approvalAdminPage.Slug = "approval"
	approvalAdminPage.BlueprintName = "approval"
	approvalAdminPage.Router = mainRouter
	err := interfaces.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(approvalAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing approval blueprint: %s", err))
	}
	approvalmodelAdminPage := interfaces.NewGormAdminPage(
		approvalAdminPage,
		func() (interface{}, interface{}) {return &models.Approval{}, &[]*models.Approval{}},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {
			fields := []string{"ContentType", "ModelPK", "ColumnName", "OldValue", "NewValue", "NewValueDescription", "ChangedBy", "ChangeDate", "ApprovalAction", "ApprovalBy", "ApprovalDate"}
			form := interfaces.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			approvalField, _ := form.FieldRegistry.GetByName("ApprovalAction")
			w := approvalField.FieldConfig.Widget.(*interfaces.SelectWidget)
			w.OptGroups = make(map[string][]*interfaces.SelectOptGroup)
			w.OptGroups[""] = make([]*interfaces.SelectOptGroup, 0)
			w.OptGroups[""] = append(w.OptGroups[""], &interfaces.SelectOptGroup{
				OptLabel: "unknown",
				Value: "0",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &interfaces.SelectOptGroup{
				OptLabel: "approved",
				Value: "1",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &interfaces.SelectOptGroup{
				OptLabel: "rejected",
				Value: "2",
			})
			approvalField.FieldConfig.Widget.SetPopulate(func(m interface{}, currentField *interfaces.Field) interface{} {
				a := m.(*models.Approval).ApprovalAction
				return strconv.Itoa(int(a))
			})
			approvalField.SetUpField = func(w interfaces.IWidget, m interface{}, v interface{}, afo interfaces.IAdminFilterObjects) error {
				approvalM := m.(*models.Approval)
				vI, _ := strconv.Atoi(v.(string))
				approvalM.ApprovalAction = models.ApprovalAction(vI)
				return nil
			}
			return form
		},
	)
	approvalmodelAdminPage.PageName = "Approval"
	approvalmodelAdminPage.Slug = "approval"
	approvalmodelAdminPage.BlueprintName = "approval"
	approvalmodelAdminPage.Router = mainRouter
	approvalActionListDisplay, _ := approvalmodelAdminPage.ListDisplay.GetFieldByDisplayName("ApprovalAction")
	approvalActionListDisplay.Populate = func(m interface{}) string {
		return models.HumanizeApprovalAction(m.(*models.Approval).ApprovalAction)
	}
	approvalDateListDisplay, _ := approvalmodelAdminPage.ListDisplay.GetFieldByDisplayName("ApprovalDate")
	approvalDateListDisplay.Populate = func(m interface{}) string {
		aD := m.(*models.Approval).ApprovalDate
		if aD == nil {
			return ""
		}
		return aD.Format(interfaces.CurrentConfig.D.Uadmin.DateTimeFormat)
	}
	contentTypeListDisplay, _ := approvalmodelAdminPage.ListDisplay.GetFieldByDisplayName("ContentType")
	contentTypeListDisplay.Populate = func(m interface{}) string {
		return m.(*models.Approval).ContentType.String()
	}
	changeDateListDisplay, _ := approvalmodelAdminPage.ListDisplay.GetFieldByDisplayName("ChangeDate")
	changeDateListDisplay.Populate = func(m interface{}) string {
		return m.(*models.Approval).ChangeDate.Format(interfaces.CurrentConfig.D.Uadmin.DateTimeFormat)
	}
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
