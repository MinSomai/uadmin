package auth

import (
	"github.com/gin-gonic/gin"
	interfaces3 "github.com/uadmin/uadmin/blueprint/auth/interfaces"
	"github.com/uadmin/uadmin/blueprint/auth/migrations"
	langmodel "github.com/uadmin/uadmin/blueprint/language/models"
	menumodel "github.com/uadmin/uadmin/blueprint/menu/models"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/utils"
)

type Blueprint struct {
	interfaces.Blueprint
	AuthAdapterRegistry *interfaces3.AuthProviderRegistry
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	for adapter := range b.AuthAdapterRegistry.Iterate() {
		adapterGroup := group.Group("/" + adapter.GetName())
		adapterGroup.POST("/signin/", adapter.Signin)
		adapterGroup.POST("/signup/", adapter.Signup)
		adapterGroup.POST("/logout/", adapter.Logout)
		adapterGroup.GET("/status/", adapter.IsAuthenticated)
	}
	mainRouter.GET(config.CurrentConfig.D.Uadmin.RootAdminURL, func(ctx *gin.Context) {
		defaultAdapter, _ := b.AuthAdapterRegistry.GetAdapter("direct-for-admin")
		userSession := defaultAdapter.GetSession(ctx)
		if userSession == nil {
			type Context struct {
				Err         string
				ErrExists   bool
				SiteName    string
				Languages   []langmodel.Language
				RootURL     string
				OTPRequired bool
				Language    *langmodel.Language
				Username    string
				Password    string
				Logo        string
				FavIcon     string
			}
			c := Context{}
			c.SiteName = config.CurrentConfig.D.Uadmin.SiteName
			c.RootURL = config.CurrentConfig.D.Uadmin.RootURL
			c.Language = utils.GetLanguage(ctx)
			c.Logo = config.CurrentConfig.D.Uadmin.Logo
			c.FavIcon = config.CurrentConfig.D.Uadmin.FavIcon
			c.Languages = utils.GetActiveLanguages()
			utils.RenderHTML(ctx, config.CurrentConfig.TemplatesFS, config.CurrentConfig.GetPathToTemplate("login"), c)
		} else {
			type Context struct {
				User     string
				Demo     bool
				Menu     []menumodel.DashboardMenu
				SiteName string
				Language *langmodel.Language
				RootURL  string
				Logo     string
				FavIcon  string
			}

			c := Context{}

			c.RootURL = config.CurrentConfig.D.Uadmin.RootURL
			c.Language = utils.GetLanguage(ctx)
			c.Logo = config.CurrentConfig.D.Uadmin.Logo
			c.FavIcon = config.CurrentConfig.D.Uadmin.FavIcon
			c.SiteName = config.CurrentConfig.D.Uadmin.SiteName
			c.User = defaultAdapter.GetUserFromRequest(ctx).Username
			c.Logo = config.CurrentConfig.D.Uadmin.Logo
			c.FavIcon = config.CurrentConfig.D.Uadmin.FavIcon

			c.Menu = make([]menumodel.DashboardMenu, 0)
			// @todo, fix
			// session.User.GetDashboardMenu()
			//for i := range c.Menu {
			//	c.Menu[i].MenuName = translation.Translate(c.Menu[i].MenuName, c.Language.Code, true)
			//}
			utils.RenderHTML(ctx, config.CurrentConfig.TemplatesFS, config.CurrentConfig.GetPathToTemplate("home"), c)
		}
	})
}

func (b Blueprint) Init(config *config.UadminConfig) {
	b.AuthAdapterRegistry.RegisterNewAdapter(&interfaces3.DirectAuthProvider{})
	b.AuthAdapterRegistry.RegisterNewAdapter(&interfaces3.TokenAuthProvider{})
	b.AuthAdapterRegistry.RegisterNewAdapter(&interfaces3.DirectAuthForAdminProvider{})
}

var ConcreteBlueprint = Blueprint{
	Blueprint: interfaces.Blueprint{
		Name:              "auth",
		Description:       "blueprint for auth functionality",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
	AuthAdapterRegistry: interfaces3.NewAuthProviderRegistry(),
}
