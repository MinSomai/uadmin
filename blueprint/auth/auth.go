package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	interfaces3 "github.com/uadmin/uadmin/blueprint/auth/interfaces"
	"github.com/uadmin/uadmin/blueprint/auth/migrations"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	"github.com/uadmin/uadmin/form"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/template"
	"github.com/uadmin/uadmin/templatecontext"
	"github.com/uadmin/uadmin/utils"
	"gorm.io/gorm/schema"
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
	mainRouter.GET(interfaces.CurrentConfig.D.Uadmin.RootAdminURL, func(ctx *gin.Context) {
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

			tr := interfaces.NewTemplateRenderer("Admin Login")
			tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("login"), c, template.FuncMap)
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
			cookieName = interfaces.CurrentConfig.D.Uadmin.AdminCookieName
			cookie, _ := ctx.Cookie(cookieName)
			session, _ := sessionAdapter.GetByKey(cookie)

			allMenu := session.GetUser().GetDashboardMenu()
			allMenus := make([]string, len(allMenu))
			for i := range allMenu {
				allMenu[i].MenuName = utils.Translate(ctx, allMenu[i].MenuName, c.Language.Code, true)
				tmpMenu, _ := json.Marshal(allMenu[i])
				allMenus[i] = string(tmpMenu)
			}
			c.Menu = strings.Join(allMenus, ",")
			c.Demo = false
			tr := interfaces.NewTemplateRenderer("Dashboard")
			tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("home"), c, template.FuncMap)
		}
	})
	mainRouter.GET(interfaces.CurrentConfig.D.Uadmin.RootAdminURL + "/profile", func(ctx *gin.Context) {
		type Context struct {
			templatecontext.AdminContext
			ID           uint
			Status       bool
			IsUpdated    bool
			Notif        string
			ProfilePhoto string
			OTPImage     string
			OTPRequired  bool
			ChangesSaved bool
			DBFields []*schema.Field
			F *form.Form
		}

		c := &Context{}
		templatecontext.PopulateTemplateContextForAdminPanel(ctx, c, templatecontext.NewAdminRequestParams())
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = interfaces.CurrentConfig.D.Uadmin.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		session, _ := sessionAdapter.GetByKey(cookie)
		user := session.GetUser()
		form1 := form.NewFormFromModelFromGinContext(c, user, make([]string, 0), []string{"Username", "FirstName", "LastName", "Email", "Photo", "LastLogin", "ExpiresOn", "OTPRequired"}, true, "")
		form1.TemplateName = interfaces.CurrentConfig.GetPathToTemplate("form/profile_form")
		c.F = form1
		if ctx.Request.Method == "POST" {
			requestForm, _ := ctx.MultipartForm()
			formError := form1.ProceedRequest(requestForm, user)
			if formError.IsEmpty() {
				db := interfaces.GetDB()
				db.Save(user).Preload("UserGroup")
				ctx.Redirect(200, ctx.Request.URL.String())
				return
			}
		}
		tr := interfaces.NewTemplateRenderer(fmt.Sprintf("%s's Profile", c.User))
		tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("profile"), c, template.FuncMap)
	})
}

func (b Blueprint) Init(config *interfaces.UadminConfig) {
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
