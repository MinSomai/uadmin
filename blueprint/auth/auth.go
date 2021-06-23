package auth

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	interfaces3 "github.com/uadmin/uadmin/blueprint/auth/interfaces"
	"github.com/uadmin/uadmin/blueprint/auth/migrations"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/templatecontext"
	"github.com/uadmin/uadmin/utils"
	"strings"
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
		if userSession == nil || userSession.GetUser().ID == 0 {
			type Context struct {
				templatecontext.AdminContext
			}
			c := &Context{}
			adminRequestParams := templatecontext.NewAdminRequestParams()
			adminRequestParams.NeedAllLanguages = true
			templatecontext.PopulateTemplateContextForAdminPanel(ctx, c, adminRequestParams)

			tr := utils.NewTemplateRenderer("Admin Login")
			tr.Render(ctx, config.CurrentConfig.TemplatesFS, config.CurrentConfig.GetPathToTemplate("login"), c)
		} else {
			type Context struct {
				templatecontext.AdminContext
				Demo     bool
				Menu     string
			}

			c := &Context{}
			templatecontext.PopulateTemplateContextForAdminPanel(ctx, c, templatecontext.NewAdminRequestParams())
			sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
			var cookieName string
			cookieName = config.CurrentConfig.D.Uadmin.AdminCookieName
			cookie, _ := ctx.Cookie(cookieName)
			session, _ := sessionAdapter.GetByKey(cookie)

			c.SetUser(defaultAdapter.GetUserFromRequest(ctx).Username)
			c.SetUserExists(defaultAdapter.GetUserFromRequest(ctx).ID != 0)
			allMenu := session.GetUser().GetDashboardMenu()
			allMenus := make([]string, len(allMenu))
			for i := range allMenu {
				allMenu[i].MenuName = utils.Translate(ctx, allMenu[i].MenuName, c.Language.Code, true)
				tmpMenu, _ := json.Marshal(allMenu[i])
				allMenus[i] = string(tmpMenu)
			}
			c.Menu = strings.Join(allMenus, ",")
			c.Demo = false
			tr := utils.NewTemplateRenderer("Dashboard")
			tr.Render(ctx, config.CurrentConfig.TemplatesFS, config.CurrentConfig.GetPathToTemplate("home"), c)
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
