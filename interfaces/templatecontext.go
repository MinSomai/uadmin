package interfaces

import (
	"net/url"
)

type IForm interface {

}

type IAdminContext interface {
	SetSiteName(siteName string)
	SetCurrentURL(currentURL string)
	SetFullURL(fullURL *url.URL)
	SetRootAdminURL(rootAdminURL string)
	SetSessionKey(sessionKey string)
	SetRootURL(rootURL string)
	SetLanguage(language *Language)
	SetLogo(logo string)
	SetFavIcon(favicon string)
	SetLanguages(langs []Language)
	SetPageTitle(pageTitle string)
	SetUser(user string)
	SetUserExists(userExists bool)
	SetDemo()
	SetError(err string)
	SetErrorExists()
	GetLanguage() *Language
	GetRootURL() string
	SetUserPermissionRegistry(permRegistry *UserPermRegistry)
	SetForm(form IForm)
	SetCurrentQuery(currentQuery string)
}

type AdminContext struct {
	Err         string
	PageTitle string
	ErrExists   bool
	SiteName    string
	Languages   []Language
	RootURL     string
	OTPRequired bool
	Language    *Language
	Username    string
	Password    string
	Logo        string
	FavIcon     string
	SessionKey string
	RootAdminURL string
	User string
	UserExists bool
	Demo bool
	UserPermissionRegistry *UserPermRegistry
	CurrentURL string
	CurrentQuery string
	FullURL *url.URL
	Form IForm
}

func (c *AdminContext) SetSiteName(siteName string) {
	c.SiteName = siteName
}

func (c *AdminContext) SetCurrentURL(currentURL string) {
	c.CurrentURL = currentURL
}

func (c *AdminContext) SetCurrentQuery(currentQuery string) {
	c.CurrentQuery = currentQuery
}

func (c *AdminContext) SetForm(form IForm) {
	c.Form = form
}

func (c *AdminContext) SetFullURL(fullURL *url.URL) {
	c.FullURL = fullURL
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

func (c *AdminContext) SetLanguage(language *Language) {
	c.Language = language
}

func (c *AdminContext) SetLogo(logo string) {
	c.Logo = logo
}

func (c *AdminContext) SetFavIcon(favicon string) {
	c.FavIcon = favicon
}

func (c *AdminContext) GetLanguage() *Language {
	return c.Language
}

func (c *AdminContext) SetLanguages(langs []Language) {
	c.Languages = langs
}

func (c *AdminContext) SetUserPermissionRegistry(permRegistry *UserPermRegistry) {
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

