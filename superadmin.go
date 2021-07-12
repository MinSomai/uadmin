package uadmin

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/jessevdk/go-flags"
	"github.com/miquella/ask"
	utils2 "github.com/uadmin/uadmin/blueprint/auth/utils"
	userblueprint "github.com/uadmin/uadmin/blueprint/user"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/utils"
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
	return commandRegistry.runAction(subaction, "", args)
}

func (c SuperadminCommand) GetHelpText() string {
	return "Manage your superusers"
}

type SuperadminCommandOptions struct {
	Username string `short:"n" required:"true" description:"Username" valid:"username-uadmin,username-unique"`
	Email string `short:"e" required:"true" description:"Email'" valid:"email,email-unique"`
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
	_, err = govalidator.ValidateStruct(opts)
	if err != nil {
		return err
	}
	db := interfaces.GetDB()
	var superuserGroup usermodels.UserGroup
	db.Model(&usermodels.UserGroup{}).Where(&usermodels.UserGroup{GroupName: "Superusers"}).First(&superuserGroup)
	if superuserGroup.ID == 0 {
		superuserGroup = usermodels.UserGroup{
			GroupName: "Superusers",
		}
		db.Create(&superuserGroup)
	}
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
			Password: password,
			ConfirmedPassword: confirmpassword,
		}
		_, err = govalidator.ValidateStruct(passwordValidationStruct)
		if err != nil {
			return err
		}
		break
	}
	salt := utils.RandStringRunes(appInstance.Config.D.Auth.SaltLength)
	// hashedPassword, err := utils2.HashPass(password, salt)
	hashedPassword, err := utils2.HashPass(password, salt)
	if err != nil {
		return err
	}
	admin := usermodels.User{
		FirstName:    opts.FirstName,
		LastName:     opts.LastName,
		Username:     opts.Username,
		Email: opts.Email,
		Password:     hashedPassword,
		Admin:        true,
		RemoteAccess: true,
		Active:       true,
		UserGroup:    superuserGroup,
		Salt: salt,
		IsSuperUser: true,
	}
	db.Create(&admin)
	return nil
}

func (command CreateSuperadmin) GetHelpText() string {
	return "Create superadmin"
}
