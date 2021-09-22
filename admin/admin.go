package admin

/*
Package admin initializes PopulateTemplateContextForAdminPanel function in the core.
It helps us to customize admin panel in case if we want, we'd like to keep core package as is, instead, all changes
have to be done outside of core package. It's non trivial task for us because we need to provide high level abstractions
to expand uadmin functionality, etc
*/

import (
	"github.com/gin-gonic/gin"
	sessionsblueprint "github.com/sergeyglazyrindev/uadmin/blueprint/sessions"
	interfaces2 "github.com/sergeyglazyrindev/uadmin/blueprint/sessions/interfaces"
	"github.com/sergeyglazyrindev/uadmin/core"
	"github.com/sergeyglazyrindev/uadmin/utils"
	"time"
)

func init() {
	core.PopulateTemplateContextForAdminPanel = func(ctx *gin.Context, context core.IAdminContext, adminRequestParams *core.AdminRequestParams) {
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = core.CurrentConfig.D.Uadmin.AdminCookieName
		// in admin panel we determine user's session using cookies
		cookie, _ := ctx.Cookie(cookieName)
		var session interfaces2.ISessionProvider
		if cookie != "" {
			session, _ = sessionAdapter.GetByKey(cookie)
		}
		if adminRequestParams.CreateSession && session == nil {
			// create session if no session found in cookies
			session = sessionAdapter.Create()
			expiresOn := time.Now().UTC().Add(time.Duration(core.CurrentConfig.D.Uadmin.SessionDuration) * time.Second)
			session.ExpiresOn(&expiresOn)
			ctx.SetCookie(core.CurrentConfig.D.Uadmin.AdminCookieName, session.GetKey(), int(core.CurrentConfig.D.Uadmin.SessionDuration), "/", ctx.Request.URL.Host, core.CurrentConfig.D.Uadmin.SecureCookie, core.CurrentConfig.D.Uadmin.HTTPOnlyCookie)
			session.Save()
		}
		if adminRequestParams.GenerateCSRFToken {
			// generate csrf token for POST requests
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
		context.SetCtx(ctx)
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
			// determine current user
			user := session.GetUser()
			if user != nil {
				context.SetUserObject(user)
				context.SetUser(user.GetUsername())
				context.SetUserExists(true)
				context.SetUserPermissionRegistry(user.BuildPermissionRegistry())
			}
		}
		// build breadcrumbs for admin panel
		breadcrumbs := core.NewAdminBreadCrumbsRegistry()
		breadcrumbs.AddBreadCrumb(&core.AdminBreadcrumb{Name: "Dashboard", URL: core.CurrentConfig.D.Uadmin.RootAdminURL, Icon: "home"})
		context.SetBreadCrumbs(breadcrumbs)
	}
	core.GenerateUserModel = func() core.IUser {
		return &core.User{}
	}
}
