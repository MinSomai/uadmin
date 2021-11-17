package estweet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/es_uadmin_example/blueprint/estweet/migrations"
	"github.com/sergeyglazyrindev/uadmin/core"
)

type Blueprint struct {
	core.Blueprint
}

type Tweet struct {
	User    string `json:"user" uadmin:"list,search"`
	Message string `json:"message" uadmin:"list,search"`
	ID      string
}

func (t *Tweet) String() string {
	return fmt.Sprintf("User %s tweeted following %s", t.User, t.Message)
}

func (t *Tweet) SetID(ID string) {
	t.ID = ID
}

func (t *Tweet) GetID() string {
	return t.ID
}

func (t *Tweet) GetIndexName() string {
	return "tweets"
}

func (b Blueprint) InitRouter(app core.IApp, group *gin.RouterGroup) {
	tweetsAdminPage := core.NewElasticSearchAdminPage(
		nil,
		nil,
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	tweetsAdminPage.PageName = "Tweets"
	tweetsAdminPage.Slug = "tweets"
	tweetsAdminPage.BlueprintName = "tweets"
	tweetsAdminPage.Router = app.GetRouter()
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(tweetsAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing tweets blueprint: %s", err))
	}
	var tweetsmodelAdminPage *core.AdminPage
	tweetsmodelAdminPage = core.NewElasticSearchAdminPage(
		tweetsAdminPage,
		&Tweet{},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"User", "Message"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			return form
		},
	)
	tweetsmodelAdminPage.PageName = "Tweets"
	tweetsmodelAdminPage.Slug = "tweet"
	tweetsmodelAdminPage.BlueprintName = "tweets"
	tweetsmodelAdminPage.Router = app.GetRouter()
	tweetsmodelAdminPage.ModelName = "tweet"
	IDListDisplayField, _ := tweetsmodelAdminPage.ListDisplay.GetFieldByDisplayName("ID")
	IDListDisplayField.SortBy.SetSortCustomImplementation(func(afo core.IAdminFilterObjects, field *core.Field, direction int) {
		directionB := true
		if direction == -1 {
			directionB = false
		}
		afo.GetPaginatedQuerySet().Order(&core.ESSortBy{
			FieldName: "_id",
			Direction: directionB,
		})
	})
	err = tweetsAdminPage.SubPages.AddAdminPage(tweetsmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing tweets blueprint: %s", err))
	}
}

func (b Blueprint) InitApp(app core.IApp) {
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &Tweet{}, &[]*Tweet{} })
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "estweet",
		Description:       "test elasticsearch",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
