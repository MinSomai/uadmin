package uadmin

import (
	"bytes"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/uadmin/uadmin/interfaces"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MigrateCommand struct {

}
func (c MigrateCommand) Proceed() {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &CommandRegistry{
		actions: make(map[string]interfaces.ICommand),
	}
	createCommand := new(CreateMigration)
	createCommand.opts = &CreateMigrationOptions{}

	commandRegistry.addAction("create", interfaces.ICommand(createCommand))
	upCommand := new(UpMigration)
	upCommand.opts = &UpMigrationOptions{}

	commandRegistry.addAction("up", interfaces.ICommand(upCommand))
	downCommand := new(DownMigration)
	downCommand.opts = &DownMigrationOptions{}

	commandRegistry.addAction("down", interfaces.ICommand(downCommand))
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
		return
	}
	commandRegistry.runAction(action)
}

func (c MigrateCommand) ParseArgs() {

}

func (c MigrateCommand) GetHelpText() string {
	return "Migrate your database"
}

var migrationTplPath = "internal/templates/migrations"
var re = regexp.MustCompile("[[:^ascii:]]")

func prepareMigrationName(message string) string {
	now := time.Now()
	sec := now.Unix()
	message = re.ReplaceAllLiteralString(message, "")
	if len(message) > 10 {
		message = message[:10]
	}
	return fmt.Sprintf("%s_%d", message, sec)
}

type CreateMigrationOptions struct {
	Message string `short:"m" required:"true" description:"Describe what is this migration for"`
	Blueprint string `short:"b" required:"true" description:"Blueprint you'd like to create migration for'"`
}

type CreateMigration struct {
	opts *CreateMigrationOptions
}

func (command CreateMigration) ParseArgs() {
	parser := flags.NewParser(command.opts, flags.Default)
	_, err := parser.ParseArgs(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
}

func (command CreateMigration) Proceed() {
	const concreteMigrationTpl = `
package migrations

type {{.MigrationName}} struct {
}

func (m {{.MigrationName}}) GetName() string {
    return "{{.ConcreteMigrationName}}"
}

func (m {{.MigrationName}}) GetId() int64 {
    return {{.ConcreteMigrationId}}
}

func (m {{.MigrationName}}) Up() {
}

func (m {{.MigrationName}}) Down() {
}

func (m {{.MigrationName}}) Deps() []string {
{{if .DependencyId}}    return []string{"{{.DependencyId}}"}{{else}}    return make([]string, 0){{end}}
}
`
	const initializeMigrationRegistryTpl = `
    BMigrationRegistry.addMigration({{.MigrationName}}{})
`
	const migrationRegistryCreationTpl = `
package migrations

import (
	"github.com/uadmin/uadmin/interfaces"
)

type MigrationRegistry struct {
	migrations map[string]interfaces.IMigration
}

func (r MigrationRegistry) addMigration(migration interfaces.IMigration) {
	r.migrations[migration.GetName()] = migration
}

func (r MigrationRegistry) FindMigrations() <-chan interfaces.IMigration{
	chnl := make(chan interfaces.IMigration)
	go func() {
		close(chnl)
	}()
	return chnl
}

var BMigrationRegistry *MigrationRegistry

func init() {
    BMigrationRegistry = &MigrationRegistry{
        migrations: make(map[string]interfaces.IMigration),
    }
    // placeholder to insert next migration
}
`
	bluePrintPath := "blueprint/" + strings.ToLower(command.opts.Blueprint)
	if _, err := os.Stat(bluePrintPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Blueprint %s doesn't exist", command.opts.Blueprint))
	}
	dirPath := "blueprint/" + strings.ToLower(command.opts.Blueprint) + "/migrations"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			panic(err)
		}
	}
	pathToBaseMigrationsFile := dirPath + "/migrations.go"
	if _, err := os.Stat(pathToBaseMigrationsFile); os.IsNotExist(err) {
		err = ioutil.WriteFile(pathToBaseMigrationsFile, []byte(migrationRegistryCreationTpl), 0755)
		if err != nil {
			panic(err)
		}
	}
	migrationName := prepareMigrationName(command.opts.Message)
	pathToConcreteMigrationsFile := dirPath + "/" + migrationName + ".go"
	var lastMigrationId int
	var err error
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		var migrationFileRegex = regexp.MustCompile(`.*?_(\d+)\.go`)
		match := migrationFileRegex.FindStringSubmatch(path)
		if len(match) > 0 {
			migrationId, _ := strconv.Atoi(match[1])
			if migrationId > lastMigrationId {
				lastMigrationId = migrationId
			}
		}
		return nil
	})
	humanizedMessage := strings.ReplaceAll(
		re.ReplaceAllLiteralString(command.opts.Message, ""),
		"\"",
		"",
	)
	var concreteTplBuffer bytes.Buffer
	now := time.Now()
	sec := now.Unix()
	concreteTpl := template.Must(template.New("concretemigration").Parse(concreteMigrationTpl))
	concreteData := struct{
		MigrationName string
		ConcreteMigrationName string
		ConcreteMigrationId string
		DependencyId string
	}{
		MigrationName: migrationName,
		ConcreteMigrationName: humanizedMessage,
		ConcreteMigrationId: strconv.Itoa(int(sec)),
		DependencyId: "",
	}
	if lastMigrationId > 0 {
		concreteData.DependencyId = strconv.Itoa(lastMigrationId)
	}
	if err = concreteTpl.Execute(&concreteTplBuffer, concreteData); err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(pathToConcreteMigrationsFile, concreteTplBuffer.Bytes(), 0755)
	if err != nil {
		panic(err)
	}
	integrateMigrationIntoRegistryTpl := template.Must(template.New("integratemigrationintoregistry").Parse(initializeMigrationRegistryTpl))
	integrateMigrationIntoRegistryData := struct{
		MigrationName string
	}{
		MigrationName: migrationName,
	}
	var integrateMigrationIntoRegistryTplBuffer bytes.Buffer
	if err = integrateMigrationIntoRegistryTpl.Execute(&integrateMigrationIntoRegistryTplBuffer, integrateMigrationIntoRegistryData); err != nil {
		panic(err)
	}
	read, err := ioutil.ReadFile(pathToBaseMigrationsFile)
	if err != nil {
		panic(err)
	}
	newContents := strings.Replace(
		string(read),
		"// placeholder to insert next migration",
		integrateMigrationIntoRegistryTplBuffer.String() + "\n    // placeholder to insert next migration", -1)
	err = ioutil.WriteFile(pathToBaseMigrationsFile, []byte(newContents), 0755)
	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"Created migration for blueprint %s with name %s\n",
		command.opts.Blueprint,
		command.opts.Message,
	)
}

func (command CreateMigration) GetHelpText() string {
	return "Create migration for your blueprint"
}

type UpMigrationOptions struct {
}

type UpMigration struct {
	opts *UpMigrationOptions
}

func (command UpMigration) ParseArgs() {
	parser := flags.NewParser(command.opts, flags.Default)
	_, err := parser.ParseArgs(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
}

func (command UpMigration) Proceed() {

}

func (command UpMigration) GetHelpText() string {
	return "Upgrade your database"
}

type DownMigrationOptions struct {
}

type DownMigration struct {
	opts *DownMigrationOptions
}

func (command DownMigration) ParseArgs() {
	parser := flags.NewParser(command.opts, flags.Default)
	_, err := parser.ParseArgs(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
}

func (command DownMigration) Proceed() {

}

func (command DownMigration) GetHelpText() string {
	return "Downgrade your database"
}
