package tests

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/database"
	"github.com/uadmin/uadmin/dialect"
	"github.com/uadmin/uadmin/interfaces"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func NewTestApp() (*uadmin.App, *gorm.DB) {
	a := new(uadmin.App)
	a.Config = config.NewConfig("configs/" + "test" + ".yaml")
	a.CommandRegistry = &uadmin.CommandRegistry{
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
	a.RegisterBaseCommands()
	a.InitializeRouter()
	dialect.CurrentDatabaseSettings = &dialect.DatabaseSettings{
		Default: a.Config.D.Db.Default,
	}
	dialectdb := dialect.NewDbDialect(a.Database.ConnectTo("default"), a.Config.D.Db.Default.Type)
	db, err := dialectdb.GetDb(
		"default",
	)
	if err != nil {
		panic(fmt.Errorf("Couldn't initialize db %s", err))
	}
	uadmin.StoreCurrentApp(a)
	return a, db
}

// Helper function to process a request and test its response
func TestHTTPResponse(t *testing.T, app *uadmin.App, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	app.Router.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}