package http

import (
	"github.com/uadmin/uadmin/config"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/uadmin/uadmin/colors"
)

const welcomeMessage = "" +
	`         ___       __          _` + "\n" +
	colors.FGBlueB + `  __  __` + colors.FGNormal + `/   | ____/ /___ ___  (_)___` + "\n" +
	colors.FGBlueB + ` / / / /` + colors.FGNormal + ` /| |/ __  / __ '__ \/ / __ \` + "\n" +
	colors.FGBlueB + `/ /_/ /` + colors.FGNormal + ` ___ / /_/ / / / / / / / / / /` + "\n" +
	colors.FGBlueB + `\__,_/` + colors.FGNormal + `_/  |_\__,_/_/ /_/ /_/_/_/ /_/` + "\n"

// const w2 = `` +
// 	`        ______      __` + "\n" +
// 	`       /\  _  \    /\ \              __` + "\n" +
// 	colors.FGBlueB + ` __  __` + colors.FGNormal + `\ \ \L\ \   \_\ \    ___ ___ /\_\    ___` + "\n" +
// 	colors.FGBlueB + `/\ \/\ \` + colors.FGNormal + `\ \  __ \  /'_' \ /' __' __'\/\ \ /' _ '\` + "\n" +
// 	colors.FGBlueB + `\ \ \_\ \` + colors.FGNormal + `\ \ \/\ \/\ \L\ \/\ \/\ \/\ \ \ \/\ \/\ \` + "\n" +
// 	colors.FGBlueB + ` \ \____/` + colors.FGNormal + ` \ \_\ \_\ \___,_\ \_\ \_\ \_\ \_\ \_\ \_\` + "\n" +
// 	colors.FGBlueB + `  \/___/ ` + colors.FGNormal + `  \/_/\/_/\/__,_ /\/_/\/_/\/_/\/_/\/_/\/_/` + "\n"

// ServerReady is a variable that is set to true once the server is ready to use
var ServerReady = false

// @todo analyze start server
// StartServer !
func StartServer(config *config.UadminConfig) {
	//InitializeDbSettingsFromConfig(config)
	//if !registered {
	//	Register()
	//}
	//if !settingsSynched {
	//	syncSystemSettings()
	//}
	//if !handlersRegistered {
	//	registerHandlers()
	//}
	//if val := getBindIP(); val != "" {
	//	BindIP = val
	//}
	//if BindIP == "" {
	//	BindIP = "0.0.0.0"
	//}
	//// Synch model translation
	//// Get Global Schema
	//stat := map[string]int{}
	//for _, v := range CustomTranslation {
	//	tempStat := syncCustomTranslation(v)
	//	for k, v := range tempStat {
	//		stat[k] += v
	//	}
	//}
	//for k := range Schema {
	//	tempStat := syncModelTranslation(Schema[k])
	//	for k, v := range tempStat {
	//		stat[k] += v
	//	}
	//}
	//for k, v := range stat {
	//	complete := float64(v) / float64(stat["en"])
	//	if complete != 1 {
	//		utils.Trail(utils.WARNING, "Translation of %s at %.0f%% [%d/%d]", k, complete*100, v, stat["en"])
	//	}
	//}
	//
	//utils.Trail(utils.OK, "Server Started: http://%s:%d", BindIP, config.D.Admin.ListenPort)
	//fmt.Println(welcomeMessage)
	//database.DbOK = true
	//ServerReady = true
	//log.Println(http.ListenAndServe(fmt.Sprintf("%s:%d", BindIP, config.D.Admin.ListenPort), nil))
}

// @todo analyze start server
// StartSecureServer !
func StartSecureServer(certFile, keyFile string, config *config.UadminConfig) {
	//InitializeDbSettingsFromConfig(config)
	//if !registered {
	//	Register()
	//}
	//if !settingsSynched {
	//	syncSystemSettings()
	//}
	//if !handlersRegistered {
	//	registerHandlers()
	//}
	//if val := getBindIP(); val != "" {
	//	BindIP = val
	//}
	//if BindIP == "" {
	//	BindIP = "0.0.0.0"
	//}
	//// Synch model translation
	//// Get Global Schema
	//stat := map[string]int{}
	//for _, v := range CustomTranslation {
	//	tempStat := syncCustomTranslation(v)
	//	for k, v := range tempStat {
	//		stat[k] += v
	//	}
	//}
	//for k := range Schema {
	//	tempStat := syncModelTranslation(Schema[k])
	//	for k, v := range tempStat {
	//		stat[k] += v
	//	}
	//}
	//for k, v := range stat {
	//	complete := float64(v) / float64(stat["en"])
	//	if complete != 1 {
	//		utils.Trail(utils.WARNING, "Translation of %s at %.0f%% [%d/%d]", k, complete*100, v, stat["en"])
	//	}
	//}
	//
	//utils.Trail(utils.OK, "Server Started: https://%s:%d\n", BindIP, config.D.Admin.SSL.ListenPort)
	//fmt.Println(welcomeMessage)
	//database.DbOK = true
	//ServerReady = true
	//log.Println(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", BindIP, config.D.Admin.SSL.ListenPort), certFile, keyFile, nil))
}

func getBindIP() string {
	// Check if there is a bind ip file in the source code
	ex, _ := os.Executable()
	buf, err := ioutil.ReadFile(path.Join(filepath.Dir(ex), ".bindip"))
	if err == nil {
		return strings.Replace(string(buf), "\n", "", -1)
	}
	return ""
}
