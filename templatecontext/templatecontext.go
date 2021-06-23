package templatecontext

import (
	"github.com/gin-gonic/gin"
	langmodel "github.com/uadmin/uadmin/blueprint/language/models"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	interfaces2 "github.com/uadmin/uadmin/blueprint/sessions/interfaces"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/utils"
	"time"
)

type IAdminContext interface {
	SetSiteName(siteName string)
	SetRootAdminURL(rootAdminURL string)
	SetSessionKey(sessionKey string)
	SetRootURL(rootURL string)
	SetLanguage(language *langmodel.Language)
	SetLogo(logo string)
	SetFavIcon(favicon string)
	SetLanguages(langs []langmodel.Language)
	SetPageTitle(pageTitle string)
	SetUser(user string)
	SetUserExists(userExists bool)
}

type AdminRequestParams struct {
	CreateSession bool
	GenerateCSRFToken bool
	NeedAllLanguages bool
}

type AdminContext struct {
	Err         string
	PageTitle string
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
	SessionKey string
	RootAdminURL string
	User string
	UserExists bool
}

func (c *AdminContext) SetSiteName(siteName string) {
	c.SiteName = siteName
}

func (c *AdminContext) SetRootAdminURL(rootAdminURL string) {
	c.RootAdminURL = rootAdminURL
}

func (c *AdminContext) SetSessionKey(sessionKey string) {
	c.SessionKey = sessionKey
}

func (c *AdminContext) SetRootURL(rootURL string) {
	c.RootURL = rootURL
}

func (c *AdminContext) SetLanguage(language *langmodel.Language) {
	c.Language = language
}

func (c *AdminContext) SetLogo(logo string) {
	c.Logo = logo
}

func (c *AdminContext) SetFavIcon(favicon string) {
	c.FavIcon = favicon
}

func (c *AdminContext) SetLanguages(langs []langmodel.Language) {
	c.Languages = langs
}

func (c *AdminContext) SetPageTitle(pageTitle string) {
	c.PageTitle = pageTitle
}

func (c *AdminContext) SetUser(user string) {
	c.User = user
}

func (c *AdminContext) SetUserExists(userExists bool) {
	c.UserExists = userExists
}


func PopulateTemplateContextForAdminPanel(ctx *gin.Context, context IAdminContext, adminRequestParams *AdminRequestParams) {
	sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
	var cookieName string
	cookieName = config.CurrentConfig.D.Uadmin.AdminCookieName
	cookie, _ := ctx.Cookie(cookieName)
	var session interfaces2.ISessionProvider
	if cookie != "" {
		session, _ = sessionAdapter.GetByKey(cookie)
	}
	if adminRequestParams.CreateSession {
		session = sessionAdapter.Create()
		expiresOn := time.Now().Add(time.Duration(config.CurrentConfig.D.Uadmin.SessionDuration)*time.Second)
		session.ExpiresOn(&expiresOn)
		if cookie == "" {
			ctx.SetCookie(config.CurrentConfig.D.Uadmin.AdminCookieName, session.GetKey(), int(config.CurrentConfig.D.Uadmin.SessionDuration), "/", ctx.Request.URL.Host, config.CurrentConfig.D.Uadmin.SecureCookie, config.CurrentConfig.D.Uadmin.HttpOnlyCookie)
		}
	}
	if adminRequestParams.GenerateCSRFToken {
		token := utils.GenerateCSRFToken()
		session.Set("csrf_token", token)
		session.Save()
	}
	context.SetSiteName(config.CurrentConfig.D.Uadmin.SiteName)
	context.SetRootAdminURL(config.CurrentConfig.D.Uadmin.RootAdminURL)
	if session != nil {
		context.SetSessionKey(session.GetKey())
	}
	context.SetRootURL(config.CurrentConfig.D.Uadmin.RootAdminURL)
	context.SetLanguage(utils.GetLanguage(ctx))
	context.SetLogo(config.CurrentConfig.D.Uadmin.Logo)
	context.SetFavIcon(config.CurrentConfig.D.Uadmin.FavIcon)
	if adminRequestParams.NeedAllLanguages {
		context.SetLanguages(utils.GetActiveLanguages())
	}
	if session != nil {
		context.SetUser(session.GetUser().Username)
		context.SetUserExists(session.GetUser().ID != 0)
	}
}
