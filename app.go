package uadmin

import (
	"fmt"
	"github.com/uadmin/uadmin/interfaces"
	"os"

	// userblueprintapi "github.com/uadmin/uadmin/blueprint/user/api"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/database"
	"github.com/uadmin/uadmin/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	Config   *config.UadminConfig
	Database *database.Database
	Router   *gin.Engine
	commandRegistry *CommandRegistry
}

var instance *App

func NewApp(environment string) *App {
	if instance == nil {
		a := new(App)
		a.Config = config.NewConfig("configs/" + environment + ".yaml")
		a.commandRegistry = &CommandRegistry{
			actions: make(map[string]interfaces.CommandInterface),
		}
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
		a.baseInitialization()
		a.registerBaseCommands()
		// a.InitializeRouter()
		instance = a
		return a
	}
	return instance
}

func (a App) Initialize() {

}

func (a App) baseInitialization() {

}

func (a App) registerBaseCommands() {
	createCommand := new(MigrateCommand)
	a.commandRegistry.addAction("migrate", interfaces.CommandInterface(createCommand))
}

func (a App) ExecuteCommand() {
	var action string
	var isCorrectActionPassed bool = false
	var help string
	if len(os.Args) > 1 {
		action = os.Args[1]
		isCorrectActionPassed = a.commandRegistry.isRegisteredCommand(action)
	}
	if !isCorrectActionPassed {
		helpText := a.commandRegistry.MakeHelpText()
		help = fmt.Sprintf(`
Please provide what do you want to do ?
%s
`, helpText)
		fmt.Print(help)
		return
	}
	a.commandRegistry.runAction(action)
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