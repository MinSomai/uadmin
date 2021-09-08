package uadmin

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/jessevdk/go-flags"
	"github.com/miquella/ask"
	utils2 "github.com/sergeyglazyrindev/uadmin/blueprint/auth/utils"
	userblueprint "github.com/sergeyglazyrindev/uadmin/blueprint/user"
	"github.com/sergeyglazyrindev/uadmin/core"
	"github.com/sergeyglazyrindev/uadmin/utils"
	"os"
)

type SuperadminCommand struct {
}

func (c SuperadminCommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &CommandRegistry{
		Actions: make(map[string]core.ICommand),
	}
	commandRegistry.AddAction("create", &CreateSuperadmin{})
	if len(os.Args) > 2 {
		action = os.Args[2]
		isCorrectActionPassed = commandRegistry.IsRegisteredCommand(action)
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
	return commandRegistry.RunAction(subaction, "", args)
}

func (c SuperadminCommand) GetHelpText() string {
	return "Manage your superusers"
}

type SuperadminCommandOptions struct {
	Username  string `short:"n" required:"true" description:"Username" valid:"username-uadmin,username-unique"`
	Email     string `short:"e" required:"true" description:"Email'" valid:"email,email-unique"`
	FirstName string `short:"f" required:"false" description:"First name'"`
	LastName  string `short:"l" required:"false" description:"Last name'"`
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
	_, err = govalidator.ValidateStruct(opts)
	if err != nil {
		return err
	}
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	db := uadminDatabase.Db
	if opts.FirstName == "" {
		opts.FirstName = "System"
	}
	if opts.LastName == "" {
		opts.LastName = "Admin"
	}
	err = ask.Print("Warning! I am about to ask you for a password!\n")
	if err != nil {
		return err
	}
	var password string
	for true {
		password, err = ask.HiddenAsk("Password: ")
		if err != nil {
			return err
		}
		confirmpassword, err := ask.HiddenAsk("Confirm Password: ")
		if err != nil {
			return err
		}
		passwordValidationStruct := &userblueprint.PasswordValidationStruct{
			Password:          password,
			ConfirmedPassword: confirmpassword,
		}
		_, err = govalidator.ValidateStruct(passwordValidationStruct)
		if err != nil {
			core.Trail(core.ERROR, fmt.Errorf("please try to to repeat password again"))
			continue
		}
		break
	}
	salt := utils.RandStringRunes(appInstance.Config.D.Auth.SaltLength)
	// hashedPassword, err := utils2.HashPass(password, salt)
	hashedPassword, err := utils2.HashPass(password, salt)
	if err != nil {
		return err
	}
	admin := core.User{
		FirstName:        opts.FirstName,
		LastName:         opts.LastName,
		Username:         opts.Username,
		Email:            opts.Email,
		Password:         hashedPassword,
		Active:           true,
		IsSuperUser:      true,
		Salt:             salt,
		IsPasswordUsable: true,
	}
	db.Create(&admin)
	core.Trail(core.INFO, "Superuser created successfully")
	return nil
}

func (command CreateSuperadmin) GetHelpText() string {
	return "Create superadmin"
}
