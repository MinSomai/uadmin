package tests

import (
	"fmt"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/dialect"
	"gorm.io/gorm"
)

func NewTestApp() (*uadmin.App, *gorm.DB) {
	app := uadmin.NewApp("test")
	dialect.CurrentDatabaseSettings = &dialect.DatabaseSettings{
		Default: app.Config.D.Db.Default,
	}
	dialectdb := dialect.NewDbDialect(app.Database.ConnectTo("default"), app.Config.D.Db.Default.Type)
	db, err := dialectdb.GetDb(
		"default",
	)
	if err != nil {
		panic(fmt.Errorf("Couldn't initialize db %s", err))
	}
	return app, db
}
