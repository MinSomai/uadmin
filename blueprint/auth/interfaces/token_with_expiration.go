package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/uadmin/core"
)

type TokenWithExpirationAuthProvider struct {
}

func (ap *TokenWithExpirationAuthProvider) GetUserFromRequest(c *gin.Context) core.IUser {
	return nil
}

// swagger:route GET /pets1 tokenwithexpiration listPets
//
// Lists pets filtered by some parameters.
//
// This will show all available pets by default.
// You can get the pets that are out of stock
//
//     Consumes:
//     - application/json
//     - application/x-protobuf
//
//     Produces:
//     - application/json
//     - application/x-protobuf
//
//     Schemes: http, https, ws, wss
//
//     Deprecated: true
//
//     Security:
//       api_key:
//       oauth: read, write
//
//     Responses:
//       default: genericError
//       200: someResponse
//       422: validationError
func (ap *TokenWithExpirationAuthProvider) Signin(c *gin.Context) {
}

func (ap *TokenWithExpirationAuthProvider) Signup(c *gin.Context) {
}

func (ap *TokenWithExpirationAuthProvider) Logout(c *gin.Context) {
}

func (ap *TokenWithExpirationAuthProvider) IsAuthenticated(c *gin.Context) {
}

func (ap *TokenWithExpirationAuthProvider) GetSession(c *gin.Context) core.ISessionProvider {
	return nil
}

func (ap *TokenWithExpirationAuthProvider) GetName() string {
	return "token-with-expiration"
}
