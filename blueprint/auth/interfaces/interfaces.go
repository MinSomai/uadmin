package interfaces

import (
	"fmt"
	"github.com/gin-gonic/gin"
	sessioninterfaces "github.com/uadmin/uadmin/blueprint/sessions/interfaces"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
)

type IAuthProvider interface {
	GetUserFromRequest (c *gin.Context) *usermodels.User
	Signin(c *gin.Context)
	Logout(c *gin.Context)
	IsAuthenticated(c *gin.Context)
	GetSession(c *gin.Context) *sessioninterfaces.ISessionProvider
	GetName() string
	Signup(c *gin.Context)
}

type AuthProviderRegistry struct {
	registeredAdapters map[string]IAuthProvider
}

func (r *AuthProviderRegistry) RegisterNewAdapter(adapter IAuthProvider) {
	r.registeredAdapters[adapter.GetName()] = adapter
}

func (r *AuthProviderRegistry) GetAdapter(name string) (IAuthProvider, error) {
	adapter, ok := r.registeredAdapters[name]
	if ok {
		return adapter, nil
	} else {
		return nil, fmt.Errorf("adapter with name %s not found", name)
	}
}

func (r *AuthProviderRegistry) Iterate() <-chan IAuthProvider {
	chnl := make(chan IAuthProvider)
	go func() {
		for _, authProvider := range r.registeredAdapters {
			chnl <- authProvider
		}
		// Ensure that at the end of the loop we close the channel!
		close(chnl)
	}()
	return chnl
}

func NewAuthProviderRegistry() *AuthProviderRegistry {
	return &AuthProviderRegistry{
		registeredAdapters: make(map[string]IAuthProvider),
	}
}
