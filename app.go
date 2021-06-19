package uadmin

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/uadmin/uadmin/dialect"
	"github.com/uadmin/uadmin/interfaces"
	"io/fs"
	nethttp "net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	abtestblueprint "github.com/uadmin/uadmin/blueprint/abtest"
	approvalblueprint "github.com/uadmin/uadmin/blueprint/approval"
	authblueprint "github.com/uadmin/uadmin/blueprint/auth"
	languageblueprint "github.com/uadmin/uadmin/blueprint/language"
	logblueprint "github.com/uadmin/uadmin/blueprint/logging"
	menublueprint "github.com/uadmin/uadmin/blueprint/menu"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	settingsblueprint "github.com/uadmin/uadmin/blueprint/settings"
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
		a := App{}
		a.Config = config.NewConfig("configs/" + environment + ".yaml")
		a.Config.TemplatesFS = templatesRoot
		a.CommandRegistry = &CommandRegistry{
			Actions: make(map[string]interfaces.ICommand),
		}
		dialect.CurrentDatabaseSettings = &dialect.DatabaseSettings{
			Default: a.Config.D.Db.Default,
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
		a.Initialize()
		a.InitializeRouter()
		appInstance = &a
		return &a
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
	a.BlueprintRegistry.Initialize(a.Config)
}

func (a App) RegisterBaseBlueprints() {
	a.BlueprintRegistry.Register(userblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(menublueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(sessionsblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(settingsblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(logblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(languageblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(approvalblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(abtestblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(authblueprint.ConcreteBlueprint)
}

func (a App) RegisterBlueprint(blueprint interfaces.IBlueprint) {
	a.BlueprintRegistry.Register(blueprint)
}

func (a App) RegisterCommand(name string, command interfaces.ICommand) {
	a.CommandRegistry.addAction(name, command)
}

func (a App) RegisterBaseCommands() {
	a.RegisterCommand("migrate", &MigrateCommand{})
	a.RegisterCommand("blueprint", &BlueprintCommand{})
	a.RegisterCommand("swagger", &SwaggerCommand{})
	a.RegisterCommand("openapi", &OpenApiCommand{})
	a.RegisterCommand("superuser", &SuperadminCommand{})
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

//go:embed templates
var templatesRoot embed.FS

//go:embed static/*
var staticRoot embed.FS

// myFS implements fs.FS
type uadminStaticFS struct {
	content embed.FS
}

func (c uadminStaticFS) Open(name string) (fs.File, error) {
	return c.content.Open(path.Join("static", name))
}

func (a App) InitializeRouter() {
	// http.FS can be used to create a http Filesystem
	staticFiles := uadminStaticFS{staticRoot}
	fs1 := nethttp.FS(staticFiles)
	// Serve static files
	a.Router.StaticFS("/static-inbuilt/", fs1)
	a.BlueprintRegistry.InitializeRouting(a.Router)
}

func (a App) BaseApiUrl() string {
	return ":" + strconv.Itoa(a.Config.D.Api.ListenPort)
}
