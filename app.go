package uadmin

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/uadmin/uadmin/interfaces"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	userblueprint "github.com/uadmin/uadmin/blueprint/user"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/database"
	"github.com/uadmin/uadmin/http"
	"strconv"
)

type App struct {
	Config            *config.UadminConfig
	Database          *database.Database
	Router            *gin.Engine
	CommandRegistry   *CommandRegistry
	BlueprintRegistry interfaces.IBlueprintRegistry
}

var appInstance *App

func NewApp(environment string) *App {
	if appInstance == nil {
		a := new(App)
		a.Config = config.NewConfig("configs/" + environment + ".yaml")
		a.CommandRegistry = &CommandRegistry{
			Actions: make(map[string]interfaces.ICommand),
		}
		a.BlueprintRegistry = interfaces.NewBlueprintRegistry()
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
		a.RegisterBaseBlueprints()
		a.RegisterBaseCommands()
		a.InitializeRouter()
		appInstance = a
		return a
	}
	return appInstance
}

func ClearApp() {
	appInstance = nil
}

func StoreCurrentApp(app *App) {
	appInstance = app
}

func (a App) Initialize() {

}

func (a App) RegisterBaseBlueprints() {
	a.BlueprintRegistry.Register(userblueprint.ConcreteBlueprint)
}

func (a App) RegisterBlueprint(blueprint interfaces.IBlueprint) {
	a.BlueprintRegistry.Register(blueprint)
}

func (a App) RegisterCommand(name string, command interfaces.ICommand) {
	a.CommandRegistry.addAction(name, command)
}

func (a App) RegisterBaseCommands() {
	migrateCommand := new(MigrateCommand)
	a.RegisterCommand("migrate", interfaces.ICommand(migrateCommand))
	blueprintCommand := new(BlueprintCommand)
	a.RegisterCommand("blueprint", interfaces.ICommand(blueprintCommand))
}

func (a App) ExecuteCommand() {
	var action string
	var isCorrectActionPassed bool = false
	var help string
	if len(os.Args) > 1 {
		action = os.Args[1]
		isCorrectActionPassed = a.CommandRegistry.isRegisteredCommand(action)
	}
	if !isCorrectActionPassed {
		helpText := a.CommandRegistry.MakeHelpText()
		help = fmt.Sprintf(`
Please provide what do you want to do ?
%s
`, helpText)
		fmt.Print(help)
		return
	}
	if len(os.Args) > 2 {
		subaction := os.Args[2]
		isCorrectActionPassed = a.CommandRegistry.isRegisteredCommand(action)
		a.CommandRegistry.runAction(action, subaction, os.Args[3:])
	} else {
		a.CommandRegistry.runAction(action,"", make([]string, 0))
	}
}

func (a App) TriggerCommandExecution(action string, subaction string, params []string) {
	a.CommandRegistry.runAction(action, subaction, params)
}

func (a App) StartAdmin() {
	a.Initialize()
	// useradmin.RegisterAdminPart()
	http.StartServer(a.Config)
}

func (a App) StartApi() {
	a.Initialize()
	// _ = a.Router.Run(":" + strconv.Itoa(a.Config.D.Api.ListenPort))
}

func (a App) InitializeRouter() {
	a.BlueprintRegistry.InitializeRouting(a.Router)

	// userblueprintapi.InitializeRouter(a.Router)
}

func (a App) BaseApiUrl() string {
	return ":" + strconv.Itoa(a.Config.D.Api.ListenPort)
}
