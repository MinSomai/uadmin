package sessions

import (
	"github.com/gin-gonic/gin"
	interfaces2 "github.com/sergeyglazyrindev/uadmin/blueprint/sessions/interfaces"
	"github.com/sergeyglazyrindev/uadmin/blueprint/sessions/migrations"
	"github.com/sergeyglazyrindev/uadmin/core"
	"github.com/sergeyglazyrindev/uadmin/utils"
	"strings"
)

type Blueprint struct {
	core.Blueprint
	SessionAdapterRegistry *interfaces2.SessionProviderRegistry
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	// function to verify CSRF
	mainRouter.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			if !core.CurrentConfig.RequiresCsrfCheck(c) {
				c.Next()
				return
			}
			if c.Request.Method != "POST" {
				c.Next()
				return
			}
			contentType := c.Request.Header.Get("Content-Type")
			if contentType == "application/json" {
				c.Next()
				return
			}
			var serverKey string
			var csrfTokenFromRequest string
			csrfTokenFromRequest = c.Request.Header.Get("X-CSRF-TOKEN")
			if csrfTokenFromRequest == "" {
				csrfTokenFromRequest, _ = c.Cookie("csrf_token")
				if csrfTokenFromRequest == "" {
					csrfTokenFromRequest = c.PostForm("csrf-token")
				}
			}
			serverKey = c.Request.Header.Get("X-" + strings.ToUpper(core.CurrentConfig.D.Uadmin.APICookieName))
			if serverKey == "" {
				if c.Query("for-uadmin-panel") == "1" {
					serverKey, _ = c.Cookie(core.CurrentConfig.D.Uadmin.AdminCookieName)
				} else {
					serverKey, _ = c.Cookie(core.CurrentConfig.D.Uadmin.APICookieName)
				}
			}
			defaultSessionAdapter, _ := b.SessionAdapterRegistry.GetDefaultAdapter()
			session, _ := defaultSessionAdapter.GetByKey(serverKey)
			if session == nil {
				c.String(400, "No user session found")
				c.Abort()
				return
			}
			// @todo, comment it out when stabilize token
			//csrfToken, err := session.Get("csrf_token")
			//if err != nil {
			//	c.String(400, err.Error())
			//	c.Abort()
			//	return
			//}

			if len(csrfTokenFromRequest) != 64 {
				c.String(400, "Incorrect length of csrf-token")
				c.Abort()
				return
			}
			// @todo, comment it out when stabilize token
			//tokenUnmasked := utils.UnmaskCSRFToken(csrfTokenFromRequest)
			//if tokenUnmasked != csrfToken {
			//	c.String(400, "Incorrect csrf-token")
			//	c.Abort()
			//	return
			//}
			c.Next()
		}
	}())
	mainRouter.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			if !strings.HasPrefix(c.Request.URL.Path, core.CurrentConfig.D.Uadmin.RootAdminURL) {
				c.Next()
				return
			}
			contentType := c.Request.Header.Get("Content-Type")
			if contentType == "application/json" {
				c.Next()
				return
			}
			serverKey := c.Request.Header.Get("X-" + strings.ToUpper(core.CurrentConfig.D.Uadmin.APICookieName))
			if serverKey == "" {
				if c.Query("for-uadmin-panel") == "1" {
					serverKey, _ = c.Cookie(core.CurrentConfig.D.Uadmin.AdminCookieName)
				} else {
					serverKey, _ = c.Cookie(core.CurrentConfig.D.Uadmin.APICookieName)
				}
			}
			defaultSessionAdapter, _ := b.SessionAdapterRegistry.GetDefaultAdapter()
			session, _ := defaultSessionAdapter.GetByKey(serverKey)
			if session.IsExpired() && c.Request.URL.Path != core.CurrentConfig.D.Uadmin.RootAdminURL {
				c.Redirect(302, core.CurrentConfig.D.Uadmin.RootAdminURL)
				return
			}
			if session.GetUser() != nil && !session.GetUser().IsStaff && !session.GetUser().IsSuperUser {
				c.Redirect(302, core.CurrentConfig.D.Uadmin.RootAdminURL)
				return
			}
			c.Next()
		}
	}())
	core.FuncMap["CSRF"] = func(Key string) string {
		sessionAdapter, _ := ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		session, _ := sessionAdapter.GetByKey(Key)
		csrfToken, _ := session.Get("csrf_token")
		return utils.MaskCSRFToken(csrfToken)
	}
}

func (b Blueprint) Init() {
	b.SessionAdapterRegistry.RegisterNewAdapter(&interfaces2.DbSession{}, true)
	core.ProjectModels.RegisterModel(func() interface{} { return &core.Session{} })
}

var ConcreteBlueprint = Blueprint{
	Blueprint: core.Blueprint{
		Name:              "sessions",
		Description:       "Sessions blueprint responsible to keep session data in database",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
	SessionAdapterRegistry: interfaces2.NewSessionRegistry(),
}
