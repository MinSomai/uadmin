package uadmin

import (
	"fmt"
	"github.com/uadmin/uadmin/admin"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/interfaces"
	"os"
)

type ContentTypeCommand struct {
}

func (c ContentTypeCommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &CommandRegistry{
		Actions: make(map[string]interfaces.ICommand),
	}

	commandRegistry.addAction("sync", &SyncContentTypes{})
	if len(os.Args) > 2 {
		action = os.Args[2]
		isCorrectActionPassed = commandRegistry.isRegisteredCommand(action)
	}
	if !isCorrectActionPassed {
		helpText := commandRegistry.MakeHelpText()
		help = fmt.Sprintf(`
Please provide what do you want to do ?
%s
`, helpText)
		fmt.Print(help)
		return nil
	}
	commandRegistry.runAction(subaction, "", args)
	return nil
}

func (c ContentTypeCommand) GetHelpText() string {
	return "Content type for uadmin project"
}

type SyncContentTypes struct {
}

func (command SyncContentTypes) Proceed(subaction string, args []string) error {
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	db := uadminDatabase.Db
	var contentType interfaces.ContentType
	var permission usermodels.Permission
	for blueprintRootAdminPage := range admin.CurrentDashboardAdminPanel.AdminPages.GetAll() {
		interfaces.Trail(interfaces.INFO, "Sync content types for blueprint %s", blueprintRootAdminPage.BlueprintName)
		for modelPage := range blueprintRootAdminPage.SubPages.GetAll() {
			db.Model(&interfaces.ContentType{}).Where(
				&interfaces.ContentType{BlueprintName: modelPage.BlueprintName, ModelName: modelPage.ModelName},
			).First(&contentType)
			if contentType.ID == 0 {
				contentType = interfaces.ContentType{BlueprintName: modelPage.BlueprintName, ModelName: modelPage.ModelName}
				db.Create(&contentType)
				interfaces.Trail(interfaces.INFO, "Created content type for blueprint %s model %s", modelPage.BlueprintName, modelPage.ModelName)
			}
			for permDescribed := range interfaces.ProjectPermRegistry.GetAllPermissions() {
				db.Model(&usermodels.Permission{}).Where(
					&usermodels.Permission{ContentTypeID: contentType.ID, PermissionBits: permDescribed.Bit},
				).First(&permission)
				if permission.ID == 0 {
					permission = usermodels.Permission{ContentTypeID: contentType.ID, PermissionBits: permDescribed.Bit}
					db.Create(&permission)
					interfaces.Trail(interfaces.INFO, "Created permission %s for blueprint %s model %s", permDescribed.Name, modelPage.BlueprintName, modelPage.ModelName)
					permission = usermodels.Permission{}
				}
				permission = usermodels.Permission{}
			}
			contentType = interfaces.ContentType{}
		}
	}
	return nil
}

func (command SyncContentTypes) GetHelpText() string {
	return "Sync your content types"
}
