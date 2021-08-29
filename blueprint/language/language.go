package language

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/language/migrations"
	"github.com/uadmin/uadmin/interfaces"
	"mime/multipart"
	"strconv"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	languageAdminPage := interfaces.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) {return nil, nil},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {return nil},
	)
	languageAdminPage.PageName = "Languages"
	languageAdminPage.Slug = "language"
	languageAdminPage.BlueprintName = "language"
	languageAdminPage.Router = mainRouter
	err := interfaces.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(languageAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing language blueprint: %s", err))
	}
	languagemodelAdminPage := interfaces.NewGormAdminPage(
		languageAdminPage,
		func() (interface{}, interface{}) {return &interfaces.Language{}, &[]*interfaces.Language{}},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {
			fields := []string{"EnglishName", "Name", "Flag", "Code", "RTL", "Default", "Active", "AvailableInGui"}
			form := interfaces.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			defaultField, _ := form.FieldRegistry.GetByName("Default")
			defaultField.Validators.AddValidator("only_one_default_language", func(i interface{}, o interface{}) error {
				isDefault := i.(bool)
				if !isDefault {
					return nil
				}
				d := o.(*multipart.Form)
				ID := d.Value["ID"][0]
				uadminDatabase := interfaces.NewUadminDatabase()
				lang := &interfaces.Language{}
				uadminDatabase.Db.Where(&interfaces.Language{Default: true}).First(lang)
				if lang.ID != 0 && ID != strconv.Itoa(int(lang.ID)) {
					return fmt.Errorf("only one default language could be configured")
				}
				return nil
			})
			return form
		},
	)
	languagemodelAdminPage.PageName = "Languages"
	languagemodelAdminPage.Slug = "language"
	languagemodelAdminPage.BlueprintName = "language"
	languagemodelAdminPage.Router = mainRouter
	languagemodelAdminPage.NoPermissionToAddNew = true
	err = languageAdminPage.SubPages.AddAdminPage(languagemodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing language blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &interfaces.Language{}})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "language",
		Description:       "Language blueprint is responsible for managing languages used in the project",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
