package templatecontext

import (
	"github.com/gin-gonic/gin"
	langmodel "github.com/uadmin/uadmin/blueprint/language/models"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	interfaces2 "github.com/uadmin/uadmin/blueprint/sessions/interfaces"
	"github.com/uadmin/uadmin/interfaces"
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
	SetDemo()
	SetError(err string)
	SetErrorExists()
	GetLanguage() *langmodel.Language
	GetRootURL() string
	SetUserPermissionRegistry(permRegistry *interfaces.UserPermRegistry)
}

type AdminRequestParams struct {
	CreateSession bool
	GenerateCSRFToken bool
	NeedAllLanguages bool
}

func NewAdminRequestParams() *AdminRequestParams {
	return &AdminRequestParams{
		CreateSession: true,
		GenerateCSRFToken: true,
		NeedAllLanguages: false,
	}
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
	Demo bool
	UserPermissionRegistry *interfaces.UserPermRegistry
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

func (c *AdminContext) GetRootURL() string {
	return c.RootURL
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

func (c *AdminContext) GetLanguage() *langmodel.Language {
	return c.Language
}

func (c *AdminContext) SetLanguages(langs []langmodel.Language) {
	c.Languages = langs
}

func (c *AdminContext) SetUserPermissionRegistry(permRegistry *interfaces.UserPermRegistry) {
	c.UserPermissionRegistry = permRegistry
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

func (c *AdminContext) SetDemo() {
	c.Demo = true
}

func (c *AdminContext) SetError(err string) {
	c.Err = err
}

func (c *AdminContext) SetErrorExists() {
	c.ErrExists = true
}

func PopulateTemplateContextForAdminPanel(ctx *gin.Context, context IAdminContext, adminRequestParams *AdminRequestParams) {
	sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
	var cookieName string
	cookieName = interfaces.CurrentConfig.D.Uadmin.AdminCookieName
	cookie, _ := ctx.Cookie(cookieName)
	var session interfaces2.ISessionProvider
	if cookie != "" {
		session, _ = sessionAdapter.GetByKey(cookie)
	}
	if adminRequestParams.CreateSession && session == nil {
		session = sessionAdapter.Create()
		expiresOn := time.Now().Add(time.Duration(interfaces.CurrentConfig.D.Uadmin.SessionDuration)*time.Second)
		session.ExpiresOn(&expiresOn)
		ctx.SetCookie(interfaces.CurrentConfig.D.Uadmin.AdminCookieName, session.GetKey(), int(interfaces.CurrentConfig.D.Uadmin.SessionDuration), "/", ctx.Request.URL.Host, interfaces.CurrentConfig.D.Uadmin.SecureCookie, interfaces.CurrentConfig.D.Uadmin.HttpOnlyCookie)
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
	context.SetSiteName(interfaces.CurrentConfig.D.Uadmin.SiteName)
	context.SetRootAdminURL(interfaces.CurrentConfig.D.Uadmin.RootAdminURL)
	if session != nil {
		context.SetSessionKey(session.GetKey())
	}
	context.SetRootURL(interfaces.CurrentConfig.D.Uadmin.RootAdminURL)
	context.SetLanguage(utils.GetLanguage(ctx))
	context.SetLogo(interfaces.CurrentConfig.D.Uadmin.Logo)
	context.SetFavIcon(interfaces.CurrentConfig.D.Uadmin.FavIcon)
	if adminRequestParams.NeedAllLanguages {
		context.SetLanguages(utils.GetActiveLanguages())
	}
	// context.SetDemo()
	if session != nil {
		user := session.GetUser()
		context.SetUser(user.Username)
		context.SetUserExists(user.ID != 0)
		if user.ID != 0 {
			context.SetUserPermissionRegistry(user.BuildPermissionRegistry())
		}
	}
}

//
//func BuildAdminHandlerForBlueprintfunc(PageTitle string) func (ctx *gin.Context) {
//
//}