package interfaces

import (
	"github.com/gin-gonic/gin"
	sessioninterfaces "github.com/uadmin/uadmin/blueprint/sessions/interfaces"
	"github.com/uadmin/uadmin/core"
)

type TokenAuthProvider struct {
}

func (ap *TokenAuthProvider) GetUserFromRequest(c *gin.Context) *core.User {
	return nil
}

func (ap *TokenAuthProvider) Signin(c *gin.Context) {
}

func (ap *TokenAuthProvider) Signup(c *gin.Context) {
}

func (ap *TokenAuthProvider) Logout(c *gin.Context) {
}

func (ap *TokenAuthProvider) IsAuthenticated(c *gin.Context) {
}

func (ap *TokenAuthProvider) GetSession(c *gin.Context) sessioninterfaces.ISessionProvider {
	return nil
}

func (ap *TokenAuthProvider) GetName() string {
	return "token"
}
