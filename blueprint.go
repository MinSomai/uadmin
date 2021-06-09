package uadmin

import (
	"bytes"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/uadmin/uadmin/interfaces"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

type BlueprintCommand struct {
}
func (c BlueprintCommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &CommandRegistry{
		Actions: make(map[string]interfaces.ICommand),
	}

	commandRegistry.addAction("create", &CreateBlueprint{})
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

func (c BlueprintCommand) GetHelpText() string {
	return "Manage your blueprints"
}

type CreateBlueprintOptions struct {
	Message string `short:"m" required:"true" description:"Describe what is this migration for"`
	Name string `short:"n" required:"true" description:"Blueprint you'd like to create migration for'"`
}

type CreateBlueprint struct {
}

func (command CreateBlueprint) Proceed(subaction string, args []string) error {
	var opts = &CreateBlueprintOptions{}
	parser := flags.NewParser(opts, flags.Default)
	var err error
	_, err = parser.ParseArgs(args)
	if len(args) == 0 {
		var help string = `
Please provide flags -n and -m which is name of blueprint and description of blueprint respectively 
`
		fmt.Printf(help)
		return nil
	}
	if err != nil {
		return err
	}
	name := re.ReplaceAllLiteralString(opts.Name, "")
	bluePrintPath := "blueprint/" + strings.ToLower(name)
	if _, err := os.Stat(bluePrintPath); os.IsExist(err) {
		panic(fmt.Sprintf("Blueprint %s does exist", name))
	}
	dirPath := "blueprint/" + strings.ToLower(name)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			panic(err)
		}
	}
	const blueprint = `package {{.Name}}

import (
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/{{.Name}}/migrations"
	"github.com/uadmin/uadmin/interfaces"
)

type Blueprint struct {
	interfaces.Blueprint
}

func (b Blueprint) InitRouter(group *gin.RouterGroup) {
}

func (b Blueprint) Init() {
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "{{.Name}}",
		Description:       "{{.Message}}",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
`
	var blueprintTplBuffer bytes.Buffer
	blueprintTpl := template.Must(template.New("blueprintmain").Parse(blueprint))
	tplData := struct{
		Name string
		Message string
	}{
		Name: name,
		Message: strings.ReplaceAll(re.ReplaceAllLiteralString(opts.Message, ""), `"`, `\"`),
	}
	if err = blueprintTpl.Execute(&blueprintTplBuffer, tplData); err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(dirPath + "/" + name + ".go", blueprintTplBuffer.Bytes(), 0755)
	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"Created blueprint with name %s\n",
		opts.Name,
	)
	migrateCommand := MigrateCommand{}
	migrateCommand.Proceed("create", []string{"-b", name, "-m", "initial"})
	return nil
}

func (command CreateBlueprint) GetHelpText() string {
	return "Create blueprint"
}
