package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	interfaces3 "github.com/uadmin/uadmin/blueprint/auth/interfaces"
	"github.com/uadmin/uadmin/blueprint/auth/migrations"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	"github.com/uadmin/uadmin/interfaces"
	"gorm.io/gorm/schema"
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
	admin.CurrentDashboardAdminPanel.ListHandler = func(ctx *gin.Context) {
		defaultAdapter, _ := b.AuthAdapterRegistry.GetAdapter("direct-for-admin")
		userSession := defaultAdapter.GetSession(ctx)
		if userSession == nil || userSession.GetUser().ID == 0 {
			type Context struct {
				interfaces.AdminContext
			}
			c := &Context{}
			adminRequestParams := interfaces.NewAdminRequestParams()
			adminRequestParams.NeedAllLanguages = true
			admin.PopulateTemplateContextForAdminPanel(ctx, c, adminRequestParams)

			tr := interfaces.NewTemplateRenderer("Admin Login")
			tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("login"), c, interfaces.FuncMap)
		} else {
			type Context struct {
				interfaces.AdminContext
				Menu     string
				CurrentPath string
			}

			c := &Context{}
			admin.PopulateTemplateContextForAdminPanel(ctx, c, interfaces.NewAdminRequestParams())
			menu := string(admin.CurrentDashboardAdminPanel.AdminPages.PreparePagesForTemplate(c.UserPermissionRegistry))
			c.Menu = menu
			c.CurrentPath = ctx.Request.URL.Path
			tr := interfaces.NewTemplateRenderer("Dashboard")
			tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("home"), c, interfaces.FuncMap)
		}
	}
	mainRouter.GET(interfaces.CurrentConfig.D.Uadmin.RootAdminURL + "/profile", func(ctx *gin.Context) {
		type Context struct {
			interfaces.AdminContext
			ID           uint
			Status       bool
			IsUpdated    bool
			Notif        string
			ProfilePhoto string
			OTPImage     string
			OTPRequired  bool
			ChangesSaved bool
			DBFields []*schema.Field
			F *interfaces.Form
		}

		c := &Context{}
		admin.PopulateTemplateContextForAdminPanel(ctx, c, interfaces.NewAdminRequestParams())
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = interfaces.CurrentConfig.D.Uadmin.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		session, _ := sessionAdapter.GetByKey(cookie)
		user := session.GetUser()
		form1 := interfaces.NewFormFromModelFromGinContext(c, user, make([]string, 0), []string{"Username", "FirstName", "LastName", "Email", "Photo", "LastLogin", "ExpiresOn", "OTPRequired"}, true, "")
		form1.TemplateName = interfaces.CurrentConfig.GetPathToTemplate("form/profile_form")
		c.F = form1
		if ctx.Request.Method == "POST" {
			requestForm, _ := ctx.MultipartForm()
			formError := form1.ProceedRequest(requestForm, user)
			if formError.IsEmpty() {
				uadminDatabase := interfaces.NewUadminDatabase()
				defer uadminDatabase.Close()
				db := uadminDatabase.Db
				db.Save(user)
				ctx.Redirect(200, ctx.Request.URL.String())
				return
			}
		}
		tr := interfaces.NewTemplateRenderer(fmt.Sprintf("%s's Profile", c.User))
		tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("profile"), c, interfaces.FuncMap)
	})
}

func (b Blueprint) Init() {
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
