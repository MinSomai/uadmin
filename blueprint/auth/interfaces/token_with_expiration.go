package interfaces

import (
	"github.com/gin-gonic/gin"
	sessioninterfaces "github.com/sergeyglazyrindev/uadmin/blueprint/sessions/interfaces"
	"github.com/sergeyglazyrindev/uadmin/core"
)

type TokenWithExpirationAuthProvider struct {
}

func (ap *TokenWithExpirationAuthProvider) GetUserFromRequest(c *gin.Context) *core.User {
	return nil
}

func (ap *TokenWithExpirationAuthProvider) Signin(c *gin.Context) {
}

func (ap *TokenWithExpirationAuthProvider) Signup(c *gin.Context) {
}

func (ap *TokenWithExpirationAuthProvider) Logout(c *gin.Context) {
}

func (ap *TokenWithExpirationAuthProvider) IsAuthenticated(c *gin.Context) {
}

func (ap *TokenWithExpirationAuthProvider) GetSession(c *gin.Context) sessioninterfaces.ISessionProvider {
	return nil
}

func (ap *TokenWithExpirationAuthProvider) GetName() string {
	return "token-with-expiration"
}
