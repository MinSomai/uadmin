package uadmin

import (
	"fmt"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/interfaces"
	"io"
	"log"
	"os"
	"os/exec"
)

type SwaggerCommand struct {
}
func (c SwaggerCommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &CommandRegistry{
		Actions: make(map[string]interfaces.ICommand),
	}

	commandRegistry.addAction("serve", &ServeSwaggerServer{})
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

func (c SwaggerCommand) GetHelpText() string {
	return "Manage your swagger integration"
}

type ServeSwaggerServerOptions struct {
}

type ServeSwaggerServer struct {
}

func (command ServeSwaggerServer) Proceed(subaction string, args []string) error {
	appInstance.Config.ApiSpec = config.NewSwaggerSpec(appInstance.Config.D.Swagger.PathToSpec)
	commandToExecute := exec.Command(
		"swagger", "serve", "--flavor=swagger", "--no-open",
		fmt.Sprintf("--port=%d", appInstance.Config.D.Swagger.ListenPort), appInstance.Config.D.Swagger.PathToSpec,
	)
	stderr, err := commandToExecute.StderrPipe()
	if err != nil {
		log.Fatal(err)
		return err
	}

	if err := commandToExecute.Start(); err != nil {
		log.Fatal(err)
		return err
	}

	slurp, _ := io.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := commandToExecute.Wait(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (command ServeSwaggerServer) GetHelpText() string {
	return "Serve your swagger api spec"
}
