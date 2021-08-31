package admin

import (
	"github.com/gin-gonic/gin"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	interfaces2 "github.com/uadmin/uadmin/blueprint/sessions/interfaces"
	"github.com/uadmin/uadmin/core"
	"github.com/uadmin/uadmin/utils"
	"time"
)

func init() {
	core.PopulateTemplateContextForAdminPanel = func(ctx *gin.Context, context core.IAdminContext, adminRequestParams *core.AdminRequestParams) {
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = core.CurrentConfig.D.Uadmin.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		var session interfaces2.ISessionProvider
		if cookie != "" {
			session, _ = sessionAdapter.GetByKey(cookie)
		}
		if adminRequestParams.CreateSession && session == nil {
			session = sessionAdapter.Create()
			expiresOn := time.Now().Add(time.Duration(core.CurrentConfig.D.Uadmin.SessionDuration) * time.Second)
			session.ExpiresOn(&expiresOn)
			ctx.SetCookie(core.CurrentConfig.D.Uadmin.AdminCookieName, session.GetKey(), int(core.CurrentConfig.D.Uadmin.SessionDuration), "/", ctx.Request.URL.Host, core.CurrentConfig.D.Uadmin.SecureCookie, core.CurrentConfig.D.Uadmin.HttpOnlyCookie)
			session.Save()
		}
		if adminRequestParams.GenerateCSRFToken {
			token := utils.GenerateCSRFToken()
			currentCsrfToken, _ := session.Get("csrf_token")
			if currentCsrfToken == "" {
				session.Set("csrf_token", token)
				session.Save()
			}
		}
		if session == nil {
			session.Save()
		}
		if ctx.Request.Method == "POST" {
			form, _ := ctx.MultipartForm()
			context.SetPostForm(form)
		}
		context.SetCurrentURL(ctx.Request.URL.Path)
		context.SetCurrentQuery(ctx.Request.URL.RawQuery)
		context.SetFullURL(ctx.Request.URL)
		context.SetSiteName(core.CurrentConfig.D.Uadmin.SiteName)
		context.SetRootAdminURL(core.CurrentConfig.D.Uadmin.RootAdminURL)
		if session != nil {
			context.SetSessionKey(session.GetKey())
		}
		context.SetRootURL(core.CurrentConfig.D.Uadmin.RootAdminURL)
		context.SetLanguage(core.GetLanguage(ctx))
		context.SetLogo(core.CurrentConfig.D.Uadmin.Logo)
		context.SetFavIcon(core.CurrentConfig.D.Uadmin.FavIcon)
		if adminRequestParams.NeedAllLanguages {
			context.SetLanguages(core.GetActiveLanguages())
		}
		// context.SetDemo()
		if session != nil {
			user := session.GetUser()
			context.SetUserObject(user)
			context.SetUser(user.Username)
			context.SetUserExists(user.ID != 0)
			if user.ID != 0 {
				context.SetUserPermissionRegistry(user.BuildPermissionRegistry())
			}
		}
		breadcrumbs := core.NewAdminBreadCrumbsRegistry()
		breadcrumbs.AddBreadCrumb(&core.AdminBreadcrumb{Name: "Dashboard", Url: core.CurrentConfig.D.Uadmin.RootAdminURL, Icon: "home"})
		context.SetBreadCrumbs(breadcrumbs)
	}
}
