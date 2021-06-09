package dialect

import (
	"encoding/json"
	"fmt"
	config2 "github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/utils"
	"gorm.io/gorm"
	"io/ioutil"
)

var Db *gorm.DB

type DatabaseSettings struct {
	Default *config2.DBSettings
}
var CurrentDatabaseSettings *DatabaseSettings

// GetDB returns a pointer to the DB
func GetDB(alias_ ...string) *gorm.DB {
	var alias string
	if len(alias_) == 0 {
		alias = "default"
	} else {
		alias = alias_[0]
	}
	if Db != nil{
		return Db
	}
	var err error

	// Check if there is a database config file
	if CurrentDatabaseSettings == nil {
		buf, err := ioutil.ReadFile(".database")
		if err == nil {
			err = json.Unmarshal(buf, CurrentDatabaseSettings)
			if err != nil {
				utils.Trail(utils.WARNING, ".database file is not a valid json file. %s", err)
			}
		}
	}
	dialect := GetDialectForDb(alias)
	Db, err = dialect.GetDb(
		alias,
	)
	if err != nil {
		utils.Trail(utils.ERROR, "unable to connect to DB. %s", err)
		Db.Error = fmt.Errorf("unable to connect to DB. %s", err)
	}
	return Db
}



func GetDialectForDb(alias string) DbDialect {
	var databaseConfig *config2.DBSettings
	if alias == "default" {
		databaseConfig = CurrentDatabaseSettings.Default
	}
	return NewDbDialect(Db, databaseConfig.Type)
}
