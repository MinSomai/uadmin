package uadmin

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/uadmin/uadmin/interfaces"
	"os"
)

type SuperadminCommand struct {
}
func (c SuperadminCommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &CommandRegistry{
		Actions: make(map[string]interfaces.ICommand),
	}
	commandRegistry.addAction("create", &CreateSuperadmin{})
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

func (c SuperadminCommand) GetHelpText() string {
	return "Manage your superusers"
}

type SuperadminCommandOptions struct {
	Username string `short:"n" required:"true" description:"Username"`
	Email string `short:"e" required:"true" description:"Email'"`
	FirstName string `short:"f" required:"false" description:"First name'"`
	LastName string `short:"l" required:"false" description:"Last name'"`
}

type CreateSuperadmin struct {
}

func (command CreateSuperadmin) Proceed(subaction string, args []string) error {
	var opts = &SuperadminCommandOptions{}
	parser := flags.NewParser(opts, flags.Default)
	var err error
	_, err = parser.ParseArgs(args)
	if len(args) == 0 {
		var help string = `
Please provide flags -n and -e which are username and email of the user respectively 
`
		fmt.Printf(help)
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (command CreateSuperadmin) GetHelpText() string {
	return "Create superadmin"
}
