package approval

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/approval/migrations"
	"github.com/uadmin/uadmin/blueprint/approval/models"
	"github.com/uadmin/uadmin/core"
	"strconv"
)

type Blueprint struct {
	core.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	approvalAdminPage := core.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) {return nil, make([]interface{}, 0)},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {return nil},
	)
	approvalAdminPage.PageName = "Approvals"
	approvalAdminPage.Slug = "approval"
	approvalAdminPage.BlueprintName = "approval"
	approvalAdminPage.Router = mainRouter
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(approvalAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing approval blueprint: %s", err))
	}
	approvalmodelAdminPage := core.NewGormAdminPage(
		approvalAdminPage,
		func() (interface{}, interface{}) {return &models.Approval{}, &[]*models.Approval{}},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"ContentType", "ModelPK", "ColumnName", "OldValue", "NewValue", "NewValueDescription", "ChangedBy", "ChangeDate", "ApprovalAction", "ApprovalBy", "ApprovalDate"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			approvalField, _ := form.FieldRegistry.GetByName("ApprovalAction")
			w := approvalField.FieldConfig.Widget.(*core.SelectWidget)
			w.OptGroups = make(map[string][]*core.SelectOptGroup)
			w.OptGroups[""] = make([]*core.SelectOptGroup, 0)
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "unknown",
				Value: "0",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "approved",
				Value: "1",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "rejected",
				Value: "2",
			})
			approvalField.FieldConfig.Widget.SetPopulate(func(m interface{}, currentField *core.Field) interface{} {
				a := m.(*models.Approval).ApprovalAction
				return strconv.Itoa(int(a))
			})
			approvalField.SetUpField = func(w core.IWidget, m interface{}, v interface{}, afo core.IAdminFilterObjects) error {
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
		return aD.Format(core.CurrentConfig.D.Uadmin.DateTimeFormat)
	}
	contentTypeListDisplay, _ := approvalmodelAdminPage.ListDisplay.GetFieldByDisplayName("ContentType")
	contentTypeListDisplay.Populate = func(m interface{}) string {
		return m.(*models.Approval).ContentType.String()
	}
	changeDateListDisplay, _ := approvalmodelAdminPage.ListDisplay.GetFieldByDisplayName("ChangeDate")
	changeDateListDisplay.Populate = func(m interface{}) string {
		return m.(*models.Approval).ChangeDate.Format(core.CurrentConfig.D.Uadmin.DateTimeFormat)
	}
	err = approvalAdminPage.SubPages.AddAdminPage(approvalmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing approval blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	core.ProjectModels.RegisterModel(func() interface{}{return &models.Approval{}})
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "approval",
		Description:       "Approval blueprint is responsible for approving things in the project",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
