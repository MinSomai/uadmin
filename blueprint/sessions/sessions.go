package sessions

import (
	"github.com/gin-gonic/gin"
	interfaces2 "github.com/uadmin/uadmin/blueprint/sessions/interfaces"
	"github.com/uadmin/uadmin/blueprint/sessions/migrations"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/utils"
	"strings"
)

type Blueprint struct {
	interfaces.Blueprint
	SessionAdapterRegistry *interfaces2.SessionProviderRegistry
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	mainRouter.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			if !config.CurrentConfig.RequiresCsrfCheck(c) {
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
			serverKey = c.Request.Header.Get("X-" + strings.ToUpper(config.CurrentConfig.D.Uadmin.ApiCookieName))
			if serverKey == "" {
				if c.Query("for-uadmin-panel") == "1" {
					serverKey, _ = c.Cookie(config.CurrentConfig.D.Uadmin.AdminCookieName)
				} else {
					serverKey, _ = c.Cookie(config.CurrentConfig.D.Uadmin.ApiCookieName)
				}
			}
			defaultSessionAdapter, _ := b.SessionAdapterRegistry.GetDefaultAdapter()
			session, _ := defaultSessionAdapter.GetByKey(serverKey)
			if session == nil {
				c.String(400, "No user session found")
				c.Abort()
				return
			}
			csrfToken, err := session.Get("csrf_token")
			if err != nil {
				c.String(400, err.Error())
				c.Abort()
				return
			}

			if len(csrfTokenFromRequest) != 64 {
				c.String(400, "Incorrect length of csrf-token")
				c.Abort()
				return
			}
			tokenUnmasked := utils.UnmaskCSRFToken(csrfTokenFromRequest)
			if tokenUnmasked != csrfToken {
				c.String(400, "Incorrect csrf-token")
				c.Abort()
				return
			}
			c.Next()
		}
	}())
	utils.FuncMap["CSRF"] = func(Key string) string {
		sessionAdapter, _ := ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		session, _ := sessionAdapter.GetByKey(Key)
		csrfToken, _ := session.Get("csrf_token")
		return utils.MaskCSRFToken(csrfToken)
	}
}

func (b Blueprint) Init(config *config.UadminConfig) {
	b.SessionAdapterRegistry.RegisterNewAdapter(&interfaces2.DbSession{}, true)
}

var ConcreteBlueprint = Blueprint{
	Blueprint: interfaces.Blueprint{
		Name:              "sessions",
		Description:       "Sessions blueprint responsible to keep session data in database",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
	SessionAdapterRegistry: interfaces2.NewSessionRegistry(),
}
