package uadmin

import (
	"fmt"
	"github.com/miquella/ask"
	"github.com/uadmin/uadmin/colors"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/debug"
	"github.com/uadmin/uadmin/interfaces"
	"log"
	"os"
	"strings"
)

type AdminCommand struct {
}
func (c AdminCommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &CommandRegistry{
		Actions: make(map[string]interfaces.ICommand),
	}

	commandRegistry.addAction("serve", &ServeAdminServer{})
	if len(os.Args) > 2 {
		action = os.Args[2]
		isCorrectActionPassed = commandRegistry.isRegisteredCommand(action)
	}
	if !isCorrectActionPassed {
		helpText := commandRegistry.MakeHelpText()
		help = fmt.Sprintf(`
Please provide what do you want to do ?
%s
`, helpText)
		fmt.Print(help)
		return nil
	}
	commandRegistry.runAction(subaction, "", args)
	return nil
}

func (c AdminCommand) GetHelpText() string {
	return "Admin functionality for uadmin project"
}

type AdminStartServerOptions struct {
}

type ServeAdminServer struct {
}

const welcomeMessage = "" +
	`         ___       __          _` + "\n" +
	colors.FGBlueB + `  __  __` + colors.FGNormal + `/   | ____/ /___ ___  (_)___` + "\n" +
	colors.FGBlueB + ` / / / /` + colors.FGNormal + ` /| |/ __  / __ '__ \/ / __ \` + "\n" +
	colors.FGBlueB + `/ /_/ /` + colors.FGNormal + ` ___ / /_/ / / / / / / / / / /` + "\n" +
	colors.FGBlueB + `\__,_/` + colors.FGNormal + `_/  |_\__,_/_/ /_/ /_/_/_/ /_/` + "\n"

func (command ServeAdminServer) Proceed(subaction string, args []string) error {
	migrateCommand := MigrateCommand{}
	err := migrateCommand.Proceed("determine-conflicts", []string{})
	if err != nil {
		debug.Trail(debug.CRITICAL, "Found problems with migrations")
		err = ask.Print("Warning! Found problems with migrations.\n")
		if err != nil {
			return err
		}
		var answer string
		for true {
			answer, err = ask.HiddenAsk("Do you want to start server ?")
			if err != nil {
				return err
			}
			if !interfaces.Contains([]string{"yes", "no"}, strings.ToLower(answer)) {
				continue
			}
			break
		}
		if answer == "no" {
			debug.Trail(debug.WARNING, "You decided to solve first migration problems, so see you next time!")
			return nil
		}
	}
	debug.Trail(debug.OK, "Server Started: http://%s:%d", config.CurrentConfig.D.Admin.BindIP, config.CurrentConfig.D.Admin.ListenPort)
	fmt.Println(welcomeMessage)
	log.Println(appInstance.Router.Run(fmt.Sprintf("%s:%d", config.CurrentConfig.D.Admin.BindIP, config.CurrentConfig.D.Admin.ListenPort)))
	return nil
}

func (command ServeAdminServer) GetHelpText() string {
	return "Serve your admin panel"
}

