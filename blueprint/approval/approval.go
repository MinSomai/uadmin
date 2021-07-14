package approval

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/blueprint/approval/migrations"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	approvalAdminPage := admin.NewAdminPage("")
	approvalAdminPage.PageName = "Approvals"
	approvalAdminPage.Slug = "approval"
	approvalAdminPage.BlueprintName = "approval"
	approvalAdminPage.Router = group
	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(approvalAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing approval blueprint: %s", err))
	}
	approvalmodelAdminPage := admin.NewAdminPage("approval")
	approvalmodelAdminPage.PageName = "Approval"
	approvalmodelAdminPage.Slug = "approval"
	approvalmodelAdminPage.BlueprintName = "approval"
	approvalmodelAdminPage.Router = group
	err = approvalAdminPage.SubPages.AddAdminPage(approvalmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing approval blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "approval",
		Description:       "Approval blueprint is responsible for approving things in the project",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
