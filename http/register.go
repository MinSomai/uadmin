package http

import (
	// "github.com/jinzhu/inflection"
	//models2 "github.com/uadmin/uadmin/blueprint/abtest/models"
	//"github.com/uadmin/uadmin/helper"
	//model2 "github.com/uadmin/uadmin/model"
	//"github.com/uadmin/uadmin/translation"
	//"github.com/uadmin/uadmin/utils"
	//"io/ioutil"
	//"net/http"
	//"os"
	//"reflect"
	//"strings"
)

// HideInDashboarder used to check if a model should be hidden in
// dashboard
type HideInDashboarder interface {
	HideInDashboard() bool
}

// CustomTranslation is where you can register custom translation files.
// To register a custom translation file, always assign it with it's key
// in the this format "category/name". For example:
//
// 		uadmin.CustomTranslation = append(uadmin.CustomTranslation, "ui/billing")
//
// This will register the file and you will be able to use it if `uadmin.Tf`.
// By default there is only one registed custom translation wich is "uadmin/system".
var CustomTranslation = []string{
	"uadmin/system",
}

// @todo analyze register
//// Register is used to register models to uadmin
//func Register(m ...interface{}) {
//	modelList := []interface{}{}
//
//	if len(models) == 0 {
//
//		// Initialize system models
//		modelList = []interface{}{
//			DashboardMenu{},
//			User{},
//			UserGroup{},
//			Session{},
//			UserPermission{},
//			GroupPermission{},
//			Language{},
//			Log{},
//			Setting{},
//			SettingCategory{},
//			Approval{},
//			models2.ABTest{},
//			models2.ABTestValue{},
//			//Builder{},
//			//BuilderField{},
//		}
//	}
//
//	// System models count
//	SMCount := len(modelList)
//
//	// Now add user defined models
//	modelList = append(modelList,
//		m...,
//	)
//
//	// Inialize the Database
//	initializeDB(modelList...)
//
//	// Setup languages
//	translation.initializeLanguage()
//
//	// Store models in Model global variable
//	// and initialize the dashboard
//	dashboardMenus := []DashboardMenu{}
//	All(&dashboardMenus)
//	var modelExists bool
//	cat := ""
//	for i := range modelList {
//		modelExists = false
//		t := reflect.TypeOf(modelList[i])
//		name := strings.ToLower(t.Name())
//		// if modelList[i] == nil || modelList[i] == "" || modelList[i] == "0" {
//		// 	if name == "user" {
//		// 		models[name] = User{}
//		// 	} else if name == "dashboardmenu" {
//		// 		models[name] = DashboardMenu{}
//		// 	} else if name == "usergroup" {
//		// 		models[name] = UserGroup{}
//		// 	} else if name == "session" {
//		// 		models[name] = Session{}
//		// 	} else if name == "userpermission" {
//		// 		models[name] = UserPermission{}
//		// 	} else if name == "grouppermission" {
//		// 		models[name] = GroupPermission{}
//		// 	} else if name == "language" {
//		// 		models[name] = Language{}
//		// 	} else if name == "log" {
//		// 		models[name] = Log{}
//		// 	} else if name == "log" {
//		// 		models[name] = Log{}
//		// 	}
//		// }
//		// Trail(ERROR, "Register model: %s - %v", name, modelList[i])
//		models[name] = modelList[i]
//
//		// Get Hidden model status
//		hideItem := false
//		if hider, ok := modelList[i].(HideInDashboarder); ok {
//			hideItem = hider.HideInDashboard()
//		}
//
//		// Register Dashboard menu
//		// First check if the model is already in dashboard
//		dashboardIndex := 0
//		for index, val := range dashboardMenus {
//			if name == val.URL {
//				modelExists = true
//				dashboardIndex = index
//				break
//			}
//		}
//
//		// If not in dashboard, then add it
//		if !modelExists {
//			// Check if the model is a system model
//			if i < SMCount {
//				cat = "System"
//			} else {
//				cat = ""
//			}
//			dashboard := DashboardMenu{
//				MenuName: inflection.Plural(strings.Join(helper.SplitCamelCase(t.Name()), " ")),
//				URL:      name,
//				Hidden:   hideItem,
//				Cat:      cat,
//			}
//			model.Save(&dashboard)
//		} else {
//			// If model exists, synchnorize it if changed
//			if hideItem != dashboardMenus[dashboardIndex].Hidden {
//				dashboardMenus[dashboardIndex].Hidden = hideItem
//				model.Save(&dashboardMenus[dashboardIndex])
//			}
//		}
//	}
//	// Check if encrypt key is there or generate it
//	if _, err := os.Stat(".key"); os.IsNotExist(err) {
//		EncryptKey = utils.generateByteArray(32)
//		ioutil.WriteFile(".key", EncryptKey, 0600)
//	} else {
//		EncryptKey, _ = ioutil.ReadFile(".key")
//	}
//
//	// Check if salt is there or generate it
//	users := []User{}
//	if _, err := os.Stat(".salt"); os.IsNotExist(err) {
//		Salt = GenerateBase64(72)
//		ioutil.WriteFile(".salt", []byte(Salt), 0600)
//		if Count(&users, "") != 0 {
//			recoveryPass := GenerateBase64(24)
//			recoverUsername := GenerateBase64(8)
//			for Count(&users, "username = ?", recoverUsername) != 0 {
//				recoverUsername = GenerateBase64(8)
//			}
//			admin := User{
//				FirstName:    "System",
//				LastName:     "Recovery Admin",
//				Username:     recoverUsername,
//				Password:     hashPass(recoveryPass),
//				Admin:        true,
//				RemoteAccess: false,
//				Active:       true,
//			}
//			admin.Save()
//			utils.Trail(utils.WARNING, "Your salt file was missing, and a new one was generated NO USERS CAN LOGIN UNTIL PASSWORDS ARE RESET.")
//			utils.Trail(utils.INFO, "uAdmin generated a recovery user for you. Username:%s Password:%s", admin.Username, recoveryPass)
//		}
//	} else {
//		saltBytes, _ := ioutil.ReadFile(".salt")
//		Salt = string(saltBytes)
//	}
//
//	// Create an admin user if there is no user in the system
//	adminUsername := "admin"
//	adminPassword := "admin"
//	if os.Getenv("UADMIN_USER") != "" {
//		adminUsername = os.Getenv("UADMIN_USER")
//	}
//	if os.Getenv("UADMIN_PASS") != "" {
//		adminPassword = os.Getenv("UADMIN_PASS")
//	}
//	if Count(&users, "") == 0 {
//		usergroup := UserGroup{
//			GroupName: "Superusers",
//		}
//		usergroup.Save()
//		admin := User{
//			FirstName:    "System",
//			LastName:     "Admin",
//			Username:     adminUsername,
//			Password:     hashPass(adminPassword),
//			Admin:        true,
//			RemoteAccess: true,
//			Active:       true,
//			UserGroup:    usergroup,
//		}
//		admin.Save()
//		utils.Trail(utils.INFO, "Auto generated admin user. Username:%s, Password:%s.", adminUsername, adminPassword)
//	}
//
//	// Register admin inlines
//	model2.RegisterInlines(UserGroup{}, map[string]string{
//		"GroupPermission": "UserGroupID",
//	})
//
//	model2.RegisterInlines(User{}, map[string]string{
//		"UserPermission": "UserID",
//	})
//
//	model2.RegisterInlines(models2.ABTest{}, map[string]string{
//		"ABTestValue": "ABTestID",
//	})
//
//	for k, v := range models {
//		Schema[k], _ = model2.GetSchema(v)
//	}
//
//	// Register JS
//	s := Schema["abtest"]
//	s.IncludeFormJS = []string{"/static/uadmin/js/abtest_form.js"}
//	Schema["abtest"] = s
//
//	// Register Limit Choices To
//	s = Schema["abtest"]
//	s.FieldByName("ModelName").LimitChoicesTo = models2.loadModels
//	s.FieldByName("Field").LimitChoicesTo = models2.loadFields
//	Schema["abtest"] = s
//
//	// Load Session data
//	if CacheSessions {
//		loadSessions()
//	}
//
//	// Load Permission data
//	if CachePermissions {
//		loadPermissions()
//	}
//
//	// Check if there are active ABTests
//	models2.abTestCount = Count([]models2.ABTest{}, "active = ?", true)
//
//	// Mark registered as true to prevent auto registeration
//	registered = true
//}
//
//// RegisterInlines is a function to register a model as an inline for another model
//// Parameters:
//// ===========
////   model (struct instance): Is the model that you want to add inlines to.
////   fk (map[interface{}]string): This is a map of the inlines to be added to the model.
////                                The map's key is the name of the model of the inline
////                                and the value of the map is the foreign key field's name.
////  Example:
////  ========
////  type Person struct {
////    uadmin.Model
////    Name string
////  }
////
////  type Card struct {
////    uadmin.Model
////    PersonID uint
////    Person   Person
////  }
////
//// func main() {
////   ...
////   uadmin.RegisterInlines(Person{}, map[string]string{
////     "Card": "PersonID",
////   })
////   ...
//// }
//
//func registerHandlers() {
//	// register static and add parameter
//	if !strings.HasSuffix(RootURL, "/") {
//		RootURL = RootURL + "/"
//	}
//	if !strings.HasPrefix(RootURL, "/") {
//		RootURL = "/" + RootURL
//	}
//
//	// Handleer for uAdmin, static and media
//	http.HandleFunc(RootURL, Handler(mainHandler))
//	http.HandleFunc("/static/", Handler(StaticHandler))
//	http.HandleFunc("/media/", Handler(mediaHandler))
//
//	// api handler
//	http.HandleFunc(RootURL+"api/", Handler(apiHandler))
//	http.HandleFunc(RootURL+"revertHandler/", Handler(revertLogHandler))
//
//	handlersRegistered = true
//}
//
