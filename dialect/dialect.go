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

var CurrentDatabaseSettings *config2.DBSettings

// GetDB returns a pointer to the DB
func GetDB() *gorm.DB {
	if Db != nil{
		return Db
	}
	var err error

	// Check if there is a database config file
	if CurrentDatabaseSettings == nil {
		buf, err := ioutil.ReadFile(".database")
		if err == nil {
			err = json.Unmarshal(buf, &CurrentDatabaseSettings)
			if err != nil {
				utils.Trail(utils.WARNING, ".database file is not a valid json file. %s", err)
			}
		}
	}

	if CurrentDatabaseSettings == nil {
		CurrentDatabaseSettings = &config2.DBSettings{
			Type: "sqlite",
		}
	}
	dialect := GetDialectForDb()
	Db, err = dialect.GetDb(
		CurrentDatabaseSettings.Host,
		CurrentDatabaseSettings.User,
		CurrentDatabaseSettings.Password,
		CurrentDatabaseSettings.Name,
		CurrentDatabaseSettings.Port,
	)
	if err != nil {
		utils.Trail(utils.ERROR, "unable to connect to DB. %s", err)
		Db.Error = fmt.Errorf("unable to connect to DB. %s", err)
	}
	return Db
}



func GetDialectForDb() DbDialect {
	return NewDbDialect(Db, CurrentDatabaseSettings.Type)
}
