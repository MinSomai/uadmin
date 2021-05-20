package uadmin

import (
	// userblueprintapi "github.com/uadmin/uadmin/blueprint/user/api"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/database"
	"strconv"
	"time"

	"github.com/uadmin/uadmin/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	Config   *config.UadminConfig
	Database *database.Database
	Router   *gin.Engine
}

var instance *App

func NewApp(environment string) *App {
	if instance == nil {
		a := new(App)
		a.Config = config.NewConfig("configs/" + environment + ".yaml")
		a.Database = database.NewDatabase(a.Config)
		a.Router = gin.Default()
		a.Router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"https://foo.com"},
			AllowMethods:     []string{"PUT", "PATCH"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return origin == "https://github.com"
			},
			MaxAge: 12 * time.Hour,
		}))
		a.InitializeRouter()
		instance = a
		return a
	}
	return instance
}

func (a App) Initialize() {

}

func (a App) StartAdmin() {
	a.Initialize()
	// useradmin.RegisterAdminPart()
	http.StartServer(a.Config)
}

func (a App) StartApi() {
	a.Initialize()
	_ = a.Router.Run(":" + strconv.Itoa(a.Config.D.Api.ListenPort))
}

func (a App) InitializeRouter() {
	// userblueprintapi.InitializeRouter(a.Router)
}

func (a App) BaseApiUrl() string {
	return ":" + strconv.Itoa(a.Config.D.Api.ListenPort)
}
