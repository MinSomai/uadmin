package proofit_example

import (
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin/core"
	"os"
)

var appForTests *uadmin.App

func NewProofItAppForTests() *uadmin.App {
	if appForTests != nil {
		return appForTests
	}
	a := uadmin.NewApp(os.Getenv("TEST_ENVIRONMENT"), true)
	a.Config.InTests = true
	uadmin.StoreCurrentApp(a)
	appForTests = a
	a.Initialize()
	microservice := &ProofItMicroservice{Microservice: core.Microservice{
		Port: 8089, AuthBackend: "auth-expert", Name: "Proof It Microservice",
		Prefix: "ProofItMicroservice", SwaggerPort: 8090,
	}}
	if microservice.AuthBackend != "" {
		authAdapter, _ := a.GetAuthAdapterRegistry().GetAdapter(microservice.AuthBackend)
		adapterGroup := a.Router.Group("/" + authAdapter.GetName())
		adapterGroup.POST("/signin/", authAdapter.Signin)
		adapterGroup.POST("/signup/", authAdapter.Signup)
		adapterGroup.POST("/logout/", authAdapter.Logout)
		adapterGroup.GET("/status/", authAdapter.IsAuthenticated)
	}
	return a
}

