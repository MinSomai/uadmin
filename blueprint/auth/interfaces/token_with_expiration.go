package interfaces

import (
	"github.com/gin-gonic/gin"
	sessioninterfaces "github.com/uadmin/uadmin/blueprint/sessions/interfaces"
	"github.com/uadmin/uadmin/interfaces"
)

type TokenWithExpirationAuthProvider struct {
}

func (ap *TokenWithExpirationAuthProvider) GetUserFromRequest(c *gin.Context) *interfaces.User {
	return nil
}

func (ap *TokenWithExpirationAuthProvider) Signin(c *gin.Context) {
}

func (ap *TokenWithExpirationAuthProvider) Signup(c *gin.Context) {
}

func (ap *TokenWithExpirationAuthProvider) Logout(c *gin.Context){
}

func (ap *TokenWithExpirationAuthProvider) IsAuthenticated(c *gin.Context) {
}

func (ap *TokenWithExpirationAuthProvider) GetSession(c *gin.Context) sessioninterfaces.ISessionProvider {
	return nil
}

func (ap *TokenWithExpirationAuthProvider) GetName() string {
	return "token-with-expiration"
}
