package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/uadmin/core"
)

type TokenAuthProvider struct {
}

func (ap *TokenAuthProvider) GetUserFromRequest(c *gin.Context) core.IUser {
	return nil
}

// swagger:route GET /pets token listPets
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
func (ap *TokenAuthProvider) Signin(c *gin.Context) {

}

func (ap *TokenAuthProvider) Signup(c *gin.Context) {
}

func (ap *TokenAuthProvider) Logout(c *gin.Context) {
}

func (ap *TokenAuthProvider) IsAuthenticated(c *gin.Context) {
}

func (ap *TokenAuthProvider) GetSession(c *gin.Context) core.ISessionProvider {
	return nil
}

func (ap *TokenAuthProvider) GetName() string {
	return "token"
}
